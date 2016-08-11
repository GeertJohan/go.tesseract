package tesseract

// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

func gobool(b C.int) bool {
	if b != 0 {
		return true
	}

	return false
}

func cStringVectorToStringslice(cStringVector **C.char) []string {
	// get pointer size to do iteration
	cPtrSize := unsafe.Sizeof(cStringVector)

	// create results string slice
	result := make([]string, 0)

	// iterate over **char
	for {
		// check for null terminator
		if *cStringVector == nil {
			return result
		}

		// add string to result slice
		result = append(result, C.GoString(*cStringVector))

		// increment pointer to next index
		cStringVectorPtr := uintptr(unsafe.Pointer(cStringVector))
		cStringVectorPtr += cPtrSize
		cStringVector = (**C.char)(unsafe.Pointer(cStringVectorPtr))
	}
}
