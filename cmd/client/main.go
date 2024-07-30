package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

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
  amqpCan,err:=connectionRabbitMQ.Channel()
  if (err!=nil){
	log.Fatalf("could not connect to RabbitMQ: %v", err)
  }
 user,err:=gamelogic.ClientWelcome()
 
 if err!=nil {
	log.Fatalf("This is the problem: %v", err)
 }
 /*pause:=routing.PauseKey*/
 /*queName:=fmt.Sprintf("%s.%s",pause,user)*/


 
//create gamestae
 value:=gamelogic.NewGameState(user)
 fmt.Printf("NEWGAMESTATE inside main client %v\n",value)
 //connect to or create Queue
 
 err=pubsub.SubscribeJSON(
	connectionRabbitMQ, 
	routing.ExchangePerilTopic,
	routing.WarRecognitionsPrefix,
	routing.WarRecognitionsPrefix+".*",
	0,
	handlerWare(value,amqpCan))

 if err!=nil {
	 log.Fatalf("This is the problem: %v", err)
  }

  err=pubsub.SubscribeGob(
	connectionRabbitMQ, 
	routing.ExchangePerilTopic,
	routing.GameLogSlug,
	routing.GameLogSlug+".*",
	0,
	handlerLog(),
)

 if err!=nil {
	 log.Fatalf("This is the problem: %v", err)
  }
 
 
  err = pubsub.SubscribeJSON(
	connectionRabbitMQ,
	routing.ExchangePerilTopic,
	routing.ArmyMovesPrefix+"."+value.GetUsername(),
	routing.ArmyMovesPrefix+".*",
	1,
	handlerMove(value, amqpCan),
)
if err != nil {
	log.Fatalf("could not subscribe to pause: %v", err)
}


/*err=pubsub.SubscribeGob(
	connectionRabbitMQ, 
	routing.ExchangePerilTopic,
	routing.GameLogSlug+"."+value.GetUsername(),
	routing.GameLogSlug+".*",
	0,
	handlerWare(value,amqpCan),
)

 if err!=nil {
	 log.Fatalf("This is the problem: %v", err)
  }*/

  err = pubsub.SubscribeJSON(
	connectionRabbitMQ,
	routing.ExchangePerilDirect,
	routing.PauseKey+"."+value.GetUsername(),
	routing.PauseKey,
	1,
	handlerPause(value),
)
if err != nil {
	log.Fatalf("could not subscribe to pause: %v", err)
}



 //infinite loop
for{
	
	fmt.Println("enter the command value of client")
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
			continue
		   }
		   err=pubsub.PublishJSON(amqpCan,routing.ExchangePerilTopic,/*queName*/routing.ArmyMovesPrefix+"."+moveValue.Player.Username,moveValue)
		if err != nil {
			log.Printf("could not publish time: %v", err)
			continue
		}	
		   fmt.Printf("the move of  %v is successful\n", moveValue.Player.Username)
		   continue
	case"statut":
		value.CommandStatus()
		
	case"help":
		gamelogic.PrintClientHelp()
		
	case"spam":
		if len(valueEnter)<2{
			log.Println("Spamming not allowed!")
			continue
		}
		intConvert, err:=strconv.Atoi(valueEnter[1])
		if err!=nil {
			log.Println("Spamming errors convert to string!")
			continue
		}
		if intConvert<=0 {
			log.Println("Spamming errors value is negative or zero")
			continue
		}
		for i := 0; i < intConvert; i++ {
			maliciousValue:=gamelogic.GetMaliciousLog()
			err := pubsub.PublishGob(
				amqpCan,
				routing.ExchangePerilTopic,
				routing.GameLogSlug+"."+value.GetUsername(),
				routing.GameLog{
					CurrentTime:time.Now(),
					Message:maliciousValue ,    
					Username : value.GetUsername(), 
		
				},
			)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				
			}
		} 
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

func publishGameLog(publishCh *amqp.Channel, username, msg string) error {
	return pubsub.PublishGob(
		publishCh,
		routing.ExchangePerilTopic,
		routing.GameLogSlug+"."+username,
		routing.GameLog{
			Username:    username,
			CurrentTime: time.Now(),
			Message:     msg,
		},
	)
}
