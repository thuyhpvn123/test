package controllers

// import (
// 	"sync"

// 	// log "github.com/sirupsen/logrus"
// 	"github.com/meta-node-blockchain/meta-node/cmd/chiabai/config"
// 	controller_client "github.com/meta-node-blockchain/meta-node/cmd/client/pkg/controllers"
// 	"github.com/meta-node-blockchain/meta-node/pkg/bls"
// 	"github.com/meta-node-blockchain/meta-node/types/network"
// 	"github.com/meta-node-blockchain/meta-node/types"
// )

import (
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/meta-node-blockchain/meta-node/cmd/client/command"
	"github.com/meta-node-blockchain/meta-node/cmd/client/pkg/client_context"
	c_config "github.com/meta-node-blockchain/meta-node/cmd/client/pkg/config"
	"github.com/meta-node-blockchain/meta-node/cmd/client/pkg/controllers"
	c_network "github.com/meta-node-blockchain/meta-node/cmd/client/pkg/network"
	client_types "github.com/meta-node-blockchain/meta-node/cmd/client/types"
	log "github.com/sirupsen/logrus"

	"github.com/meta-node-blockchain/meta-node/pkg/bls"
	p_common "github.com/meta-node-blockchain/meta-node/pkg/common"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	"github.com/meta-node-blockchain/meta-node/pkg/network"
	p_network "github.com/meta-node-blockchain/meta-node/pkg/network"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
	"github.com/meta-node-blockchain/meta-node/cmd/chiabai/core"

)
// var cli *Cli
var defaultRelatedAddress [][]byte

var (
	connectionTypesForClient = []string{p_common.NODE_CONNECTION_TYPE}
)
type Cli struct {
	clientContext *client_context.ClientContext
	mu                    sync.Mutex
	accountStateChan      chan types.AccountState
	// receiptChan          chan types.Receipt
	chData chan interface{}
	transactionController client_types.TransactionController

	server *Server
	sendChan                 chan Message1
}
// var client = Client{}

func NewCli(
	config *c_config.ClientConfig,
	server *Server,
) (*Cli, error) {
	clientContext := &client_context.ClientContext{
		Config: config,
	}
	cli := Cli{
		clientContext:    clientContext,
		accountStateChan: make(chan types.AccountState, 1),
		// receiptChan:      make(chan types.Receipt, 1),
		chData:      make(chan interface{}, 1),
		server: server,
		sendChan: 		make(chan Message1,1),
	}

	clientContext.KeyPair = bls.NewKeyPair(config.PrivateKey())
	clientContext.MessageSender = p_network.NewMessageSender(clientContext.KeyPair, config.Version())
	clientContext.ConnectionsManager = network.NewConnectionsManager(connectionTypesForClient)
	parentConn := network.NewConnection(
		common.HexToAddress(config.ParentAddress),
		config.ParentConnectionType,
		config.ParentConnectionAddress,
	)
	clientContext.Handler = c_network.NewHandler(
		cli.accountStateChan,
		// cli.receiptChan,
		cli.chData,
	)
	clientContext.SocketServer = network.NewSockerServer(
		config,
		clientContext.KeyPair,
		clientContext.ConnectionsManager,
		clientContext.Handler,
	)
	err := parentConn.Connect()
	if err != nil {
		logger.Error(fmt.Sprintf("error when connect to parent %v", err))
		return nil, err
	} else {
		// init connection
		clientContext.ConnectionsManager.AddParentConnection(parentConn)
		clientContext.SocketServer.OnConnect(parentConn)
		go clientContext.SocketServer.HandleConnection(parentConn)
	}
	cli.transactionController = controllers.NewTransactionController(
		clientContext,
	)
	// _, err = cli.SendTransaction(
	// 	common.HexToAddress("08db9fae6755cef00e98d2613686024de9bfed52"),
	// 	uint256.NewInt(10),
	// 	pb.ACTION_EMPTY,
	// 	[]byte{},
	// 	nil,
	// 	100000,
	// 	1000000000,
	// 	0,		
	// )

	// if err != nil {
	// 	fmt.Println("hahahahha:", err)
	// 	log.Warn(err)
	// } else {
	// 	logger.Info("Done send transaction from ")
	// }
	return &cli, nil
}
func (cli *Cli) TryCall(callMap map[string]interface{}) interface{} {
	fmt.Println("hello")
	i := 0
	var result interface{}
	result = "TimeOut"

	for {
		if i >= 3 {
			break
		}
		if i != 0 {
			time.Sleep(time.Second)
		}
		result = cli.call(callMap)

		if result != "TimeOut" {
			log.Info("Success time - ", i)
			log.Info(" - Result: ", result)
			return result
		}
		i++
	}

	return result
}
func (cli *Cli) call(callMap map[string]interface{}) interface{} {
	fmt.Println("hello1111111111")

	fromAddress := "45c75cfb8e20a8631c134555fa5d61fcf3e602f2"
	toAddressStr:= core.Contracts[0].Address
	toAddress := common.HexToAddress(toAddressStr)
	action  :=pb.ACTION_CALL_SMART_CONTRACT
	// action  :=pb.ACTION_EMPTY
	relatedAddress := cli.EnterRelatedAddress(callMap)
	hexAmount, _ := callMap["amount"].(string)
	if hexAmount == "" {
		hexAmount = "0"
	}
	amount := uint256.NewInt(0).SetBytes(common.FromHex(hexAmount))
	var maxGas uint64
	maxGaskq, ok := callMap["gas"].(float64)
	if !ok {
		maxGas = 2000000
	} else {
		maxGas = uint64(maxGaskq)

	}

	var maxGasPriceGwei uint64
	maxGasPriceGweikq, ok := callMap["gasPrice"].(float64)
	if !ok {
		maxGasPriceGwei = 10
	} else {
		maxGasPriceGwei = uint64(maxGasPriceGweikq)

	}
	maxGasPrice := 100000000 * maxGasPriceGwei

	var maxTimeUse uint64
	maxTimeUsekq, ok := callMap["timeUse"].(float64)
	if !ok {
		maxTimeUse = 1000
	} else {
		maxTimeUse = uint64(maxTimeUsekq)
	}
	
	input,ok :=callMap["input"].(string)
	if !ok{
		input=""
	}
	data:=common.FromHex(input)
	fmt.Println("input là:", input)
	 err := cli.SendTransaction(
		toAddress ,
		amount ,
		action ,
		data ,
		relatedAddress ,
		maxGas ,
		maxGasPrice ,
		maxTimeUse ,	
		
	)
	
	
	if err != nil {
		fmt.Println("hahahahha:", err)
		log.Warn(err)
	} else {
		logger.Info("Done send transaction from " + fromAddress)
	}
	// if hashed == nil {
	// 	fmt.Println("hashed==nil")
	// 	return "TimeOut1"
	// }
	for {

		select {
		case receiver := <-cli.chData:
			// log.Info("Hash on server", common.BytesToHash(hash.([]byte)))
			// log.Info("Hash from chain", (receiver).(network.Receipt).Hash)
			// if (receiver).(network.Receipt).Hash != common.BytesToHash(hash.([]byte)) {
			// 	continue
			// }
			// return (receiver).Value
			return (receiver).(c_network.Receipt1).Value
		case <-time.After(5 * time.Second):
			return "TimeOut"
		}
	}
		
}



