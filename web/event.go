package web

import "strings"

// HttpProcessor事件接口
type HttpProcessorEvent interface {
	// 每次出现一个新请求的时候触发
	Request(processor *HttpProcessor, context *Context)
	// 每次请求正确执行完成的时候触发
	RequestFinish(processor *HttpProcessor, context *Context, result []interface{})
	// 请求过程中出现任何错误时触发,出现错误时context需要检查是否为nil后才能使用
	Error(processor *HttpProcessor, context *Context, err error)
}

// 默认事件
type DefaultHttpProcessorEvent struct {
}

// 每次出现一个新请求的时候触发
func (this *DefaultHttpProcessorEvent) Request(processor *HttpProcessor, context *Context) {

}

// 每次请求执行完成的时候触发
func (this *DefaultHttpProcessorEvent) RequestFinish(processor *HttpProcessor, context *Context, result []interface{}) {
	if len(result) > 0 {
		var r, ok = result[0].(HttpResult)
		if ok && r.Code() == StatusCodeRedispatch {
			//StatusCodeRedispatch的返回结果需要进行重新分发
			processor.Dispatch(strings.Split(r.Message(), "/"), context.HttpContext)
			return
		} else {
			//处理其他情况
			var w, ok2 = result[0].(Result)
			if ok2 {
				var err = context.WriteResult(w)
				if err != nil {
					this.Error(processor, context, err)
				}
			}
		}

	}

}

// 出现错误时触发
func (this *DefaultHttpProcessorEvent) Error(processor *HttpProcessor, context *Context, err error) {
	if context != nil {
		var err = context.WriteResult(context.NotFound())
		if err != nil {
			processor.Logger.Error(err)
		}
	}
	processor.Logger.Error(err)
}
