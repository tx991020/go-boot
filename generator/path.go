package generator

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// ErrPath 路径错误
var ErrPath = errors.New("未找到 GOPATH , 请确定已设置GOPATH")

func GetEnvPath(aPath,tplPath string) (appPath string, templatePath string, err error) {
	var cmd *exec.Cmd
	cmd = exec.Command("/bin/sh", "-c", `go env | grep GOPATH | awk -F    '"' '{print $2}'`) // 获取本机 go env
	path, err := cmd.Output()
	if err != nil {
		return "", "", ErrPath
	}
	p := strings.TrimSpace(string(path))
	appPath = fmt.Sprintf("%s%s/src",p,aPath)
	templatePath = tplPath

	PrintInfo("项目所在位置: ", appPath)
	return appPath, templatePath, nil
}
