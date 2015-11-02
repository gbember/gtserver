package proto

const (
	SC_ROLE_INFO = uint16(10101)
)

type Sc_role_info struct {
	RoleID        int32
	RoleName      string
	Sex           int8
	Level         uint8
	Exp           int32
	VipLevel      int8
	Gold          int32
	Coin          int32
	LastLogoutSec int32
}

func (r *Sc_role_info) Read(p *Packet) error {
	value0, err := p.readInt32()
	if err != nil {
		return err
	}
	r.RoleID = value0
	value1, err := p.readString()
	if err != nil {
		return err
	}
	r.RoleName = value1
	value2, err := p.readInt8()
	if err != nil {
		return err
	}
	r.Sex = value2
	value3, err := p.readUint8()
	if err != nil {
		return err
	}
	r.Level = value3
	value4, err := p.readInt32()
	if err != nil {
		return err
	}
	r.Exp = value4
	value5, err := p.readInt8()
	if err != nil {
		return err
	}
	r.VipLevel = value5
	value6, err := p.readInt32()
	if err != nil {
		return err
	}
	r.Gold = value6
	value7, err := p.readInt32()
	if err != nil {
		return err
	}
	r.Coin = value7
	value8, err := p.readInt32()
	if err != nil {
		return err
	}
	r.LastLogoutSec = value8
	return nil
}
func (r *Sc_role_info) WriteMsgID(p *Packet) {
	p.writeUint16(SC_ROLE_INFO)
}
func (r *Sc_role_info) Write(p *Packet) {
	p.writeInt32(r.RoleID)
	p.writeString(r.RoleName)
	p.writeInt8(r.Sex)
	p.writeUint8(r.Level)
	p.writeInt32(r.Exp)
	p.writeInt8(r.VipLevel)
	p.writeInt32(r.Gold)
	p.writeInt32(r.Coin)
	p.writeInt32(r.LastLogoutSec)
}
