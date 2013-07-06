package tesseract

// #include "capi.h"
// #include <stdlib.h>
import "C"

const version = "v0.1"

func Version() string {
	libTessVersion := C.TessVersion()
	defer C.free(libTessVersion)
	return version + C.GoString(libTessVersion)
}
