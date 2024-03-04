package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

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
	Writer("./logs/subscriber_logs.txt", "name: "+data.name + " password: "+data.password + " age: "+strconv.Itoa(data.age) + " hours_spent: "+strconv.Itoa(data.hours_spent) + "\n")
}

var connectHandlerSub mqtt.OnConnectHandler = func(client mqtt.Client) {
	Writer("subscriber_logs.txt", "connected" + "\n")
}

var connectLostHandlerSub mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	var text = fmt.Sprintf("Connection lost: %v", err)
	Writer("subscriber_logs.txt",  text+ "\n")
}

func Subscriber(dbPointer *sql.DB) {
	fmt.Println("Subscriber")
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file: %s", err)
	}
	db = dbPointer

	var broker = os.Getenv("BROKER_ADDR")
	var port = 8883
	opts := mqtt.NewClientOptions()
	var connectionString = fmt.Sprintf("tls://%s%d/mqtt", broker, port)
	fmt.Println(connectionString)
	opts.AddBroker(connectionString)
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
}


