package handler

import (
	"context"
	"encoding/json"
	"net/http"

	GETAREA "sss/GetArea/proto/example"
	GETIMAGECD "sss/GetImageCd/proto/example"

	"github.com/micro/go-grpc"
	"sss/IhomeWeb/models"
	"github.com/julienschmidt/httprouter"
	"sss/IhomeWeb/utils"
	"github.com/astaxie/beego"

	"image"
	"github.com/afocus/captcha"
	"image/png"
)

// 获取地域信息
func GetArea(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("获取地区请求 GetArea  /api/v1.0/areas")

	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := GETAREA.NewExampleService("go.micro.srv.GetArea", server.Client())

	rsp, err := exampleClient.GetArea(context.TODO(), &GETAREA.Request{
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	area_list := []models.Area{}

	for _, v := range rsp.Data {
		temp := models.Area{Id: int(v.Aid), Name: v.Aname}
		area_list = append(area_list, temp)
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Error,
		"errmsg": rsp.Errmsg,
		"data":   area_list,
	}

	//回传时会直接发送数据，没有设置数据格式，前端浏览器无法识别
	//需要设置数据格式
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

//获取session
func GetSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("获取session GetSession  /api/v1.0/session")

	// we want to augment the response
	dataTemp := map[string]string{"name": "wangwen"}

	response := map[string]interface{}{
		"errno":  utils.RECODE_SESSIONERR,
		"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		"data":   dataTemp,
	}

	//回传时会直接发送数据，没有设置数据格式，前端浏览器无法识别
	//需要设置数据格式
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

//获取首页轮播图
func GetIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("获取首页轮播图 GetIndex  /api/v1.0/house/index")

	response := map[string]interface{}{
		"errno":  utils.RECODE_OK,
		"errmsg": utils.RecodeText(utils.RECODE_OK),
	}

	//回传时会直接发送数据，没有设置数据格式，前端浏览器无法识别
	//需要设置数据格式
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// 获取验证码图片
func GetImageCd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	beego.Info("获取验证码图片 GetImageCd  /api/v1.0/imagecode/:uuid")

	//创建服务
	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := GETIMAGECD.NewExampleService("go.micro.srv.GetImageCd", server.Client())

	uuid := ps.ByName("uuid")
	rsp, err := exampleClient.GetImageCd(context.TODO(), &GETIMAGECD.Request{
		Uuid:uuid,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//接收GetImageCd回传的图片
	var img image.RGBA
	img.Stride = int(rsp.Stride)
	img.Pix = []uint8(rsp.Pix)
	img.Rect.Max.X = int(rsp.Max.X)
	img.Rect.Max.Y = int(rsp.Max.Y)
	img.Rect.Min.X = int(rsp.Min.X)
	img.Rect.Min.Y = int(rsp.Min.Y)
	//转换成	/afocus/captcha 的格式
	var image captcha.Image
	image.RGBA = &img

	//将图片发送给浏览器
	png.Encode(w,image)
}
