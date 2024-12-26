package vo

import (
	"net/http"
)

type R struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

//func (r *R) GetCode() int {
//	return r.code
//}
//
//func (r *R) GetMag() string {
//	return r.msg
//}
//func (r *R) GetData() any {
//	return r.data
//}

func Success(data any) *R {
	return &R{http.StatusOK, "", data}
}

func Fail(msg string, data any) *R {
	return &R{http.StatusInternalServerError, msg, data}
}

func Error(code int, msg string, data any) *R {
	return &R{code, msg, data}
}

var result = struct {
	ok, badRequest, notFound, internalServerError *R
}{
	ok:                  &R{http.StatusOK, "成功", nil},
	badRequest:          &R{http.StatusBadRequest, "参数异常", nil},
	notFound:            &R{http.StatusNotFound, "找不到资源", nil},
	internalServerError: &R{http.StatusInternalServerError, "服务器内部错误", nil},
}

type Data int

const (
	Ok Data = iota
	BadRequest
	NotFound
	InternalServerError
)

func (d Data) Data() *R {
	switch d {
	case BadRequest:
		return result.badRequest
	case NotFound:
		return result.notFound
	case InternalServerError:
		return result.internalServerError
	default:
		return result.ok
	}
}
