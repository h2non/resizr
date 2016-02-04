package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/h2non/bimg.v0"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ServerOptions struct {
	Port             int
	Burst            int
	Concurrency      int
	HttpReadTimeout  int
	HttpWriteTimeout int
	CORS             bool
	Gzip             bool
	Address          string
	ApiKey           string
	CertFile         string
	KeyFile          string
	Placeholder      []byte
}

func Server(o ServerOptions) error {
	addr := o.Address + ":" + strconv.Itoa(o.Port)
	handler := NewServerMux(o)

	server := &http.Server{
		Addr:           addr,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    time.Duration(o.HttpReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(o.HttpWriteTimeout) * time.Second,
	}

	return listenAndServe(server, o)
}

func listenAndServe(s *http.Server, o ServerOptions) error {
	if o.CertFile != "" && o.KeyFile != "" {
		return s.ListenAndServeTLS(o.CertFile, o.KeyFile)
	}
	return s.ListenAndServe()
}

func NewServerMux(o ServerOptions) http.Handler {
	mux := httprouter.New()
	mux.GET("/", indexController)
	mux.GET("/:operation/:size/*url", resizeController(o))
	return mux
}

func resizeController(o ServerOptions) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if r.Method != "GET" {
			badRequest(w, "method not allowed")
			return
		}

		width, height, err := parseDimensions(ps.ByName("size"))
		if err != nil {
			badRequest(w, "invalid width or height path expression")
			return
		}

		debug("resize to %dx%d", width, height)
		opts := Options{Width: width, Height: height, Operation: ps.ByName("operation")}

		image, err := Fetch(ps.ByName("url")[1:])
		if err != nil {
			failed(w, opts, o, err.Error())
			return
		}

		image, err = Resize(image, opts)
		if err != nil {
			failed(w, opts, o, err.Error())
			return
		}

		mime := GetImageMimeType(bimg.DetermineImageType(image))
		w.Header().Set("Content-Type", mime)
		w.Write(image)
	}
}

func indexController(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	body, _ := json.Marshal(CurrentVersions)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func parseDimensions(value string) (int, int, error) {
	size := strings.Split(value, "x")
	width, err := strconv.Atoi(size[0])
	if err != nil {
		return 0, 0, err
	}
	height, err := strconv.Atoi(size[1])
	return width, height, err
}

func failed(w http.ResponseWriter, opts Options, o ServerOptions, msg string) {
	image := placeholder
	if len(o.Placeholder) > 1 {
		image = o.Placeholder
	}

	opts.Force = true
	image, err := Resize(image, opts)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", GetImageMimeType(bimg.DetermineImageType(image)))
	w.Header().Set("Error", msg)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(image)
}

func badRequest(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Error", msg)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(placeholder)
}
