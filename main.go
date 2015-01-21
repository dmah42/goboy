package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dominichamon/goboy/goboy"
)

const (
	// rate = time.Millisecond * 16
	rate = time.Second * 5
)

var (
	port = flag.Int("port", 8888, "Port on which to listen")
	rom = flag.String("rom", "roms/ttt.gb", "The ROM to load and run")
	run = flag.Bool("run", false, "Run the emulator automatically")
)

func main() {
	flag.Parse()

	if len(*rom) == 0 {
		log.Panic("no ROM file selected")
	}

	goboy.Run = *run

	go goboy.Loop(*rom)

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/run", runHandler)
	http.HandleFunc("/pause", pauseHandler)
	http.HandleFunc("/frame", frameHandler)
	http.HandleFunc("/keydown", keydownHandler)
	http.HandleFunc("/keyup", keyupHandler)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Panicf("failed to start listening on port %d: %v", *port, err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	rootTemplate, err := template.ParseFiles("index.html")
	if err != nil {
		log.Println("goboy: template parsing error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rootTemplate.Execute(w, nil)
}

func runHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Run")
	goboy.Run = true
	w.WriteHeader(http.StatusOK)
}

func pauseHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Pause")
	goboy.Run = false
	w.WriteHeader(http.StatusOK)
}

func frameHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: get registers as well as gpu data
	data := map[string]interface{} {
		"screen": goboy.GPU.Screen,
		"tilemap": goboy.GPU.Tilemap,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("goboy: screen marshal error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func keydownHandler(w http.ResponseWriter, r *http.Request) {
	keycode, err := strconv.ParseInt(r.URL.Query().Get("keycode"), 10, 8)
	if err != nil {
		log.Printf("goboy: keycode parse error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Println("Keydown: ", keycode)
	goboy.Key.Keydown(byte(keycode))
}

func keyupHandler(w http.ResponseWriter, r *http.Request) {
	keycode, err := strconv.ParseInt(r.URL.Query().Get("keycode"), 10, 8)
	if err != nil {
		log.Printf("goboy: keycode parse error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Println("Keyup: ", keycode)
	goboy.Key.Keyup(byte(keycode))
}
