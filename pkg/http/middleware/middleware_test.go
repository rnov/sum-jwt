package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/qredo-external/go-rnov/pkg/user"
)

type authOpMock struct {
	validateJWT func(ba string) error
}

func (am *authOpMock) CreateJWT(usr user.User) (string, error) {
	panic("Not implemented")
}

func (am *authOpMock) ValidateJWT(JWT string) error {
	if am.validateJWT != nil {
		return am.validateJWT(JWT)
	}
	panic("Not implemented")
}

type manageUsersMock struct {
	addUserToken  func(token string) bool
	isActiveToken func(token string) bool
}

func (m manageUsersMock) AddUserToken(token string) bool {
	if m.addUserToken != nil {
		return m.addUserToken(token)
	}
	panic("implement me")
}

func (m manageUsersMock) IsActiveToken(token string) bool {
	if m.isActiveToken != nil {
		return m.isActiveToken(token)
	}
	panic("implement me")
}

func TestAuthentication(t *testing.T) {
	tests := []struct {
		name           string
		auth           AuthMiddleware
		AuthHeader     bool
		Auth           string
		next           func(w http.ResponseWriter, r *http.Request)
		expectedStatus int
	}{
		{
			name: "successful validation",
			auth: AuthMiddleware{
				Operations: &authOpMock{
					validateJWT: func(ba string) error {
						return nil
					},
				},
				storage: manageUsersMock{
					isActiveToken: func(token string) bool {
						return true
					},
				},
			},
			AuthHeader:     true,
			Auth:           "Basic dXNlcm5hbWU6cGFzc3dvcmQ=",
			next:           func(w http.ResponseWriter, r *http.Request) {},
			expectedStatus: 200,
		},
		{
			name:           "error - not valid auth structure",
			AuthHeader:     true,
			Auth:           "notValidAuth123",
			expectedStatus: 401,
		},
		{
			name:           "error - missing auth header",
			expectedStatus: 401,
		},
		{
			name: "error - non existent auth",
			Auth: "basic dXNlcm5hbWU6cGFzc3dvcmQ=",
			auth: AuthMiddleware{
				Operations: &authOpMock{
					validateJWT: func(ba string) error {
						return errors.New("not valid auth")
					},
				},
			},
			AuthHeader:     true,
			expectedStatus: 401,
		},
		{
			name: "error - deactivated token",
			Auth: "basic dXNlcm5hbWU6cGFzc3dvcmQ=",
			auth: AuthMiddleware{
				Operations: &authOpMock{
					validateJWT: func(ba string) error {
						return nil
					},
				},
				storage: manageUsersMock{
					isActiveToken: func(token string) bool {
						return false
					},
				},
			},
			AuthHeader:     true,
			expectedStatus: 401,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/sum", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add(authHeader, test.Auth)

			// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			servicesRouter := mux.NewRouter()
			servicesRouter.HandleFunc("/sum", Authentication(test.auth, test.next)).Methods("POST")

			// Handlers satisfy http.Handler, so it allows to call ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			servicesRouter.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if rr.Code != test.expectedStatus {
				t.Errorf("handler returned wrong status code: expected %v got %v", test.expectedStatus, rr.Code)
			}
		})
	}
}
