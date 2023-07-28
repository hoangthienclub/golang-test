package common

import "fmt"

func AppRecover() {
	if r := recover(); r != nil {
		fmt.Println("Recovered in f", r)
	}
}
