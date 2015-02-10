package vmw

import (
	"testing"

	"github.com/coreos/coreos-cloudinit/Godeps/_workspace/src/github.com/sigma/vmw-ovflib"
)

var data_vapprun = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<Environment xmlns="http://schemas.dmtf.org/ovf/environment/1"
     xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
     xmlns:oe="http://schemas.dmtf.org/ovf/environment/1"
     oe:id="CoreOS-vmw">
   <PlatformSection>
      <Kind>vapprun</Kind>
      <Version>1.0</Version>
      <Vendor>VMware, Inc.</Vendor>
      <Locale>en_US</Locale>
   </PlatformSection>
   <PropertySection>
      <Property oe:key="guestinfo.user_data.url" oe:value="https://gist.githubusercontent.com/sigma/5a64aac1693da9ca70d2/raw/plop.yaml"/>
      <Property oe:key="guestinfo.user_data.doc" oe:value=""/>
      <Property oe:key="guestinfo.meta_data.url" oe:value=""/>
      <Property oe:key="guestinfo.meta_data.doc" oe:value=""/>
   </PropertySection>
</Environment>`)

func TestAvailable(t *testing.T) {
	env := ovf.ReadEnvironment(data_vapprun)
	gi := guestInfo{env, readVariable, readUrlBody}

	if !gi.IsAvailable() {
		t.Fatal("vmw-guestinfo datasource is unavailable")
	}
}
