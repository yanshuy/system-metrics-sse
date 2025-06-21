package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mux := http.NewServeMux()

	// fileServer := http.FileServer(http.Dir("./static/"))
	// mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "./static/index.html")
	})

	mux.HandleFunc("/stream", sseHandler)

	server := http.Server{
		Addr:    ":3333",
		Handler: mux,
	}

	fmt.Println("server started on port", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		logger.Error("error", "err", err)
	}
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	memT := time.NewTicker(time.Second)
	defer memT.Stop()

	cpuT := time.NewTicker(time.Second)
	defer cpuT.Stop()

	clientDisconnected := r.Context().Done()

	rc := http.NewResponseController(w)

	m, err := mem.VirtualMemory()
	if err == nil {
		fmt.Fprintf(w, "event: mem\ndata: Total: %d, Used: %d, Available: %d, Used(%%): %.2f\n\n",
			m.Total, m.Used, m.Available, m.UsedPercent)
		rc.Flush()
	}

	c, err := cpu.Times(false)
	if err == nil {
		fmt.Fprintf(w, "event: cpu\ndata: User: %.2f, System: %.2f, Idle: %.2f\n\n",
			c[0].User, c[0].System, c[0].Idle)
		rc.Flush()
	}

	for {
		select {
		case <-clientDisconnected:
			fmt.Println("client disconnected")
			return
		case <-memT.C:
			m, err := mem.VirtualMemory()
			if err != nil {
				log.Printf("unable to get memory stats; %v", err)
				return
			}

			fmt.Fprintf(w, "event: mem\ndata: Total: %d, Used: %d, Available: %d, Used(%%): %.2f\n\n",
				m.Total, m.Used, m.Available, m.UsedPercent)

			rc.Flush()

		case <-cpuT.C:
			c, err := cpu.Times(false)
			if err != nil {
				log.Printf("unable to get cpu stats; %v", err)
				return
			}

			fmt.Fprintf(w, "event: cpu\ndata: User: %.2f, System: %.2f, Idle: %.2f\n\n",
				c[0].User, c[0].System, c[0].Idle)

			rc.Flush()
		}
	}
}
