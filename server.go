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
	HttpCacheTtl     int
	HttpReadTimeout  int
	HttpWriteTimeout int
	CORS             bool
	Gzip             bool
	Address          string
	ApiKey           string
	CertFile         string
	KeyFile          string
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
	mux.GET("/:operation/:size/*url", resizeController)
	return mux
}

func resizeController(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != "GET" {
		badRequest(w, "method not allowed")
		return
	}

	image, err := Fetch(ps.ByName("url")[1:])
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	size := strings.Split(ps.ByName("size"), "x")
	width, err := strconv.Atoi(size[0])
	height, err2 := strconv.Atoi(size[1])
	if err != nil || err2 != nil {
		badRequest(w, "invalid width or height path expression")
		return
	}

	opts := Options{Width: width, Height: height}
	image, err = Resize(image, opts)

	mime := GetImageMimeType(bimg.DetermineImageType(image))
	w.Header().Set("Content-Type", mime)
	w.Write(image)
}

func indexController(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Version string `json:"version"`
	}{Version: Version}
	body, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func badRequest(w http.ResponseWriter, msg string) {
	msg = strings.Replace(msg, "\"", "\\\"", -1)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("{\"message\": \"" + msg + "\"}"))
}
