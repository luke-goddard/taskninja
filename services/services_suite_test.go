package services_test

import (
	"testing"

	"github.com/luke-goddard/taskninja/db"
	"github.com/luke-goddard/taskninja/interpreter"
	"github.com/luke-goddard/taskninja/services"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestServices(t *testing.T) {
	log.Logger = log.Output(zerolog.Nop())
	RegisterFailHandler(Fail)
	RunSpecs(t, "Services Suite")
}

func newTestHandler() *services.ServiceHandler {
	var store = db.NewInMemoryStore()
	var interpreter = interpreter.NewInterpreter(store)
	return services.NewServiceHandler(interpreter, store)
}

// ============================================================================
// TASK COMPLETE
// ============================================================================
var _ = Describe("Completing an incomplete task", func() {
	var services *services.ServiceHandler
	var task *db.Task
	var completed bool
	var err error

	BeforeEach(func() {
		services = newTestHandler()
		task = &db.Task{
			Title: "title",
			State: db.TaskStateIncomplete,
		}
		task, _ = services.CreateTask(task)
	})
	It("Should have no time tracking left running", func() {
		var times, err = services.GetTaskTimes(task.ID)
		Expect(err).To(BeNil())
		Expect(times).To(HaveLen(0))
	})
	Context("When the task is completed", func() {
		BeforeEach(func() {
			task, err = services.GetTaskById(task.ID)
			services.StartTimeToggleById(task.ID)
			completed, err = services.CompleteTaskById(task.ID)
			task, err = services.GetTaskById(task.ID)
		})
		It("should complete a task", func() {
			Expect(err).To(BeNil())
			Expect(completed).To(BeTrue())
			Expect(err).To(BeNil())
			Expect(task.State).To(Equal(db.TaskStateCompleted))
		})
		It("should have a completed date", func() {
			Expect(task.CompletedUtc.Valid).To(BeTrue())
		})
		It("Should have no time tracking left running", func() {
			var times, err = services.GetTaskTimes(task.ID)
			Expect(err).To(BeNil())
			Expect(times).To(HaveLen(1))
			var time = times[0]
			Expect(time.EndTimeUtc.Valid).To(BeTrue())
		})
	})
	Context("When the task is blocking another task", func() {
		BeforeEach(func() {
			services.RunProgram("add hello")
			services.RunProgram("depends 1 on 2")
		})
		It("Should delete the task dependency", func() {
			var deps []db.TaskDependency
			deps, err = services.GetDependenciesForServices(1)
			Expect(err).To(BeNil())
			Expect(deps).To(HaveLen(1))

			var _, err = services.CompleteTaskById(2)
			Expect(err).To(BeNil())

			deps, err = services.GetDependenciesForServices(1)
			Expect(err).To(BeNil())
			Expect(deps).To(HaveLen(0))
		})
	})
})

// ============================================================================
// TASK DELETE
// ============================================================================
var _ = Describe("Delete Task", func() {
	var services *services.ServiceHandler
	var task *db.Task
	Context("When the task is not being tracked", func() {
		BeforeEach(func() {
			services = newTestHandler()
			task, _ = services.CreateTask(&db.Task{Title: "title"})
		})
		It("should delete a task", func() {
			var deleted, err = services.DeleteTaskById(task.ID)
			Expect(deleted).To(BeTrue())
			Expect(err).To(BeNil())
		})
		It("should not find the task", func() {
			_, _ = services.DeleteTaskById(task.ID)
			var _, err = services.GetTaskById(task.ID)
			Expect(err).ToNot(BeNil())
		})
	})
	Context("When the task is being tracked", func() {
		BeforeEach(func() {
			services = newTestHandler()
			task, _ = services.CreateTask(&db.Task{Title: "title"})
			services.StartTimeToggleById(task.ID)
		})
		It("should delete a task", func() {
			var deleted, err = services.DeleteTaskById(task.ID)
			Expect(deleted).To(BeTrue())
			Expect(err).To(BeNil())
		})
		It("should not find the task", func() {
			_, _ = services.DeleteTaskById(task.ID)
			var _, err = services.GetTaskById(task.ID)
			Expect(err).ToNot(BeNil())
		})
		It("should not find the task in the tracking", func() {
			_, _ = services.DeleteTaskById(task.ID)
			var times, err = services.GetTaskTimes(task.ID)
			Expect(err).To(BeNil())
			Expect(times).To(BeEmpty())
			Expect(len(times)).To(Equal(0))
		})
	})
	Context("When the task has a project", func() {
		BeforeEach(func() {
			services = newTestHandler()
			services.RunProgram("add test project:home")
		})
		It("should delete a task project link", func() {
			var tasks, err = services.ListTasks()
			Expect(err).To(BeNil())
			Expect(tasks).To(HaveLen(1))
			Expect(tasks[0].ProjectNames.Value()).To(Equal("home"))
			Expect(services.Store.ListProjects()).To(HaveLen(1))
			Expect(services.Store.ProjectTasksList()).To(HaveLen(1))

			services.DeleteTaskById(tasks[0].ID)
			tasks, err = services.ListTasks()
			Expect(err).To(BeNil())
			Expect(services.Store.ProjectTasksList()).To(HaveLen(0))
		})
	})
})

// ============================================================================
// TASK CREATE
// ============================================================================
var _ = Describe("Creating a task", func() {
	var services *services.ServiceHandler
	var task *db.Task
	var err error
	BeforeEach(func() {
		services = newTestHandler()
		task, err = services.CreateTask(&db.Task{Title: "title"})
	})
	It("should create a task", func() {
		Expect(err).To(BeNil())
		Expect(task).ToNot(BeNil())
		Expect(task.Title).To(Equal("title"))
	})
	It("should have a default priority", func() {
		Expect(task.Priority).To(Equal(db.TaskPriorityNone))
	})
	It("should have a default state", func() {
		Expect(task.State).To(Equal(db.TaskStateIncomplete))
	})
})

// ============================================================================
// TASK LIST
// ============================================================================
var _ = Describe("Listing tasks", func() {
	var services *services.ServiceHandler
	var tasks []db.TaskDetailed
	var err error

	var getById = func(id int64) *db.TaskDetailed {
		for _, task := range tasks {
			if task.ID == id {
				return &task
			}
		}
		return nil

	}

	BeforeEach(func() {
		services = newTestHandler()
		_, _ = services.CreateTask(&db.Task{Title: "title", Priority: db.TaskPriorityHigh})
		_, _ = services.CreateTask(&db.Task{Title: "title", Priority: db.TaskPriorityLow})
		tasks, err = services.ListTasks()
	})
	It("should list tasks", func() {
		Expect(err).To(BeNil())
		Expect(tasks).To(HaveLen(2))
	})
	It("should sort tasks by urgency", func() {
		Expect(tasks[0].Priority).To(Equal(db.TaskPriorityHigh))
		Expect(tasks[1].Priority).To(Equal(db.TaskPriorityLow))
	})
	It("shoudld show dependencies", func() {
		services.RunProgram("depends 1 on 2")
		tasks, err = services.ListTasks()
		Expect(err).To(BeNil())
		Expect(getById(1).Dependencies.String).To(Equal("2"))
	})
	It("should show if the task is blocked", func() {
		services.RunProgram("depends 1 on 2")
		tasks, err = services.ListTasks()
		Expect(err).To(BeNil())
		Expect(getById(1).Blocked).To(BeTrue())
		Expect(getById(2).Blocked).To(BeFalse())
	})
	It("should show how many tasks as blocking", func() {
		services.RunProgram("depends 1 on 2")
		tasks, err = services.ListTasks()
		Expect(err).To(BeNil())
		Expect(getById(1).Blocking).To(Equal(0))
		Expect(getById(2).Blocking).To(Equal(1))
	})
})

// ============================================================================
// TASK COUNT
// ============================================================================
var _ = Describe("Counting tasks", func() {
	var services *services.ServiceHandler
	var count int64
	var err error
	BeforeEach(func() {
		services = newTestHandler()
		_, _ = services.CreateTask(&db.Task{Title: "title"})
		_, _ = services.CreateTask(&db.Task{Title: "title"})
		_, _ = services.CreateTask(&db.Task{Title: "title"})
		count, err = services.CountTasks()
	})
	It("should count tasks", func() {
		Expect(err).To(BeNil())
		Expect(count).To(Equal(int64(3)))
	})
})

// ============================================================================
// TASK GET BY ID
// ============================================================================
var _ = Describe("Getting a task by ID", func() {
	var services *services.ServiceHandler
	var task *db.Task
	var err error
	BeforeEach(func() {
		services = newTestHandler()
		task, _ = services.CreateTask(&db.Task{Title: "title"})
	})
	It("should get a task by ID If it exists", func() {
		task, err = services.GetTaskById(task.ID)
		Expect(err).To(BeNil())
		Expect(task).ToNot(BeNil())
	})
	It("should not get a task by ID If it does not exist", func() {
		task, err = services.GetTaskById(99)
		Expect(err).ToNot(BeNil())
		Expect(task).To(BeNil())
	})
})

// ============================================================================
// TASK INCREASE PRIORITY
// ============================================================================
var _ = Describe("Increasing the priority of a task", func() {
	var services *services.ServiceHandler
	var task *db.Task
	var err error
	Context("When the task has no priority", func() {
		BeforeEach(func() {
			services = newTestHandler()
			task, _ = services.CreateTask(&db.Task{Title: "title", Priority: db.TaskPriorityNone})
		})
		It("should increase the priority of a task to LOW", func() {
			services.IncreasePriority(task.ID)
			task, err = services.GetTaskById(task.ID)
			Expect(err).To(BeNil())
			Expect(task.Priority).To(Equal(db.TaskPriorityLow))
		})
	})
	Context("When the task has a LOW priority", func() {
		BeforeEach(func() {
			services = newTestHandler()
			task, _ = services.CreateTask(&db.Task{Title: "title", Priority: db.TaskPriorityLow})
		})
		It("should increase the priority of a task to MEDIUM", func() {
			services.IncreasePriority(task.ID)
			task, err = services.GetTaskById(task.ID)
			Expect(err).To(BeNil())
			Expect(task.Priority).To(Equal(db.TaskPriorityMedium))
		})
	})
	Context("When the task has a MEDIUM priority", func() {
		BeforeEach(func() {
			services = newTestHandler()
			task, _ = services.CreateTask(&db.Task{Title: "title", Priority: db.TaskPriorityMedium})
		})
		It("should increase the priority of a task to HIGH", func() {
			services.IncreasePriority(task.ID)
			task, err = services.GetTaskById(task.ID)
			Expect(err).To(BeNil())
			Expect(task.Priority).To(Equal(db.TaskPriorityHigh))
		})
	})
	Context("When the task has a HIGH priority", func() {
		BeforeEach(func() {
			services = newTestHandler()
			task, _ = services.CreateTask(&db.Task{Title: "title", Priority: db.TaskPriorityHigh})
		})
		It("should not increase the priority of a task", func() {
			services.IncreasePriority(task.ID)
			task, err = services.GetTaskById(task.ID)
			Expect(err).To(BeNil())
			Expect(task.Priority).To(Equal(db.TaskPriorityHigh))
		})
	})
})

// ============================================================================
// TASK DECREASE PRIORITY
// ============================================================================
var _ = Describe("Decreasing the priority of a task", func() {
	var services *services.ServiceHandler
	var task *db.Task
	var err error
	Context("When the task has no priority", func() {
		BeforeEach(func() {
			services = newTestHandler()
			task, _ = services.CreateTask(&db.Task{Title: "title", Priority: db.TaskPriorityNone})
		})
		It("should not decrease the priority of a task", func() {
			services.DecreasePriority(task.ID)
			task, err = services.GetTaskById(task.ID)
			Expect(err).To(BeNil())
			Expect(task.Priority).To(Equal(db.TaskPriorityNone))
		})
	})
	Context("When the task has a LOW priority", func() {
		BeforeEach(func() {
			services = newTestHandler()
			task, _ = services.CreateTask(&db.Task{Title: "title", Priority: db.TaskPriorityLow})
		})
		It("should decrease the priority of a task to NONE", func() {
			services.DecreasePriority(task.ID)
			task, err = services.GetTaskById(task.ID)
			Expect(err).To(BeNil())
			Expect(task.Priority).To(Equal(db.TaskPriorityNone))
		})
	})
	Context("When the task has a MEDIUM priority", func() {
		BeforeEach(func() {
			services = newTestHandler()
			task, _ = services.CreateTask(&db.Task{Title: "title", Priority: db.TaskPriorityMedium})
		})
		It("should decrease the priority of a task to LOW", func() {
			services.DecreasePriority(task.ID)
			task, err = services.GetTaskById(task.ID)
			Expect(err).To(BeNil())
			Expect(task.Priority).To(Equal(db.TaskPriorityLow))
		})
	})
	Context("When the task has a HIGH priority", func() {
		BeforeEach(func() {
			services = newTestHandler()
			task, _ = services.CreateTask(&db.Task{Title: "title", Priority: db.TaskPriorityHigh})
		})
		It("should decrease the priority of a task to MEDIUM", func() {
			services.DecreasePriority(task.ID)
			task, err = services.GetTaskById(task.ID)
			Expect(err).To(BeNil())
			Expect(task.Priority).To(Equal(db.TaskPriorityMedium))
		})
	})
})

// ============================================================================
// TASK SET PRIORITY
// ============================================================================
var _ = Describe("Setting the priority of a task", func() {
	var services *services.ServiceHandler
	var task *db.Task
	var err error
	Context("When the task has a HIGH priority", func() {
		BeforeEach(func() {
			services = newTestHandler()
			task, _ = services.CreateTask(&db.Task{Title: "title", Priority: db.TaskPriorityHigh})
		})
		It("should set the priority of a task to NONE", func() {
			services.SetPriority(task.ID, db.TaskPriorityNone)
			task, err = services.GetTaskById(task.ID)
			Expect(err).To(BeNil())
			Expect(task.Priority).To(Equal(db.TaskPriorityNone))
		})
		It("should set the priority of a task to MED", func() {
			services.SetPriority(task.ID, db.TaskPriorityMedium)
			task, err = services.GetTaskById(task.ID)
			Expect(err).To(BeNil())
			Expect(task.Priority).To(Equal(db.TaskPriorityMedium))
		})
		It("should set the priority of a task to LOW", func() {
			services.SetPriority(task.ID, db.TaskPriorityLow)
			task, err = services.GetTaskById(task.ID)
			Expect(err).To(BeNil())
			Expect(task.Priority).To(Equal(db.TaskPriorityLow))
		})
	})
})

// ============================================================================
// RUN PROGRAM
// ============================================================================
var _ = Describe("Running a program", func() {
	var services *services.ServiceHandler
	var err error
	BeforeEach(func() {
		services = newTestHandler()
	})
	It("should run a program", func() {
		_, err = services.RunProgram(`add "title"`)
		Expect(err).To(BeNil())
		var tasks, err = services.ListTasks()
		Expect(err).To(BeNil())
		Expect(tasks).To(HaveLen(1))
	})
})

// ============================================================================
// TASK START
// ============================================================================
var _ = Describe("Starting a task", func() {
	var services *services.ServiceHandler
	var task *db.Task
	var err error
	BeforeEach(func() {
		services = newTestHandler()
		task, _ = services.CreateTask(&db.Task{Title: "title", State: db.TaskStateIncomplete})
	})
	Context("When the task is not started", func() {
		BeforeEach(func() {
			err = services.StartTimeToggleById(task.ID)
			Expect(err).To(BeNil())
		})
		It("should start a task", func() {
			task, err = services.GetTaskById(task.ID)
			Expect(err).To(BeNil())
			Expect(task.State).To(Equal(db.TaskStateStarted))
		})
		It("should have a time tracking", func() {
			var times, err = services.GetTaskTimes(task.ID)
			Expect(err).To(BeNil())
			Expect(times).To(HaveLen(1))
			var time = times[0]
			Expect(time.EndTimeUtc.Valid).To(BeFalse())
			Expect(time.StartTimeUtc).ToNot(BeEmpty())
		})
	})
	Context("When the task is already started", func() {
		BeforeEach(func() {
			err = services.StartTimeToggleById(task.ID)
			Expect(err).To(BeNil())

			err = services.StartTimeToggleById(task.ID)
			Expect(err).To(BeNil())
		})
		It("should not start a task", func() {
			task, err = services.GetTaskById(task.ID)
			Expect(err).To(BeNil())
			Expect(task.State).To(Equal(db.TaskStateStarted))
		})
		It("should not stop tracking time", func() {
			var times, err = services.GetTaskTimes(task.ID)
			Expect(err).To(BeNil())
			Expect(times).To(HaveLen(1))
			var time = times[0]
			Expect(time.EndTimeUtc.Valid).To(BeFalse())
		})
	})
})

// ============================================================================
// TASK STOP
// ============================================================================
var _ = Describe("Stopping a task", func() {
	var services *services.ServiceHandler
	var task *db.Task
	var err error
	BeforeEach(func() {
		services = newTestHandler()
		task, _ = services.CreateTask(&db.Task{Title: "title", State: db.TaskStateIncomplete})
	})
	Context("When the task is not started", func() {
		BeforeEach(func() {
			err = services.StopTimeToggleById(task.ID)
			Expect(err).To(BeNil())
		})
		It("should not change the task state", func() {
			task, err = services.GetTaskById(task.ID)
			Expect(err).To(BeNil())
			Expect(task.State).To(Equal(db.TaskStateIncomplete))
		})
		It("should not create a time tracking", func() {
			var times, err = services.GetTaskTimes(task.ID)
			Expect(err).To(BeNil())
			Expect(times).To(HaveLen(0))
		})
	})
	Context("When task is started", func() {
		BeforeEach(func() {
			err = services.StartTimeToggleById(task.ID)
			Expect(err).To(BeNil())

			err = services.StopTimeToggleById(task.ID)
		})
		It("should stop a task", func() {
			task, err = services.GetTaskById(task.ID)
			Expect(err).To(BeNil())
			Expect(task.State).To(Equal(db.TaskStateIncomplete))
		})
		It("should have a time tracking row with a end time", func() {
			var times, err = services.GetTaskTimes(task.ID)
			Expect(err).To(BeNil())
			Expect(times).To(HaveLen(1))
			var time = times[0]
			Expect(time.EndTimeUtc.Valid).To(BeTrue())
		})
	})
})

// ============================================================================
// TASK GET TIMES
// ============================================================================
var _ = Describe("Getting task times", func() {
	var services *services.ServiceHandler
	var task *db.Task
	var times []db.TaskTime
	var err error
	BeforeEach(func() {
		services = newTestHandler()
		task, _ = services.CreateTask(&db.Task{Title: "title", State: db.TaskStateIncomplete})
	})
	Context("When the task has no times", func() {
		BeforeEach(func() {
			times, err = services.GetTaskTimes(task.ID)
		})
		It("should not have any times", func() {
			Expect(err).To(BeNil())
			Expect(times).To(HaveLen(0))
		})
	})
	Context("When the task has a single time that's inprogres", func() {
		BeforeEach(func() {
			err = services.StartTimeToggleById(task.ID)
			Expect(err).To(BeNil())
		})
		It("should have a single time", func() {
			times, err = services.GetTaskTimes(task.ID)
			Expect(err).To(BeNil())
			Expect(times).To(HaveLen(1))
		})
	})
	Context("When the task has a single time that's completed", func() {
		BeforeEach(func() {
			err = services.StartTimeToggleById(task.ID)
			Expect(err).To(BeNil())

			err = services.StopTimeToggleById(task.ID)
		})
		It("should have a single time", func() {
			times, err = services.GetTaskTimes(task.ID)
			Expect(err).To(BeNil())
			Expect(times).To(HaveLen(1))
		})
	})
	Context("When the task has multiple times", func() {
		BeforeEach(func() {
			err = services.StartTimeToggleById(task.ID)
			Expect(err).To(BeNil())

			err = services.StopTimeToggleById(task.ID)
			Expect(err).To(BeNil())

			err = services.StartTimeToggleById(task.ID)
			Expect(err).To(BeNil())

			err = services.StopTimeToggleById(task.ID)
			Expect(err).To(BeNil())
		})
		It("should have multiple times", func() {
			times, err = services.GetTaskTimes(task.ID)
			Expect(err).To(BeNil())
			Expect(times).To(HaveLen(2))
		})
	})
})

// ============================================================================
// TASK TAG - CREATE
// ============================================================================
var _ = Describe("Creating a task tag", func() {
	var services *services.ServiceHandler
	var tagId int64
	var err error
	BeforeEach(func() {
		services = newTestHandler()
		tagId, err = services.TagCreate("ExampleTag")
		Expect(err).To(BeNil())
	})
	It("should create a tag", func() {
		Expect(tagId).ToNot(BeZero())
		var tags, err = services.TagList()
		Expect(err).To(BeNil())
		Expect(tags).To(HaveLen(1))
		Expect(tags[0].Name).To(Equal("ExampleTag"))

	})

	It("should not create a duplicate tag", func() {
		var tagId2 int64
		tagId2, err = services.TagCreate("ExampleTag")
		Expect(err).ToNot(BeNil())
		Expect(tagId2).To(BeZero())
	})

	It("should be able to link a task to the tag and unlink it", func() {
		// Create
		var task, err = services.CreateTask(&db.Task{Title: "Example"})
		Expect(err).To(BeNil())

		// Link
		err = services.TagLinkTask(tagId, task.ID)
		Expect(err).To(BeNil())

		// Check Link
		var tasks, err2 = services.ListTasks()
		Expect(err2).To(BeNil())
		Expect(tasks).To(HaveLen(1))
		Expect(tasks[0].TagCount).To(Equal(1))
		Expect(tasks[0].TagNames.Value()).To(Equal("ExampleTag"))

		// Unlink
		err = services.Store.TagUnlinkTask(task.ID, tagId)
		Expect(err).To(BeNil())

		// Check unlink
		tasks, err2 = services.ListTasks()
		Expect(err2).To(BeNil())
		Expect(tasks).To(HaveLen(1))
		Expect(tasks[0].TagCount).To(Equal(0))
		Expect(tasks[0].TagNames.Valid).To(Equal(false))

	})
})
