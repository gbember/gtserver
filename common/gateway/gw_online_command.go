// gw_online_command.go
package gateway

import (
	"strconv"

	"github.com/gbember/gt/console/command"
)

type _gw_online struct {
}

func init() {
	command.Register(new(_gw_online))
}

func (*_gw_online) Name() string {
	return "on"
}
func (*_gw_online) Help() string {
	return "online num"
}
func (*_gw_online) Run([]string) string {
	return strconv.Itoa(GetOnlineNum())
}
