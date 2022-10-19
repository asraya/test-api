package main

import (
	"api-test/engine"
	"os"
)

func main() {
	r := engine.SetupRouter()
	r.Run(":" + os.Getenv("PORT"))
}
