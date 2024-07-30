package main

/*import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
)*/

/*func handlerPause(gameState *gamelogic.GameState)(func(gamelogic.ArmyMove)string){
	println("enter inside handle after process Subscribejson")
	return func (ps gamelogic.ArmyMove)string{
		fmt.Println("Handlepause inside client")
		defer fmt.Print("> ")
	    return gameState.HandleMove(ps)
	}
}*/

//package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerMove(gs *gamelogic.GameState, publishCh *amqp.Channel) func(gamelogic.ArmyMove) pubsub.Acktype {
	return func(move gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")
		moveOutcome := gs.HandleMove(move)
		switch moveOutcome {
		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.NackDiscard
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			err := pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilTopic,
				routing.WarRecognitionsPrefix+"."+gs.GetUsername(),
				gamelogic.RecognitionOfWar{
					Attacker: move.Player,
					Defender: gs.GetPlayerSnap(),
				},
			)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		}
		fmt.Println("error: unknown move outcome")
		return pubsub.NackDiscard
	}
}
func handlerLog() func(routing.GameLog) pubsub.Acktype {
	return func(ps routing.GameLog) pubsub.Acktype {
		defer fmt.Print("> ")
		gamelogic.WriteLog(ps)
		return pubsub.Ack
	}
}

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.Acktype {
	return func(ps routing.PlayingState) pubsub.Acktype {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
		return pubsub.Ack
	}
}

func handlerWare(gs *gamelogic.GameState,publishCh *amqp.Channel)func(gamelogic.RecognitionOfWar) pubsub.Acktype {
	return func(move gamelogic.RecognitionOfWar) pubsub.Acktype {
        log.Println("bounga bounga")
		defer fmt.Print("> ")
		var message string
		outcome,winner,loser := gs.HandleWar(move)
		
        log.Printf("OUHHHHH %v\n",outcome)
		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
			
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
			
		case gamelogic.WarOutcomeOpponentWon:
			message=fmt.Sprintf("%s won a war against %s",winner, loser)
			err := pubsub.PublishGob(
				publishCh,
				routing.ExchangePerilTopic,
				routing.GameLogSlug+"."+gs.GetUsername(),
				routing.GameLog{
					CurrentTime:time.Now(),
					Message:message ,    
					Username : gs.GetUsername(), 
		
				},
			)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeYouWon:
			message=fmt.Sprintf("%s won a war against %s",winner, loser)
			err := pubsub.PublishGob(
				publishCh,
				routing.ExchangePerilTopic,
				routing.GameLogSlug+"."+gs.GetUsername(),
				routing.GameLog{
					CurrentTime:time.Now(),
					Message:message ,    
					Username : gs.GetUsername(), 
		
				},
			)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			message=fmt.Sprintf("A war between %s and %s resulted in a draw",winner, loser)
			err := pubsub.PublishGob(
				publishCh,
				routing.ExchangePerilTopic,
				routing.GameLogSlug+"."+gs.GetUsername(),
				routing.GameLog{
					CurrentTime:time.Now(),
					Message:message ,    
					Username : gs.GetUsername(), 
		
				},
			)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		}

		
		fmt.Println("error: unknown move outcome")
		return pubsub.Ack
	}
	
	
	
}