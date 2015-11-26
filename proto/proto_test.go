// proto_test.go
package proto

import (
	"testing"
	"time"
)

func TestProto(t *testing.T) {
	r := &Sc_role_info{
		RoleID:        1,
		RoleName:      "test_proto",
		Sex:           1,
		Level:         1,
		Exp:           0,
		VipLevel:      0,
		Gold:          0,
		Coin:          0,
		LastLogoutSec: 0,
	}
	pkWriter := NewWriter()
	bs := EncodeProtoPacket(r, pkWriter)
	t.Log("bs:", bs)

	var max int64 = 1000000
	starTime := time.Now()
	for i := max; i >= 0; i-- {
		EncodeProtoPacket(r, pkWriter)
	}
	td := time.Since(starTime)
	t.Logf("encode: %v === %d", td, td.Nanoseconds()/max)

	msgID, msg, _ := DecodeProto(bs)
	t.Log(msgID, msg)
	starTime = time.Now()
	for i := max; i >= 0; i-- {
		DecodeProto(bs)
	}
	td = time.Since(starTime)
	t.Logf("decode: %v === %d", td, td.Nanoseconds()/max)
}
