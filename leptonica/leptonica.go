package leptonica

/*
#cgo pkg-config: lept
#include <allheaders.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"log/slog"
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

// Box wraps a Leptonica BOX pointer, representing a rectangular region.
type Box struct {
	ptr *C.BOX
}

// Ensure [*Box] implements [slog.LogValuer].
var _ slog.LogValuer = &Box{}

// Close releases the underlying BOX. Safe to call multiple times.
func (b *Box) Close() {
	if b == nil || b.ptr == nil {
		return
	}

	C.boxDestroy(&b.ptr)

	b.ptr = nil
	runtime.SetFinalizer(b, nil)
}

// LogValue returns a logger representation.
func (b *Box) LogValue() slog.Value {
	if b == nil || b.ptr == nil {
		return slog.Value{}
	}

	x, y, w, h := b.getGeometry()
	return slog.GroupValue(
		slog.Int("x", x),
		slog.Int("y", y),
		slog.Int("width", w),
		slog.Int("height", h),
	)
}

func (b *Box) getGeometry() (x, y, w, h int) {
	if b == nil || b.ptr == nil {
		return 0, 0, 0, 0
	}

	var (
		cx C.l_int32
		cy C.l_int32
		cw C.l_int32
		ch C.l_int32
	)

	C.boxGetGeometry(b.ptr, &cx, &cy, &cw, &ch)

	return int(cx), int(cy), int(cw), int(ch)
}
