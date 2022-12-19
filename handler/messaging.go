package handler

import (
	"context"
	"errors"
	"regexp"
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
	if err := s.Store.CreateMsg(ctx, msg); err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}

	// response
	rsp := &pb.CreateResponse{MessageThread: msg}
	return connect.NewResponse(rsp), nil
}

func (s MessagingServer) Query(ctx context.Context, req *connect.Request[pb.QueryRequest]) (*connect.Response[pb.QueryResponse], error) {
	reqMsg := req.Msg

	if reqMsg.SearchText != "" {
		pattern, err := regexp.Compile(`^[a-zA-Z@. ]+$`)
		if err != nil {
			return nil, connect.NewError(connect.CodeAborted, err)
		}
		if !pattern.MatchString(reqMsg.SearchText) {
			return nil, connect.NewError(connect.CodeInvalidArgument,
				errors.New("invalid search text format"))
		}
	}

	cur, mat, err := s.Store.QueryMsg(ctx, reqMsg)
	if err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}

	rsp := &pb.QueryResponse{
		Cursor:  cur,
		Matches: mat,
	}
	return connect.NewResponse(rsp), nil
}

func (s MessagingServer) Get(ctx context.Context, req *connect.Request[pb.GetRequest]) (*connect.Response[pb.GetResponse], error) {
	reqMsg := req.Msg.MessageId

	// store functions
	msg, err := s.Store.GetMsg(ctx, reqMsg)
	if err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}

	// response
	rsp := &pb.GetResponse{MessageThread: msg}
	return connect.NewResponse(rsp), nil
}

func (s MessagingServer) Update(ctx context.Context, req *connect.Request[pb.UpdateRequest]) (*connect.Response[pb.UpdateResponse], error) {
	reqMsg := req.Msg
	msg := reqMsg.MessageThread

	msg.Messages[len(msg.Messages)-1].Date = time.Now().Unix()

	// store functions
	if err := s.Store.UpdateMsg(reqMsg.MessageId, ctx, msg); err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}

	// response
	rsp := &pb.UpdateResponse{MessageThread: msg}
	return connect.NewResponse(rsp), nil
}
