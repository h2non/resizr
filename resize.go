package main

import "gopkg.in/h2non/bimg.v0"

type Options struct {
	Width, Height int
	Operation     string
}

func Resize(image []byte, opts Options) ([]byte, error) {
	params := bimg.Options{
		Width:  opts.Width,
		Height: opts.Height,
		Crop:   opts.Operation == "crop",
	}
	return bimg.Resize(image, params)
}

func GetImageMimeType(code bimg.ImageType) string {
	if code == bimg.PNG {
		return "image/png"
	}
	if code == bimg.WEBP {
		return "image/webp"
	}
	if code == bimg.TIFF {
		return "image/tiff"
	}
	return "image/jpeg"
}
