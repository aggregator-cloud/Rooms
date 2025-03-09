package rooms

import (
	"fmt"
	"maps"
	"slices"
	"sync"

	"github.com/google/uuid"
)

type MemberType int

const (
	MemberTypeViewer MemberType = iota
	MemberTypePresenter
)

type Member interface {
	ID() string
	Type() MemberType
	AddStream(stream Stream)
	RemoveStream(stream Stream) error
	GetStreams() map[string]Stream
	CountStreams() int
	AddOnStreamHandler(handler func(stream Stream)) *StreamHandler
	RemoveOnStreamHandler(handler *StreamHandler) (bool, error)
	GetOnStreamHandlers() []StreamHandler
}

type StreamHandler struct {
	id      string
	handler func(stream Stream)
}

func (sh *StreamHandler) ID() string {
	return sh.id
}

type RoomMember struct {
	id                 string
	memberType         MemberType
	Streams            map[string]Stream
	streamsMu          sync.Mutex
	onStreamHandlers   []StreamHandler
	onStreamHandlersMu sync.Mutex
}

func (rm *RoomMember) ID() string {
	return rm.id
}

func (rm *RoomMember) Type() MemberType {
	return rm.memberType
}

func NewRoomMember(id string, memberType MemberType) *RoomMember {
	return &RoomMember{
		id:                 id,
		memberType:         memberType,
		Streams:            make(map[string]Stream),
		streamsMu:          sync.Mutex{},
		onStreamHandlers:   []StreamHandler{},
		onStreamHandlersMu: sync.Mutex{},
	}
}

func (rm *RoomMember) AddOnStreamHandler(handler func(stream Stream)) *StreamHandler {
	rm.onStreamHandlersMu.Lock()
	defer rm.onStreamHandlersMu.Unlock()

	streamHandler := StreamHandler{
		id:      uuid.New().String(),
		handler: handler,
	}

	rm.onStreamHandlers = append(rm.onStreamHandlers, streamHandler)
	return &streamHandler
}

func (rm *RoomMember) RemoveOnStreamHandler(handler *StreamHandler) (bool, error) {
	rm.onStreamHandlersMu.Lock()
	defer rm.onStreamHandlersMu.Unlock()
	sL := len(rm.onStreamHandlers)
	rm.onStreamHandlers = slices.DeleteFunc(rm.onStreamHandlers, func(h StreamHandler) bool {
		return h.id == handler.id
	})
	if len(rm.onStreamHandlers) == sL {
		return false, fmt.Errorf("handler not found")
	}
	return true, nil
}

func (rm *RoomMember) GetOnStreamHandlers() []StreamHandler {
	rm.onStreamHandlersMu.Lock()
	defer rm.onStreamHandlersMu.Unlock()
	handlers := make([]StreamHandler, len(rm.onStreamHandlers))
	copy(handlers, rm.onStreamHandlers)
	return handlers
}

func (rm *RoomMember) AddStream(stream Stream) {
	rm.streamsMu.Lock()
	defer rm.streamsMu.Unlock()
	rm.Streams[stream.ID()] = stream
	rm.onStreamHandlersMu.Lock()
	defer rm.onStreamHandlersMu.Unlock()
	for _, handler := range rm.onStreamHandlers {
		handler.handler(stream)
	}
}

func (rm *RoomMember) RemoveStream(stream Stream) error {
	rm.streamsMu.Lock()
	defer rm.streamsMu.Unlock()

	if _, exists := rm.Streams[stream.ID()]; !exists {
		return fmt.Errorf("stream not found")
	}

	delete(rm.Streams, stream.ID())
	return nil
}

func (rm *RoomMember) GetStreams() map[string]Stream {
	rm.streamsMu.Lock()
	defer rm.streamsMu.Unlock()
	return maps.Clone(rm.Streams)
}

func (rm *RoomMember) CountStreams() int {
	rm.streamsMu.Lock()
	defer rm.streamsMu.Unlock()
	return len(rm.Streams)
}
