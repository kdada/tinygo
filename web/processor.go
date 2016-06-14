package web

import (
	"reflect"

	"github.com/kdada/tinygo/connector"
	"github.com/kdada/tinygo/log"
	"github.com/kdada/tinygo/router"
	"github.com/kdada/tinygo/session"
)

// 参数类型方法
type ParamTypeFunc func(context *Context, name string, t reflect.Type) interface{}

// HttpProcessor 用于协调http连接器和路由,并管理Http应用的所有内容
type HttpProcessor struct {
	Root             router.Router            //根路由
	Config           *HttpConfig              //http配置
	Logger           log.Logger               //日志记录
	SessionContainer session.SessionContainer //Session容器
	CSRFContainer    session.SessionContainer //Csrf容器
	Funcs            map[string]ParamTypeFunc //参数生成方法
	DefaultFunc      ParamTypeFunc            //当Funcs中不存在指定类型的方法时,使用该方法处理
}

// NewHttpProcessor 创建Http处理器
func NewHttpProcessor(root router.Router, config *HttpConfig) *HttpProcessor {
	var processor = new(HttpProcessor)
	processor.Root = root
	processor.Config = config
	//日志
	if config.Log {
		var logger, err = log.NewLogger(config.LogType, config.LogPath)
		if err != nil {
			panic(err)
		}
		logger.SetAsync(config.LogAsync)
		processor.Logger = logger
	}
	//session
	if config.Session {
		var container, err = session.NewSessionContainer(config.SessionType, config.SessionExpire, config.SessionSource)
		if err != nil {
			panic(err)
		}
		processor.SessionContainer = container
	}
	//CSRF
	if config.CSRF {
		var container, err = session.NewSessionContainer(config.CSRFType, config.CSRFExpire, config.CSRFSource)
		if err != nil {
			panic(err)
		}
		processor.CSRFContainer = container
	}

	processor.Funcs = make(map[string]ParamTypeFunc)
	register(processor.Funcs)
	processor.DefaultFunc = DefaultFunc

	return processor
}

// Dispatch 将接收到的请求进行分发
//  segments:用于进行分发的路径段信息
//  data:连接携带的数据
func (this *HttpProcessor) Dispatch(segments []string, data interface{}) {
	var ct = data.(*connector.HttpContext)
	var context = NewContext(segments, ct)
	context.Processor = this
	context.HttpContext.Request.ParseMultipartForm(int64(this.Config.MaxRequestMemory))
	var executor, ok = this.Root.Match(context)
	if ok {
		var err = executor.Execute()
		if err != nil {
			panic(err)
		}
	} else {
		panic("路由未匹配")
	}
}
