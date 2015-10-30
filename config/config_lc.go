// config_lc.go
package config

import "github.com/gbember/gt/console/command"

//配置文件更新命令
type commandLC struct{}

func init() {
	lc := new(commandLC)
	command.Register(lc)
}

func (*commandLC) Name() string {
	return "lc"
}
func (*commandLC) Help() string {
	return "config refresh command: lc [configName]"
}
func (*commandLC) Run(args []string) string {
	if len(args) == 0 {
		err := LoadAllConfig()
		if err != nil {
			return err.Error()
		}
	} else {
		err := LoadConfig(args[0])
		if err != nil {
			return err.Error()
		}
	}
	return "ok"
}
