package project_user

import (
	"fmt"
	"reflect"
	"strings"
	"test"
	_ "test_mysql"
)

/*user结构体*/
type User struct {
	Wallet         *Wallet
	Id             string
	Int_auth_token string
	Name           string
	Age            int
}

/*user认证结构体*/
type UserAuth struct {
	Int_auth_token string
	Id             string
	Device         string
	Ip             string
}

/*用户信息*/
type UserInfo struct {
	Id     string
	Name   string
	Device string
	Ip     string
}

var uMap map[string]reflect.Value

func init() {
	var wa Wallet

	uMap = make(map[string]reflect.Value)
	vf := reflect.ValueOf(&wa)
	vft := vf.Type()
	mNum := vf.NumMethod()
	for i := 0; i < mNum; i++ {
		mName := vft.Method(i).Name
		mName = strings.ToLower(mName)
		uMap[mName] = vf.Method(i)
	}
}

func (u *User) Auth(authstr string) bool {
	u.Wallet = &Wallet{}
	u.Wallet.Self = &Self{}
	uapi, err := test.NewUser("mysql")
	if err != nil {
		fmt.Println(err)
		return false
	}

	return uapi.Auth(u.Id, authstr)
}

func (u *User) GetFunctionCall(strs ...string) reflect.Value {
	var i int
	var o reflect.Value
	o = reflect.ValueOf(u)

	flens := len(strs) - 1
	for i = 0; i < flens; i++ {
		str := strings.Title(strs[i])
		o = o.Elem().FieldByName(str)
	}
	str := strings.Title(strs[i])
	return o.MethodByName(str)
}

// func (u *User) GetFunctionCall(strs ...string) reflect.Value {
// 	var i int
// 	var o reflect.Value
// 	o = reflect.ValueOf(u)

// 	flens := len(strs) - 1
// 	for i = 0; i < flens; i++ {
// 		str := strings.Title(strs[i])
// 		o = o.Elem().FieldByName(str)
// 	}
// 	str := strings.Title(strs[i])
// 	/*fmt.Println(str)
// 	o = o.MethodByName(str)
// 	fmt.Println(o.String())
// 	return o*/
// 	return o.MethodByName(str)
// }
