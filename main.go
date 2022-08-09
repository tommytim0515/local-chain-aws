package main

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	"github.com/ttvs-blockchain/local-chain-aws/internal/ledger"
)

func main() {
	// log level to INFO
	logging.SetLevel("", logging.DEBUG)
	lc := ledger.NewController()
	defer lc.Close()
	//id, err := lc.SubmitTX("test", time.Now().UnixNano())
	//handleErr(err)
	err := lc.GetAllTXs()
	handleErr(err)
	//log.Println("id: ", id)
	//log.Println("FindTX")
	//err = lc.FindTX(id)
	//handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
