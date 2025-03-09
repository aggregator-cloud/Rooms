package rooms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoomManager(t *testing.T) {
	t.Run("create room", func(t *testing.T) {
		t.Parallel()
		manager := NewRoomManager()
		rm := manager.CreateRoom("1")
		room, err := manager.GetRoom("1")
		assert.Nil(t, err)
		assert.Equal(t, rm, room)
		assert.NotNil(t, room)
		rm2 := manager.CreateRoom("1")
		assert.Equal(t, rm, rm2)
	})

	t.Run("destroy room", func(t *testing.T) {
		t.Parallel()
		manager := NewRoomManager()
		manager.CreateRoom("1")
		b, err := manager.DestroyRoom("1")
		assert.Nil(t, err)
		assert.True(t, b)
		room, err := manager.GetRoom("1")
		assert.Error(t, err)
		assert.Nil(t, room)
		room, err = manager.GetRoom("1")
		assert.Error(t, err)
		assert.Nil(t, room)
	})

	t.Run("destroy non-existent room", func(t *testing.T) {
		manager := NewRoomManager()
		b, err := manager.DestroyRoom("1")
		assert.NotNil(t, err)
		assert.False(t, b)
		room, err := manager.GetRoom("1")
		assert.Error(t, err)
		assert.Nil(t, room)
	})

	t.Run("count rooms", func(t *testing.T) {
		t.Parallel()
		manager := NewRoomManager()
		manager.CreateRoom("1")
		manager.CreateRoom("2")
		assert.Equal(t, 2, manager.CountRooms())
	})

	t.Run("get rooms", func(t *testing.T) {
		t.Parallel()
		manager := NewRoomManager()
		manager.CreateRoom("1")
		manager.CreateRoom("2")
		rooms := manager.GetRooms()
		assert.Equal(t, 2, len(rooms))
		assert.Contains(t, rooms, "1")
		assert.Contains(t, rooms, "2")
	})

	t.Run("get non-existent room", func(t *testing.T) {
		t.Parallel()
		manager := NewRoomManager()
		room, err := manager.GetRoom("1")
		assert.Error(t, err)
		assert.Nil(t, room)
	})

}
