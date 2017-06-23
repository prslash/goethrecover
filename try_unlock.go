package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

//Key variables
var passphrase string
var address string

//testPass check if pass is correct.
//Return true if pass is correct, false if pass is invalid
func testPass(pass string) bool {

	fmt.Printf("\rTesting %d password", passCount)
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
	for _, seqPre := range Conf.PreSeq {
		passTemp = seqPre + pass
		if testPass(passTemp) {
			return true
		}
	}

	//Test pass with all postSeq
	for _, seqPost := range Conf.PostSeq {
		passTemp = pass + seqPost
		if testPass(passTemp) {
			return true
		}
	}

	//Test pass with all preSeq - postSeq combinations
	for _, seqPre := range Conf.PreSeq {
		for _, seqPost := range Conf.PostSeq {
			passTemp = seqPre + pass + seqPost
			if testPass(passTemp) {
				return true
			}
		}
	}
	return false
}

func testCombinations(pass string, combs combData) bool {
	if combs.pre == true {
		for _, s := range combs.result {
			if testPass(string(s) + pass) {
				return true
			}
		}
	}
	if combs.post == true {
		for _, s := range combs.result {
			if testPass(pass + string(s)) {
				return true
			}
		}
	}
	return false
}
