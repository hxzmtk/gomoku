package manager

import "flag"

var httpPort int = 8000

func init() {
	flag.IntVar(&httpPort, "port", 8000, "example: -port 8000")
	flag.Parse()
}

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
		return err
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

var manager = &Manager{modules: make([]IModule, 0)}

func Get() *Manager {
	return manager
}

func GetRoomManager() *RoomManager {
	return manager.RoomManager
}

func GetUserManager() *UserManager {
	return manager.UserManager
}
