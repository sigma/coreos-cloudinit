package backdoor

/*
#cgo CFLAGS: -I../include
#include <stdlib.h>
#include "backdoor.h"
*/
import "C"
import "unsafe"

type BackdoorProto C.Backdoor_proto

type BackdoorProtoIn BackdoorProto

func (proto *BackdoorProto) In() *BackdoorProtoIn {
	return (*BackdoorProtoIn)(proto)
}

func (in *BackdoorProtoIn) Size() *BackdoorProtoInSize {
	return (*BackdoorProtoInSize)(unsafe.Pointer(&in[8]))
}

func (size *BackdoorProtoInSize) Set(s BackdoorProtoInSize) {
	*size = s
}
