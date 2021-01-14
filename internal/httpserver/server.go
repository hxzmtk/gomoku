package httpserver

import (
	"net/http"
	"reflect"
	"runtime"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type HandleFunc func(conn IConn, msg interface{}) (IMessage, error)
type Server struct {
	upgrader   websocket.Upgrader
	httpServer *http.Server
	Addr       string
	hub        *Hub
	engine     *gin.Engine
	handlers   map[int]HandleFunc
}

func (server *Server) handleWebsocket(c *gin.Context) {
	ws, err := server.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Info(err)
	}
	newConn := NewConn(ws, server.hub)
	newConn.Start()
	newConn.Init()
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
	go func() {
		server.httpServer.ListenAndServe()
	}()
	return nil
}

var (
	srv     *Server
	onceSrv sync.Once
)

func NewServer(Addr string) *Server {
	onceSrv.Do(func() {
		srv = &Server{
			upgrader: websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
				CheckOrigin:     func(r *http.Request) bool { return true },
			},
			Addr:     Addr,
			hub:      NewHub(),
			handlers: make(map[int]HandleFunc),
		}
	})
	return srv
}

func Register(msgId int, handle HandleFunc) {
	if _, ok := srv.handlers[msgId]; ok {
		log.Errorf("handle %d is existed", msgId)
		return
	}
	log.Infof("register msgId:%d, name:%s success", msgId, runtime.FuncForPC(reflect.ValueOf(handle).Pointer()).Name())
	srv.handlers[msgId] = handle
}

func DoHandle(conn *Conn, p IMessage) (IMessage, error) {
	handle, ok := srv.handlers[p.GetMsgId()]
	if !ok {
		log.Errorf("handle not existed,msgId:%d", p.GetMsgId())
		return nil, nil
	}
	msg, err := handle(conn, p)
	return msg, err

}
