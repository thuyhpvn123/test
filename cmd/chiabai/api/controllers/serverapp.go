package controllers

import (
	// "encoding/json"
	"fmt"
	"log"
	"net/http"

	// "net/http"
	"sync"
	// "github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"
	// c_config "github.com/meta-node-blockchain/meta-node/cmd/chiabai/config"
	c_config "github.com/meta-node-blockchain/meta-node/cmd/client/pkg/config"

	"github.com/meta-node-blockchain/meta-node/cmd/chiabai/core"
	// controller_client "github.com/meta-node-blockchain/meta-node/cmd/client/pkg/controllers"
	// "github.com/meta-node-blockchain/meta-node/pkg/bls"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	// "github.com/meta-node-blockchain/meta-node/pkg/network"
)

type Server struct {
	sync.Mutex
	contractABI map[string]*core.ContractABI
	config      *c_config.ClientConfig
	cli         Cli
}
type Message1 struct {
	Command string      `json:"command"`
	Data    interface{} `json:"data"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (server *Server) Init(config *c_config.ClientConfig) *Server {
	// init subscriber
	server.config = config
	server.contractABI = make(map[string]*core.ContractABI)
	var wg sync.WaitGroup
	for _, contract := range core.Contracts {
		wg.Add(1)
		go server.getABI(&wg, contract)
	}
	wg.Wait()

	fmt.Println("the end")
	return &Server{
		contractABI: server.contractABI,
		config:      config,
	}
}

func (server *Server) ConnectionHandler() *Cli {
	config, err := c_config.LoadConfig(c_config.CONFIG_FILE_PATH)
	if err != nil {
		logger.Error(fmt.Sprintf("error when loading config %v", err))
		panic(fmt.Sprintf("error when loading config %v", err))
	}
	cConfig := config.(*c_config.ClientConfig)

	cli, err := NewCli(
		cConfig,
		server,
	)
	if err != nil {
		logger.Error(fmt.Sprintf("error when connect to parent %v", err))
		return nil
	}

	// server.cli = cli
	logger.Info("Cli init successfully") //write on server terminal
	// defer server.clients.Remove(conn)
	return cli

}
func (server *Server) WebsocketHandler(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	config, err := c_config.LoadConfig(c_config.CONFIG_FILE_PATH)
	if err != nil {
		logger.Error(fmt.Sprintf("error when loading config %v", err))
		panic(fmt.Sprintf("error when loading config %v", err))
	}
	cConfig := config.(*c_config.ClientConfig)
	cli:= server.ConnectionHandler()
	client := Client{
		ws: conn,
		server: server,
		config: cConfig,
		cli:    cli,
	}
	client.init()
	log.Println("Client Connected successfully") //write on server terminal
	// defer server.clients.Remove(conn)

	//listen websocket
	client.handleListen()

}

func (server *Server) getABI(wg *sync.WaitGroup, contract core.Contract) {
	var temp core.ContractABI
	temp.InitContract(contract)
	server.Lock()
	server.contractABI[contract.Name] = &temp
	server.Unlock()
	wg.Done()
}
