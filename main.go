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
	Status        int
	DataFalha     string
}

func openFiles(serversFile string, downtimeFile string) (*os.File, *os.File) {
	servers, err := os.Open(serversFile)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	downtime, err := os.Open(downtimeFile)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	return servers, downtime

}

func criaListaDeServidores(data [][]string) []Server {
	var servers []Server
	for c, record := range data {
		if c > 0 {
			server := Server{
				Server:    record[0],
				ServerUrl: record[1],
			}
			servers = append(servers, server)
		}
	}
	return servers
}

func checkServers(servers []Server) {
	for _, server := range servers {
		now := time.Now()

		get, err := http.Get((server.ServerUrl))
		if err != nil {
			fmt.Println(err.Error())
		}
		server.TempoExecucao = time.Since(now).Seconds()
		server.Status = get.StatusCode

		fmt.Printf("%s Status: [%d] Tempo decorrido: [%f] segundos\n", server.Server, server.Status, server.TempoExecucao)
	}
	fmt.Println("--------------------------------------------")
}

func main() {

	serversList, downtimeList := openFiles(os.Args[1], os.Args[2])

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(serversList)
	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	servidores := criaListaDeServidores(records)

	for {
		checkServers(servidores)
		time.Sleep(5 * time.Second)
	}
}
