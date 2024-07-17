package models

type Hub struct {
	Rooms map[string]*ChatRoom
}

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*ChatRoom),
	}
}

func (h *Hub) AddNewChatRoom(roomName string) *ChatRoom {
	room := NewChatRoom(roomName)
	h.Rooms[roomName] = room
	return room
}
