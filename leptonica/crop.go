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
)

// edgeClean defines how to clean the edges during cropping.
// Leptonica expects l_int32 for this parameter.
type edgeClean int32

const (
	EDGE_CLEAN_SIDE_NOISE edgeClean = -1
	EDGE_CLEAN_NONE       edgeClean = 0
	EDGE_CLEAN_MAX        edgeClean = 15
)

// CropImage crops the image by clearing borders and adding padding as specified.
//
// It returns the cropped PIX and the bounding BOX of the cropped area.
//
// The caller owns the returned PIX and BOX.
func (p *Pix) CropImage(
	xClear int32, // pixels to clear on the left / right border
	yClear int32, // pixels to clear on the top / bottom border
	edgeClean edgeClean, // how to clean the edges
	xBorder int32, // full-res final pixels "added" padding on left and right
	yBorder int32, // full-res final pixels "added" padding on top and bottom
	maxWiden float32, // max fraction to widen the page to fit a document
) (*Pix, *Box, error) {
	if p == nil || p.ptr == nil {
		return nil, nil, fmt.Errorf("leptonica: nil PIX")
	}

	pix, box, err := pixCropImage(
		p,
		xClear,
		yClear,
		edgeClean,
		xBorder,
		yBorder,
		maxWiden,
	)
	if err != nil {
		return nil, nil, err
	}

	// Unfortunately we're returned a binarised image, so we want to re-crop the
	// original image instead.
	pix.Close()

	cropped, err := p.cropToBox(box)
	if err != nil {
		return nil, nil, err
	}

	return cropped, box, nil
}

// cropToBox crops the Pix to the specified Box.
//
// It returns the cropped Pix.
//
// The caller owns the returned Pix.
func (p *Pix) cropToBox(box *Box) (*Pix, error) {
	if p == nil || p.ptr == nil {
		return nil, fmt.Errorf("leptonica: nil PIX")
	}
	if box == nil || box.ptr == nil {
		return nil, fmt.Errorf("leptonica: nil BOX")
	}

	cropped := C.pixClipRectangle(p.ptr, box.ptr, nil)
	if cropped == nil {
		return nil, fmt.Errorf("leptonica: pixClipRectangle failed")
	}

	res := &Pix{
		ptr: cropped,
	}

	runtime.SetFinalizer(res, (*Pix).Close)

	return res, nil
}

// pixCropImage wraps pixCropImage from Leptonica.
//
// It crops the image by clearing borders and adding padding as specified.
func pixCropImage(
	p *Pix,
	xClear int32, // pixels to clear on the left / right border
	yClear int32, // pixels to clear on the top / bottom border
	edgeClean edgeClean, // how to clean the edges
	xBorder int32, // full-res final pixels "added" padding on left and right
	yBorder int32, // full-res final pixels "added" padding on top and bottom
	maxWiden float32, // max fraction to widen the page to fit a document
) (*Pix, *Box, error) {
	var croppedBounds *C.BOX

	// Leptonica pixCropImage signature includes printwiden and debugfile. We set
	// printwiden to 0 (no debug output) and debugfile to NULL.
	const printWiden C.l_int32 = 0
	var debugfile *C.char

	cropped := C.pixCropImage(
		p.ptr,
		C.l_int32(xClear),
		C.l_int32(yClear),
		C.l_int32(edgeClean),
		C.l_int32(xBorder),
		C.l_int32(yBorder),
		C.l_float32(maxWiden),
		printWiden,
		debugfile,
		&croppedBounds,
	)
	if cropped == nil {
		return nil, nil, fmt.Errorf("leptonica: pixCropImage failed")
	}

	res := &Pix{
		ptr: cropped,
	}

	box := &Box{
		ptr: croppedBounds,
	}

	runtime.SetFinalizer(res, (*Pix).Close)

	return res, box, nil
}
