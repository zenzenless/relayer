package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	ics23 "github.com/cosmos/ics23/go"

	///home/james/gopath/pkg/mod/github.com/cosmos/ibc-go/v8@v8.0.0/modules/light-clients/07-tendermint/client_state.go
	tmclient "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	"github.com/cosmos/relayer/v2/cmd"
	"github.com/cosmos/relayer/v2/relayer/chains/cosmos"
	"github.com/cosmos/relayer/v2/relayer/provider"
	zaplogfmt "github.com/jsternberg/zap-logfmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gotest.tools/v3/assert"
)

func initClient(t *testing.T, ctx context.Context) (rollapp, mehub provider.ChainProvider) {
	l, err := newRootLogger("auto", "debug")
	assert.NilError(t, err)
	// init rollapp rpc client
	var pcw cmd.ProviderConfigWrapper
	cfgFile, err := ioutil.ReadFile("./examples/demo/configs/chains/rollappevm_1234-1.json")
	assert.NilError(t, err)
	err = json.Unmarshal(cfgFile, &pcw)
	assert.NilError(t, err)

	rollapp, err = pcw.Value.NewProvider(l, "home/james/.relayer", false, "rollapp")
	assert.NilError(t, err)
	err = rollapp.Init(ctx)
	assert.NilError(t, err)

	//init mehub rpc client
	var pcwMehub cmd.ProviderConfigWrapper
	meCfgFile, err := ioutil.ReadFile("./examples/demo/configs/chains/me-hub.json")
	assert.NilError(t, err)
	err = json.Unmarshal(meCfgFile, &pcwMehub)
	assert.NilError(t, err)
	mehub, err = pcwMehub.Value.NewProvider(l, "home/james/.relayer", false, "mehub")
	assert.NilError(t, err)
	err = mehub.Init(ctx)
	assert.NilError(t, err)
	return
	//init done
}
func TestCosmosClie(t *testing.T) {
	ctx := context.Background()
	rollapp, mehub := initClient(t, ctx)
	//init done

	lightClientOnMeHubStateRes, err := mehub.QueryClientStateResponse(ctx, 17617, "07-tendermint-0")
	assert.NilError(t, err)
	fmt.Println(lightClientOnMeHubStateRes.String())

	// 查询me-hub上轻客户端的共识状态
	lightClientOnMeHubState, ok := lightClientOnMeHubStateRes.ClientState.GetCachedValue().(*tmclient.ClientState)
	assert.Assert(t, ok)
	latestHeightRollappClientOnMe := lightClientOnMeHubState.GetLatestHeight()
	meConStateRes, err := mehub.QueryClientConsensusState(ctx, 17617, "07-tendermint-0", latestHeightRollappClientOnMe)
	assert.NilError(t, err)
	csState, ok := meConStateRes.ConsensusState.GetCachedValue().(*tmclient.ConsensusState)
	assert.Assert(t, ok)

	//查询rollapp checkin key store
	cc, ok := rollapp.(*cosmos.CosmosProvider)
	if !ok {
		t.Fatal("not a cosmos provider")
	}

	req := abci.RequestQuery{
		Path:   fmt.Sprintf("store/%s/key", "checkin"),
		Height: int64(latestHeightRollappClientOnMe.GetRevisionHeight()) - 1,
		Data:   append([]byte("Dawn/value/"), []byte{0}...),
		Prove:  true,
	}

	res, err := cc.QueryABCI(ctx, req)
	assert.NilError(t, err)

	dawn := &Dawn{}
	err = dawn.Unmarshal(res.Value)
	assert.NilError(t, err)
	//

	merkleProof, err := commitmenttypes.ConvertProofs(res.ProofOps)
	assert.NilError(t, err)
	fmt.Println(merkleProof)

	path := commitmenttypes.NewMerklePath("checkin", string(append([]byte("Dawn/value/"), []byte{0}...)))
	err = merkleProof.VerifyMembership(lightClientOnMeHubState.ProofSpecs, csState.GetRoot(), path, res.Value)
	assert.NilError(t, err)

}

func newRootLogger(format string, logLevel string) (*zap.Logger, error) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(ts.UTC().Format("2006-01-02T15:04:05.000000Z07:00"))
	}
	config.LevelKey = "lvl"

	var enc zapcore.Encoder
	switch format {
	case "json":
		enc = zapcore.NewJSONEncoder(config)
	case "auto", "console":
		enc = zapcore.NewConsoleEncoder(config)
	case "logfmt":
		enc = zaplogfmt.NewEncoder(config)
	default:
		return nil, fmt.Errorf("unrecognized log format %q", format)
	}

	level := zap.InfoLevel
	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "panic":
		level = zapcore.PanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	}
	return zap.New(zapcore.NewCore(
		enc,
		os.Stderr,
		level,
	)), nil
}

