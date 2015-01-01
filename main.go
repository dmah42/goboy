package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
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

				pause();
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
