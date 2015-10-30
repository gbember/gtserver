// config_type.go
package config

type dataSetting struct {
	ServerID         int16  //服务器编号
	IsRealese        bool   //是否发布版本
	LogDir           string //日志目录
	LogLevel         int    //日志等级
	GatewayPort      int    //网关监听端口
	PProfPort        int    //pprof http端口
	MaxOnlineNum     int    //最高在线人数
	DBAddrIP         string //数据库网络地址
	DBAddrPort       int    //数据库端口
	DBName           string //数据库名
	DBUser           string //数据库用户名
	DBPassword       string //数据库密码
	LoginAuthKey     string //登陆验证Key
	HeadLen          int    //消息长度字节数
	MsgIDLen         int    //消息ID字节数
	MaxDataLen       int    //接收客户端消息最大字节数
	ConsoleAddr      string //console监听地址
	MsgMaxSpeedSec   int64  //速率控制:多少秒
	MsgMaxSpeedNum   int    //速率控制:多少个
	MsgMaxSpeedLager int    //速率超出次数控制:
}

type dataLevelExp struct {
	Level  int //等级
	MaxExp int //等级最高经验
}
