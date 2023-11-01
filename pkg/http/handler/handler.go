package handler

import (
	"encoding/json"
	"net/http"

	response "github.com/qredo-external/go-rnov/pkg/http/json"
	"github.com/qredo-external/go-rnov/pkg/service"
	"github.com/qredo-external/go-rnov/pkg/user"
)

// AuthHandler - holds the service that manages auth operation
type AuthHandler struct {
	Auth service.Authorizer
}

// NewAuthHandler - auth handler constructor
func NewAuthHandler(auth service.Authorizer) *AuthHandler {
	return &AuthHandler{
		Auth: auth,
	}
}

// CreateAuthHandler - handler for auth operation
func (a *AuthHandler) CreateAuthHandler(w http.ResponseWriter, r *http.Request) {
	usr := &user.User{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(usr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	JWTRes, err := a.Auth.CreateAuth(*usr)
	if err != nil {
		// note this could be either a Status Bad Request or a InternalError, for
		// simplicity i've left out the custom errors from the design please refer to readme.
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rBody := &response.JWT{
		JWT: JWTRes,
	}
	body, jsonErr := json.Marshal(rBody)
	if jsonErr != nil {
		// note should log error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(body)
}

// OperationHandler - holds the service that manage operations (sum)
type OperationHandler struct {
	operations service.Operations
}

// NewOperationHandler - operation handler constructor
func NewOperationHandler(op service.Operations) *OperationHandler {
	return &OperationHandler{
		operations: op,
	}
}

// SumHandler - handler for Sum operation
func (oh *OperationHandler) SumHandler(w http.ResponseWriter, r *http.Request) {
	var jsonMap interface{}
	dec := json.NewDecoder(r.Body)
	// note due the nature of jsonMap (interface{}) if we send a broken json structure it will not be detected by decode
	// it's out of scope and has no trivial solution.
	if err := dec.Decode(&jsonMap); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sumRes, err := oh.operations.Sum(jsonMap)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rBody := &response.Operation{
		Result: sumRes,
	}
	body, jsonErr := json.Marshal(rBody)
	if jsonErr != nil {
		// note should log error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)
}
