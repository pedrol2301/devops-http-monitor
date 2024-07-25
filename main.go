package main

import (
	"encoding/csv"
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

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	var servers []Server
	for c, record := range records {
		if c > 0 {
			server := Server{
				Server:    record[0],
				ServerUrl: record[1],
			}
			servers = append(servers, server)
		}
	}

	for {
		for _, server := range servers {
			now := time.Now()

			get, err := http.Get(server.ServerUrl)
			if err != nil {
				fmt.Println(err.Error())
				panic(err)
			}
			decorrido := time.Since(now).Seconds()
			status := get.Status
			fmt.Printf("%s Status: %s Tempo decorrido: %f segundos\n", server.Server, status, decorrido)
		}

		fmt.Println("--------------------------------------------")
		time.Sleep(5 * time.Second)
	}
}
