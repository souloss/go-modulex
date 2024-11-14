package pkgscan

import (
	"os"
	"path/filepath"
	"runtime"
)

func FindProjectRoot() string {
	_, dir, _, _ := runtime.Caller(0)
	for {
		// 检查当前目录是否包含go.mod
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		// 获取上级目录
		parentDir := filepath.Dir(dir)
		// 如果已经到达文件系统的根目录，则停止
		if parentDir == dir {
			break
		}
		dir = parentDir
	}
	return ""
}
