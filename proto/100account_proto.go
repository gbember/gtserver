package proto

const (
	CS_ACCOUNT_LOGIN  = uint16(10001)
	SC_ACCOUNT_LOGIN  = uint16(10002)
	CS_ACCOUNT_CREATE = uint16(10003)
	SC_ACCOUNT_CREATE = uint16(10004)
	SC_ACCOUNT_KICK   = uint16(10005)
	CS_ACCOUNT_HEART  = uint16(10006)
	SC_ACCOUNT_HEART  = uint16(10007)
	CS_ACCOUNT_LOGOUT = uint16(10008)
	CS_ACCOUNT_TEST   = uint16(10009)
	SC_ACCOUNT_TEST   = uint16(10010)
)

//请求登录
type Cs_account_login struct {
	////平台用户ID
	UserID int32
	////unix时间戳
	UnixTime int32
	////平台用户账号
	AccountName string
	////验证密钥
	Ticket string
	////服务器ID
	ServerID int16
}

func (r *Cs_account_login) Read(p *Packet) error {
	value0, err := p.readInt32()
	if err != nil {
		return err
	}
	r.UserID = value0
	value1, err := p.readInt32()
	if err != nil {
		return err
	}
	r.UnixTime = value1
	value2, err := p.readString()
	if err != nil {
		return err
	}
	r.AccountName = value2
	value3, err := p.readString()
	if err != nil {
		return err
	}
	r.Ticket = value3
	value4, err := p.readInt16()
	if err != nil {
		return err
	}
	r.ServerID = value4
	return nil
}
func (r *Cs_account_login) WriteMsgID(p *Packet) {
	p.writeUint16(CS_ACCOUNT_LOGIN)
}
func (r *Cs_account_login) Write(p *Packet) {
	p.writeInt32(r.UserID)
	p.writeInt32(r.UnixTime)
	p.writeString(r.AccountName)
	p.writeString(r.Ticket)
	p.writeInt16(r.ServerID)
}

type Sc_account_login struct {
	// 0 => 成功
	// 1 => 已经登录
	// 2 => 登录的服务器id不符
	// 3 => key超时
	// 4 => key错误
	Result int8
	////true=已经创建主公，false=未创建主公
	IsCreate bool
}

func (r *Sc_account_login) Read(p *Packet) error {
	value0, err := p.readInt8()
	if err != nil {
		return err
	}
	r.Result = value0
	value1, err := p.readBool()
	if err != nil {
		return err
	}
	r.IsCreate = value1
	return nil
}
func (r *Sc_account_login) WriteMsgID(p *Packet) {
	p.writeUint16(SC_ACCOUNT_LOGIN)
}
func (r *Sc_account_login) Write(p *Packet) {
	p.writeInt8(r.Result)
	p.writeBool(r.IsCreate)
}

//创建角色
type Cs_account_create struct {
	////名字
	RoleName string
	////0:女 1:男
	Sex int8
}

func (r *Cs_account_create) Read(p *Packet) error {
	value0, err := p.readString()
	if err != nil {
		return err
	}
	r.RoleName = value0
	value1, err := p.readInt8()
	if err != nil {
		return err
	}
	r.Sex = value1
	return nil
}
func (r *Cs_account_create) WriteMsgID(p *Packet) {
	p.writeUint16(CS_ACCOUNT_CREATE)
}
func (r *Cs_account_create) Write(p *Packet) {
	p.writeString(r.RoleName)
	p.writeInt8(r.Sex)
}

type Sc_account_create struct {
	//0 => 创建成功
	//1 => 没有登录
	//2 => 用户已经创建角色
	//3 => 性别错误
	//4 => 角色名称长度为2~6个字符
	//5 => 名字只能是字母数字和汉子组合
	//6 => 角色名称已经被使用
	Result int8
}

func (r *Sc_account_create) Read(p *Packet) error {
	value0, err := p.readInt8()
	if err != nil {
		return err
	}
	r.Result = value0
	return nil
}
func (r *Sc_account_create) WriteMsgID(p *Packet) {
	p.writeUint16(SC_ACCOUNT_CREATE)
}
func (r *Sc_account_create) Write(p *Packet) {
	p.writeInt8(r.Result)
}

//强制下线通知
type Sc_account_kick struct {
	//1 => 服务器人数已满
	//2 => 服务器关闭
	//3 => 异地登陆
	//4 => 发消息太频繁
	Reason int8
}

func (r *Sc_account_kick) Read(p *Packet) error {
	value0, err := p.readInt8()
	if err != nil {
		return err
	}
	r.Reason = value0
	return nil
}
func (r *Sc_account_kick) WriteMsgID(p *Packet) {
	p.writeUint16(SC_ACCOUNT_KICK)
}
func (r *Sc_account_kick) Write(p *Packet) {
	p.writeInt8(r.Reason)
}

//心跳包
type Cs_account_heart struct {
}

func (r *Cs_account_heart) Read(p *Packet) error {
	return nil
}
func (r *Cs_account_heart) WriteMsgID(p *Packet) {
	p.writeUint16(CS_ACCOUNT_HEART)
}
func (r *Cs_account_heart) Write(p *Packet) {
}

type Sc_account_heart struct {
	////当前服务器时间
	UnixTime int32
}

func (r *Sc_account_heart) Read(p *Packet) error {
	value0, err := p.readInt32()
	if err != nil {
		return err
	}
	r.UnixTime = value0
	return nil
}
func (r *Sc_account_heart) WriteMsgID(p *Packet) {
	p.writeUint16(SC_ACCOUNT_HEART)
}
func (r *Sc_account_heart) Write(p *Packet) {
	p.writeInt32(r.UnixTime)
}

//注销登录
type Cs_account_logout struct {
}

func (r *Cs_account_logout) Read(p *Packet) error {
	return nil
}
func (r *Cs_account_logout) WriteMsgID(p *Packet) {
	p.writeUint16(CS_ACCOUNT_LOGOUT)
}
func (r *Cs_account_logout) Write(p *Packet) {
}

type Cs_account_test struct {
}

func (r *Cs_account_test) Read(p *Packet) error {
	return nil
}
func (r *Cs_account_test) WriteMsgID(p *Packet) {
	p.writeUint16(CS_ACCOUNT_TEST)
}
func (r *Cs_account_test) Write(p *Packet) {
}

type Sc_account_test struct {
}

func (r *Sc_account_test) Read(p *Packet) error {
	return nil
}
func (r *Sc_account_test) WriteMsgID(p *Packet) {
	p.writeUint16(SC_ACCOUNT_TEST)
}
func (r *Sc_account_test) Write(p *Packet) {
}
