package entities

import "time"

type User struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt *time.Time
	Email     string
	Password  UserPassword
}

type UserPassword []byte

const hiddenPassword string = "********"

func (up UserPassword) String() string {
	return hiddenPassword
}

type UserCreatable struct {
	Email    string
	Password []byte
}
