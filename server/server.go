package server

import (
	"go-live-server/injector"
	"go-live-server/reload"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func NewRouter(root string, reloader *reload.Reloader) http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(
		http.Dir(root),
	)

	mux.HandleFunc(
		"/__live",
		reloader.HandleSSE,
	)

	mux.HandleFunc("/__live_reload.js", injector.ServeScript)

	mux.Handle("/", ServeFile(root, fileServer))

	return mux
}

func ServeFile(
	root string,
	fileServer http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func (w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, ".html") || r.URL.Path == "/" {
				path := filepath.Join(
					root,
					r.URL.Path,
				)

				if r.URL.Path == "/"{
					path = filepath.Join(root, "index.html")
				}

				content, err := os.ReadFile(path)

				if err == nil {
					html := injector.Inject(
						string(content),
					)

					w.Header().Set(
						"Content-Type",
						"text/html",
					)

					w.Write([]byte(html))

					return
				}
			}
			fileServer.ServeHTTP(w, r)
		},
	)
}