func TestABCIQuery(t *testing.T) {
	l, err := newRootLogger("auto", "debug")
	assert.NilError(t, err)
	ctx := context.Background()

	// init rollapp rpc client
	var pcw cmd.ProviderConfigWrapper
	cfgFile, err := ioutil.ReadFile("./examples/demo/configs/chains/rollappevm_1234-1.json")
	assert.NilError(t, err)
	err = json.Unmarshal(cfgFile, &pcw)
	assert.NilError(t, err)

	rollapp, err := pcw.Value.NewProvider(l, "home/james/.relayer", false, "rollapp")
	assert.NilError(t, err)
	err = rollapp.Init(ctx)
	assert.NilError(t, err)

	cc, ok := rollapp.(*cosmos.CosmosProvider)
	if !ok {
		t.Fatal("not a cosmos provider")
	}

	req := abci.RequestQuery{
		Path:   fmt.Sprintf("store/%s/subspace", "did"),
		Height: 3451,
		Data:   []byte{0x40},
		Prove:  true,
	}

	res, err := cc.QueryABCI(ctx, req)
	fmt.Printf(string(res.Value))
	assert.NilError(t, err)
	fmt.Println(res.Code)
	dawn := &Dawn{}
	err = dawn.Unmarshal(res.Value)
	assert.NilError(t, err)

	merkleProof, err := commitmenttypes.ConvertProofs(res.ProofOps)
	assert.NilError(t, err)
	fmt.Println(merkleProof.Proofs[1].GetExist().Value)
	fmt.Printf("value:%+v hash:%x", dawn, merkleProof.Proofs[1].GetExist().Value)

	req = abci.RequestQuery{
		Path:   fmt.Sprintf("store/%s/key", "checkin"),
		Height: 3452,
		Data:   append([]byte("Dawn/value/"), []byte{0}...),
		Prove:  true,
	}

	res, err = cc.QueryABCI(ctx, req)
	assert.NilError(t, err)
	fmt.Println(res.Code)
	dawn = &Dawn{}
	err = dawn.Unmarshal(res.Value)
	assert.NilError(t, err)

	merkleProof, err = commitmenttypes.ConvertProofs(res.ProofOps)
	assert.NilError(t, err)
	fmt.Println(merkleProof.Proofs[1].GetExist().Value)
	fmt.Printf("value:%+v hash:%x", dawn, merkleProof.Proofs[1].GetExist().Value)
}

func TestDID(t *testing.T) {
	l, err := newRootLogger("auto", "debug")
	assert.NilError(t, err)
	ctx := context.Background()
	// init rollapp rpc client
	var pcw cmd.ProviderConfigWrapper
	cfgFile, err := ioutil.ReadFile("./examples/demo/configs/chains/rollappevm_1234-1.json")
	assert.NilError(t, err)
	err = json.Unmarshal(cfgFile, &pcw)
	assert.NilError(t, err)

	rollapp, err := pcw.Value.NewProvider(l, "home/james/.relayer", false, "rollapp")
	assert.NilError(t, err)
	err = rollapp.Init(ctx)
	assert.NilError(t, err)

	//init mehub rpc client
	var pcwMehub cmd.ProviderConfigWrapper
	meCfgFile, err := ioutil.ReadFile("./examples/demo/configs/chains/me-hub.json")
	assert.NilError(t, err)
	err = json.Unmarshal(meCfgFile, &pcwMehub)
	assert.NilError(t, err)
	mehub, err := pcwMehub.Value.NewProvider(l, "home/james/.relayer", false, "mehub")
	assert.NilError(t, err)
	err = mehub.Init(ctx)
	assert.NilError(t, err)
	//init done

	var rollappReqHeight int64 = 47
	addr := "me19zajfxfn2wy39pdprscss8uztpxt4vtz3n0qz2"
	mehub_light_state, err := rollapp.QueryClientStateResponse(ctx, rollappReqHeight, "07-tendermint-0")
	assert.NilError(t, err)
	fmt.Println(mehub_light_state.String())

	// 查询rollapp上mehub的轻客户端的共识状态
	mehub_light_state_r, ok := mehub_light_state.ClientState.GetCachedValue().(*tmclient.ClientState)
	assert.Assert(t, ok)
	mehub_light_height := mehub_light_state_r.GetLatestHeight()
	meConStateRes, err := rollapp.QueryClientConsensusState(ctx, rollappReqHeight, "07-tendermint-0", mehub_light_height)
	t.Log("me-hub height", mehub_light_height)
	assert.NilError(t, err)
	csState, ok := meConStateRes.ConsensusState.GetCachedValue().(*tmclient.ConsensusState)
	assert.Assert(t, ok)
	key := []byte{0x10}
	address := sdk.MustAccAddressFromBech32(addr)
	key = append(key, address...)
	//查询mehub did key store
	cc, ok := mehub.(*cosmos.CosmosProvider)
	if !ok {
		t.Fatal("not a cosmos provider")
	}

	req := abci.RequestQuery{
		Path:   fmt.Sprintf("store/%s/key", "did"),
		Height: int64(mehub_light_height.GetRevisionHeight()) - 1,
		Data:   key,
		Prove:  true,
	}

	res, err := cc.QueryABCI(ctx, req)
	assert.NilError(t, err)

	//

	merkleProof, err := commitmenttypes.ConvertProofs(res.ProofOps)
	assert.NilError(t, err)
	fmt.Println(merkleProof)

	path := commitmenttypes.NewMerklePath("did", string(key))
	fmt.Printf("root:%x\n", csState.GetRoot().GetHash())
	err = merkleProof.VerifyMembership(mehub_light_state_r.ProofSpecs, csState.GetRoot(), path, res.Value)
	assert.NilError(t, err)

	proofByte, err := merkleProof.Marshal()
	assert.NilError(t, err)
	didInfo := &DidUpdateInfo{
		Address: addr,
		Did:     res.Value,
		Proof:   proofByte,
	}
	did, _ := json.Marshal(didInfo)
	fmt.Println(string(did))

}

