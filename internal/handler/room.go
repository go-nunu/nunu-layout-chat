package handler

import (
	"context"
	"fmt"
	"github.com/topfreegames/pitaya/v2"
	"github.com/topfreegames/pitaya/v2/timer"
	"go.uber.org/zap"
	v1 "nunu-layout-chat/api/v1"
	"nunu-layout-chat/internal/service"
	"strconv"
	"time"
)

type RoomHandler struct {
	*Handler
	timer       *timer.Timer
	roomService service.RoomService
}

func NewRoomHandler(
	handler *Handler,
	roomService service.RoomService,
) *RoomHandler {
	return &RoomHandler{
		Handler:     handler,
		roomService: roomService,
	}
}

// AfterInit component lifetime callback
func (h *RoomHandler) AfterInit() {
	h.logger.Debug("AfterInit")

	// TODO: You shouldn't create a room here, this line of code is just for the convenience of demonstration
	h.Create(context.Background(), nil)

	h.timer = pitaya.NewTimer(time.Second*3, func() {
		count, err := h.app.GroupCountMembers(context.Background(), "test-room")
		if err != nil {
			h.logger.Error("AfterInit error", zap.Error(err))
			return
		}
		h.logger.Debug("AfterInit", zap.Any("userCount", count))
	})
}

// Create room
func (h *RoomHandler) Create(ctx context.Context, msg []byte) (*v1.Response, error) {
	err := h.app.GroupCreate(context.Background(), "test-room")
	if err != nil {
		h.logger.WithContext(ctx)
	}
	return &v1.Response{
		Code:    0,
		Message: "success",
		Data:    nil,
	}, nil
}

// Join room
func (h *RoomHandler) Join(ctx context.Context, msg []byte) (*v1.Response, error) {
	s := h.app.GetSessionFromCtx(ctx)
	fakeUID := s.ID()                              // just use s.ID as uid !!!
	err := s.Bind(ctx, strconv.Itoa(int(fakeUID))) // binding session uid

	if err != nil {
		return nil, pitaya.Error(err, "RH-000", map[string]string{"failed": "bind"})
	}

	// notify others
	h.app.GroupBroadcast(ctx, "chat", "test-room", "onNewUser", &v1.NewUser{Content: fmt.Sprintf("New user: %s", s.UID())})
	// new user join group
	h.app.GroupAddMember(ctx, "test-room", s.UID()) // add session to group
	uids, err := h.app.GroupMembers(ctx, "test-room")
	if err != nil {
		return nil, err
	}
	s.Push("onMembers", &v1.AllMembers{Members: uids})

	// on session close, remove it from group
	s.OnClose(func() {
		h.app.GroupRemoveMember(ctx, "test-room", s.UID())
	})

	return &v1.Response{
		Code:    0,
		Message: "success",
		Data:    "",
	}, nil
}

// Message sync last message to all members
func (h *RoomHandler) Message(ctx context.Context, msg *v1.UserMessage) {
	err := h.app.GroupBroadcast(ctx, "chat", "test-room", "onMessage", msg)
	if err != nil {
		fmt.Println("error broadcasting message", err)
	}
}
