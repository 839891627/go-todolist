package web

import "net/http"

const (
	SuccessCode = 200
	FailCode    = 500
)

type JsonResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Ok(data interface{}) *JsonResult {
	return &JsonResult{
		Code: SuccessCode,
		Data: data,
	}
}

func Fail(msg string) *JsonResult {
	return &JsonResult{Code: FailCode, Msg: msg}
}

type ResponseWriterPlus struct {
	http.ResponseWriter
	Status  int
	Written bool
}

func (w *ResponseWriterPlus) WriteHeader(status int) {
	if w.Written {
		return
	}

	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

// 重写父方法，记录判断是否已完成返回，避免重复写入
func (w *ResponseWriterPlus) Write(data []byte) (int, error) {
	if w.Written {
		return 0, nil
	}
	w.Written = true
	return w.ResponseWriter.Write(data)
}
