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
	//Create connection witn server rabbitmq
	connection:="amqp://guest:guest@localhost:5672/"
 connectionRabbitMQ,error:=amqp.Dial(connection)
  if (error!=nil){
	log.Fatalf("could not connect to RabbitMQ: %v", error)
  }

  //open and return channel open
  amqpCan,err:=connectionRabbitMQ.Channel()
  if (err!=nil){
	log.Fatalf("could not connect to RabbitMQ: %v", err)
  }
  //close connection afther the end return
  defer connectionRabbitMQ.Close()


  fmt.Println("Starting Peril server, the connection was successful...")
  
// define value we send to the channel from this server
  val:=routing.PlayingState{
	IsPaused: true,
}

 // print command server help
  gamelogic.PrintServerHelp()
  

/****************************************************/


//enter the user in the server same like the client
//generate the queue name

/*_, queue, err := pubsub.DeclareAndBind(
	connectionRabbitMQ,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		0,
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)*/


/****************************************************/

//infinite loop permit continue value enter and process

for{
	
 if err!=nil {
	log.Fatalf("This is the problem: %v", err)
 }
	//enter command value give by help
	valueEnter:=gamelogic.GetInput()
/**/
	if len(valueEnter)==0{
		continue
	}
	//publish the value on the queue or do another thing  depending the command you enter
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
		return
	default:
		fmt.Println("I don't understand the command...")
		
	
}
// wait for ctrl+c
/*signalChan := make(chan os.Signal, 1)
signal.Notify(signalChan, os.Interrupt)
<-signalChan

fmt.Println("the program is shutting down ...")*/

}
}