type DidUpdateInfo struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Did     []byte `protobuf:"bytes,2,opt,name=did,proto3" json:"did,omitempty"`
	Proof   []byte `protobuf:"bytes,3,opt,name=proof,proto3" json:"proof,omitempty"`
}

func TestKycCredential(t *testing.T) {
	l, err := newRootLogger("auto", "debug")
	assert.NilError(t, err)
	ctx := context.Background()
	// init rollapp rpc client
	var pcw cmd.ProviderConfigWrapper
	cfgFile, err := ioutil.ReadFile("./examples/demo/configs/chains/rollappevm_1234-1.json")
	assert.NilError(t, err)
	err = json.Unmarshal(cfgFile, &pcw)
	assert.NilError(t, err)

	rollapp, err := pcw.Value.NewProvider(l, "home/james/.relayer", false, "rollapp")
	assert.NilError(t, err)
	err = rollapp.Init(ctx)
	assert.NilError(t, err)

	//init mehub rpc client
	var pcwMehub cmd.ProviderConfigWrapper
	meCfgFile, err := ioutil.ReadFile("./examples/demo/configs/chains/me-hub.json")
	assert.NilError(t, err)
	err = json.Unmarshal(meCfgFile, &pcwMehub)
	assert.NilError(t, err)
	mehub, err := pcwMehub.Value.NewProvider(l, "home/james/.relayer", false, "mehub")
	assert.NilError(t, err)
	err = mehub.Init(ctx)
	assert.NilError(t, err)
	//init done

	var rollappReqHeight int64 = 555
	didstr := "0000000000000099"
	mehub_light_state, err := rollapp.QueryClientStateResponse(ctx, rollappReqHeight, "07-tendermint-0")
	assert.NilError(t, err)
	fmt.Println(mehub_light_state.String())

	// 查询rollapp上mehub的轻客户端的共识状态
	mehub_light_state_r, ok := mehub_light_state.ClientState.GetCachedValue().(*tmclient.ClientState)
	assert.Assert(t, ok)
	mehub_light_height := mehub_light_state_r.GetLatestHeight()
	meConStateRes, err := rollapp.QueryClientConsensusState(ctx, rollappReqHeight, "07-tendermint-0", mehub_light_height)
	assert.NilError(t, err)
	csState, ok := meConStateRes.ConsensusState.GetCachedValue().(*tmclient.ConsensusState)
	assert.Assert(t, ok)
	key := []byte{0x40}

	key = append(key, didstr...)
	key = append(key, []byte("kyc")...)
	//查询mehub did key store
	cc, ok := mehub.(*cosmos.CosmosProvider)
	if !ok {
		t.Fatal("not a cosmos provider")
	}

	req := abci.RequestQuery{
		Path:   fmt.Sprintf("store/%s/key", "did"),
		Height: int64(mehub_light_height.GetRevisionHeight()) - 1,
		Data:   key,
		Prove:  true,
	}

	res, err := cc.QueryABCI(ctx, req)
	assert.NilError(t, err)

	//

	merkleProof, err := commitmenttypes.ConvertProofs(res.ProofOps)
	assert.NilError(t, err)
	fmt.Println(merkleProof)

	path := commitmenttypes.NewMerklePath("did", string(key))
	err = merkleProof.VerifyMembership(mehub_light_state_r.ProofSpecs, csState.GetRoot(), path, res.Value)
	assert.NilError(t, err)

	proofByte, err := merkleProof.Marshal()
	assert.NilError(t, err)
	didInfo := &CredentialUpdateInfo{
		Did:        didstr,
		Credential: res.GetValue(),
		Proof:      proofByte,
	}
	did, _ := json.Marshal(didInfo)
	fmt.Println(string(did))
	fmt.Println(string(res.Value))
	fmt.Println("size of did proof:", len(proofByte))

}

