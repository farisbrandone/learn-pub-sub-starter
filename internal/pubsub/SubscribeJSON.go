package pubsub

import (
	"encoding/json"
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)



func SubscribeJSON[T routing.PlayingState](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, //type of queue
	handler func(T),
) error {
      fmt.Println("First enter inside SubscribeJSON")
	//make sure you connected to the queeue declare on the client main 
	ch,newQueue,err:=DeclareAndBind(conn,exchange,queueName,key,simpleQueueType) 

	if err!=nil {
		return fmt.Errorf("could not declare and bind queue: %v", err)
	}
	//open the channel
	

     //connect deliver on the queue define inside channel
	 //the deliver consume the value inside channel
	deliver,err:=ch.Consume(newQueue.Name,"",false,false,false,false,nil)
	if err!=nil {
		return fmt.Errorf("could not consume messages: %v", err)
	}

	unmarshaller := func(data []byte) (T, error) {
		var target T
		err := json.Unmarshal(data, &target)
		return target, err
	}


    //fmt.Printf("Print the deliver value %v\n", deliver)
	//retrieve the value inside channel and process it



    go func() {
		defer ch.Close()
		for msg := range deliver {
			target, err := unmarshaller(msg.Body)
			if err != nil {
				fmt.Printf("could not unmarshal message: %v\n", err)
				continue
			}
			handler(target)
			msg.Ack(false)
		}
	}()
	return nil
}




	



/*func RangeDeliver[T routing.PlayingState](msg amqp.Delivery, handler func(T)){
	fmt.Println("inside Range deliver")
	 var convertMsg T
	 json.Unmarshal(msg.Body, &convertMsg)
     fmt.Printf("print the message take on the body deliver queue %v\n", convertMsg)
     handler(convertMsg)
     msg.Ack(false)

	}*/
