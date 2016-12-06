package web

import (
	"net/http"
	"path/filepath"
	"reflect"
	"time"

	"github.com/kdada/tinygo/connector"
	"github.com/kdada/tinygo/log"
	"github.com/kdada/tinygo/meta"
	"github.com/kdada/tinygo/router"
	"github.com/kdada/tinygo/session"
	"github.com/kdada/tinygo/template"
)

// 参数类型方法
type ParamTypeFunc func(context *Context, name string, t reflect.Type) interface{}

// HttpProcessor 用于协调http连接器和路由,并管理Http应用的所有内容
type HttpProcessor struct {
	Root                  router.Router                 //根路由
	Config                *HttpConfig                   //http配置
	Logger                log.Logger                    //日志记录
	SessionContainer      session.SessionContainer      //Session容器
	CSRFContainer         session.SessionContainer      //Csrf容器
	Finders               map[string]ContextValueFinder //Context单类型值查找器
	MutiTypeFinders       []ContextValueFinder          //Context多类型值查找器
	DefaultValueContainer meta.ValueContainer           //web执行器默认使用的值容器
	Templates             *template.ViewTemplates       //视图模板信息
	Event                 HttpProcessorEvent            //处理器事件
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
	processor.Finders = make(map[string]ContextValueFinder)
	processor.MutiTypeFinders = make([]ContextValueFinder, 0, 1)
	register(processor)
	//默认值提供器为meta包的全局值提供器
	processor.DefaultValueContainer = meta.GlobalValueContainer
	//创建视图模板信息
	processor.Templates = template.NewViewTemplates(config.TemplateConfig)
	if config.Precompile {
		//预编译模板
		var err = processor.Templates.CompileAll()
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
		var excutor = NewSimpleExecutor(func(r *Context) (interface{}, error) {
			return r.Redispatch(r.Processor.Config.Home), nil
		})
		r.SetRouterExcutorGenerator(func() router.RouterExcutor {
			return excutor
		})
		processor.Root.AddChild(r)
	}
	return processor, nil
}

// RegisterFinder 注册单一类型的值查找器
func (this *HttpProcessor) RegisterFinder(t reflect.Type, finder ContextValueFinder) {
	this.Finders[t.String()] = finder
}

// RegisterMutiTypeFinder 注册多类型的值查找器
func (this *HttpProcessor) RegisterMutiTypeFinder(finder ContextValueFinder) {
	this.MutiTypeFinders = append(this.MutiTypeFinders, finder)
}

// SetDefaultValueContainer 设置默认的值容器
func (this *HttpProcessor) SetDefaultValueContainer(container meta.ValueContainer) {
	this.DefaultValueContainer = container
}

// createCookie 创建cookie
func (this *HttpProcessor) createCookie(name string, id string, expire int) *http.Cookie {
	var cookieValue = new(http.Cookie)
	cookieValue.Name = name
	cookieValue.Value = id
	cookieValue.Path = "/"
	cookieValue.HttpOnly = true
	if expire > 0 {
		cookieValue.MaxAge = expire
		cookieValue.Expires = time.Now().Add(time.Second * time.Duration(expire))
	}
	return cookieValue
}

// ResolveSession 处理会话相关内容
func (this *HttpProcessor) ResolveSession(context *Context) {
	if this.SessionContainer != nil {
		//添加Session信息
		var cookieValue, exist = context.Cookie(this.Config.SessionCookieName)
		var ss session.Session
		var ok bool = false
		if exist {
			ss, ok = this.SessionContainer.Session(cookieValue.Value)
		}
		if !ok {
			ss, ok = this.SessionContainer.CreateSession()
			context.AddCookie(this.createCookie(this.Config.SessionCookieName, ss.SessionId(), this.Config.SessionCookieExpire))
		}
		if ok {
			context.Session = ss
		}
	}
	if this.CSRFContainer != nil {
		//添加CSRF Session信息,CSRF的过期时间和Session相同,使用SessionExpire设置Cookie过期时间
		var cookieValue, exist = context.Cookie(this.Config.CSRFCookieName)
		var ss session.Session
		var ok bool = false
		if exist {
			ss, ok = this.CSRFContainer.Session(cookieValue.Value)
		}
		if !ok {
			ss, ok = this.CSRFContainer.CreateSession()
			context.AddCookie(this.createCookie(this.Config.CSRFCookieName, ss.SessionId(), this.Config.CSRFCookieExpire))
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
	var context, err = NewContext(segments, ct, this)
	if err == nil {
		this.ResolveSession(context)
		if this.Event != nil {
			var ctn = this.Event.Request(this, context)
			//确定是否继续执行
			if !ctn {
				return
			}
		}
		//路由匹配
		var executor, ok = this.Root.Match(context)
		if ok {
			var r, err = executor.Execute()
			if err != nil {
				this.Event.Error(this, context, err)
			} else if this.Event != nil {
				var result = []interface{}{}
				if r != nil {
					var rs, ok = r.([]interface{})
					if ok {
						result = rs
					} else {
						result = []interface{}{r}
					}
				}
				//结果处理
				this.Event.RequestFinish(this, context, result)
			}
		} else if this.Event != nil {
			this.Event.Error(this, context, ErrorRouterNotFound.Format(context.HttpContext.Request.URL.String()).Error())
		}
	} else {
		if this.Event != nil {
			this.Event.Error(this, context, err)
		}
	}
}
