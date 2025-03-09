package rooms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoomMember(t *testing.T) {
	t.Run("add stream", func(t *testing.T) {
		t.Parallel()
		member := NewRoomMember("1", MemberTypeViewer)
		member.AddStream(NewBaseRoomStream("1"))
		assert.Equal(t, 1, member.CountStreams())
	})
	t.Run("remove stream", func(t *testing.T) {
		t.Parallel()
		member := NewRoomMember("1", MemberTypeViewer)
		member.AddStream(NewBaseRoomStream("1"))
		member.RemoveStream(NewBaseRoomStream("1"))
		assert.Equal(t, 0, member.CountStreams())
	})
	t.Run("remove non-existent stream", func(t *testing.T) {
		t.Parallel()
		member := NewRoomMember("1", MemberTypeViewer)
		err := member.RemoveStream(NewBaseRoomStream("1"))
		assert.Error(t, err)
	})
	t.Run("add on stream handler", func(t *testing.T) {
		t.Parallel()
		member := NewRoomMember("1", MemberTypeViewer)
		handler := member.AddOnStreamHandler(func(stream Stream) {})
		assert.NotNil(t, handler)
	})
	t.Run("remove on stream handler", func(t *testing.T) {
		t.Parallel()
		member := NewRoomMember("1", MemberTypeViewer)
		handler := member.AddOnStreamHandler(func(stream Stream) {})
		member.RemoveOnStreamHandler(handler)
		assert.Equal(t, 0, len(member.GetOnStreamHandlers()))
	})

	t.Run("call on stream handler", func(t *testing.T) {
		t.Parallel()
		member := NewRoomMember("1", MemberTypeViewer)
		called := 0
		member.AddOnStreamHandler(func(stream Stream) {
			called++
		})
		member.AddStream(NewBaseRoomStream("1"))
		member.AddStream(NewBaseRoomStream("2"))
		assert.Equal(t, 2, called)
	})

	t.Run("call on stream handler with multiple handlers", func(t *testing.T) {
		t.Parallel()
		member := NewRoomMember("1", MemberTypeViewer)
		called1 := 0
		called2 := 0
		handler1 := member.AddOnStreamHandler(func(stream Stream) {
			called1++
		})
		member.AddOnStreamHandler(func(stream Stream) {
			called2++
		})
		member.AddStream(NewBaseRoomStream("1"))
		assert.Equal(t, 1, called1)
		assert.Equal(t, 1, called2)
		member.RemoveOnStreamHandler(handler1)
		member.AddStream(NewBaseRoomStream("2"))
		assert.Equal(t, 1, called1)
		assert.Equal(t, 2, called2)
	})

	t.Run("remove on stream handler with invalid handler", func(t *testing.T) {
		t.Parallel()
		member := NewRoomMember("1", MemberTypeViewer)
		_, err := member.RemoveOnStreamHandler(&StreamHandler{id: "invalid-id"})
		assert.Error(t, err)
	})

	t.Run("get on stream handlers", func(t *testing.T) {
		t.Parallel()
		member := NewRoomMember("1", MemberTypeViewer)
		handler := member.AddOnStreamHandler(func(stream Stream) {})
		handlers := member.GetOnStreamHandlers()
		assert.Equal(t, 1, len(handlers))
		assert.Equal(t, handler.ID(), handlers[0].ID())
		// GetOnStreamHandlers should return a copy of the handlers
		assert.NotEqual(t, member.GetOnStreamHandlers(), member.onStreamHandlers)
	})

	t.Run("handlers are kept in order", func(t *testing.T) {
		t.Parallel()
		member := NewRoomMember("1", MemberTypeViewer)
		handler1 := member.AddOnStreamHandler(func(stream Stream) {})
		handler2 := member.AddOnStreamHandler(func(stream Stream) {})
		handlers := member.GetOnStreamHandlers()
		assert.Equal(t, 2, len(handlers))
		assert.Equal(t, handler1.ID(), handlers[0].ID())
		assert.Equal(t, handler2.ID(), handlers[1].ID())
	})

}
