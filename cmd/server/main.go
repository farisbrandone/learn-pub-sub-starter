package main

import (
	"fmt"
	"log"

	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)
func main() {
	
	connection:="amqp://guest:guest@localhost:5672/"
 connectionRabbitMQ,error:=amqp.Dial(connection)
  if (error!=nil){
	log.Fatalf("could not connect to RabbitMQ: %v", error)
  }
  amqpCan,err:=connectionRabbitMQ.Channel()
  if (err!=nil){
	log.Fatalf("could not connect to RabbitMQ: %v", err)
  }
  defer connectionRabbitMQ.Close()
  fmt.Println("Starting Peril server, the connection was successful...")
  /*signalChan := make(chan *amqp.Error, 1)
 signal:= connectionRabbitMQ.NotifyClose(signalChan)
<-signalChan*/
// wait for ctrl+c
val:=routing.PlayingState{
	IsPaused: true,
}
pubsub.PublishJSON(amqpCan,routing.ExchangePerilDirect,routing.PauseKey,val)
/*
amqpCan=cannal de rabbitmq sur lequel on se connecte
routing.EXch...=le exchange choisit dans rabbbitmq pour le transfert du messaga
pauseKey=la cle utiliser pour definir vers quels file d'attentes sera ennvoyé le message
val=valeur envoye or publish a rabbitmq pour transfert 
*/
signalChan := make(chan os.Signal, 1)
signal.Notify(signalChan, os.Interrupt)
<-signalChan
fmt.Println("the program is shutting down ...")
}
