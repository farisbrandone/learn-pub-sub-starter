package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)


func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error{
	marshVal, err := json.Marshal(val)
    fmt.Println("doncoma")
	if err!=nil {
		log.Fatalf("une erreur est survenue pendant la conversion en Json")
		return err
	}
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	 msg:= amqp.Publishing{
		ContentType:"application/json",
		Body:marshVal,
	}
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
   
	err=ch.PublishWithContext(ctx,exchange,key,false,false,msg)
	log.Printf("%v", err)
	if err!=nil{
		return err
	}
	return nil
}

