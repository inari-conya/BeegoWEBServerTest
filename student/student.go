package student

import (
	"fmt"
	"reflect"
	"strings"
	"test"
	_ "test_mysql"
)

var sMap map[string]reflect.Value

type Student struct {
	Firstname string
	Lastname  string
	Age       int
}

func init() {
	var std Student

	sMap = make(map[string]reflect.Value)
	vf := reflect.ValueOf(&std)
	vft := vf.Type()
	mNum := vf.NumMethod()
	for i := 0; i < mNum; i++ {
		mName := vft.Method(i).Name
		mName = strings.ToLower(mName)
		sMap[mName] = vf.Method(i)
	}
}

func (s *Student) List() []map[string]string {
	l := make([]map[string]string, 0)

	sapi, err := test.NewStudent("mysql")
	if err != nil {
		fmt.Println(err)
		return l
	}

	l, err = sapi.ListTutorial()
	if err != nil {
		fmt.Println(err)
		return l
	}

	return l
}

func GetFuncMap() map[string]reflect.Value {
	return sMap
}
