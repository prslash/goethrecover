package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/naoina/toml"
)

type Brute struct {
	Active bool
	Lenght int
	Chars  []byte
}

//Configuration data
type Config struct {
	CustomVariants bool
	PreSeq         []string
	PostSeq        []string
	PreBrute       struct {
		Active bool
		Lenght int
		Chars  []string
	}
	PostBrute struct {
		Active bool
		Lenght int
		Chars  []string
	}
}

type combData struct {
	n       []byte
	k       int
	indexes []int
	result  [][]byte
	pre     bool
	post    bool
}

var preComb combData
var postComb combData

//Flags
var walletFile = flag.String("wallet", "wallet.json", "Wallet file")
var passList = flag.String("passList", "passList.txt", "Text file with one passphrase per line")
var confile = flag.String("conf", "conf.toml", "TOML file with some configuration")

//Files
var keyJson []byte    //wallet
var passFile *os.File //PassList

//Configuration
var Conf Config

//Init function
func init() {

	//Parse flags
	flag.Parse()

	fmt.Print("\n--------- PASSPHRASE RECOVER TOOL ---------\n\n")
	log.Print("Initializing tool...\n\n")

	//Open Configuration file
	confData, err := ioutil.ReadFile(*confile)
	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
		return
	}
	log.Printf("Reading conf: %s", *confile)

	//Read configuration data .toml
	if err = toml.Unmarshal(confData, &Conf); err != nil {
		log.Printf("Error Configuration file: %s", err.Error())
	}

	log.Printf("Configurations %v", Conf)

	if Conf.PreBrute.Active {
		log.Println("Performing combinations [preBrute]...")
		preComb.pre = true
		for _, c := range Conf.PreBrute.Chars {
			preComb.n = append(preComb.n, c[0])
		}
		preComb.k = Conf.PreBrute.Lenght
		preComb.indexes = make([]int, Conf.PreBrute.Lenght)
		for i, _ := range preComb.indexes {
			preComb.indexes[i] = 0
		}
		preComb.result = comb(preComb.n, preComb.k, preComb.indexes, preComb.result)
	}

	if Conf.PostBrute.Active {
		log.Println("Performing combinations [postBrute]...")
		postComb.post = true
		for _, c := range Conf.PostBrute.Chars {
			postComb.n = append(postComb.n, c[0])
		}
		postComb.k = Conf.PostBrute.Lenght
		postComb.indexes = make([]int, Conf.PostBrute.Lenght)
		for i, _ := range postComb.indexes {
			postComb.indexes[i] = 0
		}
		postComb.result = comb(postComb.n, postComb.k, postComb.indexes, postComb.result)
	}

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

//on first call indexes must be an array of k elements sets to 0
//result must be a 'var result [][]byte'
//TO DO
//indexes maybe unnecessary
func comb(n []byte, k int, indexes []int, result [][]byte) [][]byte {
	//Aggiungere controlli parametri
	temp := make([]byte, k)
	for i := k - 1; i >= 0; i-- {
		temp[i] = n[indexes[i]]
	}
	result = append(result, temp)

	rcount := k - 1
	for indexes[rcount] == len(n)-1 {
		rcount--
		if rcount == -1 {
			return result
		}
		if indexes[rcount] < len(n)-1 {
			indexes[rcount]++
			for i := k - 1; i > rcount; i-- {
				indexes[i] = 0
			}
			result = comb(n, k, indexes, result)
		}
	}
	indexes[k-1]++
	result = comb(n, k, indexes, result)
	return result
}
