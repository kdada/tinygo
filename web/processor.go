package web

import (
	"net/http"
	"path/filepath"
	"reflect"
	"time"

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
	Templates        *ViewTemplates           //视图模板信息
	Event            HttpProcessorEvent       //处理器事件
}

// NewHttpProcessor 创建Http处理器
func NewHttpProcessor(root router.Router, config *HttpConfig) (*HttpProcessor, error) {
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
	//CSRF,过期时间与Session相同,CSRFExpire用于设置CSRF token的过期时间
	if config.CSRF {
		var container, err = session.NewSessionContainer(config.CSRFType, config.SessionExpire, config.CSRFSource)
		if err != nil {
			panic(err)
		}
		processor.CSRFContainer = container
	}

	//注册参数类型方法
	processor.Funcs = make(map[string]ParamTypeFunc)
	register(processor.Funcs)
	processor.DefaultFunc = DefaultFunc
	//创建视图模板信息
	processor.Templates = NewViewTemplates(filepath.Join(config.Root, config.View), config.ViewConfig, config.TemplateName, config.TemplateExt)
	if config.Precompile {
		//预编译模板
		var err = processor.Templates.CompileAll(new(CommonFunMap))
		if err != nil {
			return nil, err
		}
	}
	//注册http事件
	processor.Event = new(DefaultHttpProcessorEvent)
	//注册静态文件路由
	if processor.Config.Favicon != "" {
		processor.Root.AddChild(NewFileRouter("favicon.ico", processor.Config.Favicon))
	}
	if processor.Config.Robots != "" {
		processor.Root.AddChild(NewFileRouter("robots.txt", processor.Config.Robots))
	}
	if len(processor.Config.Static) > 0 {
		for _, s := range processor.Config.Static {
			processor.Root.AddChild(NewStaticRouter(filepath.Base(s), s))
		}
	}
	//创建首页跳转
	if processor.Config.Home != "" {
		var r = NewSpaceRouter("Get")
		var excutor = NewSimpleExecutor(func(r *Context) interface{} {
			return r.Redispatch(r.Processor.Config.Home)
		})
		r.SetRouterExcutorGenerator(func() router.RouterExcutor {
			return excutor
		})
		processor.Root.AddChild(r)
	}
	return processor, nil
}

// ParamFunc 根据类型全名获取指定的生成方法
func (this *HttpProcessor) ParamFunc(t string) ParamTypeFunc {
	var f, ok = this.Funcs[t]
	if !ok {
		f = this.DefaultFunc
	}
	return f
}

// createCookie 创建cookie(有效期为1天)
func (this *HttpProcessor) createCookie(name string, id string) *http.Cookie {
	var cookieValue = new(http.Cookie)
	cookieValue.Name = name
	cookieValue.Value = id
	cookieValue.Path = "/"
	cookieValue.MaxAge = 24 * 3600
	cookieValue.Expires = time.Now().Add(time.Hour * 24)
	cookieValue.HttpOnly = true
	return cookieValue
}

// addCookie 添加cookie
func (this *HttpProcessor) addCookie(context *Context, cookie *http.Cookie) {
	context.HttpContext.ResponseWriter.Header().Add("Set-Cookie", cookie.String())
}

// ResolveSession 处理会话相关内容
func (this *HttpProcessor) ResolveSession(context *Context) {
	if this.SessionContainer != nil {
		//添加Session信息
		var cookieValue, err = context.HttpContext.Request.Cookie(this.Config.SessionCookieName)
		var ss session.Session
		var ok bool = false
		if err == nil {
			ss, ok = this.SessionContainer.Session(cookieValue.Value)
		}
		if !ok {
			ss, ok = this.SessionContainer.CreateSession()
			this.addCookie(context, this.createCookie(this.Config.SessionCookieName, ss.SessionId()))
		}
		if ok {
			context.Session = ss
		}
	}
	if this.CSRFContainer != nil {
		//添加CSRF Session信息,CSRF的过期时间和Session相同,使用SessionExpire设置Cookie过期时间
		var cookieValue, err = context.HttpContext.Request.Cookie(this.Config.CSRFCookieName)
		var ss session.Session
		var ok bool = false
		if err == nil {
			ss, ok = this.CSRFContainer.Session(cookieValue.Value)
		}
		if !ok {
			ss, ok = this.CSRFContainer.CreateSession()
			this.addCookie(context, this.createCookie(this.Config.CSRFCookieName, ss.SessionId()))
		}
		if ok {
			context.CSRF = ss
		}
	}
}

// Dispatch 将接收到的请求进行分发
//  segments:用于进行分发的路径段信息
//  data:连接携带的数据
func (this *HttpProcessor) Dispatch(segments []string, data interface{}) {
	var ct = data.(*connector.HttpContext)
	var context, err = NewContext(segments, ct)
	if err == nil {
		context.Processor = this
		this.ResolveSession(context)
		if this.Event != nil {
			this.Event.Request(this, context)
		}
		var executor, ok = this.Root.Match(context)
		if ok {
			var r = executor.Execute()
			if this.Event != nil {
				var result = []interface{}{}
				if r != nil {
					var rs, ok = r.([]interface{})
					if ok {
						result = rs
					} else {
						result = []interface{}{r}
					}
				}
				this.Event.RequestFinish(this, context, result)
			}
		} else {
			if this.Event != nil {
				this.Event.RequestFinish(this, context, []interface{}{NewUserDefinedResult(StatusCodePageNotFound, ErrorRouterNotFound.Format(context.HttpContext.Request.URL.String()).String())})
			}
		}
	} else {
		if this.Event != nil {
			this.Event.Error(this, context, err)
		}
	}
}
