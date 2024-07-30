package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"

	//"time"

	//"fmt"

	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)


func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error{
	marshVal, err := json.Marshal(val)
     log.Println("gringo gringo1")
	if err!=nil {
		log.Fatalf("une erreur est survenue pendant la conversion en Json")
		return err
	}
	//fmt.Printf("%v\n", marshVal)
	/*var (
		ctx    context.Context
		cancel context.CancelFunc
	)*/
	 msg:= amqp.Publishing{
		ContentType:"application/json",
		Body:marshVal,
	}
	/*ctx, cancel = context.WithCancel(context.Background())
	defer cancel()*/
	log.Println("gringo gringo2")
	err=ch.PublishWithContext(context.Background()/*ctx*/,exchange,key,false,false,msg)
	log.Printf("%v\n", err)
	if err!=nil{
		return err
	}
	log.Println("gringo gringo3")
	return nil
}

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T)error{
    log.Println("grinGob1")
	encGobValue, err := Encode(val)
	if err!=nil {
		log.Fatalf("une erreur est survenue pendant la conversion en Json")
		return err
	}
	log.Println("grinGob2")
	msg:= amqp.Publishing{
		ContentType:"application/gob",
		Body:encGobValue,
	}
	err=ch.PublishWithContext(context.Background()/*ctx*/,exchange,key,false,false,msg)
	log.Printf("%v\n", err)
	if err!=nil{
		return err
	}
	log.Println("grinGob3")
	return nil

}

func Encode(gl any) ([]byte, error) {//boot dev fast encoding more than JSON
	log.Println("encode")
	var network bytes.Buffer // Stand-in for the network.

	// Create an encoder and send a value.
	enc := gob.NewEncoder(&network)
	err := enc.Encode(gl)
	a:=network.Bytes()
	return a, err
}





