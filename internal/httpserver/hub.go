package httpserver

type Hub struct {
	clients    map[string]IConn
	broadcast  chan []byte
	register   chan IConn
	unregister chan IConn
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			client := conn.(Conn)
			h.clients[client.Username] = conn
		case conn := <-h.unregister:
			client := conn.(Conn)
			delete(h.clients, client.Username)
		case <-h.broadcast:
			for _, _ = range h.clients {

			}

		}
	}
}

func NewHub() *Hub {
	hub := &Hub{
		clients:    make(map[string]IConn),
		broadcast:  make(chan []byte, 1024),
		register:   make(chan IConn, 100),
		unregister: make(chan IConn, 100),
	}
	return hub
}
