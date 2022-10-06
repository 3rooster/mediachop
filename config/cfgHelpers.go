package config

import (
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
)

// BinDir exec bin file dir
var BinDir string

// RootDir root dir of program
var RootDir string

// GetAbsPath get abs file or dir path ，
// if path start with '/' path return path, or return path with RootDir in front
func GetAbsPath(path string) string {
	if path == "" {
		return RootDir + "/"
	}
	if strings.Contains(path, ":") || strings.HasPrefix(path, "/") {
		return path
	}
	if RootDir != "" {
		return RootDir + "/" + path
	}
	return path
}

func initRootDir() {
	if isInIDE() {
		RootDir = ""
		BinDir = ""
		return
	}
	path := os.Args[0]
	dir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		panic(err)
	}

	BinDir = dir
	RootDir = filepath.Dir(BinDir)
	zap.S().Infof("bin dir is [ %s ] \nroot dir is : [ %s ]", BinDir, RootDir)
}

func isInIDE() bool {
	path := os.Args[0]
	if strings.Contains(path, "go_build") { // in jetbrains IDE
		zap.S().Infof("run by IDE, set root dir empty ")
		return true
	}
	return false
}
