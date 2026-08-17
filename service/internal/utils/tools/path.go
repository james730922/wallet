package tools

import (
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"os"
	"path/filepath"
	"runtime"
)

func InitRootFolder(path string) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(filename), path)
	err := os.Chdir(dir)
	if err != nil {
		logger.ApLog().Error(err)
	}
}
