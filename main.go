package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type Server struct {
	Server        string
	ServerUrl     string
	TempoExecucao float64
}

func main() {

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	defer f.Close()

	f.Read([]byte(""))

	now := time.Now()
	url := os.Args[1]

	get, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	decorrido := time.Since(now).Seconds()
	status := get.Status
	fmt.Printf("Status: %s Tempo decorrido: %f segundos", status, decorrido)

}
