package handlers

import (
	"backend/internal/actors"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/gin-gonic/gin"
	"github.com/lxzan/gws"
)

const (
	PingInterval = 10 * time.Second
	PingWait     = 10 * time.Second
)

var (
	Upgrader *gws.Upgrader
	usersMap map[*gws.Conn]*actor.PID
	mapMutex sync.RWMutex
)

type WSHandler struct {
	actorService actors.Service
}

func NewWSHandlers(actorService actors.Service) *WSHandler {
	handler := &WSHandler{
		actorService: actorService,
	}

	Upgrader = gws.NewUpgrader(handler, &gws.ServerOption{
		ParallelEnabled:   true,
		Recovery:          gws.Recovery,
		PermessageDeflate: gws.PermessageDeflate{Enabled: true},
	})
	usersMap = make(map[*gws.Conn]*actor.PID)

	return handler
}

func (ws *WSHandler) OnOpen(socket *gws.Conn) {
	// _ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
}

func (ws *WSHandler) OnClose(socket *gws.Conn, err error) {
	mapMutex.Lock()
	if userPID, exists := usersMap[socket]; exists {
		delete(usersMap, socket)
		ws.actorService.KillActor(userPID)
	}
	mapMutex.Unlock()

	slog.Info("user disconnected")
}

func (ws *WSHandler) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
	_ = socket.WriteString("heartbeat")
}

func (ws *WSHandler) OnPong(socket *gws.Conn, payload []byte) {}

func (ws *WSHandler) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()
	if b := message.Data.Bytes(); len(b) == 9 && string(b) == "heartbeat" {
		ws.OnPing(socket, nil)
		return
	}
}

func (ws *WSHandler) Setup(c *gin.Context) {
	userID := c.Param("user_id")

	socket, err := Upgrader.Upgrade(c.Writer, c.Request)
	if err != nil {
		slog.Error("failed upgrading connection", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed upgrading connection"})
		return
	}

	userPID := ws.actorService.CreateUser(userID, socket)

	mapMutex.Lock()
	usersMap[socket] = userPID
	mapMutex.Unlock()

	go func() {
		socket.ReadLoop()
	}()
}
