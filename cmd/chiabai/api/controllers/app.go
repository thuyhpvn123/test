package controllers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb"

	random "math/rand"
)
var IV = []byte("1234567812345678")
const deckSize = 52
var levelDb *leveldb.DB

// func (caller *CallData)test() {
// 	// Generate private keys for four players
// 	playerKeys,err := caller.GeneratePlayerKeys(4)
// 	if err != nil {
// 		fmt.Println("Error generating keys:", err)
// 		return
// 	}

// 	fmt.Println("Generated AES keys:")
// 	for i, key := range playerKeys {
// 		fmt.Printf("Key %d: %s\n", i+1, key)
// 	}
// 	// Create and shuffle a deck of 52 cards
// 	deck := caller.CreateDeck()
// 	// fmt.Println("deck:",deck)
// 	caller.ShuffleDeck(deck)
// 	// deck:=[]string{"1"}
// 	// Encrypt the deck using the private keys of each player
// 	encryptedDeck := caller.EncryptDeck(deck, playerKeys)

// 	// Print the encrypted deck
// 	fmt.Println("Encrypted Deck:",encryptedDeck)
// 	decryptedDeck := caller.DecryptDeck(encryptedDeck,playerKeys)
// 	fmt.Println("Decrypted Deck:",decryptedDeck)

// }

// func (caller *CallData)GeneratePlayerKeys(numPlayers int) ([]string, error) {
func GeneratePlayerKeys(call map[string]interface{})map[string]interface{} {
	result:=make(map[string]interface{})
	// numPlayers := int(call["numPlayers"].(float64))
	// keys := make([]string, numPlayers)
	roomid,_ := call["roomid"].(string)
	player,_ := call["player"].(string)

	keysArr:=make([]string, 52)

	for i := 0; i < 52; i++ {
		key := make([]byte, 16) // AES-256 requires a 32-byte key
		_, err := rand.Read(key)

		if err != nil {
			result=(map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})	
			return result
		}

		keysArr[i] = hex.EncodeToString(key)
	}
	keysArray=keysArr
	leveldb, err := leveldb.OpenFile("./db/device_info", nil)
	if err != nil {
		panic(err)
	}
	barray,_:=json.Marshal(keysArr)
    var strbarray string
    strbarray=string(barray)
	callmap3 :=map[string]interface{}{
		"key":roomid+player,
		"data":strbarray,

	}
	result3:= WriteValueStorage(callmap3,leveldb)
	fmt.Println("write deviceKey to storage success:",result3["success"])

	result=(map[string]interface{}{
		"success": true,
		"message": keysArr,
	})	

	return result
}

// createDeck creates a standard deck of 52 cards
// func (caller *CallData)CreateDeck() []string {
func CreateDeck() map[string]interface{} {
	result:=make(map[string]interface{})
	ranks := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

	suits := []string{"S", "H", "D", "C"}

	deck := make([]string, deckSize)

	for i, suit := range suits {
		for j, rank := range ranks {
			card := rank + suit
			deck[i*len(ranks)+j] = card
		}
	}
	result=(map[string]interface{}{
		// "success": true,
		"deck": deck,
	})	
	fmt.Println("Create deck success")
	return result
}

// shuffleDeck shuffles the given deck of cards
// func (caller *CallData)ShuffleDeck(deck []string) {
func ShuffleDeck(call map[string]interface{})[]string {
	// result:=make(map[string]interface{})
	deck := call["deck"].([]string)
	random.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	// result=(map[string]interface{}{
	// 	// "success": true,
	// 	"deck": deck,
	// })	
	fmt.Println("Shuffle deck success:",deck)
	return deck
}
// encryptDeck encrypts each card in the deck using the private keys of the players
// func (caller *CallData)EncryptDeck(deck []string, playerKeys []string) []string {
func EncryptDeck(deck []string, arrmap []string) []string {
	fmt.Println("begin encrypt deck ")
	// result:=make(map[string]interface{})
	// deck := call["deck"].([]string)
	// playerKeys,_ := call["playerKeys"].([]interface{})
	// var arrmap []string
	// for _,v := range playerKeys{
	// 	arrmap= append(arrmap,v.(string)) //arr có 52  key
	// }
	fmt.Println("arrmap:",arrmap)
	encryptedDeck := make([]string, len(deck))
	// var count=0
	// for _, card := range deck {
	// 	for _, priKey := range arrmap{
	// 		fmt.Println("begin encrypt deck2222222222 ")
	// 		fmt.Println("card: ",card)
	// 		fmt.Println("priKey: ",priKey)
			// encryptedCard:= encryption(card,priKey)
			// card =encryptedCard
			// encryptedDeck[i] = card
			// fmt.Println("card: ",card)
	// 		fmt.Println("count: ",count)
	// 		count++
			
	// 	}		
	// }
		for i, priKey := range arrmap{
			card:=deck[i]
			encryptedCard:= encryption(card,priKey)
			card =encryptedCard
			encryptedDeck[i] = card

		}
		deckKq=encryptedDeck

	// result=(map[string]interface{}{
	// 	"success": true,
	// 	"message": encryptedDeck,
	// })	
	// fmt.Println("Encrypted deck success")
	// return result
	return encryptedDeck
}
func createCipher(key string) cipher.Block {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Fatalf("Failed to create the AES cipher: %s", err)
	}
	return c
}
func encryption(plainText string,key string) string{
	bytes := []byte(plainText)
	blockCipher := createCipher(key)
	stream := cipher.NewCTR(blockCipher, IV)
	// Buffer for storing decrypted data
	encryptedData := make([]byte, len(bytes))
	stream.XORKeyStream(encryptedData, bytes)
	result:=base64.StdEncoding.EncodeToString(encryptedData)
	return result
}
// func (caller *CallData)DecryptDeck(encrytedDeck []string, playerKeys []string) []string {
func DecryptDeck(call map[string]interface{}) map[string]interface{}  {
	fmt.Println("-----------")
	result:=make(map[string]interface{})
	encrytedDeck := call["encrytedDeck"].([]interface{})
	playerKeys := call["playerKeys"].([]interface{})
	decryptedDeck := make([]string, len(encrytedDeck))
	// for i, encryptedcard := range encrytedDeck {
	// 	for j:=len(playerKeys)-1;j>=0;j--{
	// 			decryptedBlockBytes:= decryption(encryptedcard.(string),playerKeys[j].(string))
	// 			encryptedcard =string(decryptedBlockBytes)
	// 			decryptedDeck[i] = encryptedcard.(string)
	// 	}
	// }
	for i, encryptedcard := range encrytedDeck {
		decryptedBlockBytes:= decryption(encryptedcard.(string),playerKeys[i].(string))
		encryptedcard =string(decryptedBlockBytes)
		decryptedDeck[i] = encryptedcard.(string)

	}
	result=(map[string]interface{}{
		"success": true,
		"message": decryptedDeck,
	})	
	return result
} 

func decryption(encrypted string,key string) []byte {
	bytes,err:=base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		fmt.Println("Error decoding the encrypted string:", err)
		return nil
	}
	blockCipher := createCipher(key)
	stream := cipher.NewCTR(blockCipher, IV)
	stream.XORKeyStream(bytes, bytes)
	return bytes
}
