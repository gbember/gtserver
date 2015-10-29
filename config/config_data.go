// config_data.go
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var (
	_register_config_map = make(map[string]func() error)
	//setting
	DataSetting  *dataSetting
	DataLevelExp map[int]int
)

func init() {
	register("setting", load_data_setting)
	register("load_data_level_exp", load_data_level_exp)
}

//加载所有配置
func LoadAllConfig() error {
	for _, fun := range _register_config_map {
		err := fun()
		if err != nil {
			return err
		}
	}
	return nil
}

//加载单个配置
func LoadConfig(fileName string) error {
	suffix := ".config"
	if strings.HasSuffix(fileName, ".config") {
		fileName = fileName[:len(fileName)-len(suffix)]
	}
	if fun, ok := _register_config_map[fileName]; ok {
		return fun()
	}
	return errors.New(fmt.Sprintf("没有找到注册的配置:%s", fileName))
}

func load_data_setting() error {
	_dataSetting := &dataSetting{}
	err := load_setting("data_setting", _dataSetting)
	if err != nil {
		return err
	}
	DataSetting = _dataSetting
	return nil
}

func load_data_level_exp() error {
	_dataLevelExpList := make([]dataLevelExp, 200)
	err := load_config("data_level_exp", &_dataLevelExpList)
	if err != nil {
		return err
	}
	_dataLevelExp := make(map[int]int)
	for _, v := range _dataLevelExpList {
		_dataLevelExp[v.Level] = v.MaxExp
	}
	DataLevelExp = _dataLevelExp
	return nil
}

func register(fileName string, fun func() error) {
	_register_config_map[fileName] = fun
}

func load_config(fileName string, v interface{}) error {
	fileName = fileName + ".config"
	fileName1 := filepath.Join(GetConfigDir(), fileName)
	bs, err := ioutil.ReadFile(fileName1)
	if err != nil {
		return errors.New(fmt.Sprintf("加载config配置[%s]错误:%s", fileName, err.Error()))
	}
	err = json.Unmarshal(bs, v)
	if err != nil {
		return errors.New(fmt.Sprintf("加载config配置[%s]错误:%s", fileName, err.Error()))
	}
	return nil
}

func load_setting(fileName string, v interface{}) error {
	fileName = fileName + ".config"
	fileName1 := filepath.Join(GetSettingDir(), fileName)
	bs, err := ioutil.ReadFile(fileName1)
	if err != nil {
		return errors.New(fmt.Sprintf("加载setting配置[%s]错误:%s", fileName, err.Error()))
	}
	err = json.Unmarshal(bs, v)
	if err != nil {
		return errors.New(fmt.Sprintf("加载setting配置[%s]错误:%s", fileName, err.Error()))
	}
	return nil
}
