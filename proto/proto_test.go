// proto_test.go
package proto

import (
	"testing"
	"time"
)

func TestProto1(t *testing.T) {
	d := &Cs_account_login{AccountName: "sahojwkohrjqohskagopjsojh", Ticket: "sghiahiwngineur9uhnng	HGWUIHGWH"}
	t.Log(EncodeProto(d))
	var max int64 = 1000000
	startTime := time.Now()
	for i := max; i > 0; i-- {
		EncodeProto(d)
	}
	t.Log(time.Since(startTime).Nanoseconds() / max)
}

//func TestProto2(t *testing.T) {
//	d := &Cs_account_login{AccountName: "sahojwkohrjqohskagopjsojh", Ticket: "sghiahiwngineur9uhnng	HGWUIHGWH"}
//	p := NewWriter()
//	t.Log(EncodeProtoPacket(d, p))
//	var max int64 = 10000000
//	startTime := time.Now()
//	for i := max; i > 0; i-- {
//		EncodeProtoPacket(d, p)
//	}
//	t.Log(time.Since(startTime).Nanoseconds() / max)
//}
