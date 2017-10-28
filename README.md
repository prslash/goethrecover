## goethrecover
### What is goethrecover?
goethrecover is a very simple tool written in golang to recover lost passphrase of Ethereum Wallets using a list of passphrases (dictionary attack). You can also add some prefix and suffix strings to these passphrases.

## How to install goethrecover?
To run goethrecover you need to install golang from [Install Golang](https://golang.org/doc/install)

Add some packages using 'go get'
```
go get "github.com/ethereum/go-ethereum/accounts/keystore"
go get "github.com/naoina/toml"
```

Then build goethrecover
```
go build goethrecover.go try_unlock.go config.go
```

## How does it works?
It is simple...
```
./goethrecover -h

Usage of ./goethrecover:
  -conf string
    	TOML file with some configuration (default "conf.toml")
  -passList string
    	Text file with one passphrase per line (default "passList.txt")
  -wallet string
    	Wallet file (default "wallet.json")
```

In conf.toml there are some configurations option. Lines with # are comments<br />
See conf.toml example below:
```
#---CONFIGURATION FILE---
#Modify values as you wish

#testVariants is boolean value
#true: Prefix and Suffix variations are tested
#false: Prefix and Suffix variations are not tested
testVariants = true

#preSeq is an array of prefix strings to test
preSeq = [ "123", "2008", "FOX" ]
#postSeq is an array of suffix strings to test
postSeq = [ "456", "1995", "512"]

#preBrute section
[preBrute]
active = true
lenght = 3
chars = ['a', 'b', 'c']

[postBrute]
active = true
lenght = 3
chars = ['1', '2', '3']
```
If you want to try prefix and suffix you need to set: `testVariants = true`<br />
-Set prefix and suffix strings<br />
Otherwise set: `testVariants = false`<br /><br />
If you want to try brute force of prefix (preBrute) and/or suffix (postBrute) set:<br />
`active = true`<br />
`lenght = 3` lenght of string to bruteforce<br />
`chars = ['a', 'b', 'c']` set of custom character<br />

### Example
```
./goethrecover -wallet "myWallet.json" -passList "myFavoritesPasswords.txt" -conf "myOwnConf.toml"
```

## Conclusion
This is a very simple tool. I hope it will help you to recover your ethereum wallet.
Feel free to suggest new features or improvements.

Please make a donation to developer:

ETH: 0x2feD76d5abE6c001D259eC769c28f6068E0166CB<br />
BTC: 1HTpxVw6KkDakhjqL3bgkYtM7Gsxxzmjw5
