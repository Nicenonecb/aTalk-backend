package main

import (
	"aTalkBackEnd/internal/app/router"
)

func main() {
	r := router.SetupRoutes()
	r.Run()
}
