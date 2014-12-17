package vmw

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/coreos/coreos-cloudinit/Godeps/_workspace/src/github.com/sigma/vmw-guestinfo/rpcvmx"
	"github.com/coreos/coreos-cloudinit/Godeps/_workspace/src/github.com/sigma/vmw-guestinfo/vmcheck"
	"github.com/coreos/coreos-cloudinit/Godeps/_workspace/src/github.com/sigma/vmw-ovflib"
)

type guestInfo struct {
	user_data []byte
	meta_data []byte
}

func readVariable(var_name string, ovf_env *ovf.OvfEnvironment) (string, bool) {
	if val, ok := ovf_env.Properties["guestinfo."+var_name]; ok {
		return val, ok
	} else if vmcheck.IsVirtualWorld() {
		val := rpcvmx.ConfigGetString(var_name, "")
		return val, val != ""
	}
	return "", false
}

func readUrlBody(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		return make([]byte, 0)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return make([]byte, 0)
	}
	return body
}

func NewDatasource(filename string) *guestInfo {
	var ovf_env []byte
	gi := &guestInfo{}

	if filename == "" {
		if vmcheck.IsVirtualWorld() {
			ovf_env = []byte(rpcvmx.ConfigGetString("ovfenv", ""))
		} else {
			ovf_env = make([]byte, 0)
		}
	} else {
		var err error
		ovf_env, err = ioutil.ReadFile(filename)
		if err != nil {
			ovf_env = make([]byte, 0)
		}
	}

	env := &ovf.OvfEnvironment{}
	if len(ovf_env) != 0 {
		env = ovf.ReadEnvironment(ovf_env)
	}

	val, ok := readVariable("user_data.doc", env)
	if ok {
		gi.user_data = []byte(val)
	} else if val, ok = readVariable("user_data.url", env); ok {
		gi.user_data = readUrlBody(val)
	}

	val, ok = readVariable("meta_data.doc", env)
	if ok {
		gi.user_data = []byte(val)
	} else if val, ok = readVariable("meta_data.url", env); ok {
		gi.meta_data = readUrlBody(val)
	}

	return gi
}

func (gi *guestInfo) IsAvailable() bool {
	return len(gi.user_data) != 0 || len(gi.meta_data) != 0
}

func (gi *guestInfo) AvailabilityChanges() bool {
	return false
}

func (gi *guestInfo) ConfigRoot() string {
	return ""
}

func (gi *guestInfo) FetchMetadata() ([]byte, error) {
	var err error
	if len(gi.meta_data) == 0 {
		err = errors.New("No metadata")
	} else {
		err = nil
	}
	return gi.meta_data, err
}

func (gi *guestInfo) FetchUserdata() ([]byte, error) {
	var err error
	if len(gi.user_data) == 0 {
		err = errors.New("No metadata")
	} else {
		err = nil
	}
	return gi.user_data, err
}

func (gi *guestInfo) FetchNetworkConfig(filename string) ([]byte, error) {
	return nil, nil
}

func (gi *guestInfo) Type() string {
	return "vmw-guestinfo"
}
