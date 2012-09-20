package api

import (
	"encoding/json"
	"fmt"
	"github.com/bmizerany/pat"
	"github.com/fd/heroku-keepalive/model"
	"log"
	"net/http"
)

var router = pat.New()

func ListenAndServe(addr string) {
	if addr == "" || addr == ":" {
		addr = ":3000"
	}
	http.ListenAndServe(addr, router)
}

func init() {
	router.Get("/assets/", http.FileServer(http.Dir("public")))
	router.Get("/apps", http.HandlerFunc(get_apps))
	router.Get("/", http.FileServer(http.Dir("public")))
}

func get_apps(w http.ResponseWriter, req *http.Request) {
	apps := model.Get()
	arry := make([]*model.Application, 0, len(apps))
	for _, app := range apps {
		arry = append(arry, app)
	}

	data, err := json.Marshal(arry)
	if err != nil {
		log.Printf("[E]: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.WriteHeader(200)
	w.Write(data)
}
