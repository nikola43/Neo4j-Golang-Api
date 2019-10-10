package utils

import (
	"fmt"
	"log"
)

func CheckIfErrorExists(err error) {
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal()
	}
}
