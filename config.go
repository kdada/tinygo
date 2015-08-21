package tinygo

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"
	"tinygo/info"
)

// Tiny配置
var tinyConfig = struct {
	Port         uint16   //监听端口
	StaticPath   []string //静态文件目录
	ViewPath     string   //视图文件目录
	PageError    string   //默认出错页面
	CompilePages bool     //是否预编译页面
	Api          string   //使用Api返回的数据的解析格式 Json 或 Xml

}{
	Port:         80,                       //默认为80端口
	StaticPath:   []string{"Content"},      //默认为Content目录
	ViewPath:     "View",                   //默认为View目录
	PageError:    "",                       //无默认页面
	CompilePages: true,                     //默认编译视图
	Api:          string(info.ApiTypeJson), //默认为Json
}

// loadConfig 加载配置
func loadConfig() {
	var config = info.DefaultConfigPath
	var content, err = ioutil.ReadFile(config)
	if err == nil {
		err = json.Unmarshal(content, &tinyConfig)
	} 
	if err != nil {
		fmt.Println(config,err)
	}
}

// getViewFilePath 返回视图文件路径的正确路径
// layout/layout.html ==> tinyConfig.ViewPath/layout/layout.html
func getViewFilePath(filePath string) string {
	return tinyConfig.ViewPath + "/" + filePath
}

// generateViewFilePath 返回视图文件相对tinyConfig.ViewPath的路径
// tinyConfig.ViewPath/layout/layout.html ==> layout/layout.html
func generateViewFilePath(filePath string) string {
	filePath = strings.TrimPrefix(filePath, tinyConfig.ViewPath)
	filePath = strings.TrimLeft(filePath, `\/`)
	filePath = translatePath(filePath)
	return filePath
}

// translatePath 将\风格的路径转换为/风格的路径
func translatePath(path string) string {
	// 即 test\some.html  ==>  test/some.html
	path = strings.Replace(path, `\`, `/`, -1)
	return path
}

// 视图布局配置
// 为了顺利解析,即使在不同目录中
// layout文件不得与任何视图文件重名,layout文件之间也不允许重名
var layoutConfig = struct {
	LayoutMap     map[string]string //为布局文件定义别名
	DefaultLayout string            //默认布局文件别名
	LayoutSpec    map[string]string //特别指定目录或文件对应的布局文件,map[目录或文件]布局文件别名
}{
	LayoutMap:     map[string]string{},
	DefaultLayout: "",
	LayoutSpec:    map[string]string{},
}

// loadLayoutConfig 加载布局配置
func loadLayoutConfig() {
	var config = tinyConfig.ViewPath + "/" + info.DefaultLayoutConfigFileName
	var content, err = ioutil.ReadFile(config)
	if err == nil {
		err = json.Unmarshal(content, &layoutConfig)
	}
	if err != nil {
		fmt.Println(config,err)
	}
}

// isLayoutFile 判断文件是否是布局文件
// filePath:相对于tinyConfig.ViewPath的文件路径,分隔符必须为/
func isLayoutFile(filePath string) bool {
	for _, layoutPath := range layoutConfig.LayoutMap {
		if layoutPath == filePath {
			return true
		}
	}
	return false
}

// getLayout 获取指定文件的Layout文件
// filePath:相对于tinyConfig.ViewPath的文件路径,分隔符必须为/
func getLayoutFile(filePath string) (string, bool) {
	var layout, ok = layoutConfig.LayoutSpec[filePath]
	if !ok {
		var dir = path.Dir(filePath)
		layout, ok = layoutConfig.LayoutSpec[dir]
		if !ok {
			layout = layoutConfig.DefaultLayout
		}
	}
	layout, ok = layoutConfig.LayoutMap[layout]
	if filePath == layout {
		//避免循环布局
		return "", false
	}
	return layout, ok
}
