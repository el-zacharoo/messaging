package main

import (
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/el-zacharoo/messaging/handler"
	pbcnn "github.com/el-zacharoo/messaging/internal/gen/messaging/v1/messagingv1connect"
	"github.com/el-zacharoo/messaging/store"
)

const port = "localhost:8080"

func main() {
	svc := &handler.MessagingServer{
		Store: store.Connect(),
	}
	mux := http.NewServeMux()

	path, h := pbcnn.NewMessagingServiceHandler(svc)
	mux.Handle(path, h)

	http.ListenAndServe(
		port,
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
