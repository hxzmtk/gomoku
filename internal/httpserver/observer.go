package httpserver


type IObserver interface {
	Do(subject ISubject, msg IMessage) error
}

type ISubject interface {
	Attach(observers ...IObserver)
	Detach(observer IObserver)
	Notify(msg IMessage) error
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
func (s *concreteSubject) Notify(msg IMessage) error {
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
