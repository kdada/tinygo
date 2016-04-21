package config

import (
	"strconv"
	"strings"
)

// iniConfigPaser ini类型文件配置解析
func iniConfigPaser(data []byte) (Config, error) {
	var config = newIniConfig()
	var currentSection = config.GlobalSection().(*IniSection)
	var dataLength = len(data)
	if dataLength >= 3 && data[0] == 239 && data[1] == 187 && data[2] == 191 {
		//去除BOM头
		data = data[3:]
		dataLength -= 3
	}
	for i := 0; i < dataLength; i++ {
		var char = string(data[i])
		//忽略空白字符
		if char == " " || char == "\n" || char == "\r" {
			continue
		}
		switch char {
		case "#", ";":
			{
				//处理以#和;开头的注释文本
				for j := i + 1; j < dataLength; j++ {
					if string(data[j]) == "\n" {
						i = j
						break
					}
				}
			}
		case "[":
			{
				//处理[]段
				var found = false
				for j := i + 1; j < dataLength; j++ {
					if !found && string(data[j]) == "]" {
						var sectionName = string(data[i+1 : j])
						var section = newIniSection(sectionName)
						config.sections[sectionName] = section
						currentSection = section
						found = true
					}
					if string(data[j]) == "\n" || j == dataLength-1 {
						if !found {
							return nil, ErrorNotMatch.Format("]").Error()
						}
						i = j
						break
					}
				}
			}
		default:
			{
				//处理key=value形式的内容
				var keyFound = false
				var keyEndPos = 0
				var key = ""
				var value = ""
				for j := i + 1; j < dataLength; j++ {
					if !keyFound && string(data[j]) == "=" {
						key = string(data[i:j])
						keyEndPos = j
						keyFound = true
					}

					if string(data[j]) == "\n" || j == dataLength-1 {
						if keyFound {
							value = string(data[keyEndPos+1 : j])
							//忽略两头的空白字符
							key = strings.TrimSpace(key)
							value = strings.TrimSpace(value)
							currentSection.add(key, value)
						} else {
							return nil, ErrorNotMatch.Format("=").Error()
						}
						i = j
						break
					}
				}
			}
		}
	}
	return config, nil
}

// Ini配置
type IniConfig struct {
	globalSection *IniSection            //全局段
	sections      map[string]*IniSection //命名段
}

// newIniConfig 创建ini配置
func newIniConfig() *IniConfig {
	return &IniConfig{
		newIniSection(""),
		make(map[string]*IniSection),
	}
}

// GlobalSection 获取全局配置段
func (this *IniConfig) GlobalSection() Section {
	return this.globalSection
}

// Section 根据name获取指定名称的配置段
func (this *IniConfig) Section(name string) (Section, bool) {
	var section, ok = this.sections[name]
	return section, ok
}

// Ini配置段
type IniSection struct {
	name string            //段名称
	kvs  map[string]string //键值map
}

// newIniSection 创建ini配置段
func newIniSection(name string) *IniSection {
	var section = &IniSection{
		name,
		make(map[string]string),
	}
	return section
}

// Name 配置段名称
func (this *IniSection) Name() string {
	return this.name
}

// add 添加键值对,对于同名key,后添加的有效
func (this *IniSection) add(key, value string) {
	this.kvs[key] = value
}

// String 获取字符串
func (this *IniSection) String(key string) (string, error) {
	var value, ok = this.kvs[key]
	if ok {
		return value, nil
	}
	return "", ErrorInvalidKey.Format(key).Error()
}

// Int 获取整数
func (this *IniSection) Int(key string) (int64, error) {
	var value, ok = this.kvs[key]
	if ok {
		var result, err = strconv.ParseInt(value, 0, 64)
		if err == nil {
			return result, nil
		}
		return 0, ErrorInvalidTypeConvertion.Format(value, "int64").Error()
	}
	return 0, ErrorInvalidKey.Format(key).Error()
}

// Bool 获取布尔值
func (this *IniSection) Bool(key string) (bool, error) {
	var value, ok = this.kvs[key]
	if ok {
		var result, err = strconv.ParseBool(value)
		if err == nil {
			return result, nil
		}
		return false, ErrorInvalidTypeConvertion.Format(value, "bool").Error()
	}
	return false, ErrorInvalidKey.Format(key).Error()
}

// Float 获取浮点值
func (this *IniSection) Float(key string) (float64, error) {
	var value, ok = this.kvs[key]
	if ok {
		var result, err = strconv.ParseFloat(value, 64)
		if err == nil {
			return result, nil
		}
		return 0.0, ErrorInvalidTypeConvertion.Format(value, "float64").Error()
	}
	return 0.0, ErrorInvalidKey.Format(key).Error()
}