func (cli *Cli) SendTransaction(
	toAddress common.Address,
	amount *uint256.Int,
	action pb.ACTION,
	data []byte,
	relatedAddress [][]byte,
	maxGas uint64,
	maxGasPrice uint64,
	maxTimeUse uint64,
) ( error) {
	fmt.Println("hello222222222222")

	cli.mu.Lock()
	defer cli.mu.Unlock()
	// get account state
	parentConn := cli.clientContext.ConnectionsManager.ParentConnection()
	cli.clientContext.MessageSender.SendBytes(
		parentConn,
		command.GetAccountState,
		cli.clientContext.KeyPair.Address().Bytes(),
		p_common.Sign{},
	)
	fmt.Println("hello33333333")

	as := <-cli.accountStateChan
	lastHash := as.LastHash()
	pendingBalance := as.PendingBalance()

	// bRelatedAddresses := make([][]byte, len(relatedAddress))
	// for i, v := range relatedAddress {
	// 	bRelatedAddresses[i] = v.Bytes()
	// }
	transaction, err := cli.transactionController.SendTransaction(
		lastHash,
		toAddress,
		pendingBalance,
		amount,
		maxGas,
		maxGasPrice,
		maxTimeUse,
		action,
		data,
		relatedAddress,
	)
	
	logger.Info("Sending transaction", transaction)
	// if err != nil {
	// 	return nil, err
	// }
	if err != nil {
		logger.Error(err)
	} else {
		logger.Info("Done send transaction from " )
	}

	// receipt := <-cli.receiptChan
	// fmt.Println("receipt:",receipt)
	return  nil
}

func (cli *Cli) AccountState(address common.Address) (types.AccountState, error) {
	cli.mu.Lock()
	defer cli.mu.Unlock()
	// get account state
	parentConn := cli.clientContext.ConnectionsManager.ParentConnection()
	cli.clientContext.MessageSender.SendBytes(
		parentConn,
		command.GetAccountState,
		address.Bytes(),
		p_common.Sign{},
	)
	as := <-cli.accountStateChan
	return as, nil
}

func (cli *Cli) Subcribe(storageHost string, smartContractAddress common.Address) (chan types.EventLogs, error) {
	storageConnection := network.NewConnection(common.Address{}, p_common.STORAGE_CONNECTION_TYPE, storageHost)
	err := storageConnection.Connect()
	if err != nil {
		logger.Error("Unable to connect to storage", err)
		return nil, fmt.Errorf("unable to connect to storage")
	}
	go cli.clientContext.SocketServer.HandleConnection(storageConnection)

	err = cli.clientContext.MessageSender.SendBytes(storageConnection, command.SubscribeToAddress, smartContractAddress.Bytes(), p_common.Sign{})
	if err != nil {
		return nil, fmt.Errorf("unable to send subscribe")
	}
	evenLogsChan := make(chan types.EventLogs)
	cli.clientContext.Handler.(*c_network.Handler).SetEventLogsChan(evenLogsChan)
	return evenLogsChan, nil
}
func (cli *Cli) EnterRelatedAddress(call map[string]interface{}) [][]byte {
	var arrmap []map[string]interface{}
	arr, _ := call["relatedAddresses"].([]interface{})
	if call["relatedAddresses"] == nil || len(arr) == 0 {
		return defaultRelatedAddress
	} else {
		for _, v := range arr {
			arrmap = append(arrmap, v.(map[string]interface{}))
		}

		var relatedAddStr []string

		for _, v := range arrmap {
			relatedAddStr = append(relatedAddStr, v["address"].(string))
		}
		var relatedAddress [][]byte

		// temp := strings.Split(relatedAddStr, ",")
		logger.Info("Temp Related Address")
		for _, addr := range relatedAddStr {
			addressHex := common.HexToAddress(addr)
			logger.Info(addressHex)
			relatedAddress = append(relatedAddress, addressHex.Bytes())
		}
		defaultRelatedAddress = append(defaultRelatedAddress, relatedAddress...)
		return relatedAddress

	}
}
func (cli *Cli) sentToClient(command string, data interface{}) {
	cli.sendChan <- Message1{command, data}
}

