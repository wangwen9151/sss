package handler

import (
	"context"

	"github.com/micro/go-log"

	example "sss/GetArea/proto/example"
	"github.com/astaxie/beego"
	"sss/IhomeWeb/utils"
	"github.com/astaxie/beego/orm"
	"sss/IhomeWeb/models"

	_ "github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/gomodule/redigo/redis"
	_ "github.com/garyburd/redigo/redis"

	"encoding/json"
	"github.com/astaxie/beego/cache"
	"time"
)

type Example struct{}

// 获取地域信息
func (e *Example) GetArea(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("获取地区信息 url:api/v1.0/areas")
	//初始化错误码
	rsp.Error = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Error)
	//1从缓存中获取数据
	//1.1准备连接redis信息
	redis_conf := map[string]string{
		"key":   utils.G_server_name,
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
	beego.Info(redis_conf)
	//1.2将连接redis信息map转换成json
	redis_conf_js, _ := json.Marshal(redis_conf)
	//1.3创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_js))
	if err != nil{
		beego.Info("获取地区信息 url:api/v1.0/areas 数据库redis查询失败", err)
		rsp.Error = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
	}
	//1.4定制一个固定的存储areas的key为  area_info
	area_value := bm.Get("area_info")
	if area_value != nil{
		beego.Info("获取地区信息 url:api/v1.0/areas 查询到缓存数据")
		area_map := []map[string]interface{}{}
		json.Unmarshal(area_value.([]byte),&area_map)
		beego.Info("获取地区信息 url:api/v1.0/areas 缓存中数据为",area_map )
		for _, v := range area_map {
			temp := example.ResponseArea{
				Aid:   int32(v["aid"].(float64)),
				Aname: v["aname"].(string),
			}
			rsp.Data = append(rsp.Data, &temp)
		}
		return nil
	}
	//2缓存中有数据，发送给前端
	//3缓存中无数据，直接在mysql中获取
	o := orm.NewOrm()
	var areas []models.Area
	qs := o.QueryTable("area")
	n, err := qs.All(&areas)
	if err != nil {
		beego.Info("获取地区信息 url:api/v1.0/areas 数据库mysql查询失败", err)
		rsp.Error = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}
	if n == 0 {
		beego.Info("获取地区信息 url:api/v1.0/areas 数据库mysql 空数据表")
		rsp.Error = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}
	beego.Info("获取地区信息 url:api/v1.0/areas 数据库mysql数据")
	//查询到的数据存入缓存
	area_json,_ := json.Marshal(areas)
	err = bm.Put("area_info",area_json,time.Second*3600)
	if err != nil {
		beego.Info("获取地区信息 url:api/v1.0/areas 缓存存储失败", err)
		rsp.Error = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}
	//将查到的数据发送给前端
	for _, v := range areas {
		temp := example.ResponseArea{
			Aid:   int32(v.Id),
			Aname: v.Name,
		}
		rsp.Data = append(rsp.Data, &temp)
	}

	return nil
}


// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Example) Stream(ctx context.Context, req *example.StreamingRequest, stream example.Example_StreamStream) error {
	log.Logf("Received Example.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&example.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Example) PingPong(ctx context.Context, stream example.Example_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&example.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
