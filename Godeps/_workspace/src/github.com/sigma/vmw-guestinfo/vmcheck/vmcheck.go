package vmcheck

import (
	"github.com/coreos/coreos-cloudinit/Godeps/_workspace/src/github.com/sigma/vmw-guestinfo/bridge"
)

func IsVirtualWorld() bool {
	return bridge.VmCheckIsVirtualWorld()
}

func GetVersion() (version uint32, typ uint32) {
	return bridge.VmCheckGetVersion()
}
