package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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
  gamelogic.PrintServerHelp()
  
// wait for ctrl+c
val:=routing.PlayingState{
	IsPaused: true,
}
//pubsub.PublishJSON(amqpCan,routing.ExchangePerilDirect,routing.PauseKey,val)
/*
amqpCan=cannal de rabbitmq sur lequel on se connecte
routing.EXch...=le exchange choisit dans rabbbitmq pour le transfert du messaga
pauseKey=la cle utiliser pour definir vers quels file d'attentes sera ennvoyé le message
val=valeur envoye or publish a rabbitmq pour transfert 
*/
for{
	valueEnter:=gamelogic.GetInput()
/**/
	if len(valueEnter)==0{
		continue
	}
	switch valueEnter[0] {
	case "pause":
		fmt.Println("the message pause is sending ...")
		val.IsPaused=true
		err=pubsub.PublishJSON(amqpCan,routing.ExchangePerilDirect,routing.PauseKey,val)
		if err != nil {
			log.Printf("could not publish time: %v", err)
		}	
	case "resume":
		fmt.Println("the message resume is sending ...")
		val.IsPaused=false
		err=pubsub.PublishJSON(amqpCan,routing.ExchangePerilDirect,routing.PauseKey,val)
		if err != nil {
			log.Printf("could not publish time: %v", err)
		}
	case"quit":
		log.Println("your are exiting ...")
		break
	default:
		fmt.Println("I don't understand the command...")
		
	
}

/*signalChan := make(chan os.Signal, 1)
signal.Notify(signalChan, os.Interrupt)
<-signalChan

fmt.Println("the program is shutting down ...")*/

}
}