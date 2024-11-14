package main

import (
	"github.com/souloss/go-clean-arch/pkg/logger"
)

func main() {
	app, err := initializeApp("")
	if err != nil {
		logger.L().Error("initialize Application error: %v", err)
		panic(err)
	}
	app.Run()
}
