package tinygo

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kdada/tinygo/info"
	"github.com/kdada/tinygo/router"
	"github.com/kdada/tinygo/session"
)

//Session提供器,默认为内存Session
var sessionProvider session.SessionProvider

func initSession(sessionType string, expire int64) {
	//暂时只有内存Session
	sessionProvider = session.NewMemSessionProvider(expire)
}

// 生成静态路由
func generateStaticRouters() {
	for _, path := range tinyConfig.static {
		var components = strings.Split(path, `\/`)
		var count = len(components)
		var lastRouter = RootRouter
		if count > 2 {
			for _, c := range components {
				var space = router.NewSpaceRouter(c)
				lastRouter.AddChild(space)
				lastRouter = space
			}
		}
		lastRouter.AddChild(router.NewStaticRouter(components[count-1], path))
	}
}

// handler 统一路由处理方法
func handler(w http.ResponseWriter, r *http.Request) {
	var oldTime = time.Now().UnixNano()
	dispatch(w, r)
	var duration = (time.Now().UnixNano() - oldTime) / 1000000
	fmt.Println("[Info]", duration, "ms ", r.URL.Path)
}

// dispatch 路由查询处理
func dispatch(w http.ResponseWriter, r *http.Request) {
	var context = HttpContext{}
	context.urlParts = strings.Split(r.URL.Path, "/")[1:]
	context.request = r
	context.responseWriter = w
	//添加Session信息
	var cookie, err = context.request.Cookie(info.DefaultSessionCookieName)
	var ss session.Session
	var ok bool = false
	if err == nil {
		ss, ok = sessionProvider.Session(cookie.Value)
	}
	if !ok {
		ss, ok = sessionProvider.CreateSession()
		context.responseWriter.Header().Set("Set-Cookie", info.DefaultSessionCookieName+"="+ss.SessionId())
	}
	if ok {
		context.session = ss
	}

	var result = RootRouter.Pass(&context)
	if !result {
		//页面不存在
		HttpNotFound(w, r)
	}
}
