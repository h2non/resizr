# resizr [![Go Report Card](http://goreportcard.com/badge/h2non/resizr)](http://goreportcard.com/report/h2non/resizr) [![Heroku](https://img.shields.io/badge/Heroku-Deploy_Now-blue.svg)](https://heroku.com/deploy)

A dead simple HTTP service developed in 1 hour for image resizing/cropping. 
Designed for easy integration from web apps and programmatic HTTP clients.

Behind the scenes it's powered by [libvips](https://github.com/jcupitt/libvips) and [bimg](https://github.com/h2non/bimg).

`resizr` is currently used in production processing thoundands of images per day.

## Rationale

I just created this to cover a very particular need. 
If you need more versatility, I would recommend you start building something custom by your own.

## Features

- Fast: written in Go and uses libvips, a powerful image processing library in C.
- Simple (for now): just type the URL
- Supports image resize with crop calculus.
- Supports JPEG, PNG and WEBP formats and conversion between them.
- Automatic image rotation based on EXIF orientation metadata.
- Default image placeholder in case of processing error.
- Image fetching and resizing.
- No cache. Build your cache layer in front of it.

## Upcoming features

- gzip responses
- CORS support
- Traffic throttle strategy

## Installation

```bash
go get github.com/h2non/resizr
```

## Usage

```bash
resizr 0.1.0

Usage:
  resizr -p 80
  resizr -placeholder image.jpg

Options:
  -a <addr>                 bind address [default: *]
  -p <port>                 bind port [default: 9000]
  -h, -help                 output help
  -v, -version              output version
  -placeholder <path>       placeholder image to use on error
  -cors                     Enable CORS support [default: false]
  -gzip                     Enable gzip compression [default: false]
  -key <key>                Define API key for authorization
  -http-read-timeout <num>  HTTP read timeout in seconds [default: 30]
  -http-write-timeout <num> HTTP write timeout in seconds [default: 30]
  -certfile <path>          TLS certificate file path
  -keyfile <path>           TLS private key file path
  -concurreny <num>         Throttle concurrency limit per second [default: disabled]
  -burst <num>              Throttle burst max cache size [default: 100]
  -mrelease <num>           OS memory release inverval in seconds [default: 30]
  -cpus <num>               Number of used cpu cores.
                            (default for current machine is 8 cores)
```

Start the server:
```bash
resizr -p 8080
```

Then, from a web browser, try opening the following URL:
```bash
http://localhost:8080/crop/200x200/http://imgsv.imaging.nikon.com/lineup/lens/zoom/normalzoom/af-s_dx_18-300mmf_35-56g_ed_vr/img/sample/sample4_l.jpg
```

Using it from HTML `img` tag is as simple as:
```html
<img src="http://localhost:8080/crop/200x200/http://imgsv.imaging.nikon.com/lineup/lens/zoom/normalzoom/af-s_dx_18-300mmf_35-56g_ed_vr/img/sample/sample4_l.jpg" />
```

## HTTP API

### Handling errors

Since `resizr` has been designed to be used as public HTTP service, including web pages, the response MIME type must be respected in most scenarios,
so the server will always reply with a placeholder image in case of error. 

You can customize the placeholder image passing the `-placeholder` flag when starting `resizr`.

If image resizing fails for some reason, a 400 Bad Request will be used as response status, but the `Content-Type` will always `image/*`.
If you want to see the error details, you have it in the `Error` header field.

### GET /
Content-Type: `application/json`

Returns versions info. 

### GET /crop/{width}x{height?}/{imageUrl}
Content-Type: `image/*`

Performs an image resize with implicit crop calculus to automatically fit to the desired resolution.

`height` value is optional.

### GET /resize/{width}x{height?}/{imageUrl}
Content-Type: `image/*`

Performs an image resize with implicit crop calculus and enlarge, if necessary, to the desired resolution.

`height` value is optional.

## License

MIT