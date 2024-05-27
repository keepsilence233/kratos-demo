package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewGreeterUsecase, NewStudentUsercase, NewGinUsecase)

type PageInfo struct {
	PageNo     int32
	PageSize   int32
	TotalCount int64
}

// GinUsecase
type GinUsecase struct {
	log *log.Helper

	ginRepo BizGinRepo
}

func NewGinUsecase(logger log.Logger, ginRepo BizGinRepo) *GinUsecase {
	return &GinUsecase{
		log:     log.NewHelper(logger),
		ginRepo: ginRepo,
	}
}
