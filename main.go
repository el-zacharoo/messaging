package main

import (
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/el-zacharoo/messaging/handler"
	pbcnn "github.com/el-zacharoo/messaging/internal/gen/messaging/v1/messagingv1connect"
	"github.com/el-zacharoo/messaging/store"
	"github.com/rs/cors"
)

const port = "localhost:8080"

func main() {
	s := store.Connect()

	svc := &handler.MessagingServer{
		Store: s,
	}

	c := setCORS()

	mux := http.NewServeMux()
	path, h := pbcnn.NewMessagingServiceHandler(svc)
	mux.Handle(path, h)
	handler := c.Handler(mux)

	http.ListenAndServe(
		port,
		h2c.NewHandler(handler, &http2.Server{}),
	)
}

func setCORS() *cors.Cors {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Access-Control-Allow-Origin", "Content-Type"},
		AllowedMethods:   []string{"POST"},
		AllowCredentials: true,
	})
	return c
}
