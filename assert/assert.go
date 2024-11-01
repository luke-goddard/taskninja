package assert

import (
	"fmt"
	"os"
	"reflect"
	"runtime/debug"

	"github.com/rs/zerolog/log"
)

var TaskNinjaSkipAssert = "TASK_NINJA_SKIP_ASSERTS" // #nosec G101

func True(truth bool, msg string, data ...any) {
	if !truth {
		failAssert(msg, data...)
	}
}

func Nil(item any, msg string, data ...any) {
	log.Info().Interface("item", item).Msg("Nil Check")
	if item == nil {
		return
	}

	log.Error().Interface("item", item).Msg("Is Nil Check Failed")
	failAssert(msg, data...)
}

func NotNil(item any, msg string, data ...any) {
	if item == nil || reflect.ValueOf(item).Kind() == reflect.Ptr && reflect.ValueOf(item).IsNil() {
		log.Error().Interface("item", item).Msg("Not Nil Check Failed")
		failAssert(msg, data...)
	}
}

func Fail(msg string, data ...any) {
	failAssert(msg, data...)
}

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
		return
	}
	os.Exit(1)
}
