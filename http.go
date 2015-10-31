package tinygo

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/kdada/tinygo/session"
)

// Session提供器,默认为内存Session
var sessionProvider session.SessionProvider

// Csrf提供器,使用内存Session
var csrfProvider session.SessionProvider

// initSession 初始化session
func initSession(sessionType string, expire int64) {
	var err error
	sessionProvider, err = session.NewSessionProvider(session.SessionType(sessionType), expire)
	if err != nil {
		//会话创建失败 严重错误
		panic(err)
	}
}

// cleanAllDeadSession 清理过期的session和csrf session
func cleanAllDeadSession(expire int64) {
	defer func() {
		if err := recover(); err != nil {
			Error(err)
			go cleanAllDeadSession(expire)
		}
	}()
	for true {
		//清理间隔
		time.Sleep(time.Duration(expire) * time.Second)
		if sessionProvider != nil {
			sessionProvider.Clean()
		}
		if csrfProvider != nil {
			csrfProvider.Clean()
		}
	}
}

// initCsrfSession 初始化csrf session,csrf有效期与session相同
func initCsrfSession(expire int64) {
	var err error
	csrfProvider, err = session.NewSessionProvider(session.SessionTypeMemory, expire)
	if err != nil {
		//会话创建失败 严重错误
		panic(err)
	}
}

// handler 统一路由处理方法
func handler(w http.ResponseWriter, r *http.Request) {
	var oldTime = time.Now().UnixNano()
	var found, static = dispatch(w, r)
	if !static {
		var duration = (time.Now().UnixNano() - oldTime) / 1000000
		var foundstr = "  found  "
		if !found {
			foundstr = "not found"
		}
		fmt.Println("["+r.Method+"]", foundstr, duration, "ms ", r.URL.Path)
	}
}

