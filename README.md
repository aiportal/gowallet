# gowallet

A console bitcoin wallet application written in golang. 

![GoWallet Account View](https://raw.githubusercontent.com/aiportal/gowallet/master/_doc/account.png)

**GoWallet uses a secret phrase and a salt phrase to generate your safe wallets.**  
Project location: https://github.com/aiportal/gowallet

**GoWallet is a safe brain wallet for bitcoin.**  
  Secret phrase at least 16 characters.  
  Salt phrase at least 6 characters.  

  Secret phrases should contain uppercase letters, lowercase letters, numbers, and special characters.  
  Both secret phrases and salt phrases can use hexadecimal notation such as \xff or \xFF to represent a character.   

**It is recommended that use a more complex secret and put it on paper.**  
**It's also recommended that keep your salt in mind.**  

Donations are welcome at <code>[<b>1BTC</b>zvzTn7QYBwFkRRkXGcVPodwrYoQyAq](https://blockchain.info/address/1BTCzvzTn7QYBwFkRRkXGcVPodwrYoQyAq)</code>

![GoWallet Encryption Process](https://raw.githubusercontent.com/aiportal/gowallet/master/_doc/encryption.png)

#### Installation 

requires installation of golang version 1.10 or later

```
go get github.com/aiportal/gowallet
dep ensure
go install
```

#### Advanced usage

You can export bulk wallets using the command line.

  -n or -number uint  
    &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
    Number of wallets to generate.   
  -v or -vanity string  
    &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
    Find vanity wallet address matching. (prefix)  
  -e or -export string  
    &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
    Export wallets(child number, private key and address) in WIF format.
