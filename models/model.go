package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SendMessageReq struct {
	UserId      primitive.ObjectID `json:"userId"`
	RoomId      primitive.ObjectID `json:"roomId"`
	Content     string             `json:"content"`
	ContentType string             `json:"content_type"`
}

type SendMessageRes struct {
	UserId      primitive.ObjectID `json:"userId"`
	RoomId      primitive.ObjectID `json:"roomId"`
	Content     string             `json:"content"`
	ContentType string             `json:"content_type"`
	MessageId   string             `json:"message_Id"`
}

type ReplyMessageReq struct {
	UserId      primitive.ObjectID `json:"userId"`
	RoomId      primitive.ObjectID `json:"roomId"`
	Content     string             `json:"content"`
	ContentType string             `json:"content_type"`
	PrevMessage string             `json:"prev_message"`
}

type ReplyMessageRes struct {
	UserId      primitive.ObjectID `json:"userId"`
	RoomId      primitive.ObjectID `json:"roomId"`
	Content     string             `json:"content"`
	ContentType string             `json:"content_type"`
	PrevMessage string             `json:"prev_message"`
	MessageId   string             `json:"message_Id"`
}

type ForwardMessageReq struct {
	UserId        primitive.ObjectID `json:"userId"`
	RoomId        primitive.ObjectID `json:"roomId"`
	Content       string             `json:"content"`
	ContentType   string             `json:"content_type"`
	ForwardRoomId primitive.ObjectID `json:"forward_room_id"`
}

type DeleteMessageReq struct {
	UserId      primitive.ObjectID `json:"userId"`
	RoomId      primitive.ObjectID `json:"roomId"`
	Content     string             `json:"content"`
	ContentType string             `json:"content_type"`
	MessageId   primitive.ObjectID `json:"forward_room_id"`
}

type SendMessage struct {
	UserId      primitive.ObjectID `json:"userId"`
	RoomId      primitive.ObjectID `json:"roomId"`
	MessageId   string `json:"messageId"`
	Content     string             `json:"content"`
	ContentType string             `json:"content_type"`
	PrevMessage string             `json:"prev_message"`
	TimeStamp   int64              `json:"timestamp"`
}
