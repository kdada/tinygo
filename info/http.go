package info

// Http方法常量
type HttpMethod string

const (
	HttpMethodGet     HttpMethod = "GET"    //GET方法
	HttpMethodPost    HttpMethod = "POST"   //POST方法
	HttpMethodPut     HttpMethod = "PUT"    //PUT方法
	HttpMethodHead    HttpMethod = "HEAD"   //HEAD方法
	HttpMethodOptions HttpMethod = "OPTION" //OPTIONS方法
	HttpMethodDelete  HttpMethod = "DELETE" //DELETE方法
)

// 默认配置文件路径
const DefaultConfigPath = "web.config"

// 默认布局配置文件名
const DefaultLayoutConfigFileName = "Layout.config"

//默认模板文件扩展名
const DefaultTemplateExt = ".html"

//默认模板文件内模板名,用于返回部分视图时使用
const DefaultTemplateName = "Content"

//默认Cookie名
const DefaultSessionCookieName = "SessionId"

// Api格式
type ApiType string

const (
	ApiTypeJson ApiType = "Json" //Json
	ApiTypeXml  ApiType = "Xml"  //Xml
)
