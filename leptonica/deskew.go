package leptonica

/*
#cgo pkg-config: lept
#include <allheaders.h>
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

var (
	// ErrNilPix is returned when a nil Pix is provided.
	ErrNilPix = errors.New("leptonica: nil PIX")

	// ErrFailedToDetectSkew is returned when skew detection fails.
	ErrFailedToDetectSkew = errors.New("leptonica: failed to detect skew")

	// ErrRotationFailed is returned when image rotation fails.
	ErrRotationFailed = errors.New("leptonica: rotation failed")
)

// DeskewOptions configures pixFindSkewSweepAndSearch and rotation.
type DeskewOptions struct {
	SweepRange                  float32 // degrees to sweep when searching for skew angle
	SweepDelta                  float32 // angular resolution of sweep, in tenths of a degree
	MinConfidence               float64 // minimum confidence required to apply rotation
	SweepReductionFactor        int32   // reduction factor for sweep
	BinarySearchReductionFactor int32   // reduction factor for binary search
	Threshold                   int32   // threshold for binarization; 0 for default
}

// Deskew creates a rotated copy based on detected skew.
//
// It returns the new PIX, the detected angle in degrees, the confidence value,
// and an error when the detection or rotation fails.
//
// The caller owns the returned PIX.
func (p *Pix) Deskew(opts DeskewOptions) (*Pix, float64, float64, error) {
	if p == nil || p.ptr == nil {
		return nil, 0, 0, errors.New("leptonica: nil PIX")
	}

	deskewed, angle, conf, err := pixDeskewGeneral(
		p,
		opts.SweepReductionFactor,
		opts.SweepRange,
		opts.SweepDelta,
		opts.BinarySearchReductionFactor,
		opts.Threshold,
	)
	if err != nil {
		return nil, 0, 0, err
	}

	runtime.SetFinalizer(deskewed, (*Pix).Close)

	if opts.MinConfidence > 0 && float64(conf) < opts.MinConfidence {
		deskewed.Close()

		return nil, 0, float64(conf), fmt.Errorf(
			"leptonica: skew confidence %.2f below threshold %.2f",
			float64(conf), opts.MinConfidence,
		)
	}

	return deskewed, float64(angle), float64(conf), nil
}

// WriteToFile encodes the PIX to disk. If format is UNKNOWN or DEFAULT, Leptonica
// will choose a suitable encoder based on file extension.
func (p *Pix) WriteToFile(filename string, format ImageType) error {
	if p == nil || p.ptr == nil {
		return ErrNilPix
	}

	cfn := C.CString(filename)
	defer C.free(unsafe.Pointer(cfn))

	if format == UNKNOWN {
		format = DEFAULT
	}

	rc := C.pixWrite(cfn, p.ptr, C.l_int32(format))
	if rc != 0 {
		return fmt.Errorf("leptonica: failed to write image %q", filename)
	}

	return nil
}

// pixDeskewGeneral wraps the Leptonica pixDeskewGeneral function.
//
// It finds the skew angle of the provided Pix and returns the deskewed image,
// the angle, and the confidence.
//
// Note: If the confidence is below opts.MinConfidence, a cloned copy of the
// original Pix is returned instead.
func pixDeskewGeneral(
	p *Pix,
	sweepReductionFactor int32,
	sweepRange float32,
	sweepDelta float32,
	binarySearchReductionFactor int32,
	threshold int32,
) (*Pix, float64, float64, error) {
	var (
		angle C.l_float32
		conf  C.l_float32
	)

	deskewedPix := C.pixDeskewGeneral(
		p.ptr,
		C.l_int32(sweepReductionFactor),
		C.l_float32(sweepRange),
		C.l_float32(sweepDelta),
		C.l_int32(binarySearchReductionFactor),
		C.l_int32(threshold),
		&angle,
		&conf,
	)
	if deskewedPix == nil {
		return nil, 0, 0, ErrFailedToDetectSkew
	}

	res := &Pix{
		ptr: deskewedPix,
	}

	return res, float64(angle), float64(conf), nil
}
