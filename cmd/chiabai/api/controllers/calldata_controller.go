package controllers

// import (
// 	"encoding/hex"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"math/big"
// 	"os"
// 	"reflect"
// 	"strconv"
// 	"strings"
// 	"sync"
// 	"time"

// 	log "github.com/sirupsen/logrus"

// 	// "github.com/meta-node-blockchain/meta-node/cmd/chiabai/utils"
// 	// "github.com/meta-node-blockchain/meta-node/cmd/chiabai/core"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/crypto"
// 	"github.com/holiman/uint256"
// 	"github.com/meta-node-blockchain/meta-node/cmd/client/command"
// 	"github.com/meta-node-blockchain/meta-node/pkg/bls"
// 	cm "github.com/meta-node-blockchain/meta-node/pkg/common"
// 	"github.com/meta-node-blockchain/meta-node/types"

// 	. "github.com/ethereum/go-ethereum/accounts/abi"
// 	"github.com/meta-node-blockchain/meta-node/pkg/network"
// 	"github.com/meta-node-blockchain/meta-node/pkg/logger"
// 	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
// 	t "github.com/meta-node-blockchain/meta-node/pkg/transaction"
// )

// type ActionListenerCallback map[string]interface{}

// var defaultRelatedAddress [][]byte
// var (
// 	ErrorGetAccountStateTimedOut = errors.New("get account state timed out")
// 	ErrorInvalidAction           = errors.New("invalid action")
// )

// type WalletKey struct {
// 	PriKey []byte
// 	PubKey []byte
// }

// func (cli *Cli) TryCall(callMap map[string]interface{}) interface{} {
// 	i := 0
// 	var result interface{}
// 	result = "TimeOut"

// 	for {
// 		if i >= 3 {
// 			break
// 		}
// 		if i != 0 {
// 			time.Sleep(time.Second)
// 		}
// 		result = cli.call(callMap)

// 		if result != "TimeOut" {
// 			log.Info("Success time - ", i)
// 			log.Info(" - Result: ", result)
// 			return result
// 		}
// 		i++
// 	}

// 	return result
// }

// func (cli *Cli) call(callMap map[string]interface{}) interface{} {

// 	fromAddress := "45c75cfb8e20a8631c134555fa5d61fcf3e602f2"
// 	_, err := cli.SendTransaction2(callMap)
// 	if err != nil {
// 		fmt.Println("hahahahha:", err)
// 		log.Warn(err)
// 	} else {
// 		logger.Info("Done send transaction from " + fromAddress)
// 	}
// 	// if hashed == nil {
// 	// 	fmt.Println("hashed==nil")
// 	// 	return "TimeOut1"
// 	// }
// 	for {

// 		select {
// 		case receiver := <-cli.tcpServerMap[fromAddress].GetHandler():
// 			// log.Info("Hash on server", common.BytesToHash(hash.([]byte)))
// 			// log.Info("Hash from chain", (receiver).(network.Receipt).Hash)
// 			// if (receiver).(network.Receipt).Hash != common.BytesToHash(hash.([]byte)) {
// 			// 	continue
// 			// }
// 			return (receiver).(network.Receipt1).Value
// 		case <-time.After(5 * time.Second):
// 			return "TimeOut"
// 		}
// 	}

// }

// var transferMu sync.Mutex

// //chỉ owner gọi
// func (cli *Cli) SendTransaction2(call map[string]interface{}) (t.Transaction, error) {
// 	// fmt.Println("call là:", call)
// 	inputStr, _ := call["input"].(string)
// 	relatedAddress := cli.EnterRelatedAddress(call)
// 	// fromAddress:= cli.config.Address
// 	fromAddress := "45c75cfb8e20a8631c134555fa5d61fcf3e602f2"
// 	// toAddressStr:= core.Contracts[0].Address
// 	toAddressStr := "fdd11471417109d88c48030e579f3523e485f6fa"

// 	toAddress := common.HexToAddress(toAddressStr)
// 	hexAmount, _ := call["amount"].(string)
// 	if hexAmount == "" {
// 		hexAmount = "0"
// 	}
// 	amount := uint256.NewInt(0).SetBytes(common.FromHex(hexAmount))
// 	var maxGas uint64
// 	maxGaskq, ok := call["gas"].(float64)
// 	if !ok {
// 		maxGas = 2000000
// 	} else {
// 		maxGas = uint64(maxGaskq)

// 	}

// 	var maxGasPriceGwei uint64
// 	maxGasPriceGweikq, ok := call["gasPrice"].(float64)
// 	if !ok {
// 		maxGasPriceGwei = 10
// 	} else {
// 		maxGasPriceGwei = uint64(maxGasPriceGweikq)