type CredentialUpdateInfo struct {
	Did        string `protobuf:"bytes,1,opt,name=did,proto3" json:"did,omitempty"`
	Credential []byte `protobuf:"bytes,2,opt,name=credential,proto3" json:"credential,omitempty"`
	Proof      []byte `protobuf:"bytes,3,opt,name=proof,proto3" json:"proof,omitempty"`
}

func TestScanCredential(t *testing.T) {
	ctx := context.Background()
	_, mehub := initClient(t, ctx)
	cc, ok := mehub.(*cosmos.CosmosProvider)
	if !ok {
		t.Fatal("not a cosmos provider")
	}
	req := abci.RequestQuery{
		Path: "/metaearth.kyc.Query/KYCs",
		//Path:   fmt.Sprintf("store/%s/subspace", "did"),
		Height: 0,
		Data:   nil,
		Prove:  false,
	}

	res, err := cc.QueryABCI(ctx, req)
	fmt.Printf(string(res.Value))
	assert.NilError(t, err)
	fmt.Println(res.Code)

}

type QueryKYCs struct {
	RegionId   string             `protobuf:"bytes,1,opt,name=regionId,proto3" json:"regionId,omitempty"`
	Pagination *query.PageRequest `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func TestSubscribe(t *testing.T) {
	mehub, err := client.NewClientFromNode("http://localhost:46657")
	assert.NilError(t, err)
	mehub.Start()
	ch, err := mehub.Subscribe(context.TODO(), "test", `tm.event='NewBlock'`)
	assert.NilError(t, err)
	get := <-ch
	t.Log(get)
}

func TestBatchVerify(t *testing.T) {
	ctx := context.Background()
	rollapp, mehub := initClient(t, ctx)

	//
	var rollappReqHeight int64 = 8
	mehub_light_state, err := rollapp.QueryClientStateResponse(ctx, rollappReqHeight, "07-tendermint-0")
	assert.NilError(t, err)
	fmt.Println(mehub_light_state.String())

	// 查询rollapp上mehub的轻客户端的共识状态
	mehub_light_state_r, ok := mehub_light_state.ClientState.GetCachedValue().(*tmclient.ClientState)
	assert.Assert(t, ok)
	mehub_light_height := mehub_light_state_r.GetLatestHeight()
	meConStateRes, err := rollapp.QueryClientConsensusState(ctx, rollappReqHeight, "07-tendermint-0", mehub_light_height)
	assert.NilError(t, err)
	csState, ok := meConStateRes.ConsensusState.GetCachedValue().(*tmclient.ConsensusState)
	assert.Assert(t, ok)

	getProofs := func() (didRoot []byte, didProof *ics23.CommitmentProof, items map[string][]byte, proofs []*ics23.CommitmentProof) {

		prefix := []byte{0x40}
		items = make(map[string][]byte)
		for i := 0; i < 10; i++ {
			key := append(prefix, []byte(fmt.Sprintf("%016d", i))...)
			key = append(key, []byte("kyc")...)
			//查询mehub did key store
			cc, ok := mehub.(*cosmos.CosmosProvider)
			if !ok {
				t.Fatal("not a cosmos provider")
			}
			req := abci.RequestQuery{
				Path:   fmt.Sprintf("store/%s/key", "did"),
				Height: 281512,
				//Height: int64(mehub_light_height.GetRevisionHeight()) - 1,
				Data:  key,
				Prove: true,
			}

			res, err := cc.QueryABCI(ctx, req)
			assert.NilError(t, err)
			if len(res.Value) == 0 {
				//	continue
			}
			merkleProof, err := commitmenttypes.ConvertProofs(res.ProofOps)
			assert.NilError(t, err)
			proofs = append(proofs, merkleProof.Proofs[0])
			didRoot = merkleProof.Proofs[1].GetExist().Value
			didProof = merkleProof.Proofs[1]
			items[string(key)] = res.Value
		}
		return
	}
	didRoot, didProof, items, proofs := getProofs()
	batchProofs, err := ics23.CombineProofs(proofs)
	assert.NilError(t, err)
	pass := ics23.BatchVerifyMembership(mehub_light_state_r.ProofSpecs[0], didRoot, batchProofs, items)
	assert.Equal(t, pass, true)

	pass = ics23.VerifyMembership(mehub_light_state_r.ProofSpecs[1], csState.Root.GetHash(), didProof, []byte("did"), didRoot)
	assert.Equal(t, pass, true)

	fmt.Println(batchProofs.Size())
	aa, _ := json.Marshal(batchProofs)

	ioutil.WriteFile("testData", aa, 0666)

	fmt.Println("tx command")
	var itemList []item
	for k, v := range items {
		itemList = append(itemList, item{Key: []byte(k), Value: v})
	}
	itemsStr, _ := json.Marshal(itemList)
	storeHashHex := hex.EncodeToString(didRoot)
	fmt.Println(didProof.String())
	didProofBz, _ := didProof.Marshal()
	storeProofHex := hex.EncodeToString(didProofBz)
	proofsBz, _ := batchProofs.Marshal()
	proofsHex := hex.EncodeToString(proofsBz)
	fmt.Printf("tx kyc update-credential %d '%s' %s %s %s\n", int64(mehub_light_height.GetRevisionHeight()), string(itemsStr), storeHashHex, storeProofHex, proofsHex)
}

type item struct {
	Key   []byte
	Value []byte
}

func TestVerifyNotMemberShip(t *testing.T) {
	ctx := context.Background()
	rollapp, mehub := initClient(t, ctx)

	//
	var rollappReqHeight int64 = 43
	mehub_light_state, err := rollapp.QueryClientStateResponse(ctx, rollappReqHeight, "07-tendermint-0")
	assert.NilError(t, err)
	fmt.Println(mehub_light_state.String())

	//查询rollapp上mehub的轻客户端的共识状态
	mehub_light_state_r, ok := mehub_light_state.ClientState.GetCachedValue().(*tmclient.ClientState)
	assert.Assert(t, ok)
	mehub_light_height := mehub_light_state_r.GetLatestHeight()
	meConStateRes, err := rollapp.QueryClientConsensusState(ctx, rollappReqHeight, "07-tendermint-0", mehub_light_height)
	assert.NilError(t, err)
	csState, ok := meConStateRes.ConsensusState.GetCachedValue().(*tmclient.ConsensusState)
	assert.Assert(t, ok)

	getProofs := func() (didRoot []byte, didProof *ics23.CommitmentProof, items [][]byte, proofs []*ics23.CommitmentProof) {

		prefix := []byte{0x40}
		items = make([][]byte, 0)
		for i := 0; i < 10; i++ {
			key := append(prefix, []byte(fmt.Sprintf("%016d", i))...)
			key = append(key, []byte("kyc")...)
			//查询mehub did key store
			cc, ok := mehub.(*cosmos.CosmosProvider)
			if !ok {
				t.Fatal("not a cosmos provider")
			}
			req := abci.RequestQuery{
				Path: fmt.Sprintf("store/%s/key", "did"),
				//Height: 281512,
				Height: int64(mehub_light_height.GetRevisionHeight()) - 1,
				Data:   key,
				Prove:  true,
			}

			res, err := cc.QueryABCI(ctx, req)
			assert.NilError(t, err)
			if len(res.Value) == 0 {
				// not member ship
				merkleProof, err := commitmenttypes.ConvertProofs(res.ProofOps)
				assert.NilError(t, err)
				proofs = append(proofs, merkleProof.Proofs[0])
				didRoot = merkleProof.Proofs[1].GetExist().Value
				didProof = merkleProof.Proofs[1]
				items = append(items, key)
			}
		}
		return
	}
	didRoot, didProof, items, proofs := getProofs()
	batchProofs, err := ics23.CombineProofs(proofs)
	assert.NilError(t, err)

	// clientState, err := mehub.QueryClientState(ctx, 281512, "09-localhost")
	// assert.NilError(t, err)
	// clientS := clientState.(*tmclient.ClientState)
	pass := ics23.BatchVerifyNonMembership(mehub_light_state_r.ProofSpecs[0], didRoot, batchProofs, items)
	assert.Equal(t, pass, true)


	pass = ics23.VerifyMembership(mehub_light_state_r.ProofSpecs[1], csState.Root.GetHash(), didProof, []byte("did"), didRoot)
	assert.Equal(t, pass, true)
}
