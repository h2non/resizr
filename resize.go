package main

import "gopkg.in/h2non/bimg.v0"

type Options struct {
	Width, Height int
	Force         bool
	Operation     string
}

func Resize(image []byte, opts Options) ([]byte, error) {
	params := bimg.Options{
		Enlarge:      true,
		NoAutoRotate: true,
		Width:        opts.Width,
		Height:       opts.Height,
		Force:        opts.Force,
		Crop:         opts.Operation == "crop" || opts.Operation == "resize",
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
	return "image/jpeg"
}
