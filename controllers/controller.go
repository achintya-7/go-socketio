package controllers

import (
	"context"
	"fmt"
	"log"

	"github.com/achintya-7/go_socketio/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var MONGO_URI = utils.GetDotEnvVariable("MONGO_URI")
// var DB_NAME = utils.GetDotEnvVariable("DB_NAME")

var MONGO_URI string = "mongodb+srv://Achintya:X7hunter@cluster0.wuvef.mongodb.net/test"
var DB_NAME string = "test"

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

	data := models.SendMessageDB{
		UserId:      req.UserId,
		RoomId:      req.RoomId,
		Content:     req.Content,
		PrevMessage: "",
		ContentType: req.ContentType,
		TimeStamp:   req.Timestamp,
		MessageId:   req.MessageId,
	}

	_, err := database.Collection("messages").InsertOne(context.Background(), data)
	if err != nil {
		fmt.Println("Insert Error : ", err)
	}

	roomId, _ := primitive.ObjectIDFromHex(req.RoomId.Hex())
	filter := bson.D{{Key: "_id", Value: roomId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "latestmessageid", Value: req.MessageId}, {Key: "latestmessage", Value: req.Content}}}}

	val, err := database.Collection("chats").UpdateOne(context.Background(), filter, update)
	fmt.Println(val, err)
	fmt.Println("Message Sent")

}

func ReplyMessage(req models.ReplyMessageRes) {

	data := models.SendMessageDB{
		UserId:      req.UserId,
		RoomId:      req.RoomId,
		Content:     req.Content,
		ContentType: req.ContentType,
		TimeStamp:   req.Timestamp,
		PrevMessage: req.PrevMessage,
		MessageId:   req.MessageId,
	}

	roomId, _ := primitive.ObjectIDFromHex(req.RoomId.Hex())
	filter := bson.D{{Key: "_id", Value: roomId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "latestmessageid", Value: req.MessageId}, {Key: "latestmessage", Value: req.Content}}}}

	_, err := database.Collection("messages").InsertOne(context.Background(), data)
	if err != nil {
		fmt.Println("Insert Error : ", err)
	}

	val, err := database.Collection("chats").UpdateOne(context.Background(), filter, update)
	fmt.Println(val, err)
	fmt.Println("Message Sent")

}

func DeleteMessage(req models.DeleteMessageReq) {

	fmt.Println(req.MessageId)
	filter := bson.D{{Key: "messageid", Value: req.MessageId}}
	val, err := database.Collection("messages").DeleteOne(context.Background(), filter)
	if err != nil {
		fmt.Println("Delete Message Error", err)
	}
	fmt.Println(val)

	filter2 := bson.D{{Key: "_id", Value: req.RoomId}, {Key: "latestmessageid", Value: req.MessageId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "latestmessageid", Value: ""}, {Key: "latestmessage", Value: ""}}}}
	val2, err := database.Collection("chats").UpdateOne(context.Background(), filter2, update)
	if err != nil {
		fmt.Println("Delete Message Error", err)
	}

	fmt.Println(val2, err)

}

func UpdateMessage(req models.UpdateMessageReq) {

	filter := bson.D{{Key: "messageid", Value: req.MessageId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "content", Value: req.Content}, {Key: "contenttype", Value: req.ContentType}}}}
	_, err := database.Collection("messages").UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println("Delete Message Error", err)
	}

	filter2 := bson.D{{Key: "_id", Value: req.RoomId}, {Key: "latestmessageid", Value: req.MessageId}}
	update2 := bson.D{{Key: "$set", Value: bson.D{{Key: "latestmessageid", Value: req.MessageId}, {Key: "latestMessage", Value: req.Content}}}}
	val, err := database.Collection("chats").UpdateOne(context.Background(), filter2, update2)
	if err != nil {
		fmt.Println("Update Message Error", err)
	}

	fmt.Println(val, err)

}
