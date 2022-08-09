package ledger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

const (
	configProviderName = "connection-config.yaml"
	channelName        = "mychannel"
	contractType       = "basic"
	walletPath         = "wallet"
	walletLabel        = "appUser"
	org1MSPid          = "Org1MSP"
	createTXFuncName   = "CreateTX"
	findTXFuncName     = "ReadTX"
	getAllTXFuncName   = "GetAllTXs"
)

var (
	//configPath = filepath.Join(
	//	"admin-msp",
	//)
	configPath   = "/Users/tommytian/codebase/admin-msp"
	org1CertPath = filepath.Join(configPath, "signcerts", "cert.pem")
	org1KeyDir   = filepath.Join(configPath, "keystore")
)

func init() {
	log.Println("============ application-golang starts ============")
	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "false")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environment variable: %v", err)
	}
	err = os.RemoveAll(walletPath)
	if err != nil {
		log.Fatalf("Error removing wallet directory: %v", err)
	}
}

type Controller struct {
	gw *gateway.Gateway
	ct *gateway.Contract
}

// NewController starts a new service instance
func NewController() *Controller {
	service := new(Controller)
	wallet, err := gateway.NewFileSystemWallet(walletPath)
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}
	if !wallet.Exists(walletLabel) {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(configProviderName))),
		gateway.WithIdentity(wallet, walletLabel),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	service.gw = gw
	network, err := gw.GetNetwork(channelName)
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}
	contract := network.GetContract(contractType)
	service.ct = contract
	// log.Println("--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger")
	// result, err := contract.SubmitTransaction("InitLedger")
	// if err != nil {
	// 	log.Fatalf("Failed to Submit transaction: %v", err)
	// }
	// log.Println(string(result))
	return service
}

func (s *Controller) Close() {
	s.gw.Close()
}

func (s *Controller) SubmitTX(binding string, timestamp int64) (string, error) {
	// log.Println("--> Submit Transaction: Invoke, function that adds a new asset")
	txID, err := s.ct.SubmitTransaction(createTXFuncName, binding, strconv.FormatInt(timestamp, 10))
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
		return "", err
	}
	// log.Println(string(result))
	return string(txID), nil
}

func (s *Controller) FindTX(txID string) error {
	result, err := s.ct.EvaluateTransaction(findTXFuncName, txID)
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(result))
	return nil
}

func (s *Controller) GetAllTXs() error {
	results, err := s.ct.EvaluateTransaction(getAllTXFuncName)
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(results))
	return nil
}

func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")

	// read the certificate pem
	cert, err := os.ReadFile(filepath.Clean(org1CertPath))
	if err != nil {
		return err
	}

	// there's a single file in this dir containing the private key
	fmt.Println(org1KeyDir)
	files, err := os.ReadDir(org1KeyDir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		//return fmt.Errorf("keystore folder should have contain one file")
		return fmt.Errorf("no file in keystore folder")
	}
	keyPath := filepath.Join(org1KeyDir, files[0].Name())
	key, err := os.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity(org1MSPid, string(cert), string(key))

	return wallet.Put(walletLabel, identity)
}
