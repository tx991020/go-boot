package generator

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// PrintErr 打印错误
func PrintErr(err error) {
	fmt.Printf("%c[%d;%d;%dm%s(错误: %s )%c[0m \n", 0x1B, 41, 37, 1, "", err, 0x1B)
}

// PrintInfo 打印信息
func PrintInfo(a ...interface{}) {
	fmt.Printf("%c[%d;%d;%dm%s%s %c[0m \n", 0x1B, 46, 37, 1, "", a, 0x1B)
}

// PrintInfoMap 打印成对信息
func PrintInfoMap(a map[string]interface{}) {
	for k, v := range a {
		fmt.Printf("%c[%d;%d;%dm%s%s: %s;  %c[0m \n", 0x1B, 46, 37, 1, "", k, Interface2String(v), 0x1B)
	}
}

// Interface2String 接口转字符串
func Interface2String(i interface{}) string {
	switch i.(type) {
	case string:
		return i.(string)
	case int:
		return strconv.Itoa(i.(int))
	case int64:
		return strconv.FormatInt(i.(int64), 10)
	case float64:
		return strconv.FormatFloat(i.(float64), 'f', 0, 64)
	case bool:
		return strconv.FormatBool(i.(bool))
	default:
		return ""
	}
}

// Struct2Map struct -> map
func Struct2Map(a interface{}) (m map[string]interface{}, err error) {
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &m)
	return
}

// MergeMap 合并多个map
func MergeMap(a ...map[string]interface{}) map[string]interface{} {
	m := map[string]interface{}{}
	for _, v := range a {
		for k, val := range v {
			m[k] = val
		}
	}
	return m
}

