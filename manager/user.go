package manager

type UserManager struct {
}

func (UserManager) Init() error {
	return nil
}

func NewUserManager() *UserManager {
	return &UserManager{}
}
