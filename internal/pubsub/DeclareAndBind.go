package pubsub

import (
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
		return nil,amqp.Queue{},err
	  }
	  if simpleQueueType==0{
		dur,auto,exclu=true,false,false
		
	  }else{
		dur,auto,exclu=false,true,true
		
	  }
	  newQueue,err:=amqpChan.QueueDeclare(queueName,dur,auto, exclu,false,nil)

	 if (err!=nil){
		return nil,amqp.Queue{},err
	  }
	  amqpChan.QueueBind(queueName,key,exchange,false,nil)

	  return amqpChan,newQueue,nil
}