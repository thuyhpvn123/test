package core

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	// . "github.com/ethereum/go-ethereum/accounts/abi"
)

type Account struct {
	Address string
	Private string
}

var PORT int
var Contracts = [...]Contract{
	{Name: "chiabai", Address: "4acca9e4c68560e16ebb5cbdaed79b3029540b77"},
}

// var accounts = [...]Account{
// 	{
// 		Address: "45c75cfb8e20a8631c134555fa5d61fcf3e602f2",
// 		Private: "36e1aa979f98c7154fb2491491ec044ccac099651209ccfbe2561746dbe29ebb",
// 	},

// }
//f8eaba3eb679f6defbe78ce8dd5229ec3622f2a7
func GetPORT() int {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	// log.Info("PORT: ", os.Getenv("PORT"))
	PORT, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
	return PORT
}

type Contract struct {
	Name    string
	Address string
}
