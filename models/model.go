package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SendMessageReq struct {
	UserId      primitive.ObjectID `json:"userId"`
	RoomId      primitive.ObjectID `json:"roomId"`
	Content     string             `json:"content"`
	ContentType string             `json:"contentType"`
}

type SendMessageRes struct {
	UserId      primitive.ObjectID `json:"userId"`
	RoomId      primitive.ObjectID `json:"roomId"`
	Content     string             `json:"content"`
	ContentType string             `json:"contentType"`
	MessageId   string             `json:"messageId"`
}

type ReplyMessageReq struct {
	UserId        primitive.ObjectID `json:"userId"`
	RoomId        primitive.ObjectID `json:"roomId"`
	Content       string             `json:"content"`
	ContentType   string             `json:"contentType"`
	PrevMessage   string             `json:"prevMessage"`
	PrevMessageId string             `json:"prevMessageId"`
}

type ReplyMessageRes struct {
	UserId        primitive.ObjectID `json:"userId"`
	RoomId        primitive.ObjectID `json:"roomId"`
	Content       string             `json:"content"`
	ContentType   string             `json:"contentType"`
	PrevMessage   string             `json:"prevMessage"`
	MessageId     string             `json:"messageId"`
	PrevMessageId string             `json:"prevMessageId"`
}

type DeleteMessageReq struct {
	RoomId    primitive.ObjectID `json:"roomId"`
	MessageId primitive.ObjectID `json:"messageId"`
	TimeStamp int64              `json:"timestamp"`
}

type DeleteMessageRes struct {
	RoomId        primitive.ObjectID `json:"roomId"`
	MessageId     primitive.ObjectID `json:"messageId"`
	TimeStamp     int64              `json:"timestamp"`
	DeleteMessage bool               `json:"deleteMessage"`
}

type UpdateMessageReq struct {
	RoomId      primitive.ObjectID `json:"roomId"`
	MessageId   primitive.ObjectID `json:"messageId"`
	TimeStamp   int64              `json:"timestamp"`
	Content     string             `json:"content"`
	ContentType string             `json:"contentType"`
}

type UpdateMessageRes struct {
	RoomId        primitive.ObjectID `json:"roomId"`
	MessageId     primitive.ObjectID `json:"messageId"`
	TimeStamp     int64              `json:"timestamp"`
	Content       string             `json:"content"`
	ContentType   string             `json:"contentType"`
	UpdateMessage bool               `json:"updateMessage"`
}

type SendMessageDB struct {
	UserId      primitive.ObjectID `json:"userId"`
	RoomId      primitive.ObjectID `json:"roomId"`
	MessageId   string             `json:"messageId"`
	Content     string             `json:"content"`
	ContentType string             `json:"contentType"`
	PrevMessage string             `json:"prevMessage"`
	TimeStamp   int64              `json:"timestamp"`
}
