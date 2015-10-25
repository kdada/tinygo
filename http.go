package tinygo

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kdada/tinygo/session"
)

//Session提供器,默认为内存Session
var sessionProvider session.SessionProvider

func initSession(sessionType string, expire int64) {
	var err error
	sessionProvider, err = session.NewSessionProvider(session.SessionType(sessionType), expire)
	if err != nil {
		//会话创建失败 严重错误
		panic(err)
	}
}

// handler 统一路由处理方法
func handler(w http.ResponseWriter, r *http.Request) {
	var oldTime = time.Now().UnixNano()
	dispatch(w, r)
	var duration = (time.Now().UnixNano() - oldTime) / 1000000
	fmt.Println("["+r.Method+"]", duration, "ms ", r.URL.Path)
}

// dispatch 路由查询处理
func dispatch(w http.ResponseWriter, r *http.Request) {
	var context = HttpContext{}
	context.urlParts = strings.Split(r.URL.Path, "/")
	var i = len(context.urlParts) - 1
	for ; i > 0; i-- {
		if context.urlParts[i] != "" {
			break
		}
	}
	context.urlParts = context.urlParts[:i+1]
	context.request = r
	context.responseWriter = w
	if sessionProvider != nil {
		//添加Session信息
		var cookie, err = context.request.Cookie(DefaultSessionCookieName)
		var ss session.Session
		var ok bool = false
		if err == nil {
			ss, ok = sessionProvider.Session(cookie.Value)
		}
		if !ok {
			ss, ok = sessionProvider.CreateSession()
			context.responseWriter.Header().Set("Set-Cookie", DefaultSessionCookieName+"="+ss.SessionId()+";Max-Age="+strconv.Itoa(int(tinyConfig.sessionexpire))+";Path=/")
		}
		if ok {
			context.session = ss
		}
	}
	// 检索路由信息
	var result = RootRouter.Pass(&context)
	if result {
		//执行
		context.execute()
	} else {
		//页面不存在
		HttpNotFound(w, r)
	}
}
