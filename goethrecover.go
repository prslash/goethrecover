package main

import (
	"bufio"
	"fmt"
	"log"
	"time"
)

//Counter password tested
var passCount int
var elapsed time.Duration

//end print some outputs
//success true if passphrase FOUND
//success false if not FOUND
func end(success bool) {
	//Some output
	fmt.Print("\n")
	log.Printf("Tested %v passphrases in %v seconds.", passCount, elapsed.Seconds())

	if success {
		log.Print("PASSPHRASE FOUND")
		fmt.Printf("\nWallet Address: %s\n\n------------ PASSPHRASE ------------\n\n%s\n\n------------------------------------\n\n", address, passphrase)
		fmt.Print("Please make a donation to developer:\n\nETH: 0x2feD76d5abE6c001D259eC769c28f6068E0166CB\nBTC: 1HTpxVw6KkDakhjqL3bgkYtM7Gsxxzmjw5\n\n")
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
	if Conf.CustomVariants { //With Password Variants
		fmt.Print("\n")
		log.Println("Searching Passphrase with variants and without variants... Please wait")
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
		fmt.Print("\n")
		log.Println("\nSearching Passphrase without variants... Please wait")
		for scanner.Scan() {
			pass := scanner.Text()
			if testPass(pass) {
				elapsed = time.Since(startT)
				end(true)
				return
			}
		}
	}
	if Conf.PreBrute.Active {
		fmt.Print("\n")
		log.Println("Searching Passphrase with preBrute... Please wait")
		passFile.Seek(0, 0)
		scanner = bufio.NewScanner(passFile)
		for scanner.Scan() {
			pass := scanner.Text()
			if testCombinations(pass, preComb) {
				elapsed = time.Since(startT)
				end(true)
				return
			}
		}
	}
	if Conf.PostBrute.Active {
		fmt.Print("\n")
		log.Println("Searching Passphrase with postBrute... Please wait")
		passFile.Seek(0, 0)
		scanner = bufio.NewScanner(passFile)
		for scanner.Scan() {
			pass := scanner.Text()
			if testCombinations(pass, postComb) {
				elapsed = time.Since(startT)
				end(true)
				return
			}
		}
	}
	//if conf.T

	//if passphrase not found
	elapsed = time.Since(startT)

	// check scanner errors
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	end(false)
	return
}
