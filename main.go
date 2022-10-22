package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	// "log"
	// "net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
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

		token := strings.Split(strings.Split(c.URL().RawQuery, "&")[0], "=")[1]
		if token == "" {
			return fmt.Errorf("cant get the token")
		}

		var keyfunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
			return []byte("D191359BDAA07A792E10C1916746897FC16729F88D50E5A55D7E11D6C7AE8BE4"), nil
		}

		parsed, err := jwt.Parse(token, keyfunc)
		if err != nil {
			return fmt.Errorf("failed to parse JWT.\nError: %s", err.Error())
		}

		if !parsed.Valid {
			return fmt.Errorf("token is not valid")
		}

		claims, _ := parsed.Claims.(jwt.MapClaims)
		userId := claims["id"].(string)

		isUser := controllers.GetUser(userId)

		if !isUser {
			return fmt.Errorf("user not found")
		}

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
	log.Println("Serving on localhost:4001")
	log.Fatal(http.ListenAndServe("127.0.0.1:4001", nil))

}
