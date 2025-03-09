package rooms

import (
	"errors"
	"maps"
	"slices"
	"sync"
)

type RoomStream struct {
	Stream
	RecipientTypes []MemberType
}

func NewRoomStream(stream Stream, recipientTypes ...MemberType) *RoomStream {
	return &RoomStream{
		Stream:         stream,
		RecipientTypes: recipientTypes,
	}
}

type Room struct {
	id        string
	streams   map[string]RoomStream
	streamsMu sync.Mutex
	members   map[string]Member
	membersMu sync.Mutex
}

func NewRoom(id string) *Room {
	return &Room{
		id:        id,
		streams:   make(map[string]RoomStream),
		streamsMu: sync.Mutex{},
		members:   make(map[string]Member),
		membersMu: sync.Mutex{},
	}
}

func (r *Room) ID() string {
	return r.id
}

func (r *Room) AddMember(member Member) {
	r.membersMu.Lock()
	defer r.membersMu.Unlock()
	r.members[member.ID()] = member

	// add streams to viewers
	for _, stream := range r.GetStreams() {
		if !slices.Contains(stream.RecipientTypes, member.Type()) && len(stream.RecipientTypes) > 0 {
			continue
		}
		member.AddStream(stream.Stream)
	}
}

func (r *Room) RemoveMember(member Member) {
	r.membersMu.Lock()
	defer r.membersMu.Unlock()
	delete(r.members, member.ID())
}

func (r *Room) GetMembers() []Member {
	r.membersMu.Lock()
	defer r.membersMu.Unlock()
	members := make([]Member, 0, len(r.members))
	for _, member := range r.members {
		members = append(members, member)
	}
	return members
}

func (r *Room) CountMembers() int {
	r.membersMu.Lock()
	defer r.membersMu.Unlock()
	return len(r.members)
}

func (r *Room) AddStream(stream Stream, memberTypes ...MemberType) {
	r.streamsMu.Lock()
	defer r.streamsMu.Unlock()
	roomStream := NewRoomStream(stream, memberTypes...)
	r.streams[stream.ID()] = *roomStream
	r.membersMu.Lock()
	defer r.membersMu.Unlock()
	lenMemberTypes := len(memberTypes)
	for _, member := range r.members {
		mType := member.Type()
		if lenMemberTypes > 0 && !slices.Contains(memberTypes, mType) {
			continue
		}
		member.AddStream(stream)
	}
}

func (r *Room) RemoveStream(stream Stream) error {
	r.membersMu.Lock()
	r.streamsMu.Lock()
	defer r.streamsMu.Unlock()
	defer r.membersMu.Unlock()

	if _, exists := r.streams[stream.ID()]; !exists {
		return errors.New("stream not found")
	}

	var errs []error
	for _, member := range r.members {
		if err := member.RemoveStream(stream); err != nil {
			errs = append(errs, err)
		}
	}

	delete(r.streams, stream.ID())

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (r *Room) GetStreams() map[string]RoomStream {
	r.streamsMu.Lock()
	defer r.streamsMu.Unlock()
	return maps.Clone(r.streams)
}

func (r *Room) Close() {
	r.streamsMu.Lock()
	defer r.streamsMu.Unlock()
	r.streams = make(map[string]RoomStream)

	r.membersMu.Lock()
	defer r.membersMu.Unlock()
	r.members = make(map[string]Member)
}
