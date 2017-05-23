package main

import (
	"github.com/geobe/gostip/go/controller"
	"github.com/geobe/gostip/go/model"
	"net/http"
	"log"
	//"net"
	"net/url"
	"strings"
	"regexp"
)

const httpport = ":8070"
const tlsport = ":8090"
const schema = "http"

func main() {
	// prepare database
	model.Setup("")
	db := model.Db()
	defer db.Close()
	model.ClearTestDb(db)
	model.InitTestDb(db)

	mux := controller.SetRouting()

	// konfiguriere server
	server := &http.Server{
		Addr:    "0.0.0.0" + tlsport,
		Handler: mux,
	}

	// switching redirect handler
	handlerSwitch := &HandlerSwitch{
		Mux:    mux,
		Redirect: http.HandlerFunc(RedirectHTTP),
	}

	// konfiguriere redirect server
	redirectserver := &http.Server{
		Addr:    "0.0.0.0" + httpport,
		Handler: handlerSwitch, //http.HandlerFunc(RedirectHTTP),
	}
	// starte den redirect server
	go redirectserver.ListenAndServe()

	// und starte den primären server
	log.Printf("server starting\n")
	server.ListenAndServe()
}


// RedirectHTTP is an HTTP handler (suitable for use with http.HandleFunc)
// that responds to all requests by redirecting to the same URL served over HTTPS.
// It should only be invoked for requests received over HTTP.
func RedirectHTTP(w http.ResponseWriter, r *http.Request) {
	if r.TLS != nil || r.Host == "" {
		http.Error(w, "not found", 404)
	}

	var u *url.URL
	u = r.URL
	host := r.Host
	u.Host = strings.Split(host, ":")[0] + tlsport
	u.Scheme = schema
	log.Printf("redirect to u.host  %s -> %s\n", r.Host, u.String())
	http.Redirect(w, r, u.String(), 302)
}

type HandlerSwitch struct {
	Mux      http.Handler
	Redirect http.Handler
}

// nicht richtig für 172.16.0.0/12
var matcher = regexp.MustCompile("(192\\.168.*)|(localhost)|(10\\..*)|(172\\..*)")

func (h *HandlerSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	local := matcher.MatchString(host)
	if local {
		h.Mux.ServeHTTP(w, r)
	} else {
		h.Redirect.ServeHTTP(w, r)
	}
}
