package handler

import (
	"context"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"

	pb "github.com/el-zacharoo/messaging/internal/gen/messaging/v1"
	pbcnn "github.com/el-zacharoo/messaging/internal/gen/messaging/v1/messagingv1connect"
	"github.com/el-zacharoo/messaging/store"
)

type MessagingServer struct {
	Store store.Storer
	pbcnn.UnimplementedMessagingServiceHandler
}

func (s MessagingServer) Create(ctx context.Context, req *connect.Request[pb.CreateRequest]) (*connect.Response[pb.CreateResponse], error) {
	reqMsg := req.Msg
	msg := reqMsg.MessageThread

	msg.Id = uuid.NewString()
	msg.Messages[0].Date = time.Now().Unix()

	// store functions
	if err := s.Store.Create(ctx, msg); err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}

	// response
	rsp := &pb.CreateResponse{MessageThread: msg}
	return connect.NewResponse(rsp), nil
}


func (s MessagingServer) Update(ctx context.Context, req *connect.Request[pb.UpdateRequest]) (*connect.Response[pb.UpdateResponse], error) {
	reqMsg := req.Msg
	msg := reqMsg.MessageThread

	msg.Messages[len(msg.Messages)-1].Date = time.Now().Unix()

	// store functions
	if err := s.Store.Update(reqMsg.MessageId, ctx, msg); err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}

	// response
	rsp := &pb.UpdateResponse{MessageThread: msg}
	return connect.NewResponse(rsp), nil

}
