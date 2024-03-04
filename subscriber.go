package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"database/sql"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	godotenv "github.com/joho/godotenv"
)
type Data struct {
	name string
	password string
	age int
	hours_spent int
}
var data *Data;

var db *sql.DB

var messagePubHandlerSub mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var text = fmt.Sprintf("Recebido: %s do tópico: %s com QoS: %d\n", msg.Payload(), msg.Topic(), msg.Qos())
	Writer("./logs/subscriber_logs.txt",  text+ "\n")
	result := strings.Split(string(msg.Payload()), ",")
	age,_ := strconv.Atoi(result[2])
	hours_spent_value, _ := strconv.Atoi(result[3])
	data = &Data{name: result[0], password: result[1], age: age, hours_spent: hours_spent_value}
	Insert(db, *data)
	fmt.Printf("name: "+data.name, "password: "+data.password, "age: "+strconv.Itoa(data.age), "hours_spent: "+strconv.Itoa(data.hours_spent))
}

var connectHandlerSub mqtt.OnConnectHandler = func(client mqtt.Client) {
	Writer("subscriber_logs.txt", "connected" + "\n")
}

var connectLostHandlerSub mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	var text = fmt.Sprintf("Connection lost: %v", err)
	Writer("subscriber_logs.txt",  text+ "\n")
}

func Subscriber(db *sql.DB) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file: %s", err)
	}
	db = SetPointer(db)

	var broker = os.Getenv("BSetPointerROKER_ADDR")
	var port = 8883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tls://%s%d/mqtt", broker, port))
	opts.SetClientID("Subscriber")
	opts.SetUsername(os.Getenv("HIVE_USER"))
	opts.SetPassword(os.Getenv("HIVE_PSWD"))
	opts.SetDefaultPublishHandler(messagePubHandlerSub)
	opts.OnConnect = connectHandlerSub
	opts.OnConnectionLost = connectLostHandlerSub

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe("test/topic", 1, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return
	}

	select {}
}


func SetPointer(dbPointer *sql.DB) *sql.DB{
	
	return dbPointer;
}

