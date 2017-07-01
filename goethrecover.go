package main

import (
	"bufio"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type worker struct {
	name     string
	routines *int32
	in       chan string
	done     chan *worker
	found    bool
	f        func(*worker)
	wg       *sync.WaitGroup
}

func (w *worker) start(f func(*worker), wg *sync.WaitGroup, name string) *worker {
	w.name = name
	w.routines = new(int32)
	*w.routines = int32(1)
	w.in = make(chan string, 250)
	w.done = make(chan *worker)
	w.f = f
	wg.Add(1)
	w.wg = wg
	go w.f(w)
	return w
}

func (w *worker) add(n int32) {
	r := atomic.LoadInt32(w.routines)
	r += n
	atomic.SwapInt32(w.routines, r)
	var i int32
	for i = 0; i < n; i++ {
		w.wg.Add(1)
		go w.f(w)
	}
}

func (w *worker) finished(n int32) {
	r := atomic.LoadInt32(w.routines)
	r -= n
	atomic.SwapInt32(w.routines, r)
}

var elapsed time.Duration
var startT time.Time

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
		log.Print("Sorry. Passphrase not found!\n\n")
		fmt.Print("------------------------------------\n\nPlease make a donation to developer:\n\nETH: 0x2feD76d5abE6c001D259eC769c28f6068E0166CB\nBTC: 1HTpxVw6KkDakhjqL3bgkYtM7Gsxxzmjw5\n\n")

	}
}

func maxLoads(workers []*worker) *worker {
	max := len(workers[0].in)
	maxW := workers[0]
	for _, w := range workers {
		if len(w.in) > max {
			max = len(w.in)
			maxW = w
		}
	}
	return maxW
}

func manager(workers []*worker) {

	defer passFile.Close()

	procs := runtime.GOMAXPROCS(runtime.NumCPU())
	if len(workers) < procs {
		n := procs - len(workers)
		for i := 0; i < n; i++ {
			time.Sleep(1 * time.Second)
			t := maxLoads(workers)
			t.add(1)
		}
	}
	// Create a scanner to read passList line by line
	scanner := bufio.NewScanner(passFile)
	//Start time for elapsed
	startT = time.Now()

	for scanner.Scan() {
		pass := scanner.Text()
		for _, w := range workers {
			w.in <- pass
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	for _, w := range workers {
		close(w.in)
	}

	for {
		for i := 0; i < len(workers); i++ {
			w := workers[i]
			select {
			case wd := <-w.done:
				if wd.found {
					elapsed = time.Since(startT)
					end(true)
					return
				} else {
					//delete worker
					if len(workers) != 1 {
						t := maxLoads(workers)
						t.add(1)
					}
					n := atomic.LoadInt32(wd.routines)
					if n == 1 {
						workers = append(workers[:i], workers[i+1:]...)
						if len(workers) == 0 {
							elapsed = time.Since(startT)
							end(false)
							return
						}
					} else {
						w.finished(1)
					}
				}
			}
		}
	}
}

func onlyPass(w *worker) {
	defer w.wg.Done()
	for {
		select {
		case pass, ok := <-w.in:
			if ok && testPass(pass) {
				w.found = true
				w.done <- w
				return
			} else if !ok {
				w.found = false
				w.done <- w
				return
			}
		}
	}
}

func customVariants(w *worker) {
	defer w.wg.Done()
	for {
		select {
		case pass, ok := <-w.in:
			if ok && testPassVariants(pass) {
				w.found = true
				w.done <- w
				return
			} else if !ok {
				w.found = false
				w.done <- w
				return
			}
		}
	}
}

func preBrute(w *worker) {
	defer w.wg.Done()
	for {
		select {
		case pass, ok := <-w.in:
			if ok && testCombinations(pass, preComb) {
				w.found = true
				w.done <- w
				return
			} else if !ok {
				w.found = false
				w.done <- w
				return
			}
		}
	}
}

func postBrute(w *worker) {
	defer w.wg.Done()
	for {
		select {
		case pass, ok := <-w.in:
			if ok && testCombinations(pass, postComb) {
				w.found = true
				w.done <- w
				return
			} else if !ok {
				w.found = false
				w.done <- w
				return
			}
		}
	}
}

//Main function
func main() {

	var wg sync.WaitGroup
	var w []*worker

	if Conf.CustomVariants { //With Password Variants
		work1 := new(worker).start(onlyPass, &wg, "onlyPass")
		work2 := new(worker).start(customVariants, &wg, "CustomVariants")
		log.Println("Searching Passphrase without variants... Please wait")
		w = append(w, work1)
		log.Println("Searching Passphrase with variants... Please wait")
		w = append(w, work2)
	} else { //Without Password Variants
		work := new(worker).start(onlyPass, &wg, "OnlyPass")
		log.Println("Searching Passphrase without variants... Please wait")
		w = append(w, work)
	}

	if Conf.PreBrute.Active {
		work := new(worker).start(preBrute, &wg, "preBrute")
		log.Println("Searching Passphrase with preBrute... Please wait")
		w = append(w, work)
	}

	if Conf.PostBrute.Active {
		work := new(worker).start(postBrute, &wg, "PostBrute")
		log.Println("Searching Passphrase with postBrute... Please wait")
		w = append(w, work)
	}

	manager(w)
	//wg.Wait()
}
