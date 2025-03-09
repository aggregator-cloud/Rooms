package rooms

import (
	"fmt"
	"maps"
	"sync"
)

type RoomManager struct {
	rooms   map[string]*Room
	roomsMu sync.Mutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Room),
	}
}

/* CreateRoom creates a new room if it doesn't exist and returns the room */
func (rm *RoomManager) CreateRoom(id string) *Room {
	rm.roomsMu.Lock()
	defer rm.roomsMu.Unlock()
	if _, ok := rm.rooms[id]; ok {
		return rm.rooms[id]
	}
	room := NewRoom(id)
	rm.rooms[id] = room
	return room
}

/* DestroyRoom destroys a room if it exists */
func (rm *RoomManager) DestroyRoom(id string) (bool, error) {
	rm.roomsMu.Lock()
	defer rm.roomsMu.Unlock()
	room := rm.rooms[id]
	if room == nil {
		return false, fmt.Errorf("room %s not found", id)
	}
	room.Close()
	delete(rm.rooms, id)
	return true, nil
}

func (rm *RoomManager) GetRoom(id string) (*Room, error) {
	rm.roomsMu.Lock()
	defer rm.roomsMu.Unlock()
	room := rm.rooms[id]
	if room == nil {
		return nil, fmt.Errorf("room %s not found", id)
	}
	return room, nil
}

func (rm *RoomManager) GetRooms() map[string]*Room {
	rm.roomsMu.Lock()
	defer rm.roomsMu.Unlock()
	return maps.Clone(rm.rooms)
}

func (rm *RoomManager) CountRooms() int {
	rm.roomsMu.Lock()
	defer rm.roomsMu.Unlock()
	return len(rm.rooms)
}
