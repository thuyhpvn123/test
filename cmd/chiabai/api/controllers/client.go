package controllers

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	// "github.com/meta-node-blockchain/meta-node/cmd/chiabai/config"
	log "github.com/sirupsen/logrus"
	"github.com/meta-node-blockchain/meta-node/cmd/client/pkg/config"

	// controller_client "github.com/meta-node-blockchain/meta-node/cmd/client/pkg/controllers"
	// "github.com/meta-node-blockchain/meta-node/pkg/bls"
	// "github.com/meta-node-blockchain/meta-node/pkg/network"
	// "github.com/meta-node-blockchain/meta-node/pkg/state"
)

type Client struct {
	ws     *websocket.Conn
	server *Server
	cli *Cli
	sync.Mutex
	config *config.ClientConfig
}

func (client *Client) init() {
	go client.handleMessage()
	log.Info("End init client")
}
func (client *Client) handleListen() {
	for {
		// Read in a new message as JSON and map it to a Message object
		var msg map[string]interface{}
		err := client.ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			client.ws.Close()
			break
		}
		// log.Info("Message from client: ", msg)
		client.handleCallChain(msg)
	}
}

// handle message struct tu chain tra ve va chuyen qua dang JSON gui toi cac client
func (client *Client) handleMessage() {
	for {
		msg := <-client.cli.sendChan
		// msg1 := <-sendDataC
		log.Info(msg)
		err := client.ws.WriteJSON(msg)

		if err != nil {
			log.Printf("error: %v", err)
			client.ws.Close()
		}
	}
}
func (client *Client) handleCallChain(msg map[string]interface{}) {
	// handle call
	switch msg["command"] {
	case "deal-cards":
		call := msg["value"].(map[string]interface{})
		kq := GeneratePlayerKeys(call)
		call1 := CreateDeck()
		deck := ShuffleDeck(call1)
		// call2:=msg["value"].(map[string]interface{})
		playerkeys := kq["message"].([]string)
		kq2 := EncryptDeck(deck, playerkeys)
		encryptDeck := kq2
		fmt.Println("encryptDeck:", encryptDeck)
		roomNumber:=call["roomNumber"].(string)
		player:="fdd11471417109d88c48030e579f3523e485f6fa"
		// client.cli.SetDeck(roomNumber,encryptDeck)
		client.cli.SetDeck(roomNumber,player)

		// client.cli.SetDeck()

		// client.cli.Caller.sentToClient("shuffle-and-encrypt-cards","success")

	// case "decrypt-cards":
	// 	call:=msg["value"].(map[string]interface{})
	// 	kq:=DecryptDeck(call)
	// 	client.cli.Caller.sentToClient("decrypt-cards",kq)
	// case "get-key-for-player":
	// 	call:=msg["value"].(map[string]interface{})
	// 	fmt.Println("get-key-for-player callmap:",call)
	// 	client.cli.Caller.GetKeyForPlayer(call)
	// case "get-cards":
	// 	call:=msg["value"].(map[string]interface{})
	// 	client.cli.Caller.GetCards(call)
	// case "get-sign":
	// 	call:=msg["value"].(map[string]interface{})
	// 	client.cli.Caller.GetSign(call)
	// case "verify-sign":
	// 	call:=msg["value"].(map[string]interface{})

	// 	verifyKq:=client.cli.Caller.VerifySign(call)
	// 	if(verifyKq["data"].(bool)==true){
	// 		// client.Caller.GetKeyForPlayer(call)
	// 		client.cli.Caller.sentToClient("verified-sign", verifyKq["address"].(string))

	// 	}else{
	// 		client.cli.Caller.sentToClient("get-key-for-player", "Not Authorised Address")
	// 	}

	default:
		log.Warn("Require call not match: ", msg)
	}
}
