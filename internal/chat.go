package internal

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/sehogas/gosocket/models"
	"github.com/sehogas/gosocket/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     check,
}

type messageChannel chan *models.Message
type userChannel chan *UserChat
type Channel struct {
	messageChannel messageChannel
	leaveChannel   userChannel
}
type WebSocketChat struct {
	users       map[string]*UserChat
	joinChannel userChannel
	channel     *Channel
}

func NewWebSocketChat() *WebSocketChat {
	return &WebSocketChat{
		users:       make(map[string]*UserChat),
		joinChannel: make(userChannel),
		channel: &Channel{
			messageChannel: make(messageChannel),
			leaveChannel:   make(userChannel),
		},
	}
}
func check(r *http.Request) bool {
	log.Printf("%s %s%s %v", r.Method, r.Host, r.RequestURI, r.Proto)
	return r.Method == http.MethodGet
}

func (w *WebSocketChat) HandlerConnections(rw http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Error de conexión: %s", err.Error())
		return
	}
	keys := r.URL.Query()
	username := strings.TrimSpace(keys.Get("username"))
	log.Println(username)
	if strings.TrimSpace(username) == "" {
		username = fmt.Sprintf("user-%d", utils.GetRandonInt())
	}
	u := NewUserChat(w.channel, username, connection)
	w.joinChannel <- u
	u.OnLine()
}

func (w *WebSocketChat) UsersManager() {
	for {
		select {
		case userChat := <-w.joinChannel:
			w.AddUser(userChat)
		case message := <-w.channel.messageChannel:
			w.SendMessage(message)
		case user := <-w.channel.leaveChannel:
			w.DisconnectUser(user.UserName)
		}
	}
}

func (w *WebSocketChat) AddUser(userchat *UserChat) {
	if user, ok := w.users[userchat.UserName]; ok {
		user.Connection = userchat.Connection
		log.Printf("Usuario reconectado: %s \n", userchat.UserName)
	} else {
		w.users[userchat.UserName] = userchat
		log.Printf("Usuario conectado: %s \n", userchat.UserName)
	}
}
func (w *WebSocketChat) DisconnectUser(username string) {
	if user, ok := w.users[username]; ok {
		defer user.Connection.Close()
		delete(w.users, username)
		log.Printf("Usuario: %s, ha dejado el chat.", username)
	}
}

func (w *WebSocketChat) SendMessage(message *models.Message) {
	if user, ok := w.users[message.Target]; ok {
		if err := user.SendMessage(message); err != nil {
			log.Printf("No se pudo enviar el mensaje a [%s]", message.Target)
		}

	}
}

func StartWebSocket(port string) {
	log.Printf("Chat escuchando en http://localhost:%s", port)
	ws := NewWebSocketChat()
	http.HandleFunc("/chat", ws.HandlerConnections)
	go ws.UsersManager()
	log.Fatalln(http.ListenAndServe(fmt.Sprintf("localhost:%s", port), nil))
}
