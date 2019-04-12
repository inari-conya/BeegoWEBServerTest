package main

import (
	//"fmt"
	//"io/ioutil"
	"fmt"
	"io/ioutil"
	_ "routers"

	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
)

func main() {
	b, err := ioutil.ReadFile("server.config")
	if err != nil {
		fmt.Println("打开配置文件失败")
	}

	js, err := simplejson.NewJson(b)
	if err != nil {
		panic(err.Error())
	}

	port, _ := js.Get("port").Int()

	beego.BConfig.Listen.HTTPPort = port
	beego.Run()
}
