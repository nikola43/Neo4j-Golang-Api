package main

import (
	"github.com/nikola43/ecodadys_api/controllers"
)

func main() {
	a := controllers.App{}
	a.Initialize()
	a.Run(":8080")
}
