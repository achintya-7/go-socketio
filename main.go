package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/achintya-7/go_socketio/controllers"
	"github.com/achintya-7/go_socketio/models"
	socketio "github.com/googollee/go-socket.io"
	gonanoid "github.com/matoous/go-nanoid/v2"
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
		var req models.SendMessageReq
		err := json.Unmarshal([]byte(r), &req)
		if err != nil {
			fmt.Println("Unable to parse :", err)
			server.BroadcastToRoom("/", req.RoomId.Hex(), "send", err)
			return
		}

		uid, err := gonanoid.New()
		if err != nil {
			fmt.Printf("unable to generate nanoid: %v\n", err)
			server.BroadcastToRoom("/", req.RoomId.Hex(), "send", err)
			return
		}

		res := models.SendMessageRes{
			UserId:      req.UserId,
			RoomId:      req.RoomId,
			Content:     req.Content,
			ContentType: req.ContentType,
			MessageId:   uid,
		}

		server.BroadcastToRoom("/", req.RoomId.Hex(), "send", res)
		
		controllers.SendMessage(res)
	})

	server.OnEvent("/", "reply", func(c socketio.Conn, r string) {
		var req models.ReplyMessageReq

		err := json.Unmarshal([]byte(r), &req)
		if err != nil {
			fmt.Printf("unable to generate nanoid: %v\n", err)
			server.BroadcastToRoom("/", req.RoomId.Hex(), "send", err)
			return
		}

		uid, err := gonanoid.New()
		if err != nil {
			fmt.Printf("unable to generate nanoid: %v\n", err)
			server.BroadcastToRoom("/", req.RoomId.Hex(), "send", err)
			return
		}

		res := models.ReplyMessageRes{
			UserId:      req.UserId,
			RoomId:      req.RoomId,
			Content:     req.Content,
			ContentType: req.ContentType,
			PrevMessage: req.PrevMessage,
			MessageId:   uid,
		}

		server.BroadcastToRoom("/", req.RoomId.Hex(), "send", res)

		controllers.ReplyMessage(res)


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
