package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	keystore "github.com/ethereum/go-ethereum/accounts/keystore"
	toml "github.com/naoina/toml"
)

//Configuration data
type config struct {
	TestVariants bool
	PreSeq       []string
	PostSeq      []string
}

//Flags
var walletFile = flag.String("wallet", "wallet.json", "Wallet file")
var passList = flag.String("passList", "passList.txt", "Text file with one passphrase per line")
var confile = flag.String("conf", "conf.toml", "TOML file with some configuration")

//Files
var keyJson []byte    //wallet
var passFile *os.File //PassList
var conf config       //Configuration

//Counter password tested
var passCount int
var elapsed time.Duration

//Key variables
var passphrase string
var address string

//Init function
func init() {

	fmt.Print("\n--------- PASSPHRASE RECOVER TOOL ---------\n\n")
	log.Print("Initializing tool...\n\n")

	//Parse flags
	flag.Parse()

	//Open Configuration file
	confData, err := ioutil.ReadFile(*confile)
	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
		return
	}
	log.Printf("Reading conf: %s", *confile)

	//Read configuration data
	if err = toml.Unmarshal(confData, &conf); err != nil {
		log.Printf("Error Configuration file: %s", err.Error())
	}

	log.Printf("Configurations preSeq: %v\tpostSeq: %v", conf.PreSeq, conf.PostSeq)

	log.Printf("Opening wallet: %s\n", *walletFile)
	//Open and Read Json from walletFile
	keyJson, err = ioutil.ReadFile(*walletFile)
	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
		return
	}

	log.Printf("Opening passList: %s\n\n", *passList)
	//Open passList.txt
	if passFile, err = os.Open(*passList); err != nil {
		log.Fatal(err)
	}

	//Set counter to 0
	passCount = 0
}

//testPass check if pass is correct.
//Return true if pass is correct, false if pass is invalid
func testPass(pass string) bool {
	if len(pass) <= 7 {
		return false
	} else {
		passCount++
		if test, err := keystore.DecryptKey(keyJson, pass); err != nil {
			//fmt.Printf("Errore: %v\n", err)
			return false
		} else {
			passphrase = pass
			address = test.Address.String()
			return true
		}
	}
}

//testPassVariants check all prefix and suffix variations
//of pass. Remember to set preSeq and postSeq in your conf.toml
func testPassVariants(pass string) bool {
	var passTemp string

	//Test pass with all preSeq
	for _, seqPre := range conf.PreSeq {
		passTemp = seqPre + pass
		//fmt.Println(passTemp)
		if testPass(passTemp) {
			return true
		}
	}

	//Test pass with all postSeq
	for _, seqPost := range conf.PostSeq {
		passTemp = pass + seqPost
		//fmt.Println(passTemp)
		if testPass(passTemp) {
			return true
		}
	}

	//Test pass with all preSeq - postSeq combinations
	for _, seqPre := range conf.PreSeq {
		for _, seqPost := range conf.PostSeq {
			passTemp = seqPre + pass + seqPost
			//fmt.Println(passTemp)
			if testPass(passTemp) {
				return true
			}
		}
	}
	return false
}

//end print some outputs
//success true if passphrase FOUND
//success false if not FOUND
func end(success bool) {
	//Some output
	log.Printf("Tested %v passphrases in %v seconds.", passCount, elapsed.Seconds())

	if success {
		log.Print("PASSPHRASE FOUND")
		fmt.Printf("\nWallet Address: %s\n\n------------ PASSPHRASE ------------\n\n%s\n\n------------------------------------\n\n", address, passphrase)
		fmt.Print("Please make a donation to developer:\nETH: 0x2feD76d5abE6c001D259eC769c28f6068E0166CB\nBTC: 1HTpxVw6KkDakhjqL3bgkYtM7Gsxxzmjw5\n\n")
	} else {
		log.Print("Sorry. Passphrase not found!")
	}

}

//Main function
func main() {

	defer passFile.Close()
	// Create a scanner to read passList line by line
	scanner := bufio.NewScanner(passFile)

	//Start time for elapsed
	startT := time.Now()

	//Read line by line passList
	if conf.TestVariants { //With Password Variants
		log.Print("Searching Passphrase with variants... Please wait\n\n")
		for scanner.Scan() {
			pass := scanner.Text()
			if testPass(pass) {
				elapsed = time.Since(startT)
				end(true)
				return
			} else {
				if testPassVariants(pass) {
					elapsed = time.Since(startT)
					end(true)
					return
				}
			}
		}
	} else { //Without Password Variants
		log.Print("Searching Passphrase without variants... Please wait\n\n")
		for scanner.Scan() {
			pass := scanner.Text()
			if testPass(pass) {
				elapsed = time.Since(startT)
				end(true)
				return
			}
		}
	}

	//if passphrase not found
	elapsed = time.Since(startT)

	// check scanner errors
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	end(false)
	return
}
