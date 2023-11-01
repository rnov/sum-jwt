package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	response "github.com/qredo-external/go-rnov/pkg/http/json"
	"github.com/qredo-external/go-rnov/pkg/user"
)

type OperationServiceMock struct {
	sum func(data interface{}) (string, error)
}

func (osm OperationServiceMock) Sum(data interface{}) (string, error) {
	if osm.sum != nil {
		return osm.sum(data)
	}
	panic("Not implemented")
}

func TestNewOperationHandler(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		requestPayload []byte
		service        OperationServiceMock
		status         int
		expectedRes    response.Operation
	}{
		{
			name:           "Successful request",
			url:            "/sum",
			requestPayload: json.RawMessage(`[1,2,3,4]`),
			service: OperationServiceMock{
				sum: func(data interface{}) (string, error) {
					return "qwertyHasBeenHashed", nil
				},
			},
			status: 200,
			expectedRes: response.Operation{
				Result: "qwertyHasBeenHashed",
			},
		},
		{
			name:           "Successful request",
			url:            "/sum",
			requestPayload: json.RawMessage(`[1,2,3,4]`),
			service: OperationServiceMock{
				sum: func(data interface{}) (string, error) {
					return "", errors.New("error in sum")
				},
			},
			status: 500,
		},
		//{
		//	name:           "error special case incoming body is not a json - error unmarshal",
		//	requestPayload: json.RawMessage(`[1,2,3,4`),
		//	url:            "/sum",
		//	status:         400,
		//},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var jsonBody []byte
			jsonBody, _ = json.Marshal(&test.requestPayload)
			req, err := http.NewRequest("POST", test.url, bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatal(err)
			}

			rh := NewOperationHandler(&test.service)

			// We sum a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			servicesRouter := mux.NewRouter()
			servicesRouter.HandleFunc("/sum", rh.SumHandler).Methods("POST")

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			servicesRouter.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if rr.Code != test.status {
				t.Errorf("handler returned wrong status code: expectedRes %v got %v", test.status, rr.Code)
			}
			if rr.Body.Len() > 0 {
				res := response.Operation{}
				dec := json.NewDecoder(rr.Body)
				dec.DisallowUnknownFields()
				if err := dec.Decode(&res); err != nil {
					t.Error("unable to decode body response")
				}
				if res.Result != test.expectedRes.Result {
					t.Errorf("error expectedRes response %s got %s", test.expectedRes.Result, res.Result)
				}
			}
		})
	}
}

type AuthorizerServiceMock struct {
	createAuth func(usr user.User) (string, error)
}

func (asm AuthorizerServiceMock) CreateAuth(usr user.User) (string, error) {
	if asm.createAuth != nil {
		return asm.createAuth(usr)
	}
	panic("Not implemented")
}

func TestNewAuthHandler(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		requestPayload user.User
		service        AuthorizerServiceMock
		status         int
		expectedRes    string
	}{
		{
			name: "Successful request",
			url:  "/auth",
			requestPayload: user.User{
				UserName: "qwerty",
				Password: "z1x2c3",
			},
			service: AuthorizerServiceMock{
				createAuth: func(usr user.User) (string, error) {
					return "validAuthToken", nil
				},
			},
			status:      201,
			expectedRes: "validAuthToken",
		},
		{
			name: "error - empty password",
			url:  "/auth",
			requestPayload: user.User{
				UserName: "qwerty",
				Password: "",
			},
			service: AuthorizerServiceMock{
				createAuth: func(usr user.User) (string, error) {
					return "", errors.New("error: invalid user")
				},
			},
			status: 400,
		},
		//{
		//	name:   "error special case incoming body is not a user - error unmarshal",
		//	url:    "/auth",
		//	status: 400,
		//},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var jsonBody []byte
			jsonBody, _ = json.Marshal(&test.requestPayload)
			req, err := http.NewRequest("POST", test.url, bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatal(err)
			}

			rh := NewAuthHandler(&test.service)

			// We sum a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			servicesRouter := mux.NewRouter()
			servicesRouter.HandleFunc("/auth", rh.CreateAuthHandler).Methods("POST")

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			servicesRouter.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if rr.Code != test.status {
				t.Errorf("handler returned wrong status code: expectedRes %v got %v", test.status, rr.Code)
			}
			if rr.Body.Len() > 0 {
				res := &response.JWT{}
				dec := json.NewDecoder(rr.Body)
				dec.DisallowUnknownFields()
				if err := dec.Decode(&res); err != nil {
					t.Error("unable to decode body response")
				}
				if res.JWT != test.expectedRes {
					t.Errorf("error expectedRes response %s got %s", test.expectedRes, res)
				}
			}
		})
	}
}
