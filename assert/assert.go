package assert

import (
	"fmt"
	"os"
	"reflect"
	"runtime/debug"

	"github.com/rs/zerolog/log"
)

// Used to disable asserts panicing
var TaskNinjaSkipAssert = "TASK_NINJA_SKIP_ASSERTS" // #nosec G101
var TotalPanics = 0                                 // Total number of panics

// True causes a panic if the truth is false
func True(truth bool, msg string, data ...any) {
	if !truth {
		failAssert(msg, data...)
	}
}

// Nil asserts that the item passed is nil
func Nil(item any, msg string, data ...any) {
	if item == nil {
		return
	}

	log.Error().Interface("item", item).Msg("Is Nil Check Failed")
	failAssert(msg, data...)
}

// NotNil asserts that the item passed is not nil
func NotNil(item any, msg string, data ...any) {
	if item == nil || reflect.ValueOf(item).Kind() == reflect.Ptr && reflect.ValueOf(item).IsNil() {
		log.Error().Interface("item", item).Msg("Not Nil Check Failed")
		failAssert(msg, data...)
	}
}

// Fail causes a panic with the message
func Fail(msg string, data ...any) {
	failAssert(msg, data...)
}

// NoError asserts that the error is nil
func NoError(err error, msg string, data ...any) {
	if err != nil {
		data = append(data, "error", err)
		failAssert(msg, data...)
	}
}

func failAssert(msg string, args ...interface{}) {
	var err = fmt.Errorf(msg, args...)
	log.Error().
		Str("stack", string(debug.Stack())).
		Err(err).
		Str("msg", msg).
		Interface("args", args).
		Msg("Assert Failed")

	if os.Getenv(TaskNinjaSkipAssert) == "true" {
		TotalPanics++
		return
	}
	os.Exit(1)
}
