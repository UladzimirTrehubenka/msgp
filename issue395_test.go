package main

import "testing"

func TestIssue395TupleAllowNil(t *testing.T) {
	runTest(t, issue395TupleAllowNil)
}

const issue395TupleAllowNil = `package main

import (
	"fmt"
	"os"
)

//go:generate msgp

//msgp:tuple User
type User struct {
	ID       []byte ` + "`msgpack:\",allownil\"`" + `
	Name     string
	Email    string
	IsActive bool
}

func main() {
	u := User{ID: nil, Name: "user"}
	data, err := u.MarshalMsg(nil)
	if err != nil {
		fmt.Println("User MarshalMsg failed:", err)
		os.Exit(1)
	}

	var user User
	if _, err := user.UnmarshalMsg(data); err != nil {
		fmt.Println("User UnmarshalMsg failed:", err)
		os.Exit(1)
	}
	if user.ID != nil || user.Name != "user" {
		fmt.Println("User UnmarshalMsg bad user:", user)
		os.Exit(1)
	}
}
`
