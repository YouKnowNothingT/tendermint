package client_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/tendermint/tendermint/abci/example/kvstore"
	nm "github.com/tendermint/tendermint/node"
	rpctest "github.com/tendermint/tendermint/rpc/test"

	dkgOffChain "github.com/corestario/dkglib/lib/offChain"
)

var node *nm.Node

func TestMain(m *testing.M) {
	// start a tendermint node (and kvstore) in the background to test against
	dir, err := ioutil.TempDir("/tmp", "rpc-client-test")
	if err != nil {
		panic(err)
	}
	app := kvstore.NewPersistentKVStoreApplication(dir)
	node = rpctest.StartTendermint(app)

	node.GetConsensusState().SetVerifier(dkgOffChain.GetVerifier(1, 1)("rpc_client_tests", 0))
	code := m.Run()

	// and shut down proper at the end
	rpctest.StopTendermint(node)
	os.Exit(code)
}
