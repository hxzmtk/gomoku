package conn

type IMsg interface {
	Receive()
	ToBytes() []byte
}

type IObserver interface {
	Do(subject ISubject, msg IMsg) error
}

type ISubject interface {
	Attach(observers ...IObserver)
	Detach(observer IObserver)
	Notify(msg IMsg) error
}

type concreteSubject struct {
	observers []IObserver
}

func (s *concreteSubject) Attach(observers ...IObserver) {
	s.observers = append(s.observers, observers...)
}

func (s *concreteSubject) Detach(observer IObserver) {
	for k, item := range s.observers {
		if item == observer {
			s.observers = append(s.observers[:k], s.observers[k+1:]...)
		}
	}
}
func (s *concreteSubject) Notify(msg IMsg) error {
	for _, item := range s.observers {
		if err := item.Do(s, msg); err != nil {
			return err
		}
	}
	return nil
}

func NewSubject() ISubject {
	return &concreteSubject{observers: []IObserver{}}
}
