package rooms

/*
	RoomStream is a stream that is shared between members.
	Streams may be webrtc streams or other types of streams.
*/

type Stream interface {
	ID() string
}

type BaseRoomStream struct {
	Stream
	id     string
	closed bool
}

func NewBaseRoomStream(id string) *BaseRoomStream {
	return &BaseRoomStream{
		id:     id,
		closed: false,
	}
}

func (rt *BaseRoomStream) ID() string {
	return rt.id
}

func (rt *BaseRoomStream) Close() {
	rt.closed = true
}
