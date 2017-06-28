package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

//Counter password tested
var passCount int32
var elapsed time.Duration
var startT time.Time
var found bool

//end print some outputs
//success true if passphrase FOUND
//success false if not FOUND
func end(success bool) {
	//Some output
	fmt.Print("\n")
	log.Printf("Tested %v passphrases in %v seconds.", passCount, elapsed.Seconds())

	if success {
		found = true
		log.Print("PASSPHRASE FOUND")
		fmt.Printf("\nWallet Address: %s\n\n------------ PASSPHRASE ------------\n\n%s\n\n------------------------------------\n\n", address, passphrase)
		fmt.Print("Please make a donation to developer:\n\nETH: 0x2feD76d5abE6c001D259eC769c28f6068E0166CB\nBTC: 1HTpxVw6KkDakhjqL3bgkYtM7Gsxxzmjw5\n\n")
		os.Exit(0)
	} else {
		found = false
		log.Print("Sorry. Passphrase not found!")
		os.Exit(0)
	}
}

func manager(chans []chan string, wg *sync.WaitGroup) {
	//defer fmt.Print("Manager: done.")
	found = false

	defer passFile.Close()
	defer wg.Done()
	time.Sleep(500 * time.Millisecond)
	// Create a scanner to read passList line by line
	scanner := bufio.NewScanner(passFile)
	//Start time for elapsed
	startT = time.Now()

	for scanner.Scan() {
		pass := scanner.Text()
		for _, c := range chans {
			c <- pass
		}
	}
	for _, c := range chans {
		close(c)
	}

	// check scanner errors
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

func onlyPass(ch <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	//defer fmt.Print("onlyPass: done.")
	log.Println("Searching Passphrase without variants... Please wait")
	time.Sleep(500 * time.Millisecond)
	for {
		select {
		case pass, ok := <-ch:
			if ok && testPass(pass) {
				elapsed = time.Since(startT)
				end(true)
				return
			} else if !ok {
				return
			}
		}
	}
}

func customVariants(ch <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	//defer fmt.Print("customVariants: done.")
	log.Println("Searching Passphrase with variants... Please wait")
	time.Sleep(500 * time.Millisecond)
	for {
		select {
		case pass, ok := <-ch:
			if ok && testPassVariants(pass) {
				elapsed = time.Since(startT)
				end(true)
				return
			} else if !ok {
				return
			}
		}
	}
}

func preBrute(ch <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	//defer fmt.Print("preBrute: done.")
	log.Println("Searching Passphrase with preBrute... Please wait")
	time.Sleep(500 * time.Millisecond)
	for {
		select {
		case pass, ok := <-ch:
			if testCombinations(pass, preComb) {
				elapsed = time.Since(startT)
				end(true)
				return
			} else if !ok {
				return
			}
		}
	}
}

func postBrute(ch <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	//defer fmt.Print("postBrute: done.")
	log.Println("Searching Passphrase with postBrute... Please wait")
	time.Sleep(500 * time.Millisecond)
	for {
		select {
		case pass, ok := <-ch:
			if testCombinations(pass, postComb) {
				elapsed = time.Since(startT)
				end(true)
				return
			} else if !ok {
				return
			}
		}
	}
}

//Main function
func main() {

	var wg sync.WaitGroup

	chans := make([]chan string, 0)

	if Conf.CustomVariants { //With Password Variants
		ch1 := make(chan string, 250)
		ch2 := make(chan string, 250)
		chans = append(chans, ch1)
		chans = append(chans, ch2)
		wg.Add(2)
		go onlyPass(ch1, &wg)
		go customVariants(ch2, &wg)
	} else { //Without Password Variants
		ch := make(chan string, 250)
		chans = append(chans, ch)
		wg.Add(1)
		go onlyPass(ch, &wg)
	}

	if Conf.PreBrute.Active {
		ch := make(chan string, 250)
		chans = append(chans, ch)
		wg.Add(1)
		go preBrute(ch, &wg)
	}

	if Conf.PostBrute.Active {
		ch := make(chan string, 250)
		chans = append(chans, ch)
		wg.Add(1)
		go postBrute(ch, &wg)
	}

	wg.Add(1)
	manager(chans, &wg)
	wg.Wait()
	if !found {
		//if passphrase not found
		elapsed = time.Since(startT)
		end(false)
	}
}
