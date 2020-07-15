package hub

type IObserver interface {
	Do(subject ISubject, msg Msg) error
}

type ISubject interface {
	Attach(observers ...IObserver)
	Detach(observer IObserver)
	Notify(msg Msg) error
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
func (s *concreteSubject) Notify(msg Msg) error {
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

type ObserverChessWalk struct {
	client IClient
}

func (o *ObserverChessWalk) Do(subject ISubject, msg Msg) error {
	switch v := o.client.(type) {
	case *HumanClient:

		// 如果发生错误，移除自己，下次将收不到推送的消息
		defer func() {
			if err := recover(); err != nil {
				subject.Detach(o)
			}
		}()

		v.Send <- &msg
	}

	return nil
}
