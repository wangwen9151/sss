package handler

import (
	"context"

	"github.com/micro/go-log"

	example "sss/GetImageCd/proto/example"
	"github.com/astaxie/beego"

	"github.com/afocus/captcha"
	"image/color"

	"sss/IhomeWeb/utils"
	_ "github.com/astaxie/beego/orm"
	_ "sss/IhomeWeb/models"
	"encoding/json"

	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/gomodule/redigo/redis"
	_ "github.com/garyburd/redigo/redis"
	"time"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetImageCd(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("获取验证码图片 GetImageCd  /api/v1.0/imagecode/:uuid")

	//生成验证码图片
	cap := captcha.New()
	if err := cap.SetFont("comic.ttf"); err != nil {
		panic(err.Error())
	}
	cap.SetSize(90, 41)
	cap.SetDisturbance(captcha.MEDIUM)
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	//生成图片
	img, str := cap.Create(4, captcha.NUM)

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
	if err != nil {
		beego.Info("获取地区信息 url:api/v1.0/areas 数据库redis查询失败", err)
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
	}
	//1.4存储图片至redis
	bm.Put(req.Uuid, str, time.Second*300)
	beego.Info(str)

	//图片解引用
	imgTemp := *((*img).RGBA)
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	//返回图片拆分
	rsp.Pix = []byte(imgTemp.Pix)
	rsp.Stride = int64(imgTemp.Stride)
	rsp.Max = &example.ResponsePoint{X: int64(imgTemp.Rect.Max.X), Y: int64(imgTemp.Rect.Max.Y)}
	rsp.Min = &example.ResponsePoint{X: int64(imgTemp.Rect.Min.X), Y: int64(imgTemp.Rect.Min.Y)}

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
