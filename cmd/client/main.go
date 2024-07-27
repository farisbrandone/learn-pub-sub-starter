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
	fmt.Println("Starting Peril client...")
	connection:="amqp://guest:guest@localhost:5672/"
 connectionRabbitMQ,error:=amqp.Dial(connection)
 if (error!=nil){
	log.Fatalf("could not connect to RabbitMQ: %v", error)
  }
  defer connectionRabbitMQ.Close()
 user,err:=gamelogic.ClientWelcome()
 pause:=routing.PauseKey
 queName:=fmt.Sprintf("%s.%s",pause,user)
 if err!=nil {
	log.Fatalf("This is the problem: %v", err)
 }
_,newQueue,err:= pubsub.DeclareAndBind(connectionRabbitMQ,routing.ExchangePerilDirect,queName,routing.PauseKey,1)



if err!=nil {
	log.Fatalf("This is the problem: %v", err)

 }
 fmt.Printf("Queue %v declared and bound!\n", newQueue.Name)
 value:=gamelogic.NewGameState(user)
for{
	valueEnter:=gamelogic.GetInput()
	
	if len(valueEnter)==0{
		continue
	}
	switch valueEnter[0] {
	case "spawn":
		fmt.Println("the message spawn value is sending ...")
		err=value.CommandSpawn(valueEnter)
       if err!=nil {
		log.Printf("This is the problem: %v\n", err)
        }
		continue
		
	case "move":
		fmt.Println("the message move value is sending ...")
		moveValue,err:=value.CommandMove(valueEnter)
		if err!=nil {
			log.Printf("This is the problem: %v", err)
		   }
		   fmt.Printf("the move of  %v is successful\n", moveValue.Player.Username)
		   continue
	case"statut":
		value.CommandStatus()
		
	case"help":
		gamelogic.PrintClientHelp()
		
	case"spam":
		log.Println("Spamming not allowed yet!")
		
	case"quit":
		gamelogic.PrintQuit()
		return
	default:
		fmt.Println("unknown command ....")
}
}	


// wait for ctrl+c
	/*signalChan := make(chan os.Signal, 1)
signal.Notify(signalChan, os.Interrupt)
<-signalChan*/

}
