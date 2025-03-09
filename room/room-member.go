package rooms

import (
	"maps"
	"slices"
	"sync"

	"github.com/google/uuid"
)

type Member interface {
	ID() string
}

type RoomMember struct {
	id string
}

func (rm *RoomMember) ID() string {
	return rm.id
}

type StreamHandler struct {
	id      string
	handler func(stream Stream)
}

type RoomViewerMember struct {
	RoomMember
	Streams            map[string]Stream
	streamsMu          sync.Mutex
	onStreamHandlers   []StreamHandler
	onStreamHandlersMu sync.Mutex
}

func NewRoomViewerMember(id string) *RoomViewerMember {
	return &RoomViewerMember{
		RoomMember: RoomMember{
			id: id,
		},
		Streams:            make(map[string]Stream),
		streamsMu:          sync.Mutex{},
		onStreamHandlers:   []StreamHandler{},
		onStreamHandlersMu: sync.Mutex{},
	}
}

func (rm *RoomViewerMember) AddOnStreamHandler(handler func(stream Stream)) *StreamHandler {
	rm.onStreamHandlersMu.Lock()
	defer rm.onStreamHandlersMu.Unlock()

	streamHandler := StreamHandler{
		id:      uuid.New().String(),
		handler: handler,
	}

	rm.onStreamHandlers = append(rm.onStreamHandlers, streamHandler)
	return &streamHandler
}

func (rm *RoomViewerMember) RemoveOnStreamHandler(handler *StreamHandler) {
	rm.onStreamHandlersMu.Lock()
	defer rm.onStreamHandlersMu.Unlock()

	rm.onStreamHandlers = slices.DeleteFunc(rm.onStreamHandlers, func(h StreamHandler) bool {
		return h.id == handler.id
	})
}

func (rm *RoomViewerMember) GetOnStreamHandlers() []StreamHandler {
	rm.onStreamHandlersMu.Lock()
	defer rm.onStreamHandlersMu.Unlock()
	handlers := make([]StreamHandler, len(rm.onStreamHandlers))
	copy(handlers, rm.onStreamHandlers)
	return handlers
}

func (rm *RoomViewerMember) AddStream(stream Stream) {
	rm.streamsMu.Lock()
	defer rm.streamsMu.Unlock()
	rm.Streams[stream.ID()] = stream
	rm.onStreamHandlersMu.Lock()
	defer rm.onStreamHandlersMu.Unlock()
	for _, handler := range rm.onStreamHandlers {
		handler.handler(stream)
	}
}

func (rm *RoomViewerMember) RemoveStream(stream Stream) {
	rm.streamsMu.Lock()
	defer rm.streamsMu.Unlock()

	delete(rm.Streams, stream.ID())
}

func (rm *RoomViewerMember) GetStreams() map[string]Stream {
	rm.streamsMu.Lock()
	defer rm.streamsMu.Unlock()
	return maps.Clone(rm.Streams)
}

type RoomSourceMember struct {
	RoomMember
}

func NewRoomSourceMember(id string) *RoomSourceMember {
	return &RoomSourceMember{
		RoomMember: RoomMember{
			id: id,
		},
	}
}
