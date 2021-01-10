package httpserver

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	upgrader   websocket.Upgrader
	httpServer *http.Server
	Addr       string
	hub        *Hub
	engine     *gin.Engine
}


func (server *Server) handleWebsocket(c *gin.Context) {
	conn, err := server.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Info(err)
	}
	_ = conn
}

func (server *Server) init() {
	if runtime.GOOS == "windows" {
		gin.DisableConsoleColor()
	}
	engine := gin.Default()

	// load static
	engine.Static("/static", "web/static")

	// load html template
	engine.LoadHTMLFiles("web/chess.html", "web/ai.html", "web/index.html")

	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"debug": gin.IsDebugging(),
		})
	})
	engine.GET("/ws", server.handleWebsocket)
	server.engine = engine
}

func (server *Server) Start() error {
	server.init()
	server.httpServer = &http.Server{Addr: server.Addr, Handler: server.engine}
	return server.httpServer.ListenAndServe()
}

func NewServer(Addr string) Server {
	return Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		Addr: Addr,
		hub:  NewHub(),
	}
}
