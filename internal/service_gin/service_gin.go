package service_gin

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"kratos-demo/internal/biz"
	"kratos-demo/internal/service"
)

var ProviderSet = wire.NewSet(NewGinService)

const (
	OK            = 200
	ErrCodeCommon = 450
)

type Response struct {
	Code    int64       `json:"code"` // 业务响应码
	Reason  string      `json:"reason"`
	Message interface{} `json:"message"` // 返回数据
	//ServerTime int64       `json:"server_time"` // 1655210662 秒级时间戳
	//TraceId    string      `json:"trace_id"`    // trace_id
}

type GinService struct {
	// 可以调用GreeterService中封装好的属性(biz层的usecase)、方法
	srv    *service.GreeterService
	GinBiz *biz.GinUsecase

	log *log.Helper
}

func NewGinService(us *service.GreeterService, logger log.Logger, ginBiz *biz.GinUsecase) *GinService {
	return &GinService{
		srv:    us,
		GinBiz: ginBiz,
		log:    log.NewHelper(logger),
	}
}