// 	}
// 	maxGasPrice := 100000000 * maxGasPriceGwei

// 	var maxTimeUse uint64
// 	maxTimeUsekq, ok := call["timeUse"].(float64)
// 	if !ok {
// 		maxTimeUse = 1000
// 	} else {
// 		maxTimeUse = uint64(maxTimeUsekq)
// 	}
// 	var action pb.ACTION
// 	// action = pb.ACTION_CALL_SMART_CONTRACT
// 	action = 0

// 	sign := cli.GetSignGetAccountState()

// 	as, err := cli.GetAccountState(fromAddress, sign)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var data []byte
// 	if len(inputStr) > 0 {
// 		// data = common.FromHex(inputStr)
// 		data = common.FromHex(inputStr)

// 	} else {

// 		data, err = cli.GetDataForCallSmartContract(call)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// 	transaction, err := cli.transactionControllerMap[fromAddress].SendTransaction(
// 		as.LastHash(),
// 		toAddress,
// 		as.PendingBalance(),
// 		amount,
// 		maxGas,
// 		maxGasPrice,
// 		maxTimeUse,
// 		action,
// 		data,
// 		relatedAddress,
// 	)
// 	logger.Debug("Sending transaction", transaction)
// 	if err != nil {
// 		logger.Warn(err)
// 	}
// 	fmt.Printf("Send transaction %v", transaction)

// 	return transaction, err
// }

// // func (caller *CallData) SendTransaction1(call map[string]interface{}) error {
// // 	fmt.Println("call là:", call)
// // 	inputStr, _ := call["input"].(string)
// // 	relatedAddress := caller.EnterRelatedAddress(call)
// // 	fromAddress, _ := call["from"].(string)
// // 	toAddressStr, _ := call["to"].(string)
// // 	toAddress := common.HexToAddress(toAddressStr)
// // 	hexAmount, _ := call["amount"].(string)
// // 	if hexAmount == "" {
// // 		hexAmount = "0"
// // 	}
// // 	amount := uint256.NewInt(0).SetBytes(common.FromHex(hexAmount))
// // 	var maxGas uint64
// // 	maxGaskq, ok := call["gas"].(float64)
// // 	if !ok {
// // 		maxGas = 500000
// // 	}
// // 	maxGas = uint64(maxGaskq)

// // 	var maxGasPriceGwei uint64
// // 	maxGasPriceGweikq, ok := call["gasPrice"].(float64)
// // 	if !ok {
// // 		maxGasPriceGwei = 10
// // 	}
// // 	maxGasPriceGwei = uint64(maxGasPriceGweikq)
// // 	maxGasPrice := 1000000000 * maxGasPriceGwei

// // 	var maxTimeUse uint64
// // 	maxTimeUsekq, ok := call["timeUse"].(float64)
// // 	if !ok {
// // 		maxTimeUse = 60000
// // 	}
// // 	maxTimeUse = uint64(maxTimeUsekq)
// // 	var action pb.ACTION
// // 	action = pb.ACTION_CALL_SMART_CONTRACT

// // 	sign := caller.GetSignGetAccountState(call)

// // 	as, err := caller.GetAccountState(fromAddress, sign)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	var data []byte
// // 	if len(inputStr) > 0 {
// // 		data = common.FromHex(inputStr)
// // 	} else {

// // 		data, err = caller.GetDataForCallSmartContract(call)
// // 		if err != nil {
// // 			panic(err)
// // 		}
// // 	}
// // 	transaction, err := caller.client.transactionControllerMap[fromAddress].SendTransaction(
// // 		as.GetLastHash(),
// // 		toAddress,
// // 		as.GetPendingBalance(),
// // 		amount,
// // 		maxGas,
// // 		maxGasPrice,
// // 		maxTimeUse,
// // 		action,
// // 		data,
// // 		relatedAddress,
// // 	)
// // 	logger.Debug("Sending transaction", transaction)
// // 	if err != nil {
// // 		logger.Warn(err)
// // 	}
// // 	fmt.Printf("Send transaction %v", transaction)

// // 	return err
// // }

// func (cli *Cli) GetSignGetAccountState() cm.Sign {
// 	hash := crypto.Keccak256(common.FromHex("36e1aa979f98c7154fb2491491ec044ccac099651209ccfbe2561746dbe29ebb"))

// 	// if call["priKey"] == nil {
// 	// 	logger.Error(fmt.Sprintf("error when get wallet key "))
// 	// }
// 	privateKey := cli.config.PrivateKey
// 	keyPair := bls.NewKeyPair(common.FromHex(privateKey))
// 	prikey := keyPair.PrivateKey()
// 	sign := bls.Sign(prikey, hash)
// 	return sign
// }

