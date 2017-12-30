package api

import (
	"testing"
	"time"
)

func Test_publishEventValidator_works(t *testing.T) {
	base := BaseEvent{
		ID:        "ID",
		Type:      "eventType",
		Version:   "1",
		Timestamp: time.Now(),
		Source:    "self",
	}
	e := PublishEvent{
		BaseEvent: &base,
	}
	if err := e.Validate(); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_publishEventValidator_missingID(t *testing.T) {
	base := BaseEvent{
		Type:      "eventType",
		Version:   "1",
		Timestamp: time.Now(),
		Source:    "self",
	}
	e := PublishEvent{
		BaseEvent: &base,
	}
	if err := e.Validate(); err == nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_publishEventValidator_missingType(t *testing.T) {
	base := BaseEvent{
		ID:        "ID",
		Version:   "1",
		Timestamp: time.Now(),
		Source:    "self",
	}
	e := PublishEvent{
		BaseEvent: &base,
	}
	if err := e.Validate(); err == nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_publishEventValidator_missingVersion(t *testing.T) {
	base := BaseEvent{
		Type:      "eventType",
		ID:        "ID",
		Timestamp: time.Now(),
		Source:    "self",
	}
	e := PublishEvent{
		BaseEvent: &base,
	}
	if err := e.Validate(); err == nil {
		t.Log(err)
		t.Fail()
	}
}
