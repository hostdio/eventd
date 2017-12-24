package api

import (
	"testing"
	"time"
)

func Test_publishEventValidator_works(t *testing.T) {
	e := PublishEvent{
		ID: []byte("ID"),
		Type: []byte("eventType"),
		Version: []byte("1"),
		Timestamp: time.Now(),
		Source: []byte("self"),
	}
	if err := e.Validate(); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_publishEventValidator_missingID(t *testing.T) {
	e := PublishEvent{
		Type: []byte("eventType"),
		Version: []byte("1"),
		Timestamp: time.Now(),
		Source: []byte("self"),
	}
	if err := e.Validate(); err == nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_publishEventValidator_missingType(t *testing.T) {
	e := PublishEvent{
		ID: []byte("ID"),
		Version: []byte("1"),
		Timestamp: time.Now(),
		Source: []byte("self"),
	}
	if err := e.Validate(); err == nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_publishEventValidator_missingVersion(t *testing.T) {
	e := PublishEvent{
		Type: []byte("eventType"),
		ID: []byte("ID"),
		Timestamp: time.Now(),
		Source: []byte("self"),
	}
	if err := e.Validate(); err == nil {
		t.Log(err)
		t.Fail()
	}
}