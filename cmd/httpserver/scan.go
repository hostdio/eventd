package httpserver

import (
	"net/http"
	"time"
	"github.com/hostdio/eventd/api"
	"encoding/json"
	"errors"
	"log"
	"strconv"
)

var (
	defaultLimit = 10
)

type scanResponse struct {
	Events []api.PersistedEvent
}

func (resp scanResponse) JSON() []byte {
	byt, err :=json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	return byt
}

var missingTimestamp  = errors.New("Timestamp is missing. Please try to specify timestamp using query parameter \"from\"")

func scanHandler(scanner api.Scanner) http.HandlerFunc {
	return SimpleHandler(func(byt []byte, r*http.Request) (*Response, error) {
		val := r.URL.Query()
		fromStr := val.Get("from")

		if fromStr == "" {
			return nil, missingTimestamp
		}

		from, parseErr := time.Parse(time.RFC3339, fromStr)
		if parseErr != nil {
			return nil, parseErr
		}

		limit := defaultLimit
		limitStr := val.Get("limit")
		if limitStr != "" {
			parseLimit, err := strconv.Atoi(limitStr)
			if err != nil {
				return nil, errors.New("Illegal limit value")
			}
			limit = parseLimit
		}

		persistedEvents, scanErr := scanner.Scan(r.Context(), from, limit)

		if scanErr != nil {
			return nil, scanErr
		}

		resp := scanResponse{Events:persistedEvents}

		return &Response{
			Payload: resp.JSON(),
			StatusCode: 200,
		}, nil

	}, func(err error) ErrorPayload {
		if err == missingTimestamp {
			return ErrorPayload{
				Status: 400,
				Message: "Something was wrong with the request",
				DeveloperMessage: missingTimestamp.Error(),
			}
		}
		log.Println(err)
		return ErrorPayload{
			Status: 500,
			Message: "An unknown error occured",
			DeveloperMessage: "An unknown error occured. Please check with the administrator what went wrong",
		}
	})
}

