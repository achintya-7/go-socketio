package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/achintya-7/go_socketio/controllers"
	"github.com/achintya-7/go_socketio/models"
	socketio "github.com/googollee/go-socket.io"
)

func main() {

	server := socketio.NewServer(nil)

	// on connect
	server.OnConnect("/", func(c socketio.Conn) error {
		c.SetContext("")
		fmt.Println("connected", c.ID())
		return nil

	})

	server.OnEvent("/", "join", func(c socketio.Conn, room string) {
		c.Context()
		c.Join(room)
		c.Emit("join", room)
	})

	server.OnEvent("/", "send", func(c socketio.Conn, r string) {
		// TODO : generate a random number or uuid and add it to the req struct
		// can use the key and value payer as {messageId: "dkbkfnbfbn"}

		var req models.SendMessageReq
		err := json.Unmarshal([]byte(r), &req)
		if err != nil {
			fmt.Println("Unable to parse :", err)
		}

		res := models.SendMessageRes{
			UserId: req.UserId,
			RoomId: req.RoomId,
			Content: req.Content,
			ContentType: req.ContentType,
			MessageId: "323424234",
		}

		server.BroadcastToRoom("/", req.RoomId.String(), "send", res)
		controllers.SendMessage(req)
	})

	server.OnError("/", func(c socketio.Conn, e error) {
		fmt.Println("met error:", e)
	})

	server.OnDisconnect("/", func(c socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	log.Println("Serving on localhost:4000")
	log.Fatal(http.ListenAndServe("127.0.0.1:4000", nil))

}
