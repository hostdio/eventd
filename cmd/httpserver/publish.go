package httpserver

import (
	"net/http"

	"github.com/hostdio/eventd/api"
	"encoding/json"
	"log"
)

func publishHandler(publisher api.Publisher) http.HandlerFunc {
	return SimpleHandler(func(byt []byte, r*http.Request) (*Response, error){
		var v api.PublishEvent
		if err := json.Unmarshal(byt, &v); err != nil {
			return nil, err
		}
		if err := v.Validate(); err != nil {
			log.Println(err)
			return nil, err
		}
		log.Println(v)
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
