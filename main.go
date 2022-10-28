package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/achintya-7/go_socketio/controllers"
	"github.com/achintya-7/go_socketio/models"
	socketio "github.com/googollee/go-socket.io"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func main() {

	server := socketio.NewServer(nil)
	// SECRET_KEY := utils.GetDotEnvVariable("SECRET_TOKEN")

	server.OnConnect("/", func(c socketio.Conn) error {
		return nil
	})

	// on connect with jwt sign in
	// server.OnConnect("/", func(c socketio.Conn) error {
	// 	c.SetContext("")

	// 	token := strings.Split(strings.Split(c.URL().RawQuery, "&")[0], "=")[1]
	// 	if token == "" {
	// 		return fmt.Errorf("cant get the token")
	// 	}

	// 	var keyfunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
	// 		return []byte(SECRET_KEY), nil
	// 	}

	// 	parsed, err := jwt.Parse(token, keyfunc)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to parse JWT.\nError: %s", err.Error())
	// 	}

	// 	if !parsed.Valid {
	// 		return fmt.Errorf("token is not valid")
	// 	}

	// 	claims, _ := parsed.Claims.(jwt.MapClaims)
	// 	userId := claims["id"].(string)

	// 	isUser := controllers.GetUser(userId)

	// 	if !isUser {
	// 		return fmt.Errorf("user not found")
	// 	}

	// 	return nil

	// })

	// joining of room
	server.OnEvent("/", "join", func(c socketio.Conn, room string) {
		c.Context()
		c.Join(room)
		c.Emit("join", room)
	})

	// sending message and forwarding message
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

		timestamp := time.Now().UTC().Unix()

		res := models.SendMessageRes{
			UserId:      req.UserId,
			RoomId:      req.RoomId,
			Content:     req.Content,
			ContentType: req.ContentType,
			MessageId:   uid,
			Timestamp:   timestamp,
		}

		server.BroadcastToRoom("/", req.RoomId.Hex(), "send", res)

		controllers.SendMessage(res)
	})

	// replying a previous message
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

		timestamp := time.Now().UTC().Unix()

		res := models.ReplyMessageRes{
			UserId:        req.UserId,
			RoomId:        req.RoomId,
			Content:       req.Content,
			ContentType:   req.ContentType,
			PrevMessage:   req.PrevMessage,
			PrevMessageId: req.PrevMessageId,
			MessageId:     uid,
			Timestamp:     timestamp,
		}

		server.BroadcastToRoom("/", req.RoomId.Hex(), "send", res)

		controllers.ReplyMessage(res)

	})

	server.OnEvent("/", "delete", func(c socketio.Conn, r string) {
		var req models.DeleteMessageReq

		if err := json.Unmarshal([]byte(r), &req); err != nil {
			fmt.Println("Unable to unmarshal the json", err)
			return
		}

		timestamp := time.Now().UTC().Unix() - 1800

		if timestamp <= req.TimeStamp {
			res := models.DeleteMessageRes{
				RoomId:        req.RoomId,
				MessageId:     req.MessageId,
				TimeStamp:     timestamp,
				DeleteMessage: true,
			}

			server.BroadcastToRoom("/", req.RoomId.Hex(), "modify", res)

			controllers.DeleteMessage(req)
		} else {
			res := models.DeleteMessageRes{
				RoomId:        req.RoomId,
				MessageId:     req.MessageId,
				TimeStamp:     timestamp,
				DeleteMessage: false,
			}

			server.BroadcastToRoom("/", req.RoomId.Hex(), "modify", res)
		}
	})

	server.OnEvent("/", "update", func(c socketio.Conn, r string) {
		var req models.UpdateMessageReq

		if err := json.Unmarshal([]byte(r), &req); err != nil {
			fmt.Println("Unable to unmarshal the json", err)
			return
		}

		timestamp := time.Now().UTC().Unix() - 1800
		if timestamp <= req.TimeStamp {
			res := models.UpdateMessageRes{
				RoomId:        req.RoomId,
				MessageId:     req.MessageId,
				TimeStamp:     timestamp,
				Content:       req.Content,
				ContentType:   req.ContentType,
				UpdateMessage: true,
			}

			server.BroadcastToRoom("/", req.RoomId.Hex(), "modify", res)

			controllers.UpdateMessage(req)
		} else {
			res := models.UpdateMessageRes{
				RoomId:        req.RoomId,
				MessageId:     req.MessageId,
				TimeStamp:     timestamp,
				Content:       req.Content,
				ContentType:   req.ContentType,
				UpdateMessage: false,
			}

			server.BroadcastToRoom("/", req.RoomId.Hex(), "modify", res)
		}
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
	log.Fatal(http.ListenAndServe(":4001", nil)) // windows
	// log.Fatal(http.ListenAndServe(":4001", nil)) // linux

}
