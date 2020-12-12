package infrastructure

import (
	"fmt"
	"net/http"

	"github.com/LinMAD/InTweets/domain"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// ServerAPI ...
type ServerAPI struct {
	router    *mux.Router
	wsUpgrade websocket.Upgrader
	log       *logrus.Logger
}

// InitServerAPI returns prepared API for use
func InitServerAPI(l *logrus.Logger) *ServerAPI {
	ws := &ServerAPI{
		log:    l,
		router: mux.NewRouter(),
		wsUpgrade: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}

	return ws
}

// loadHandlers api routes
func (a *ServerAPI) LoadHandlers() {
	a.router.Handle("/ws", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		conn, err := a.wsUpgrade.Upgrade(w, r, nil) // Upgrade GET to a websocket
		if err != nil {
			panic(err)
		}

		go a.handleConnection(conn)

	})).Methods(http.MethodGet)
}

// handleConnection for each client under websocket
func (a *ServerAPI) handleConnection(conn *websocket.Conn) {
	var exit = make(chan bool)

	go func() {
		for {
			client := fmt.Sprintf("Client(%s)", conn.RemoteAddr().String())
			msg := domain.WebSocketEvent{}
			err := conn.ReadJSON(&msg)

			if ce, ok := err.(*websocket.CloseError); ok {
				switch ce.Code {
				case websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived:
					a.log.Infof("%s: Closed websocket...", client)
					exit <- true

					return
				}
			}

			a.log.Debugf("%s: Send an event message: %s", client, msg)
			// TODO Get tweets via stream...
		}
	}()
}

// Run server
func (a *ServerAPI) Run(addr string) {
	loggedRouter := handlers.LoggingHandler(a.log.Writer(), a.router)

	a.log.Infof("HTTP Server started on %s", addr)
	if err := http.ListenAndServe(addr, loggedRouter); err != nil {
		panic(err)
	}
}