// func (cli *Cli) GetAccountState(address string, sign cm.Sign) (types.AccountState, error) {
// 	parentConn := cli.connectionsManager.ParentConnection()
// 	cli.messageSenderMap[address].SendBytes(parentConn, command.GetAccountState, common.FromHex(address), sign)

// 	select {
// 	case accountState := <-cli.accountStateChan:
// 		return accountState, nil
// 	case <-time.After(5 * time.Second):
// 		return nil, ErrorGetAccountStateTimedOut
// 	}

// }

// // func (cli *Cli) GetWalletInfo(call map[string]interface{}) {

// // 	sign := cli.GetSignGetAccountState(call)
// // 	as, err := cli.GetAccountState(call["from"].(string), sign)
// // 	if err != nil {
// // 		logger.Error(fmt.Sprintf("error when GetAccountState %", err))
// // 		panic(fmt.Sprintf("error when GetAccountState %v", err))
// // 	}
// // 	result := map[string]interface{}{
// // 		"address":         as.GetAddress(),
// // 		"last_hash":       as.GetLastHash(),
// // 		"balance":         as.GetBalance(),
// // 		"pending_balance": as.GetPendingBalance(),
// // 	}
// // 	fmt.Println("result:", result)
// // }

// func (cli *Cli) sentToClient(command string, data interface{}) {
// 	cli.sendChan <- Message1{command, data}
// 	// sendQueue[cli.client.ws] <- Message{msgType, value}
// }

// func (cli *Cli) EnterRelatedAddress(call map[string]interface{}) [][]byte {
// 	var arrmap []map[string]interface{}
// 	arr, _ := call["relatedAddresses"].([]interface{})
// 	if call["relatedAddresses"] == nil || len(arr) == 0 {
// 		return defaultRelatedAddress
// 	} else {
// 		for _, v := range arr {
// 			arrmap = append(arrmap, v.(map[string]interface{}))
// 		}

// 		var relatedAddStr []string

// 		for _, v := range arrmap {
// 			relatedAddStr = append(relatedAddStr, v["address"].(string))
// 		}
// 		var relatedAddress [][]byte

// 		// temp := strings.Split(relatedAddStr, ",")
// 		logger.Info("Temp Related Address")
// 		for _, addr := range relatedAddStr {
// 			addressHex := common.HexToAddress(addr)
// 			logger.Info(addressHex)
// 			relatedAddress = append(relatedAddress, addressHex.Bytes())
// 		}
// 		defaultRelatedAddress = append(defaultRelatedAddress, relatedAddress...)
// 		return relatedAddress

// 	}
// }
// func (cli *Cli) GetDataForCallSmartContract(call map[string]interface{}) ([]byte, error) {
// 	kq := cli.EncodeAbi(call)
// 	callData := t.NewCallData(kq)
// 	return callData.Marshal()
// }

// func (cli *Cli) EncodeAbi(call map[string]interface{}) []byte {
// 	var inputArray []interface{}
// 	if call["inputArray"] == nil {
// 		inputArray = []interface{}{}
// 	} else {
// 		inputArray, _ = call["inputArray"].([]interface{})
// 	}
// 	functionName, _ := call["function-name"].(string)

// 	reader, err := os.Open("./abi/chiabai.json")
// 	fmt.Println("1111111111111111111")
// 	if err != nil {
// 		log.Fatalf("Error occured while reading %s", "./abi/chiabai.json")
// 	}
// 	abiJson, err := JSON(reader)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var abiTypes []interface{}
// 	for _, item := range inputArray {
// 		itemArr := encodeAbiItem(item)
// 		for _, v := range itemArr {
// 			abiTypes = append(abiTypes, v)
// 		}

// 	}

// 	fmt.Println("kkkkkkkkkk")
// 	fmt.Printf("type: %T", abiTypes)
// 	out, err := abiJson.Pack(functionName, abiTypes[:]...)

// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("out:", hex.EncodeToString(out))
// 	return out
// }

// func encodeAbiItem(item interface{}) []interface{} {
// 	var result []interface{}
// 	var itemMap map[string]interface{}
// 	fmt.Println("222222222222222")

// 	if err := json.Unmarshal([]byte(item.(string)), &itemMap); err != nil {
// 		log.Fatal(err)
// 	}
// 	itemType, _ := itemMap["type"].(string)
// 	fmt.Println("itemType:", itemType)
// 	switch itemType {
// 	case "tuple":
// 		fmt.Println("3333333333")

// 		var value []interface{}
// 		if err := json.Unmarshal([]byte(fmt.Sprintf("%v", itemMap["value"])), &value); err != nil {
// 			log.Fatal(err)
// 		}

