package users

import "fmt"

type UserNotFoundErr struct {
	userID int64
}

func NewUserNotFoundErr(userID int64) UserNotFoundErr {
	return UserNotFoundErr{userID}
}

func (err UserNotFoundErr) Error() string {
	return fmt.Sprintf("user with ID %v not found", err.userID)
}
