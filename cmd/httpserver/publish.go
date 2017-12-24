package httpserver

import (
	"net/http"

	"github.com/hostdio/eventd/api"
	"encoding/json"
	"log"
	"context"
	"time"
	"errors"
)

func publishHandler(publisher api.Publisher) http.HandlerFunc {
	return SimpleHandler(func(byt []byte, r*http.Request) (*Response, error){
		var event api.PublishEvent
		if err := json.Unmarshal(byt, &event); err != nil {
			return nil, err
		}
		if err := event.Validate(); err != nil {
			log.Println(err)
			return nil, err
		}
		
		ctx, cancel := context.WithTimeout(r.Context(), 10 * time.Second)
		defer cancel()
		if _, err := publisher.Publish(ctx, event); err != nil {
			return nil, errors.New("Publishing event timedout. Please try again")
		}
		return &Response{
			Payload: []byte(`{"status":"Payload received"}`),
			StatusCode: 201,
		}, nil
	}, func(error) ErrorPayload {
		return ErrorPayload{
			Status: 500,
			Message: "An unknown error occured",
			DeveloperMessage: "An unknown error occured. Please check with the administrator what went wrong",
		}
	})
}
