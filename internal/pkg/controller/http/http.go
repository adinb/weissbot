package http

import (
	"fmt"
	"net/http"

	"github.com/adinb/weissbot/internal/pkg/meta"
)

func CreateAndStartHTTPServer(port string, metac chan<- meta.Meta, errc chan<- error) *http.Server {
	metaHandler := createWeissMetaHandler(metac)

	srv := &http.Server{Addr: ":" + port}
	http.HandleFunc("/", handleMainPage)
	http.HandleFunc("/meta", metaHandler)
	http.HandleFunc("/line_webhook", handleLineEvent)

	go func() {
		fmt.Printf("HTTP Server listening at port %s\n", port)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			errc <- err
		}

	}()

	return srv
}
