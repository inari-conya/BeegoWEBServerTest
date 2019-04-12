package controllers

import (
	"fmt"
	"project_user"
	"reflect"
	"strconv"
	"strings"

	"student"

	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
)

var vMap map[string]reflect.Value

/*json返回类型*/
type Jsonret struct {
	Ret     int    `json:"ret,omitempty"`
	Version int    `json:"version,omitempty"`
	Action  string `json:"action,omitempty"`
	Result  int    `json:"result,omitempty"`
}

func init() {
	var vs Vers

	vMap = make(ControllerMapsType, 0)
	vf := reflect.ValueOf(&vs)
	vft := vf.Type()
	mNum := vf.NumMethod()
	for i := 0; i < mNum; i++ {
		mName := vft.Method(i).Name
		mName = strings.ToLower(mName)
		vMap[mName] = vf.Method(i)
	}
}

//定义控制器函数Map类型
type ControllerMapsType map[string]reflect.Value

//声明控制器函数Map类型变量
var ControllerMaps ControllerMapsType

//定义路由器结构类型
type Vers struct {
}

func (i *Vers) Testapi() (ret Jsonret) {
	ret = Jsonret{}
	ret.Ret = 1000
	return ret
}

func (i *Vers) Plus(v, a, b int, m string) (ret Jsonret) {
	ret = Jsonret{}
	ret.Ret = 1000
	ret.Version = v
	ret.Action = m
	ret.Result = a + b
	return ret
}

//API1
type API1 struct {
	beego.Controller
}

/*
curl "http://127.0.0.1:8808/v3/test-api"
curl "http://127.0.0.1:8808/v3/plus?a=3&b=4"
*/
func (c *API1) Get() {
	a, _ := c.GetInt("a")
	b, _ := c.GetInt("b")
	vers := c.Ctx.Input.Param(":ver")
	ver, _ := strconv.Atoi(vers)
	Method := c.Ctx.Input.Param(":splat")

	parms := make([]reflect.Value, 4)
	parms[0] = reflect.ValueOf(ver)
	parms[1] = reflect.ValueOf(a)
	parms[2] = reflect.ValueOf(b)
	parms[3] = reflect.ValueOf(Method)

	if Method == "test-api" {
		res := vMap["testapi"].Call(nil)
		ret := res[0].Interface().(Jsonret)
		c.Data["json"] = &ret
		c.ServeJSON()
	} else {
		_, ok := vMap[Method]
		if ok {
			res := vMap[Method].Call(parms)
			ret := res[0].Interface().(Jsonret)
			c.Data["json"] = &ret
			c.ServeJSON()
		} else {
			c.Abort("404")
		}
	}

	return
}

/*
curl -H "Content-Type: application/x-www-form-urlencoded" -X POST  --data "{\"a\":3, \"b\":4}" http://127.0.0.1:8808/v3/plus
*/
func (c *API1) Post() {
	jsonstr := c.Ctx.Input.RequestBody

	js, err := simplejson.NewJson(jsonstr)
	if err != nil {
		fmt.Println("NewJson failed", err.Error())
		c.Abort("400")
	}

	a, err := js.Get("a").Int()
	if err != nil {
		fmt.Println("get a failed", err.Error())
		c.Abort("400")
	}

	b, err := js.Get("b").Int()
	if err != nil {
		fmt.Println("get b failed", err.Error())
		c.Abort("400")
	}

	vers := c.Ctx.Input.Param(":ver")
	ver, _ := strconv.Atoi(vers)
	Method := c.Ctx.Input.Param(":splat")

	parms := make([]reflect.Value, 4)
	parms[0] = reflect.ValueOf(ver)
	parms[1] = reflect.ValueOf(a)
	parms[2] = reflect.ValueOf(b)
	parms[3] = reflect.ValueOf(Method)

	_, ok := vMap[Method]
	if ok {
		res := vMap[Method].Call(parms)
		ret := res[0].Interface().(Jsonret)
		c.Data["json"] = &ret
		c.ServeJSON()
	} else {
		c.Abort("404")
	}

	return
}

type Tutorial struct {
	beego.Controller
}

/*
curl -H "Content-Type: application/x-www-form-urlencoded" -X GET http://127.0.0.1:8808/tutorial/student/list
*/
func (c *Tutorial) Get() {
	ns := c.Ctx.Input.Param(":namespace")
	rsc := c.Ctx.Input.Param(":resource")
	act := c.Ctx.Input.Param(":action")

	sMap := student.GetFuncMap()

	if ns == "tutorial" {
		if rsc == "student" {
			_, ok := sMap[act]
			if ok == false {
				fmt.Println("sMap[", act, "] does not exist")
				c.Abort("404")
			}
			res := sMap[act].Call(nil)
			data := res[0].Interface().([]map[string]string)
			fmt.Println(data)
			c.Data["list"] = data
			c.TplName = "list.html"
		}
	}
	return
}

type UserInfo struct {
	beego.Controller
}

/*
curl -H "Content-Type: application/x-www-form-urlencoded" -X POST  --data "{\"intAuthToken\":\"yuqbajnnr\"}" http://127.0.0.1:8808/user/sp100032/wallet/self/detail
curl -H "Content-Type: application/x-www-form-urlencoded" -X POST  --data "{}" http://127.0.0.1:8808/user/sp100032/wallet/self/detail
*/
func (c *UserInfo) Post() {
	var user project_user.User
	authstr := ""
	uid := c.Ctx.Input.Param(":userid")
	obj := c.Ctx.Input.Param(":obj")
	act := c.Ctx.Input.Param(":act")
	function := c.Ctx.Input.Param(":function")

	jsonstr := c.Ctx.Input.RequestBody

	js, err := simplejson.NewJson(jsonstr)
	if err != nil {
		fmt.Println("NewJson failed", err.Error())
	} else {
		authstr, err = js.Get("intAuthToken").String()
		if err != nil {
			fmt.Println("get intAuthToken failed", err.Error())
		}
	}

	user.Id = uid
	auth_res := user.Auth(authstr)
	if auth_res == false {
		ret := Jsonret{}
		ret.Ret = 1001
		c.Data["json"] = &ret
		c.ServeJSON()
		return
	}
	user.Int_auth_token = authstr

	parms := make([]reflect.Value, 2)
	parms[0] = reflect.ValueOf(user.Id)
	parms[1] = reflect.ValueOf(user.Int_auth_token)

	v := user.GetFunctionCall(function, obj, act)
	res := v.Call(parms)
	data := res[0].Interface().(map[string]string)

	c.Data["list"] = data
	c.TplName = "user.html"
}
