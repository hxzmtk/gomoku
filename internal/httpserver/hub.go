package httpserver

type Hub struct {
	clients    map[string]Conn
	broadcast  chan []byte
	register   chan Conn
	unregister chan Conn
}

func (h *Hub) Run() {
	for {
		select {
		case <-h.register:

		case <-h.unregister:
		case <-h.broadcast:
			for _, _ = range h.clients {

			}

		}
	}
}

func NewHub() *Hub {
	hub := &Hub{clients: make(map[string]Conn, 0)}
	return hub
}
