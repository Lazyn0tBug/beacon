// main.go

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Lazyn0tBug/beacon/server/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create the logger based on the configuration
	Logger=utils.GetLogger() 

	r := gin.Default()

	// Now you can use the logger
	Logger.Info("This is an info message")
	Logger.Error("This is an error message")

	for i := 0; i < 12; i++ {
		go Logger.Info(fmt.Sprintf("test log: %d", i))
	}
	time.Sleep(time.Second * 3)
}
