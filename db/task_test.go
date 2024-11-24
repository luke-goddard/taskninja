package db

import (
	"fmt"
	"testing"
	"time"
)

func TestPrettyAge(t *testing.T) {
	var task = Task{}
	var payloads = []struct {
		Time   time.Duration
		Expect string
	}{
		{Time: time.Duration(0) * time.Second, Expect: "0m"},
		{Time: time.Duration(1) * time.Second, Expect: "0m"},
		{Time: time.Duration(60) * time.Second, Expect: "1m"},
		{Time: time.Duration(61) * time.Second, Expect: "1m"},
		{Time: time.Duration(60) * time.Minute, Expect: "1h0m"},
		{Time: time.Duration(61) * time.Minute, Expect: "1h1m"},
		{Time: time.Duration(24) * time.Hour, Expect: "1d0h"},
		{Time: time.Duration(25) * time.Hour, Expect: "1d1h"},
		{Time: time.Duration(28) * time.Hour, Expect: "1d4h"},
		{Time: time.Duration(24) * time.Hour * 7, Expect: "1w0d"},
		{Time: time.Duration(1) * time.Hour * 526, Expect: "3w3d"},
	}

	for _, payload := range payloads {
		var name = fmt.Sprintf("pretty-age-%s", payload.Expect)
		t.Run(name, func(t *testing.T) {
			var result = task.PrettyAge(payload.Time)
			if result != payload.Expect {
				t.Errorf("Expected %s, got %s", payload.Expect, result)
			}
		})
	}
}
