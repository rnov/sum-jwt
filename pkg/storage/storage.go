package storage

// ManageUsers - defines all the operations that need to be supported by any type of storage solutions used.
type ManageUsers interface {
	AddUserToken(token string) bool
	IsActiveToken(token string) bool
}
