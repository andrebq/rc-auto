package main

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"net/http"
	"rcauto/bus"
	"rcauto/ui"
)

func main() {
	addr := "127.0.0.1:8080"
	mainBus := bus.New()
	mainBus.RegisterHandlerFunc(printMessages)
	http.Handle("/assets/", ui.Assets("/assets"))
	versionHttp, versionBus := versionHandlers()
	http.Handle("/version.json", versionHttp)
	mainBus.RegisterHandler(versionBus)
	http.Handle("/dispatch", ui.Dispatch(mainBus))
	http.HandleFunc("/edit-file", fileEditor)
	http.HandleFunc("/", mainPage)
	log.Printf("Starting at %v", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func mainPage(w http.ResponseWriter, req *http.Request) {
	ui.RenderNode(w, req, ui.MainLayout("Controlador", ui.Controls()))
}

func versionHandlers() (http.Handler, bus.Handler) {
	version := struct {
		Version int `json:"version"`
	}{
		Version: rand.Int(),
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(version)
		}), bus.HandlerFunc(func(output bus.Transmitter, from uint64, p bus.Payload) error {
			if p.Meta.Kind != "core.app.get-version" {
				return nil
			}
			return bus.SendData(output, from, "core.app.version", version)
		})
}

func printMessages(t bus.Transmitter, from uint64, p bus.Payload) error {
	log.Printf("%#v", p)
	return nil
}

func fileEditor(w http.ResponseWriter, req *http.Request) {
	// TODO: super insecure, should provide some form of validation here
	filepath := req.FormValue("filepath")
	ui.RenderNode(w, req, ui.TextEditorLayout("Text editor", filepath))
}
