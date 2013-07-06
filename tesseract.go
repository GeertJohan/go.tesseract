package tesseract

// #include "tesseract/capi.h"
// #include <stdlib.h>
// #cgo LDFLAGS: -ltesseract
import "C"

const version = "0.1"

func Version() string {
	libTessVersion := C.TessVersion()
	// defer C.free(unsafe.Pointer(libTessVersion))
	return "go.tesseract:" + version + "  tesseract lib:" + C.GoString(libTessVersion)
}
