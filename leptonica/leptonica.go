package leptonica

/*
#cgo pkg-config: lept
#include <allheaders.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

// Pix wraps a Leptonica PIX pointer.
// The underlying memory must be released with Close.
type Pix struct {
	ptr *C.PIX
}

// ImageType maps to Leptonica IFF_* constants for encoding.
type ImageType int32

const (
	UNKNOWN ImageType = ImageType(C.IFF_UNKNOWN)
	DEFAULT ImageType = ImageType(C.IFF_DEFAULT)
	BMP     ImageType = ImageType(C.IFF_BMP)
	JPEG    ImageType = ImageType(C.IFF_JFIF_JPEG)
	PNG     ImageType = ImageType(C.IFF_PNG)
	TIFF    ImageType = ImageType(C.IFF_TIFF)
	WEBP    ImageType = ImageType(C.IFF_WEBP)
	JP2     ImageType = ImageType(C.IFF_JP2)
	PNM     ImageType = ImageType(C.IFF_PNM)
)

// NewPixFromFile loads an image from disk into a PIX.
func NewPixFromFile(filename string) (*Pix, error) {
	cfn := C.CString(filename)
	defer C.free(unsafe.Pointer(cfn))

	pix := C.pixRead(cfn)
	if pix == nil {
		return nil, fmt.Errorf("leptonica: failed to read image %q", filename)
	}

	p := &Pix{ptr: pix}
	runtime.SetFinalizer(p, (*Pix).Close)
	return p, nil
}

// Close releases the underlying PIX. Safe to call multiple times.
func (p *Pix) Close() {
	if p == nil || p.ptr == nil {
		return
	}

	C.pixDestroy(&p.ptr)
	p.ptr = nil
	runtime.SetFinalizer(p, nil)
}
