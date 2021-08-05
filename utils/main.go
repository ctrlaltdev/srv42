package utils

import (
	"fmt"
	"log"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Debug(debug *bool, msg string) {
	if *debug {
		fmt.Printf("[debug] %s\n", msg)
	}
}
