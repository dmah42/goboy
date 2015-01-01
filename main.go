package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dominichamon/goboy/goboy"
)

const (
	// rate = time.Millisecond * 16
	rate = time.Second * 5

	rootHTML = `<!DOCTYPE html>
	<html>
		<head>
			<title>goboy</title>
		</head>
		<body>
			<h1>goboy</h1>
			<canvas id="screen" width="160" height="144">
			</canvas>
			<span id="run">Run</span>
			<!-- TODO: reset -->
			<!-- TODO: load input -->
			<script type="text/javascript">
				var run = function() {
					var req = new XMLHttpRequest();
					req.open('GET', '/run', false);
					req.send();

					document.getElementById('run').innerHTML = 'Pause';
					document.getElementById('run').onclick = pause;
				}

				var pause = function() {
					var req = new XMLHttpRequest();
					req.open('GET', '/pause', false);
					req.send();

					document.getElementById('run').innerHTML = 'Run';
					document.getElementById('run').onclick = run;
				}

				var frame = function() {
					var req = new XMLHttpRequest();
					req.onload = render
					req.open('GET', '/frame', false);
					req.send();

					window.requestAnimationFrame(frame);
				}

				var render = function() {
					var screen = document.getElementById('screen');
					var ctx = screen.getContext('2d');
					var b = ctx.getImageData(0, 0, screen.width, screen.height);

					var screenData = JSON.parse(this.responseText);

					// TODO: scale up canvas
					for (var i = 0; i < screenData.length; ++i) {
						b.data[i] = screenData[i];
					}

					ctx.putImageData(b, 0, 0);
				}

				pause();

				window.onkeydown = function(e) {
					var req = new XMLHttpRequest();
					req.open('GET', '/keydown?keycode=' + e.keyCode, true);
					req.send();
				}

				window.onkeyup = function(e) {
					var req = new XMLHttpRequest();
					req.open('GET', '/keyup?keycode=' + e.keyCode, true);
					req.send();
				}

				window.requestAnimationFrame(frame);
			</script>
		</body>
	</html>`
)

var (
	port = flag.Int("port", 8888, "Port on which to listen")

	rootTemplate = template.Must(template.New("root").Parse(rootHTML))
)

func main() {
	if len(os.Args) < 2 {
		log.Panic("no ROM file selected")
	}
	go goboy.Loop(os.Args[1])

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
	b, err := json.Marshal(goboy.GPU.Screen)
	if err != nil {
		log.Println("goboy: screen marshal error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
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
