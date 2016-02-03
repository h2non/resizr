# resizr

A dead simple HTTP service developed in 1 hour for image resizing/cropping. 
Designed for easy integration from web apps and programmatic HTTP clients.

Behind the scenes it's powered by [libvips](https://github.com/jcupitt/libvips) and [bimg](https://github.com/h2non/bimg).

## Rationale

I just created this to cover a very particular need. 
If you need more versatility, I would recommend you start building something custom by your own.

## Features

- Fast: written in Go and uses libvips, a powerful image processing library in C.
- Simple (for now): just type the URL
- Supports image resize with crop calculus.
- Supports JPEG, PNG and WEBP formats and conversion between them.
- Automatic image rotation based on EXIF orientation metadata.
- Image fetching and resizing.
- Default image placeholder in case of error.
- No cache. Build your cache layer in front of it.

## Upcoming features

- API token based authorization
- gzip responses
- CORS support
- Traffic throttle strategy

## Installation

```bash
go get github.com/h2non/resizr
```

## Usage

```bash

```

## HTTP API

### GET /
Content-Type: `application/json`

Returns the server version. 

### GET /resize/{width}x{height}/{imageUrl}
Content-Type: `image/*`

Performs an image resize with implicit crop calculus to automatically fit to the desired resolution.

## License

MIT