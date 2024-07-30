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
	servers, err := os.OpenFile(serversFile, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	downtime, err := os.OpenFile(downtimeFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	return servers, downtime

}

func criaListaDeServidores(file *os.File) []Server {
	csvReader := csv.NewReader(file)
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

	return servers
}

func checkServers(servers []Server) []Server {
	var downServers []Server

	for _, server := range servers {
		agora := time.Now()

		get, err := http.Get((server.ServerUrl))
		if err != nil {
			fmt.Println("Erro ao tentar acessar o servidor: ", server.Server, err.Error())
			server.Status = 0
			server.DataFalha = agora.Format("02/01/2006 15:04:05")
			downServers = append(downServers, server)
			continue
		}
		server.Status = get.StatusCode

		if server.Status != 200 {
			server.DataFalha = agora.Format("02/01/2006 15:04:05")
			downServers = append(downServers, server)
		}
		server.TempoExecucao = time.Since(agora).Seconds()
		fmt.Printf("%s Status: [%d] Tempo decorrido: [%f] segundos\n", server.Server, server.Status, server.TempoExecucao)
	}

	fmt.Println("--------------------------------------------")

	return downServers
}

func generateDowntimeReport(file *os.File, servers []Server) {
	csvWriter := csv.NewWriter(file)

	for _, server := range servers {
		record := []string{server.Server, server.ServerUrl, fmt.Sprintf("%d", server.Status), server.DataFalha}
		err := csvWriter.Write(record)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
	}
	csvWriter.Flush()
}

func main() {

	serversList, downtimeList := openFiles(os.Args[1], os.Args[2])
	defer serversList.Close()
	defer downtimeList.Close()

	servidores := criaListaDeServidores(serversList)

	for {
		downServers := checkServers(servidores)
		generateDowntimeReport(downtimeList, downServers)
		time.Sleep(5 * time.Second)
	}
}
