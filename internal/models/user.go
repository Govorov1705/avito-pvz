package models

import "github.com/google/uuid"

type UserRole string

const (
	Employee  UserRole = "employee"
	Moderator UserRole = "moderator"
)

type User struct {
	ID       uuid.UUID
	Email    string
	Password string
	Role     UserRole
}
