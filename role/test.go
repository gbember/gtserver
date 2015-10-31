// test.go
package role

import "github.com/gbember/gtserver/proto"

func init() {
	register_handle(proto.CS_ACCOUNT_TEST, test)
}

func test(rc *roleClient, msg proto.Messager) {
	m := new(proto.Sc_account_test)
	rc.gw.Send(m)
}
