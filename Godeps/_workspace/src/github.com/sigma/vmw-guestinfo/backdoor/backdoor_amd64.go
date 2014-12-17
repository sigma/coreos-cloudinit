package backdoor

/*
#cgo CFLAGS: -I../include
#include <stdlib.h>
#include "backdoor.h"
*/
import "C"

type BackdoorProtoInSize C.size_t

type BackdoorProtoReg uint64
