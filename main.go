package main

import (
	"context"
	"errors"
	"log/slog"
	"path"
	"strings"

	"github.com/Crystalix007/deskew/leptonica"
	"github.com/spf13/cobra"
)

type flags struct {
	leptonica.DeskewOptions

	OutputFilename string
}

var (
	ErrNoInputFilename        = errors.New("no input filename provided")
	ErrMultipleInputFilenames = errors.New("more than one input filename provided")
)

func main() {
	var flags flags

	cmd := &cobra.Command{
		Use:   "deskew <input-filename>",
		Short: "A tool to deskew images",
		PreRunE: func(_ *cobra.Command, args []string) error {
			if len(args) < 1 {
				return ErrNoInputFilename
			}

			if len(args) > 1 {
				return ErrMultipleInputFilenames
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), flags, args[0])
		},
	}

	cmd.Flags().Float32VarP(
		&flags.SweepRange,
		"sweep-range",
		"r",
		10,
		"The angle range, in degrees, to sweep when detecting skew",
	)

	cmd.Flags().Float32VarP(
		&flags.SweepDelta,
		"sweep-delta",
		"d",
		1,
		"Angular resolution (tenths of a degree) used during the sweep",
	)

	cmd.Flags().Float64VarP(
		&flags.MinConfidence,
		"min-confidence",
		"c",
		2.0,
		"Minimum confidence required to apply the detected skew correction",
	)

	cmd.Flags().Int32VarP(
		&flags.Threshold,
		"threshold",
		"t",
		0,
		"Binarization threshold; 0 for default",
	)

	cmd.Flags().StringVarP(
		&flags.OutputFilename,
		"output-filename",
		"o",
		"deskewed.jpg",
		"The output filename for the deskewed image",
	)

	cmd.Execute()
}

// run takes an image file path, and deskews the image to correct for any skew.
func run(_ context.Context, args flags, inputFilename string) error {
	input, err := leptonica.NewPixFromFile(inputFilename)
	if err != nil {
		return err
	}
	defer input.Close()

	opts := leptonica.DeskewOptions{
		SweepRange:                  args.SweepRange,
		SweepDelta:                  args.SweepDelta,
		MinConfidence:               args.MinConfidence,
		SweepReductionFactor:        4,
		BinarySearchReductionFactor: 1,
	}

	deskewed, angle, confidence, err := input.Deskew(opts)
	if err != nil {
		return err
	}

	defer deskewed.Close()

	outputFormat := filenameToFormat(args.OutputFilename)
	if err := deskewed.WriteToFile(args.OutputFilename, outputFormat); err != nil {
		return err
	}

	slog.Info(
		"deskew complete",
		slog.Float64("angle", angle),
		slog.Float64("confidence", confidence),
		slog.String("output", args.OutputFilename),
	)
	return nil
}

func filenameToFormat(filename string) leptonica.ImageType {
	ext := strings.ToLower(path.Ext(filename))

	switch ext {
	case ".jpg", ".jpeg":
		return leptonica.JPEG
	case ".png":
		return leptonica.PNG
	case ".webp":
		return leptonica.WEBP
	case ".tiff", ".tif":
		return leptonica.TIFF
	case ".bmp":
		return leptonica.BMP
	case ".jp2":
		return leptonica.JP2
	case ".pnm", ".pbm", ".pgm", ".ppm":
		return leptonica.PNM
	default:
		return leptonica.DEFAULT
	}
}
