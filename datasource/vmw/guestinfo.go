package vmw

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/coreos/coreos-cloudinit/Godeps/_workspace/src/github.com/sigma/vmw-guestinfo/rpcvmx"
	"github.com/coreos/coreos-cloudinit/Godeps/_workspace/src/github.com/sigma/vmw-guestinfo/vmcheck"
	"github.com/coreos/coreos-cloudinit/Godeps/_workspace/src/github.com/sigma/vmw-ovflib"
)

type guestInfo struct {
	env       *ovf.OvfEnvironment
	varReader func(string, *ovf.OvfEnvironment) (string, bool)
	urlReader func(string) []byte
}

func readVariable(var_name string, ovf_env *ovf.OvfEnvironment) (string, bool) {
	if val, ok := ovf_env.Properties["guestinfo."+var_name]; ok {
		return val, ok && val != ""
	} else if vmcheck.IsVirtualWorld() {
		val := rpcvmx.ConfigGetString(var_name, "")
		return val, val != ""
	}
	return "", false
}

func readUrlBody(url string) []byte {
	log.Printf("Reading from url %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Url unavailable")
		return make([]byte, 0)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body")
		return make([]byte, 0)
	}
	return body
}

func NewDatasource(filename string) *guestInfo {
	var ovf_env []byte

	if filename == "" {
		if vmcheck.IsVirtualWorld() {
			log.Println("Trying to read from VMware backdoor")
			ovf_env = []byte(rpcvmx.ConfigGetString("ovfenv", ""))
		} else {
			log.Println("Not in a VMware world, giving up")
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

	return &guestInfo{env, readVariable, readUrlBody}
}

func (gi *guestInfo) IsAvailable() bool {
	vars := []string{
		"user_data.doc", "user_data.url",
		"meta_data.doc", "meta_data.url"}
	for _, v := range vars {
		_, ok := readVariable(v, gi.env)
		if ok {
			return true
		}
	}
	log.Println("vmw-guestinfo datasource is not available")
	return false
}

func (gi *guestInfo) AvailabilityChanges() bool {
	return false
}

func (gi *guestInfo) ConfigRoot() string {
	return ""
}

func (gi *guestInfo) fetchData(key string) ([]byte, error) {
	val, ok := readVariable(key+".doc", gi.env)
	if ok {
		log.Println("Direct document available")
		return []byte(val), nil
	} else if val, ok = readVariable(key+".url", gi.env); ok {
		log.Println("Url available")
		return readUrlBody(val), nil
	}
	return make([]byte, 0), nil
}

func (gi *guestInfo) FetchMetadata() ([]byte, error) {
	log.Println("Reading metadata")
	return gi.fetchData("meta_data")
}

func (gi *guestInfo) FetchUserdata() ([]byte, error) {
	log.Println("Reading user data")
	return gi.fetchData("user_data")
}

func (gi *guestInfo) FetchNetworkConfig(filename string) ([]byte, error) {
	return nil, nil
}

func (gi *guestInfo) Type() string {
	return "vmw-guestinfo"
}
