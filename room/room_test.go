package rooms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoom(t *testing.T) {
	t.Run("add member", func(t *testing.T) {
		t.Parallel()
		room := NewRoom("1")
		member := NewRoomMember("1", MemberTypeViewer)
		room.AddMember(member)

		assert.Equal(t, 1, room.CountMembers())

		member2 := NewRoomMember("2", MemberTypePresenter)
		room.AddMember(member2)

		assert.Equal(t, 2, room.CountMembers())
	})

	t.Run("remove member", func(t *testing.T) {
		t.Parallel()
		room := NewRoom("1")
		member := NewRoomMember("1", MemberTypeViewer)
		room.AddMember(member)
		assert.Equal(t, 1, room.CountMembers())
		room.RemoveMember(member)
		assert.Equal(t, 0, room.CountMembers())
	})

	t.Run("add stream", func(t *testing.T) {
		t.Parallel()
		room := NewRoom("1")

		videoStream := NewBaseRoomStream("1")
		audioStream := NewBaseRoomStream("2")

		room.AddStream(videoStream)
		room.AddStream(audioStream)

		streams := room.GetStreams()
		assert.Equal(t, 2, len(streams))
		assert.Equal(t, videoStream, streams["1"].Stream)
		assert.Equal(t, audioStream, streams["2"].Stream)

		member := NewRoomMember("1", MemberTypeViewer)
		room.AddMember(member)

		mstreams := member.GetStreams()
		assert.Equal(t, 2, len(mstreams))
		assert.Equal(t, videoStream, mstreams["1"])
		assert.Equal(t, audioStream, mstreams["2"])
	})

	t.Run("add streams to late comers", func(t *testing.T) {
		t.Parallel()
		room := NewRoom("room1")
		room.AddStream(NewBaseRoomStream("video1"))
		room.AddStream(NewBaseRoomStream("audio1"), MemberTypePresenter)

		latecomer := NewRoomMember("member1", MemberTypeViewer)
		room.AddMember(latecomer)

		streams := latecomer.GetStreams()
		assert.Equal(t, 1, len(streams))
		assert.Equal(t, "video1", streams["video1"].ID())
	})

	t.Run("stream handler lifecycle", func(t *testing.T) {
		t.Parallel()
		room := NewRoom("1")
		member := NewRoomMember("1", MemberTypeViewer)
		room.AddMember(member)

		streamEvents := []string{}
		handlerFunc := func(stream Stream) {
			streamEvents = append(streamEvents, "stream_added:"+stream.ID())
		}

		// Test adding handler
		handler := member.AddOnStreamHandler(handlerFunc)
		assert.NotEmpty(t, handler, "handler ID should not be empty")

		// Test handler is called when stream is added
		videoStream := NewBaseRoomStream("video1")
		room.AddStream(videoStream)
		assert.Equal(t, []string{"stream_added:video1"}, streamEvents)

		// Test multiple streams
		audioStream := NewBaseRoomStream("audio1")
		room.AddStream(audioStream)
		assert.Equal(t, []string{
			"stream_added:video1",
			"stream_added:audio1",
		}, streamEvents)

		// Test removing handler stops future notifications
		member.RemoveOnStreamHandler(handler)
		newStream := NewBaseRoomStream("video2")
		room.AddStream(newStream)
		assert.Equal(t, []string{
			"stream_added:video1",
			"stream_added:audio1",
		}, streamEvents, "no new events should be recorded after handler removal")

		// Verify handlers list is empty
		assert.Empty(t, member.GetOnStreamHandlers(), "handlers should be empty after removal")
	})

	t.Run("multiple handlers with different member types", func(t *testing.T) {
		t.Parallel()
		room := NewRoom("1")
		member := NewRoomMember("1", MemberTypeViewer)
		member2 := NewRoomMember("2", MemberTypePresenter)
		room.AddMember(member)
		room.AddMember(member2)

		handler1Called := 0
		handler2Called := 0
		handler3Called := 0

		member.AddOnStreamHandler(func(stream Stream) {
			handler1Called++
		})
		member.AddOnStreamHandler(func(stream Stream) {
			handler2Called++
		})
		member2.AddOnStreamHandler(func(stream Stream) {
			handler3Called++
		})

		// Both handlers should be called
		stream := NewBaseRoomStream("1")
		room.AddStream(stream)
		assert.Equal(t, 1, handler1Called, "handler1 should be called")
		assert.Equal(t, 1, handler2Called, "handler2 should be called")
		assert.Equal(t, 1, handler3Called, "handler3 should be called")
		room.AddStream(stream, MemberTypePresenter)
		assert.Equal(t, 1, handler1Called, "handler1 should be called")
		assert.Equal(t, 1, handler2Called, "handler2 should be called")
		assert.Equal(t, 2, handler3Called, "handler3 should be called")
		room.AddStream(stream, MemberTypeViewer)
		assert.Equal(t, 2, handler1Called, "handler1 should be called")
		assert.Equal(t, 2, handler2Called, "handler2 should be called")
		assert.Equal(t, 2, handler3Called, "handler3 should be called")
	})

	t.Run("remove stream", func(t *testing.T) {
		t.Parallel()
		room := NewRoom("1")
		member := NewRoomMember("1", MemberTypeViewer)
		room.AddMember(member)
		stream := NewBaseRoomStream("1")
		room.AddStream(stream)
		err := room.RemoveStream(stream)
		assert.NoError(t, err)
		assert.Equal(t, 0, member.CountStreams())
	})

	t.Run("handler removal with invalid ID", func(t *testing.T) {
		t.Parallel()
		member := NewRoomMember("1", MemberTypeViewer)

		handler := member.AddOnStreamHandler(func(stream Stream) {})
		assert.Equal(t, 1, len(member.GetOnStreamHandlers()))

		// Remove with invalid ID should not affect valid handlers
		member.RemoveOnStreamHandler(&StreamHandler{id: "invalid-id"})
		assert.Equal(t, 1, len(member.GetOnStreamHandlers()),
			"valid handler should remain after removing invalid ID")

		// Remove with valid ID should work
		member.RemoveOnStreamHandler(handler)
		assert.Equal(t, 0, len(member.GetOnStreamHandlers()),
			"handler should be removed with valid ID")
	})

	t.Run("get members", func(t *testing.T) {
		t.Parallel()
		room := NewRoom("1")
		member := NewRoomMember("1", MemberTypeViewer)
		room.AddMember(member)
		members := room.GetMembers()
		assert.Contains(t, members, member)

		member2 := NewRoomMember("2", MemberTypePresenter)
		room.AddMember(member2)
		members = room.GetMembers()
		assert.Contains(t, members, member2)
	})
}
