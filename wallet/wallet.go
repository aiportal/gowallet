package wallet

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"strconv"
)

type Wallet struct {
	No      uint32
	Private string
	Address string
}


func (w *Wallet) Balance() (satoshi uint64, err error) {

	_url := fmt.Sprintf("https://blockchain.info/q/addressbalance/%s", w.Address)
	resp, err := http.Get(_url)
	if err != nil {
		return
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	satoshi, err = strconv.ParseUint(string(bs), 10, 64)
	if err != nil {
		return
	}
	return

	//https://blockchain.info/q/addressbalance/1DogeKd9JrUNzFaLEyWAVxCVXSvWxe6sAm
	//https://blockexplorer.com/api/addr/
	//https://www.bitgo.com/api/v1/address/
	//https://bitcoin.toshi.io/api/v0/addresses/
	//https://chain.api.btc.com/v3/address/
	//https://api.blocktrail.com/v1/btc/address/1NcXPMRaanz43b1kokpPuYDdk6GGDvxT2T?api_key=MY_APIKEY
	//https://api.blockcypher.com/v1/btc/main/addrs/1DEP8i3QJCsomS4BSMY2RpU1upv62aGvhD/balance
	//https://api-r.bitcoinchain.com/v1/address/1Chain4asCYNnLVbvG6pgCLGBrtzh4Lx4b
	//https://api.kaiko.com/v1/addresses/3Nt1smucEdFks8uYQhyGvXGBuocTcMSmsT
	//https://chainflyer.bitflyer.jp/v1/address/1LDWeSRJukN7zWXDBpuvB2WGsMxYE7UTnQ
	//https://insight.bitpay.com/api/addr/1NcXPMRaanz43b1kokpPuYDdk6GGDvxT2T/?noTxList=1
	//https://api.coinprism.com/v1/addresses/1dice97ECuByXAvqXpaYzSaQuPVvrtmz6
	//http://btc.blockr.io/api/v1/address/info/
}
