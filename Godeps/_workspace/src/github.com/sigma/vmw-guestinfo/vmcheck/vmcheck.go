package vmcheck

import (
	"github.com/coreos/coreos-cloudinit/Godeps/_workspace/src/github.com/sigma/vmw-guestinfo/bridge"
)

func IsVirtualWorld() bool {
	return bridge.VmCheckIsVirtualWorld()
}

func GetVersion() (uint32, uint32) {
	return bridge.VmCheckGetVersion()
}
