package controllers

import "net/http"

// Template 定义一个统一的模板方法签名,接收HTTP请求和响应writer,以及任意数据和错误
// 实现该接口的具体模板方法,需要读取请求,处理数据,并通过writer写入HTTP响应
// 在HTTP处理函数中,可以通过该接口调用不同的模板方法来处理请求
// templates可以方便地重用,而不需要重复处理请求和响应对象
// 通过interface{}传入数据,可以灵活使用任意数据类型
// 通过errs返回错误,可以方便处理失败情况
type Template interface {
	Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error)
}
