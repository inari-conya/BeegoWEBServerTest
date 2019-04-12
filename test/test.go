package test

import (
	"dbapi"
	"errors"
	"fmt"
	"sync"
)

type T_user struct {
	api dbapi.UserAPI
}

type T_student struct {
	api dbapi.StudentAPI
}

var (
	regMu sync.RWMutex
	udrv  = make(map[string]dbapi.UserAPI)
	sdrv  = make(map[string]dbapi.StudentAPI)
)

func UserReg(name string, user dbapi.UserAPI) {
	regMu.Lock()
	defer regMu.Unlock()
	if udrv == nil {
		fmt.Println("sql: Register driver is nil")
		return
	}
	if _, dup := udrv[name]; dup {
		fmt.Println("sql: Register called twice for driver " + name)
		return
	}
	udrv[name] = user
}

func StudentReg(name string, std dbapi.StudentAPI) {
	regMu.Lock()
	defer regMu.Unlock()
	if sdrv == nil {
		fmt.Println("sql: Register driver is nil")
		return
	}
	if _, dup := sdrv[name]; dup {
		fmt.Println("sql: Register called twice for driver " + name)
		return
	}
	sdrv[name] = std
}

/*创建T_user对象*/
func NewUser(reg string) (*T_user, error) {
	_, ok := udrv[reg]
	if ok == false {
		return nil, errors.New("user driver: " + reg + " does not exist")
	}
	var p T_user
	p.api = udrv[reg]
	return &p, nil
}

/*创建T_student对象*/
func NewStudent(reg string) (*T_student, error) {
	_, ok := sdrv[reg]
	if ok == false {
		return nil, errors.New("student driver: " + reg + " does not exist")
	}
	var p T_student
	p.api = sdrv[reg]
	return &p, nil
}

func (s *T_student) ListTutorial() ([]map[string]string, error) {
	return s.api.ListTutorial()
}

func (u *T_user) Auth(id, authstr string) bool {
	return u.api.Auth(id, authstr)
}

func (u *T_user) Detail(id, authstr string) (map[string]string, error) {
	return u.api.Detail(id, authstr)
}
