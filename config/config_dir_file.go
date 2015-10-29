// config_dir.go
package config

import (
	"os"
	"path/filepath"
)

//得到服务器运行目录
func GetServerDir() string {
	dir, _ := os.Getwd()
	return dir
}

//得到config目录
func GetConfigDir() string {
	return filepath.Join(GetServerDir(), "config", "config_file")
}

//得到setting目录
func GetSettingDir() string {
	return filepath.Join(GetServerDir(), "config", "setting")
}

func getConfigFileName(fileShorName string) string {
	return filepath.Join(GetConfigDir(), fileShorName)
}

func getSettingFileName(fileShorName string) string {
	return filepath.Join(GetConfigDir(), fileShorName)
}