// dispatch 路由查询处理
//  return:(是否查找到路由,是否是静态路由)
func dispatch(w http.ResponseWriter, r *http.Request) (bool, bool) {
	var context = HttpContext{}
	var url = filepath.Clean("/" + strings.Replace(r.URL.Path, `\`, `/`, -1))
	context.urlParts = strings.Split(url, "/")
	var i = len(context.urlParts) - 1
	for ; i > 0; i-- {
		if context.urlParts[i] != "" {
			break
		}
	}
	context.urlParts = context.urlParts[:i+1]
	context.request = r
	context.responseWriter = w

	// 检索路由信息
	var result = RootRouter.Pass(&context)
	if result {
		//执行
		if !context.static {
			//只有非静态的上下文才能设置session和csrf
			if sessionProvider != nil {
				//添加Session信息
				var cookieValue, err = context.Cookie(DefaultSessionCookieName)
				var ss session.Session
				var ok bool = false
				if err == nil {
					ss, ok = sessionProvider.Session(cookieValue)
				}
				if !ok {
					ss, ok = sessionProvider.CreateSession()
				}
				if ok {
					//更新cookie有效期
					var cookie = &http.Cookie{}
					cookie.Name = DefaultSessionCookieName
					cookie.Value = ss.SessionId()
					cookie.Path = "/"
					cookie.MaxAge = int(tinyConfig.sessionexpire)
					cookie.Expires = time.Now().Add(time.Duration(cookie.MaxAge) * time.Second)
					cookie.HttpOnly = true
					context.AddCookie(cookie)
					context.session = ss
				}
			}
			if csrfProvider != nil {
				//添加Csrf Session信息,csrf有效期与session相同
				var cookieValue, err = context.Cookie(DefaultCSRFCookieName)
				var ss session.Session
				var ok bool = false
				if err == nil {
					ss, ok = csrfProvider.Session(cookieValue)
				}
				if !ok {
					ss, ok = csrfProvider.CreateSession()
				}
				if ok {
					//更新cookie有效期
					var cookie = &http.Cookie{}
					cookie.Name = DefaultCSRFCookieName
					cookie.Value = ss.SessionId()
					cookie.Path = "/"
					cookie.MaxAge = int(tinyConfig.sessionexpire)
					cookie.Expires = time.Now().Add(time.Duration(cookie.MaxAge) * time.Second)
					cookie.HttpOnly = true
					context.AddCookie(cookie)
					context.csrf = ss
				}
			}
		}
		context.execute()
	} else {
		//页面不存在
		HttpNotFound(w, r)
	}
	return result, context.static
}

// HttpNotFound 返回页面不存在(404)错误
func HttpNotFound(w http.ResponseWriter, r *http.Request) {
	if tinyConfig.pageerr != "" {
		w.WriteHeader(404)
		ParseTemplate(&HttpContext{responseWriter: w, request: r}, tinyConfig.pageerr, nil)
	} else {
		http.NotFound(w, r)
	}
}

// Redirect 302重定向
func Redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, 302)
}

// RedirectPermanently 301重定向
func RedirectPermanently(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, 301)
}

// ParseUrlValueToStruct 将url值解析到结构体中
//  urlValues:url值
//  value:结构体的反射值
func ParseUrlValueToStruct(urlValues url.Values, value reflect.Value) {
	if value.Kind() == reflect.Struct {
		for i := 0; i < value.NumField(); i++ {
			var fieldValue = value.Field(i)
			var fieldType = value.Type().Field(i)
			if fieldType.Anonymous {
				//匿名组合字段,进行递归解析
				ParseUrlValueToStruct(urlValues, fieldValue)
			} else {
				//非匿名字段
				if fieldValue.CanSet() {
					var fieldName = fieldType.Tag.Get("from")
					if fieldName == "-" {
						//如果是-,则忽略当前字段
						continue
					}
					if fieldName == "" {
						//如果为空,则使用字段名
						fieldName = fieldType.Name
					}
					var urlValue = urlValues.Get(fieldName)
					switch fieldType.Type.Kind() {
					case reflect.Bool:
						result, err := strconv.ParseBool(urlValue)
						fieldValue.SetBool(result && err == nil)
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						result, err := strconv.ParseInt(urlValue, 10, 64)
						if err == nil {
							fieldValue.SetInt(result)
						} else {
							fieldValue.SetInt(0)
						}

					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						result, err := strconv.ParseUint(urlValue, 10, 64)
						if err == nil {
							fieldValue.SetUint(result)
						} else {
							fieldValue.SetUint(0)
						}
					case reflect.Float32, reflect.Float64:
						result, err := strconv.ParseFloat(urlValue, 64)
						if err == nil {
							fieldValue.SetFloat(result)
						} else {
							fieldValue.SetFloat(0.0)
						}
					case reflect.Interface:
						fieldValue.Set(reflect.ValueOf(urlValue))
					case reflect.String:
						fieldValue.SetString(urlValue)
					case reflect.Struct:
						if fieldType.Type.String() == "time.Time" {
							//使用本地时区解析时间
							result, err := time.ParseInLocation("2006-01-02 15:04:05", urlValue, time.Local)
							if err == nil {
								fieldValue.Set(reflect.ValueOf(result))
							}
						}
					case reflect.Slice:
						//可以解析的数组类型包括bool,int,float64,string类型
						switch fieldType.Type.String() {
						case "[]bool":
							stringValue := urlValues[fieldName]
							boolValue := make([]bool, len(stringValue), len(stringValue))
							for i := 0; i < len(boolValue); i++ {
								boolValue[i], _ = strconv.ParseBool(stringValue[i])
							}
							fieldValue.Set(reflect.ValueOf(boolValue))
						case "[]int":
							stringValue := urlValues[fieldName]
							intValue := make([]int, len(stringValue), len(stringValue))
							for i := 0; i < len(intValue); i++ {
								intValue[i], _ = strconv.Atoi(stringValue[i])
							}
							fieldValue.Set(reflect.ValueOf(intValue))
						case "[]float64":
							stringValue := urlValues[fieldName]
							floatValue := make([]float64, len(stringValue), len(stringValue))
							for i := 0; i < len(floatValue); i++ {
								floatValue[i], _ = strconv.ParseFloat(stringValue[i], 64)
							}
						case "[]string":
							stringValue := urlValues[fieldName]
							fieldValue.Set(reflect.ValueOf(stringValue))
						}
					}
				}

			}
		}
	}
}
