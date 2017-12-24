package httpserver

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func jsonSeralizer(p interface{}) []byte {
	byt, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return byt
}

// Serializer defines the interface of implementing seralizers
type Serializer func(p interface{}) []byte

var (
	SerializerJSON = jsonSeralizer
)

type ErrorPayload struct {
	Status int
	Message string
	DeveloperMessage string
}

func (p ErrorPayload) Serialize(seralizer Serializer) []byte {
	return seralizer(p)
}

type Response struct {
	Payload []byte
	StatusCode int
}

func SimpleHandler(
	handler func([]byte, *http.Request) (*Response, error),
	errorHandler func(error) ErrorPayload) http.HandlerFunc {
		errors := func(err error, w http.ResponseWriter) {
			payload := errorHandler(err)
			w.WriteHeader(payload.Status)
			output := payload.Serialize(SerializerJSON)
			w.Write(output)
		}
	return func(w http.ResponseWriter, r *http.Request) {
		byt, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errors(err, w)
			return
		}
		resp, err := handler(byt, r)
		if err != nil {
			errors(err, w)
			return
		}

		w.WriteHeader(resp.StatusCode)
		w.Write(resp.Payload)
	}
}
