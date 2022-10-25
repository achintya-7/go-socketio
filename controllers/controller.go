package controllers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/achintya-7/go_socketio/models"
	"github.com/achintya-7/go_socketio/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MONGO_URI = utils.GetDotEnvVariable("MONGO_URI")
var DB_NAME = utils.GetDotEnvVariable("DB_NAME")

var database *mongo.Database

func init() {
	clientOptions := options.Client().ApplyURI(MONGO_URI)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB Connected")

	database = client.Database(DB_NAME)

	fmt.Println("Name of Db :", database.Name())
}

func GetUser(userId string) bool {
	userIdHex, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		fmt.Println("Error", err)
	}

	filter := bson.D{{Key: "_id", Value: userIdHex}}
	err2 := database.Collection("users").FindOne(context.Background(), filter)
	if err2 != nil {
		if err2.Err() == mongo.ErrNoDocuments {
			fmt.Println("No user")
			return false
		}
	}
	fmt.Println(userId)
	return true
}

func SendMessage(req models.SendMessageRes) {

	timestamp := time.Now().UTC().Unix()

	data := models.SendMessageDB{
		UserId:      req.UserId,
		RoomId:      req.RoomId,
		Content:     req.Content,
		PrevMessage: "",
		ContentType: req.ContentType,
		TimeStamp:   timestamp,
		MessageId:   req.MessageId,
	}

	_, err := database.Collection("messages").InsertOne(context.Background(), data)
	if err != nil {
		fmt.Println("Insert Error : ", err)
	}

	roomId, _ := primitive.ObjectIDFromHex(req.RoomId.Hex())
	filter := bson.D{{Key: "_id", Value: roomId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "latestMessageId", Value: req.MessageId}, {Key: "latestMessage", Value: req.Content}}}}

	val, err := database.Collection("chats").UpdateOne(context.Background(), filter, update)
	fmt.Println(val, err)
	fmt.Println("Message Sent")

}

func ReplyMessage(req models.ReplyMessageRes) {

	timestamp := time.Now().UTC().Unix()

	data := models.SendMessageDB{
		UserId:      req.UserId,
		RoomId:      req.RoomId,
		Content:     req.Content,
		ContentType: req.ContentType,
		TimeStamp:   timestamp,
		PrevMessage: req.PrevMessage,
		MessageId:   req.MessageId,
	}

	roomId, _ := primitive.ObjectIDFromHex(req.RoomId.Hex())
	filter := bson.D{{Key: "_id", Value: roomId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "latestMessageId", Value: req.MessageId}, {Key: "latestMessage", Value: req.Content}}}}

	_, err := database.Collection("messages").InsertOne(context.Background(), data)
	if err != nil {
		fmt.Println("Insert Error : ", err)
	}

	val, err := database.Collection("chats").UpdateOne(context.Background(), filter, update)
	fmt.Println(val, err)
	fmt.Println("Message Sent")

}

func DeleteMessage(req models.DeleteMessageReq) {

	filter := bson.D{{Key: "messageId", Value: req.MessageId}}
	_, err := database.Collection("messages").DeleteOne(context.Background(), filter)
	if err != nil {
		fmt.Println("Delete Message Error", err)
	}

	filter2 := bson.D{{Key: "_id", Value: req.RoomId}, {Key: "latestMessageId", Value: req.MessageId}}
	update := bson.D{{Key: "latestMessageId", Value: ""}, {Key: "latestMessage", Value: ""}}
	val, err := database.Collection("chats").UpdateOne(context.Background(), filter2, update)
	if err != nil {
		fmt.Println("Update Message Error", err)
	}

	fmt.Println(val, err)

}

func UpdateMessage(req models.UpdateMessageReq) {

	filter := bson.D{{Key: "messageId", Value: req.MessageId}}
	update := bson.D{{Key: "content", Value: req.Content}, {Key: "contentType", Value: req.ContentType}}
	_, err := database.Collection("messages").UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println("Delete Message Error", err)
	}

	filter2 := bson.D{{Key: "_id", Value: req.RoomId}, {Key: "latestMessageId", Value: req.MessageId}}
	update2 := bson.D{{Key: "latestMessageId", Value: req.MessageId}, {Key: "latestMessage", Value: req.Content}}
	val, err := database.Collection("chats").UpdateOne(context.Background(), filter2, update2)
	if err != nil {
		fmt.Println("Update Message Error", err)
	}

	fmt.Println(val, err)

}
