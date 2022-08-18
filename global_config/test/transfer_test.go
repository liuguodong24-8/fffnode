package test

import (
	"github.com/fff-chain/go-fff/global_config"
	"github.com/fff-chain/go-fff/global_config/wallet"
	"log"
	"testing"
)

func TestTransfer(t *testing.T) {
	wl := wallet.Init("http://87.118.86.2:8488", "http://127.0.0.1:8545")
	hash := wl.Transfer("20a442166fda1598b32d85d0cf5b30fe5867bf23c943f8584013eab0ad49a88c", "FFF3rPN3k7M1EASWb6QU1QvJHXmrUGhPF8QS5U7hpDfkoWwGyb6aBoZhD1", global_config.EthToWei(10000), "", 0)
	log.Println(hash)
}
