package service

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/qredo-external/go-rnov/pkg/auth"
	"github.com/qredo-external/go-rnov/pkg/storage"
	"github.com/qredo-external/go-rnov/pkg/user"
)

type authOperationsMock struct {
	createJWT   func(usr user.User) (string, error)
	validateJWT func(JWT string) error
}

func (a authOperationsMock) CreateJWT(usr user.User) (string, error) {
	if a.createJWT != nil {
		return a.createJWT(usr)
	}
	panic("implement me")
}

func (a authOperationsMock) ValidateJWT(JWT string) error {
	if a.validateJWT != nil {
		return a.validateJWT(JWT)
	}
	panic("implement me")
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

func TestOperationManager_Sum(t *testing.T) {
	tests := []struct {
		input        []byte
		expectedHash string
	}{
		{
			input:        json.RawMessage(`[1,2,3,4]`),
			expectedHash: "4a44dc15364204a80fe80e9039455cc1608281820fe2b24f1e5233ade6af1dd5",
		},
		{
			input:        json.RawMessage(`{"a":6,"b":4}`),
			expectedHash: "4a44dc15364204a80fe80e9039455cc1608281820fe2b24f1e5233ade6af1dd5",
		},
		{
			input:        json.RawMessage(`[[[2]]]`),
			expectedHash: "d4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35",
		},
		{
			input:        json.RawMessage(`{"a":{"b":4},"c":-2}`),
			expectedHash: "d4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35",
		},
		{
			input:        json.RawMessage(`{"a":[-1,1,"dark"]}`),
			expectedHash: "5feceb66ffc86f38d952786c6d696c79c2dbc239dd4e91b46729d73a27fb57e9",
		},
		{
			input:        json.RawMessage(`[-1,{"a":1, "b":"light"}]`),
			expectedHash: "5feceb66ffc86f38d952786c6d696c79c2dbc239dd4e91b46729d73a27fb57e9",
		},
		{
			input:        json.RawMessage(`{}`),
			expectedHash: "5feceb66ffc86f38d952786c6d696c79c2dbc239dd4e91b46729d73a27fb57e9",
		},
		{
			input:        json.RawMessage(`[]`),
			expectedHash: "5feceb66ffc86f38d952786c6d696c79c2dbc239dd4e91b46729d73a27fb57e9",
		},
	}

	var jsonMap interface{}
	for i, test := range tests {
		t.Run(fmt.Sprintf("generic test %d", i+1), func(t *testing.T) {
			if err := json.Unmarshal(test.input, &jsonMap); err != nil {
				t.Errorf("error unmarshaling body: %s", err.Error())
			}
			os := NewOperationsService()
			res, err := os.Sum(jsonMap)
			if err != nil {
				t.Errorf("non-nil error : %s", err.Error())
			}
			if test.expectedHash != res {
				t.Errorf("error in expectedHash value: expected %s got %s", test.expectedHash, res)
			}
		})
	}
}

func Test_getSum(t *testing.T) {
	tests := []struct {
		input          []byte
		ExpectedResult int
	}{
		{
			input: json.RawMessage(`[1,2,3,4]`), ExpectedResult: 10,
		},
		{
			input: json.RawMessage(`{"a":6,"b":4}`), ExpectedResult: 10,
		},
		{
			input: json.RawMessage(`[[[2]]]`), ExpectedResult: 2,
		},
		{
			input: json.RawMessage(`{"a":{"b":4},"c":-2}`), ExpectedResult: 2,
		},
		{
			input: json.RawMessage(`{"a":[-1,1,"dark"]}`), ExpectedResult: 0,
		},
		{
			input: json.RawMessage(`[-1,{"a":1, "b":"light"}]`), ExpectedResult: 0,
		},
		{
			input: json.RawMessage(`{}`), ExpectedResult: 0,
		},
		{
			input: json.RawMessage(`[]`), ExpectedResult: 0,
		},
	}

	var jsonMap interface{}
	for i, test := range tests {
		t.Run(fmt.Sprintf("generic test %d", i+1), func(t *testing.T) {
			if err := json.Unmarshal(test.input, &jsonMap); err != nil {
				t.Errorf("error unmarshaling body: %s", err.Error())
			}
			res := getSum(jsonMap)
			if test.ExpectedResult != res {
				t.Errorf("error in expectedHash value: expected %d got %d", test.ExpectedResult, res)
			}
		})
	}
}

func TestAuthManager_CreateAuth(t *testing.T) {
	tests := []struct {
		name           string
		usr            user.User
		authOp         auth.Operations
		storage        storage.ManageUsers
		expectedResult string
		expectedErr    error
	}{
		{
			name: "successful auth create",
			usr: user.User{
				UserName: "qwerty",
				Password: "mnbvc",
			},
			authOp: authOperationsMock{
				createJWT: func(usr user.User) (string, error) {
					return "aValidToken", nil
				},
			},
			storage: manageUsersMock{
				addUserToken: func(token string) bool {
					return true
				},
			},
			expectedResult: "aValidToken",
		},
		{
			name: "error - invalid user input",
			usr: user.User{
				UserName: "qwerty",
				Password: "",
			},
			expectedErr: fmt.Errorf("error validating user: invalid user data, empty fields"),
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("generic test %d", i+1), func(t *testing.T) {
			ah := NewAuthService(test.authOp, test.storage)
			res, err := ah.CreateAuth(test.usr)
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("expected: '%s' instead got: '%s'", test.expectedErr, err)
			}
			if test.expectedResult != res {
				t.Errorf("error in expectedHash value: expected %s got %s", test.expectedResult, res)
			}
		})
	}
}
