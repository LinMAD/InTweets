package infrastructure

import (
	"fmt"
	"net/http"

	"github.com/LinMAD/InTweets/core"
	"github.com/LinMAD/InTweets/domain"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// ServerAPI ...
type ServerAPI struct {
	router    *mux.Router
	wsUpgrade websocket.Upgrader
	log       *core.Logger
	twitCred  *domain.TwitterCredential
}

// InitServerAPI returns prepared API for use
func InitServerAPI(c *domain.TwitterCredential, l *core.Logger) *ServerAPI {
	ws := &ServerAPI{
		log:      l,
		twitCred: c,
		router:   mux.NewRouter(),
		wsUpgrade: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}

	l.Info("HTTP server initialized...")

	return ws
}

// LoadRouteHandlers api routes
func (sApi *ServerAPI) LoadRouteHandlers() {
	sApi.router.Handle("/ws", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		conn, err := sApi.wsUpgrade.Upgrade(w, r, nil) // Upgrade GET to sApi websocket
		if err != nil {
			panic(err)
		}

		go sApi.handleConnection(conn)

	})).Methods(http.MethodGet)
}

// handleConnection for each client under websocket
func (sApi *ServerAPI) handleConnection(conn *websocket.Conn) {
	go func() {
		var exit = make(chan bool)
		client := fmt.Sprintf("Client(%s)", conn.RemoteAddr().String())

		sApi.log.Infof("%s: Connected to websocket", client)
		twitterClient := InitClientTwitter(client, sApi.twitCred, sApi.log)

		for {
			msg := domain.WebSocketEvent{}
			err := conn.ReadJSON(&msg)

			if ce, ok := err.(*websocket.CloseError); ok {
				switch ce.Code {
				// check for client closed connections
				case websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived:
					sApi.log.Infof("%s: Closed websocket...", client)
					twitterClient.StopStream()
					exit <- true
					return
				}
			}

			sApi.log.Debugf("%s: Send an event message: %s", client, msg)

			if ok := twitterClient.StartStream(msg.Data); !ok {
				sApi.log.Errorf("Twitter client: Unexpected error while creating stream for tweets...")
				return
			}

			sApi.log.Debugf("Twitter client: Tweets feed started...")

			go func(exit chan bool) {
				tweetCh := make(chan string)
				go twitterClient.FetchTweet(tweetCh, exit)

				for t := range tweetCh {
					if err := conn.WriteJSON(t); err != nil {
						sApi.log.Error(err)
						return
					}
				}
			}(exit)
		}
	}()
}

// Run server
func (sApi *ServerAPI) Run(addr string) {
	loggedRouter := handlers.LoggingHandler(sApi.log.Writer(), sApi.router)

	sApi.log.Infof("HTTP Server started on http://%s", addr)
	if err := http.ListenAndServe(addr, loggedRouter); err != nil {
		panic(err)
	}
}
