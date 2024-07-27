package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error){

	var dur,auto,exclu bool

	amqpChan,err:=conn.Channel()
	if (err!=nil){
		return nil, amqp.Queue{}, fmt.Errorf("could not create channel: %v", err)
	  }
	  //durable queue type
	  if simpleQueueType==0{
		dur,auto,exclu=true,false,false
		
	 //transiant (transitoire queue type)
	  }else{
		dur,auto,exclu=false,true,true
		
	  }
	  
	  newQueue,err:=amqpChan.QueueDeclare(queueName,dur,auto, exclu,false,nil)

	 if (err!=nil){
		return nil, amqp.Queue{}, fmt.Errorf("could not declare queue: %v", err)
	  }
     fmt.Println("zinzin")
	  err=amqpChan.QueueBind(queueName,key,exchange,false,nil)
	  if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not bind queue: %v", err)
	}
	  return amqpChan,newQueue,nil
}