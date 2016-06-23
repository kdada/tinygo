package web

// HttpProcessor事件接口
type HttpProcessorEvent interface {
	// 每次出现一个新请求的时候触发
	Request(processor *HttpProcessor, context *Context)
	// 每次请求执行完成的时候触发
	RequestFinish(processor *HttpProcessor, context *Context, result []interface{})
	// 出现错误时触发,出现错误时context需要检查是否为nil后才能使用
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
		var r, ok = result[0].(Result)
		if ok {
			var err = r.WriteTo(context.HttpContext.ResponseWriter)
			if err != nil {
				processor.Logger.Error(err)
			}
		}
	}

}

// 出现错误时触发
func (this *DefaultHttpProcessorEvent) Error(processor *HttpProcessor, context *Context, err error) {

}
