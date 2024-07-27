package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gameState *gamelogic.GameState)func(routing.PlayingState){
	println("enter inside handle after process Subscribejson")
	return func (ps routing.PlayingState){
		fmt.Println("The beginning return handle value")
		defer fmt.Print("> ")
	gameState.HandlePause(ps)
	}
}