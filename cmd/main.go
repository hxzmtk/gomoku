package main

import (
	"github.com/bzyy/gomoku/api"
)

func main() {
	g := api.InitRouter()
	g.Run(":8000")
}
