package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewGreeterUsecase, NewStudentUsercase)

type PageInfo struct {
	PageNo     int32
	PageSize   int32
	TotalCount int64
}
