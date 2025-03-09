package rooms

/*
	RoomStream is a stream that is shared between members.
	Streams may be webrtc streams or other types of streams.
*/

type Stream interface {
	ID() string
}

type RoomStream struct {
	Stream
	id string
}

func (rt *RoomStream) ID() string {
	return rt.id
}

type RoomVideoStream struct {
	RoomStream
}

func NewRoomVideoStream(id string) *RoomVideoStream {
	return &RoomVideoStream{
		RoomStream: RoomStream{
			id: id,
		},
	}
}

type RoomAudioStream struct {
	RoomStream
}

func NewRoomAudioStream(id string) *RoomAudioStream {
	return &RoomAudioStream{
		RoomStream: RoomStream{
			id: id,
		},
	}
}
