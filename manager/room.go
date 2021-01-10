package manager

type RoomManager struct {
}

func (RoomManager) Init() error {
	return nil
}

func NewRoomManager() *RoomManager {
	return &RoomManager{}
}
