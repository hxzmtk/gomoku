package main

import (
	"github.com/bzyy/gobang/api"
)

func main() {
	g := api.InitRouter()
	g.Run(":8000")
}
