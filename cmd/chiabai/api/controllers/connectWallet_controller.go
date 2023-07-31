package controllers

// import (

// 	// "go.mongodb.org/mongo-driver/mongo"
// 	"fmt"

// 	"github.com/ethereum/go-ethereum/common"
// 	c_config "github.com/meta-node-blockchain/meta-node/cmd/chiabai/config"
// 	c_network "github.com/meta-node-blockchain/meta-node/cmd/chiabai/network"
// 	"github.com/meta-node-blockchain/meta-node/cmd/client/pkg/controllers"
// 	"github.com/meta-node-blockchain/meta-node/pkg/bls"
// 	p_common "github.com/meta-node-blockchain/meta-node/pkg/common"
// 	"github.com/meta-node-blockchain/meta-node/pkg/logger"
// 	"github.com/meta-node-blockchain/meta-node/pkg/network"
// 	"github.com/meta-node-blockchain/meta-node/types"
// )

// type CallData struct {
// 	server *Server
// 	// client *Client
// 	cli *Cli
// }

// func (cli *Cli) ConnectWallet(callMap map[string]interface{}) map[string]interface{} {
// 	result := make(map[string]interface{})

// 	address, _ := callMap["address"].(string)
// 	privatekey, _ := callMap["priKey"].(string)

// 	call := map[string]interface{}{
// 		"address": address,
// 		"priKey":  privatekey,
// 	}
// 	kq := cli.ConnectSocket(call)

// 	result = (map[string]interface{}{
// 		"success": true,
// 		"message": kq,
// 	})
// 	return result
// }

// var (
// 	connectionTypesForClient = []string{p_common.NODE_CONNECTION_TYPE}
// )

// func (cli *Cli) ConnectSocket(walletKey map[string]interface{}) map[string]interface{} {
// 	fmt.Println("ConnectSocket")
// 	config, err := c_config.LoadConfig(c_config.CONFIG_FILE_PATH)
// 	if err != nil {
// 		logger.Error(fmt.Sprintf("error when loading config %v", err))
// 		panic(fmt.Sprintf("error when loading config %v", err))
// 	}
// 	cConfig := config.(*c_config.Config)
// 	result := make(map[string]interface{})
// 	// connect to parent
// 	connectionsManager := network.NewConnectionsManager(connectionTypesForClient)

// 	parentConn := network.NewConnection(
// 		common.HexToAddress(cConfig.ParentAddress),
// 		cConfig.ParentConnectionType,
// 		cConfig.ParentConnectionAddress,
// 	)
// 	accountStateChan := make(chan types.AccountState)
// 	chData := make(chan interface{})
// 	handler := c_network.NewHandler(accountStateChan, chData)

// 	err = parentConn.Connect()
// 	if err != nil {
// 		logger.Error(fmt.Sprintf("error when connect to parent %v", err))
// 	} else {
// 		// init connection
// 		connectionsManager.AddParentConnection(parentConn)

// 		if walletKey["priKey"] == nil {
// 			logger.Error(fmt.Sprintf("error when GetWalletKeyFromAddress %", err))
// 			panic(fmt.Sprintf("error when GetWalletKeyFromAddress %v", err))
// 		} else {
// 			priKey := common.FromHex(walletKey["priKey"].(string))
// 			keyPair := bls.NewKeyPair(priKey)

// 			logger.Info("Running with key pair: " + "\n" + keyPair.String())
// 			messageSender := network.NewMessageSender(keyPair, config.GetVersion())
// 			tcpServer := network.NewSockerServer(config, keyPair, connectionsManager, handler)
// 			tcpServer.OnConnect(parentConn)

// 			go tcpServer.HandleConnection(parentConn)

// 			// init controller
// 			transactionCtl := controllers.NewTransactionController(keyPair, messageSender, connectionsManager)
// 			// init and start client
// 			// fmt.Println("addressString:",addressString)
// 			addressString := walletKey["address"].(string)
// 			cli.keyPairMap[addressString] = keyPair
// 			cli.messageSenderMap[addressString] = messageSender
// 			cli.transactionControllerMap[addressString] = transactionCtl
// 			cli.tcpServerMap[addressString] = tcpServer
// 			cli.accountStateChan = accountStateChan

// 		}

// 		cli.connectionsManager = connectionsManager
// 		cli.config = cConfig
// 	}
// 	fmt.Println("init connection")
// 	return result
// }
