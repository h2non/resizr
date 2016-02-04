package main

import "gopkg.in/h2non/bimg.v0"

const Version = "0.1.0"

type Versions struct {
	Version     string `json:"resizr"`
	BimgVersion string `json:"bimg"`
	VipsVersion string `json:"libvips"`
}

var CurrentVersions = Versions{Version, bimg.Version, bimg.VipsVersion}
