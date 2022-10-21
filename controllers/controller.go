package controllers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/achintya-7/go_socketio/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error oading .env file")
	}

	return os.Getenv(key)
}

var mongoURI = getDotEnvVariable("MONGO_URI")
var dbName = getDotEnvVariable("DB_NAME")


var database *mongo.Database

func init() {
	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB Connected")

	database = client.Database(dbName)

	fmt.Println("Name of Db :", database.Name())
}

func SendMessage(req models.SendMessageReq) {

	timestamp := time.Now().Unix()

	data := models.SendMessage{
		UserId:      req.UserId,
		RoomId:      req.RoomId,
		Content:     req.Content,
		PrevMessage: "",
		ContentType: req.ContentType,
		TimeStamp:   timestamp,
	}

	_, err := database.Collection("messages").InsertOne(context.Background(), data)
	if err != nil {
		fmt.Println("Insert Error : ", err)
	}

	roomId, _ := primitive.ObjectIDFromHex(req.RoomId.String())

	filter := bson.D{{Key: "_id", Value: roomId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "latestMessage", Value: req.Content}}}}

	val, err := database.Collection("chats").UpdateOne(context.Background(), filter, update)
	fmt.Println(val, err)
	fmt.Println("Message Sent")

}

func ReplyMessage(req models.ReplyMessageReq) {

	timestamp := time.Now().Unix()

	data := models.SendMessage{
		UserId:      req.UserId,
		RoomId:      req.RoomId,
		Content:     req.Content,
		ContentType: req.ContentType,
		TimeStamp:   timestamp,
		PrevMessage: req.PrevMessage,
	}

	roomId, _ := primitive.ObjectIDFromHex(req.RoomId.String())
	filter := bson.D{{Key: "_id", Value: roomId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "latestMessage", Value: req.Content}}}}

	_, err := database.Collection("messages").InsertOne(context.Background(), data)
	if err != nil {
		fmt.Println("Insert Error : ", err)
	}

	val, err := database.Collection("chats").UpdateOne(context.Background(), filter, update)
	fmt.Println(val, err)
	fmt.Println("Message Sent")

}

func ForwardMessage(req models.SendMessageReq) {

	timestamp := time.Now().Unix()

	data := models.SendMessage{
		UserId:      req.UserId,
		RoomId:      req.RoomId, // different room id
		Content:     req.Content,
		ContentType: req.ContentType,
		TimeStamp:   timestamp,
	}

	_, err := database.Collection("messages").InsertOne(context.Background(), data)
	if err != nil {
		fmt.Println("Insert Error : ", err)
	}

	forwardRoomId, _ := primitive.ObjectIDFromHex(req.RoomId.String())

	filter := bson.D{{Key: "_id", Value: forwardRoomId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "latestMessage", Value: req.Content}}}}

	val, err := database.Collection("chats").UpdateOne(context.Background(), filter, update)
	fmt.Println(val, err)
	fmt.Println("Message Sent")

}

// func DeleteMessage(req models.DeleteMessageReq) {

// 	var message primitive.M
// 	timestamp := time.Now().Unix()
// 	messageId, _ := primitive.ObjectIDFromHex(req.MessageId.String())

// 	filter := bson.D{{Key: "_id", Value: messageId}}

// 	cur := database.Collection("chats").FindOne(context.TODO(), filter)
// 	if cur == nil {
// 		log.Fatal("No Id Found")
// 	}

// 	err := cur.Decode(&message)
// 	if err != nil {
// 		log.Fatal("Error Decoding in data : ", err)
// 	}

// 	message



// 	// data := models.SendMessage{
// 	// 	UserId:      req.UserId,
// 	// 	RoomId:      req.ForwardRoomId, // different room id
// 	// 	Content:     req.Content,
// 	// 	ContentType: req.ContentType,
// 	// 	TimeStamp:   timestamp,
// 	// }

// 	// _, err := database.Collection("messages").InsertOne(context.Background(), data)
// 	// if err != nil {
// 	// 	fmt.Println("Insert Error : ", err)
// 	// }

// 	forwardRoomId, _ := primitive.ObjectIDFromHex(req.ForwardRoomId.String())

	

// 	val, err := database.Collection("chats").UpdateOne(context.Background(), filter, update)
// 	fmt.Println(val, err)
// 	fmt.Println("Message Sent")

// }
