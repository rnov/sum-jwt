package service

import (
	"crypto/sha256"
	"fmt"
	"strconv"

	"github.com/qredo-external/go-rnov/pkg/auth"
	"github.com/qredo-external/go-rnov/pkg/storage"
	"github.com/qredo-external/go-rnov/pkg/user"
)

// note in a prod ready it would have a couple of things, a big miss the lack of data that encapsulate `operation` just like
// auth but due the simplicity and the specifics of hashing is not worthy. I kept it due the design specifics.
type OperationManager struct {
	//	note it would have a logger
}

func NewOperationsService() *OperationManager {
	return &OperationManager{}
}

type Operations interface {
	Sum(data interface{}) (string, error)
}

// Sum - finds all the numbers throughout a valid (json) document and adds them together and hashes the expectedHash.
func (om OperationManager) Sum(data interface{}) (string, error) {
	sumRes := getSum(data)
	res := strconv.Itoa(sumRes)
	hashRes := fmt.Sprintf("%x", sha256.Sum256([]byte(res)))
	return hashRes, nil
}

// getSum - finds all the integers in a json structure and adds
// reason to use float64: https://golang.org/pkg/encoding/json/#Unmarshal same with decode
func getSum(data interface{}) int {
	var sum = 0
	switch data := data.(type) {
	case []interface{}:
		for _, v := range data {
			sum += getSum(v)
		}
	case map[string]interface{}:
		for _, v := range data {
			sum += getSum(v)
		}
	case float64:
		toInt := int(data)
		sum += toInt
	default:
		break
	}

	return sum
}

// AuthManager -
type AuthManager struct {
	// note add logger
	AuthOps auth.Operations
	Storage storage.ManageUsers
}

func NewAuthService(auth auth.Operations, storage storage.ManageUsers) *AuthManager {
	return &AuthManager{
		AuthOps: auth,
		Storage: storage,
	}
}

type Authorizer interface {
	CreateAuth(usr user.User) (string, error)
}

func (am AuthManager) CreateAuth(usr user.User) (string, error) {
	if err := usr.ValidateUser(); err != nil {
		return "", fmt.Errorf("error validating user: %s", err.Error())
	}
	jwt, err := am.AuthOps.CreateJWT(usr)
	if err != nil {
		return "", fmt.Errorf("error creting JWT: %s", err.Error())
	}
	if !am.Storage.AddUserToken(jwt) {
		// note different ways to handle it, as right now for simplicity just log...
		return jwt, nil
	}

	return jwt, nil
}
