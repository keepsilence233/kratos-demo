package common

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/errors"
)

// 定义错误码
type Error struct {
	Code    int
	Message string
}

func (err Error) Error() string {
	return err.Message
}

// 定义错误
type Errord struct {
	Code    int    // 错误码
	Message string // 展示给用户看的
	Err     error  // 保存内部错误信息
}

func (err *Errord) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.Message, err.Err)
}

func Parse(err error) (int, string) {
	e, ok := err.(*Error)
	if ok {
		return e.Code, e.Message
	}
	if se := errors.FromError(err); se.Code == 400 {
		return ErrBind.Code, ErrBind.Message
	}
	return InternalServerError.Code, InternalServerError.Message
}

/*
错误码设计
第一位表示错误级别, 1 为系统错误, 2 为普通错误
第二三位表示服务模块代码
第四五位表示具体错误代码
*/

var (
	OK = &Error{Code: 000000, Message: "success"}

	// 系统内部错误, 前缀为 100
	InternalServerError = &Error{Code: 10001, Message: "内部服务器错误"}
	ErrBind             = &Error{Code: 10002, Message: "请求参数错误"}
)
