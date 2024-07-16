package main

import (
	"goqmc5883"
	"log"
	"time"
)

func main() {
	m := goqmc5883.New()

	for i := 0; i < 10000; i++ {
		time.Sleep(100 * time.Millisecond)

		azmth, err := m.GetAzimuth()
		if err != nil {
			log.Printf("%v\n", err)
		}
		log.Printf("azuimuth: %d", azmth)
	}
}
