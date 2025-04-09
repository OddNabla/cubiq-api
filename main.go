package main

import (
	"github.com/DouglasValerio/cubiq-api/router"
	"github.com/DouglasValerio/cubiq-api/setup"
)

func main() {
	setup.Init()
	router.Initialize()
}
