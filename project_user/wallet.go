package project_user

import (
	"fmt"
	"test"
	_ "test_mysql"
)

type Wallet struct {
	Self *Self
}

type Self struct {
}

func (s *Self) Detail(id, authstr string) map[string]string {
	m := make(map[string]string)
	uapi, err := test.NewUser("mysql")
	if err != nil {
		fmt.Println(err)
		return m
	}

	m, err = uapi.Detail(id, authstr)
	if err != nil {
		fmt.Println(err)
		return m
	}

	return m
}
