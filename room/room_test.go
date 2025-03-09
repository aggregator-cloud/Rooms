package rooms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoomMembers(t *testing.T) {
	room := NewRoom("1")
	member := NewRoomViewerMember("1")
	room.AddMember(member)

	members := room.GetMembers()
	assert.Equal(t, 1, len(members))
	assert.Equal(t, member, members[0])

	member2 := NewRoomViewerMember("2")
	room.AddMember(member2)

	members = room.GetMembers()
	assert.Equal(t, 2, len(members))
	assert.Equal(t, member2, members[1])
}

func TestAddStream(t *testing.T) {
	room := NewRoom("1")

	videoStream := NewRoomVideoStream("1")
	audioStream := NewRoomAudioStream("2")

	room.AddStream(videoStream)
	room.AddStream(audioStream)

	streams := room.GetStreams()
	assert.Equal(t, 2, len(streams))
	assert.Equal(t, videoStream, streams["1"])
	assert.Equal(t, audioStream, streams["2"])

	member := NewRoomViewerMember("1")
	room.AddMember(member)

	streams = member.GetStreams()
	assert.Equal(t, 2, len(streams))
	assert.Equal(t, videoStream, streams["1"])
	assert.Equal(t, audioStream, streams["2"])
}

func TestRoomStreamHandlers(t *testing.T) {
	t.Run("stream handler lifecycle", func(t *testing.T) {
		room := NewRoom("1")
		member := NewRoomViewerMember("1")
		room.AddMember(member)

		streamEvents := []string{}
		handlerFunc := func(stream Stream) {
			streamEvents = append(streamEvents, "stream_added:"+stream.ID())
		}

		// Test adding handler
		handler := member.AddOnStreamHandler(handlerFunc)
		assert.NotEmpty(t, handler, "handler ID should not be empty")

		// Test handler is called when stream is added
		videoStream := NewRoomVideoStream("video1")
		room.AddStream(videoStream)
		assert.Equal(t, []string{"stream_added:video1"}, streamEvents)

		// Test multiple streams
		audioStream := NewRoomVideoStream("audio1")
		room.AddStream(audioStream)
		assert.Equal(t, []string{
			"stream_added:video1",
			"stream_added:audio1",
		}, streamEvents)

		// Test removing handler stops future notifications
		member.RemoveOnStreamHandler(handler)
		newStream := NewRoomVideoStream("video2")
		room.AddStream(newStream)
		assert.Equal(t, []string{
			"stream_added:video1",
			"stream_added:audio1",
		}, streamEvents, "no new events should be recorded after handler removal")

		// Verify handlers list is empty
		assert.Empty(t, member.GetOnStreamHandlers(), "handlers should be empty after removal")
	})

	t.Run("multiple handlers", func(t *testing.T) {
		room := NewRoom("1")
		member := NewRoomViewerMember("1")
		room.AddMember(member)

		handler1Called := false
		handler2Called := false

		handler1 := member.AddOnStreamHandler(func(stream Stream) {
			handler1Called = true
		})
		handler2 := member.AddOnStreamHandler(func(stream Stream) {
			handler2Called = true
		})

		// Both handlers should be called
		stream := NewRoomVideoStream("1")
		room.AddStream(stream)
		assert.True(t, handler1Called, "handler1 should be called")
		assert.True(t, handler2Called, "handler2 should be called")

		// Remove one handler
		member.RemoveOnStreamHandler(handler1)
		assert.Equal(t, 1, len(member.GetOnStreamHandlers()),
			"should have one handler remaining")

		member.RemoveOnStreamHandler(handler2)
		assert.Equal(t, 0, len(member.GetOnStreamHandlers()),
			"should have no handlers remaining")
	})

	t.Run("handler removal with invalid ID", func(t *testing.T) {
		member := NewRoomViewerMember("1")

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
}
