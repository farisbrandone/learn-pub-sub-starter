package pubsub

import (
	"fmt"
	"log"

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
	  log.Printf("poupoute %d\n",simpleQueueType)
	  if simpleQueueType==0{
		log.Println("Banko Banko ohoho")
		dur,auto,exclu=true,false,false
		
	 //transiant (transitoire queue type)
	  }else{
		dur,auto,exclu=false,true,true
		
	  }
	  ///now i pass amqp.Table inside QueueDeclare with x-dead-letter-exchange key 
	  //who take comme value the name of the dead exchange define inside ui rabbitMQ
	  newQueue,err:=amqpChan.QueueDeclare(queueName,dur,auto, exclu,false,amqp.Table{"x-dead-letter-exchange":"peril_dlx"})

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