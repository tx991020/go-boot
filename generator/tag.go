package generator

import "github.com/iancoleman/strcase"

// FormatName 转驼峰&特有名词
func FormatName(name string) (strcutName string, tagName string) {

	return strcase.ToCamel(name), strcase.ToLowerCamel(name)
}

