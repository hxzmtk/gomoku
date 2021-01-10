package human

import "github.com/zqhhh/gomoku/service/hub"

type ObserverChessWalk struct {
	client *Client
}

func (o *ObserverChessWalk) Do(subject hub.ISubject, msg hub.IMsg) error {
	// 如果发生错误，移除自己，下次将收不到推送的消息
	defer func() {
		if err := recover(); err != nil {
			subject.Detach(o)
		}
	}()

	o.client.Send <- msg

	return nil
}
