package rooms

import (
	"maps"
	"sync"
)

type Room struct {
	ID        string
	Streams    map[string]Stream
	streamsMu  sync.Mutex
	Members   map[string]Member
	membersMu sync.Mutex
}

func NewRoom(id string) *Room {
	return &Room{
		ID:        id,
		Streams:    make(map[string]Stream),
		streamsMu:  sync.Mutex{},
		Members:   make(map[string]Member),
		membersMu: sync.Mutex{},
	}
}

func (r *Room) AddMember(member Member) {
	r.membersMu.Lock()
	defer r.membersMu.Unlock()
	r.Members[member.ID()] = member

	// add streams to viewers
	if viewer, ok := member.(*RoomViewerMember); ok {
		for _, stream := range r.GetStreams() {
			viewer.AddStream(stream)
		}
	}
}

func (r *Room) RemoveMember(member Member) {
	r.membersMu.Lock()
	defer r.membersMu.Unlock()
	delete(r.Members, member.ID())
}

func (r *Room) GetMembers() []Member {
	r.membersMu.Lock()
	defer r.membersMu.Unlock()
	members := make([]Member, 0, len(r.Members))
	for _, member := range r.Members {
		members = append(members, member)
	}
	return members
}

func (r *Room) AddStream(stream Stream) {
	r.streamsMu.Lock()
	defer r.streamsMu.Unlock()
	r.Streams[stream.ID()] = stream
	r.membersMu.Lock()
	defer r.membersMu.Unlock()
	for _, member := range r.Members {
		if viewer, ok := member.(*RoomViewerMember); ok {
			viewer.AddStream(stream)
		}
	}
}

func (r *Room) RemoveStream(stream Stream) {
	r.streamsMu.Lock()
	defer r.streamsMu.Unlock()
	delete(r.Streams, stream.ID())
}

func (r *Room) GetStreams() map[string]Stream {
	r.streamsMu.Lock()
	defer r.streamsMu.Unlock()
	return maps.Clone(r.Streams)
}

func (r *Room) Close() {
	r.streamsMu.Lock()
	defer r.streamsMu.Unlock()
	r.Streams = make(map[string]Stream)

	r.membersMu.Lock()
	defer r.membersMu.Unlock()
	r.Members = make(map[string]Member)
}
