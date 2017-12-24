package httpserver

import (
	"net/http"

	"github.com/hostdio/eventd/api"
)

func publishHandler(publisher api.Publisher) http.HandlerFunc {
	return SimpleHandler(func(byt []byte, r*http.Request) (Response, error){
		return Response{
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
