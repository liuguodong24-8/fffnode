package nft

import (
	"context"
	"github.com/fff-chain/go-fff/accounts/abi/bind"
	"github.com/fff-chain/go-fff/common"
	"github.com/fff-chain/go-fff/crypto"
	"github.com/fff-chain/go-fff/ethclient"
	"github.com/fff-chain/go-fff/global_config/wallet"
	"log"
	"testing"
)

func TestDeployAddress(t *testing.T) {


	d, _ := ethclient.Dial("http://87.118.86.2:8488")

	pri := "20a442166fda1598b32d85d0cf5b30fe5867bf23c943f8584013eab0ad49a88c"


	priKeyECD, err := crypto.HexToECDSA(pri)

	ownerAddress:=wallet.GetPublicAddressFromKey(pri)

	if err != nil {
		log.Println("私钥异常")

	}
	chianId, err := d.ChainID(context.Background())
	if err != nil {
		log.Println("节点chianId异常", err)

	}
	param, err := bind.NewKeyedTransactorWithChainID(priKeyECD, chianId)
	if err != nil {
		log.Println("节点Transactor异常", err)

	}

	param.GasPrice,_= d.SuggestGasPrice(context.Background())


	addr,_,_,_:=DeployFFFNFT(param,d)

	log.Println(addr,ownerAddress)



}

func TestFFFNFTTransactor_Mint(t *testing.T) {
	d, _ := ethclient.Dial("http://87.118.86.2:8488")

	pri := "20a442166fda1598b32d85d0cf5b30fe5867bf23c943f8584013eab0ad49a88c"


	priKeyECD, err := crypto.HexToECDSA(pri)

	ownerAddress:=wallet.GetPublicAddressFromKey(pri)

	contractAddress:=common.HexToAddress("FFF3kbeRvCXCQT4eJ4VZ6mRZDu9RzWaNCXK8csoCGiCBNNf8iu1QMpMnt1")

	cl,_:=NewFFFNFTTransactor(contractAddress,d)
	chianId, err := d.ChainID(context.Background())
	if err != nil {
		log.Println("节点chianId异常", err)

	}
	param, err := bind.NewKeyedTransactorWithChainID(priKeyECD, chianId)
	if err != nil {
		log.Println("节点Transactor异常", err)

	}

	param.GasPrice,_= d.SuggestGasPrice(context.Background())

	z,_:=cl.Mint(param,common.HexToAddress(ownerAddress),"1","1","1",nil)

	log.Println(z.Hash())
}

func TestFFFNFTCaller_GetApproved(t *testing.T) {

	d, _ := ethclient.Dial("http://87.118.86.2:8488")

	rc,_:=d.TransactionReceipt(context.Background(),common.HexToHash("0xa6438cfe68f349743b9d7ed85df14510bea7938f36a7b1b467a124d1a6efbbbf"))
	log.Println(rc.BlockNumber)

}
func TestDeployEnumerableMap(t *testing.T) {
	pri := "218ec8f94b29a96d38a21b19e0b866faf9963f48a1a0d575d91990fdf08f8936"

	ownerAddress:=wallet.GetPublicAddressFromKey(pri)

	log.Println(ownerAddress)

}