// 		var components []interface{}
// 		fmt.Println("444444444444")

// 		if err := json.Unmarshal([]byte(fmt.Sprintf("%v", itemMap["components"])), &components); err != nil {
// 			log.Fatal(err)
// 		}

// 		var abiTypes []interface{}
// 		for i, component := range components {
// 			componentBytes, _ := json.Marshal(component)
// 			componentType, _ := component.(map[string]interface{})["type"].(string)
// 			if componentType == "tuple" || componentType == "tuple[]" {
// 				components[i].(map[string]interface{})["value"] = value[i]
// 				abiTypes = append(abiTypes, encodeAbiItem(componentBytes))
// 			} else {
// 				abiTypes = append(abiTypes, getAbiType(componentType, value[i]))
// 			}
// 		}
// 		result = abiTypes
// 	case "tuple[]":
// 		var value []interface{}
// 		fmt.Println("555555555555")

// 		if err := json.Unmarshal([]byte(fmt.Sprintf("%v", itemMap["value"])), &value); err != nil {
// 			log.Fatal(err)
// 		}

// 		fmt.Println("66666666666666")
// 		var components []interface{}
// 		if err := json.Unmarshal([]byte(fmt.Sprintf("%v", itemMap["components"])), &components); err != nil {
// 			log.Fatal(err)
// 		}

// 		var tuples []interface{}
// 		for _, v := range value {
// 			vArray := v.([]interface{})
// 			var abiTypes []interface{}
// 			for j, component := range components {
// 				componentBytes, _ := json.Marshal(component)
// 				componentType, _ := component.(map[string]interface{})["type"].(string)
// 				components[j].(map[string]interface{})["value"] = vArray[j]
// 				if componentType == "tuple" || componentType == "tuple[]" {
// 					abiTypes = append(abiTypes, encodeAbiItem(componentBytes))
// 				} else {
// 					abiTypes = append(abiTypes, getAbiType(componentType, vArray[j]))
// 				}
// 			}
// 			tuples = append(tuples, abiTypes...)
// 		}
// 		result = tuples
// 	default:
// 		fmt.Println("77777777777777")

// 		value := itemMap["value"]

// 		var arr []interface{}

// 		result1 := getAbiType(itemType, value)
// 		result = append(arr, result1)
// 		fmt.Println("jjjjjjjjjjjjjj")
// 	}
// 	return result
// }
// func getAbiType(dataType string, data interface{}) interface{} {
// 	fmt.Println("888888888888")

// 	if strings.Contains(dataType, "int") {
// 		params := big.NewInt(0)
// 		params, ok := params.SetString(fmt.Sprintf("%v", int64(data.(float64))), 10)

// 		if !ok {
// 			log.Warn("Format big int: error")
// 			return nil
// 		}
// 		return params

// 	} else {
// 		fmt.Println("dataType:", dataType)
// 		switch dataType {

// 		case "string":
// 			fmt.Println("aaaaaaaaaaa")
// 			return data.(string)
// 		case "bool":
// 			return data.(bool)
// 		case "address":
// 			return common.HexToAddress(data.(string))
// 		case "uint8":
// 			intVar, err := strconv.Atoi(data.(string))
// 			if err != nil {
// 				log.Warn("Conver Uint8 fail", err)
// 				return nil
// 			}
// 			return uint8(intVar)
// 		// case "uint", "uint256":
// 		// 	nubmer := big.NewInt(0)
// 		// 	nubmer, ok := nubmer.SetString(data.(string), 10)
// 		// 	if !ok {
// 		// 		log.Warn("Format big int: error")
// 		// 		return nil
// 		// 	}
// 		// 	return nubmer
// 		case "array", "slice":
// 			fmt.Println("999999999")

// 			fmt.Println("array nè")
// 			fmt.Println("data:", data)
// 			rv := reflect.ValueOf(data)
// 			var out []interface{}
// 			for i := 0; i < rv.Len(); i++ {
// 				out = append(out, rv.Index(i).Interface())
// 			}

// 			return out
// 		case "string[]":
// 			fmt.Println("ppppppp")
// 			var out []string
// 			for i := 0; i < len(data.([]interface{})); i++ {
// 				out = append(out, data.([]interface{})[i].(string))
// 			}

// 			return out
// 		case "address[]":
// 			fmt.Println("kkkkkk")
// 			var out []common.Address
// 			for i := 0; i < len(data.([]interface{})); i++ {
// 				out = append(out, common.HexToAddress(data.([]interface{})[i].(string)))
// 			}

// 			return out

// 		default:
// 			fmt.Println("1000000000")
// 			return data
// 		}
// 	}
// }
