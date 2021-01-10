package manager

type IModule interface {
	Init() error
}
type Manager struct {
	ClientManager *ClientManager
	UserManager   *UserManager
	RoomManager   *RoomManager
	modules       []IModule
}

func (m *Manager) init() error {
	return nil
}
func (m *Manager) appendModule(module IModule) IModule {
	m.modules = append(m.modules, module)
	return module
}

func (m *Manager) Init() error {
	if err := m.init(); err != nil {
		return nil
	}
	m.ClientManager = m.appendModule(NewClientManager()).(*ClientManager)
	m.UserManager = m.appendModule(NewUserManager()).(*UserManager)
	m.RoomManager = m.appendModule(NewRoomManager()).(*RoomManager)
	for _, m := range m.modules {
		if err := m.Init(); err != nil {
			return nil
		}
	}
	return nil
}

func (m *Manager) Stop() {
}

var m = &Manager{modules: make([]IModule, 0)}

func Get() *Manager {
	return m
}

func GetRoomManager() *RoomManager {
	return m.RoomManager
}

func GetUserManager() *UserManager {
	return m.UserManager
}
