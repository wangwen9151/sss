package main

import (
	"github.com/micro/go-log"
	"net/http"

	"github.com/micro/go-web"
	"github.com/julienschmidt/httprouter"

	_"sss/IhomeWeb/models"
	"sss/IhomeWeb/handler"
)

func main() {
	// create new web service
	service := web.NewService(
		web.Name("go.micro.web.IhomeWeb"),
		web.Version("latest"),
		web.Address(":8080"),
	)

	// initialise service
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	//使用路由中间件来印射页面
	rou := httprouter.New()
	rou.NotFound = http.FileServer(http.Dir("html"))

	//获取地区请求
	rou.GET("/api/v1.0/areas",handler.GetArea)
	//获取session
	rou.GET("/api/v1.0/session",handler.GetSession)
	//获取首页轮播图
	rou.GET("/api/v1.0/house/index",handler.GetIndex)
	//获取验证码图片
	rou.GET("/api/v1.0/imagecode/:uuid",handler.GetImageCd)

	// register html handler
	service.Handle("/", rou)

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
