package producer_plugin

import (
	"fmt"
	"github.com/zhangsifeng92/geos/chain/types"
	. "github.com/zhangsifeng92/geos/chain/types/generated_containers"
	. "github.com/zhangsifeng92/geos/plugins/chain_interface"
	"github.com/zhangsifeng92/geos/plugins/chain_plugin"

	//Chain "github.com/zhangsifeng92/geos/plugins/producer_plugin/testing" /*test model*/
	"encoding/json"
	"github.com/urfave/cli"
	Chain "github.com/zhangsifeng92/geos/chain" /*real chain*/
	"github.com/zhangsifeng92/geos/common"
	"github.com/zhangsifeng92/geos/crypto"
	"github.com/zhangsifeng92/geos/crypto/ecc"
	. "github.com/zhangsifeng92/geos/exception"
	. "github.com/zhangsifeng92/geos/exception/try"
	"github.com/zhangsifeng92/geos/libraries/asio"
	"github.com/zhangsifeng92/geos/log"
	. "github.com/zhangsifeng92/geos/plugins/appbase/app"
	"strings"
	"time"
)

const ProducerPlug = PluginTypeName("ProducerPlugin")

var (
	producerPlugin = App().RegisterPlugin(ProducerPlug, NewProducerPlugin(App().GetIoService()))
	ppLog          = log.GetLoggerMap()["producer_plugin"]
	trxTraceLog    = log.GetLoggerMap()["transaction_tracing"]
)

type ProducerPlugin struct {
	AbstractPlugin
	//ConfirmedBlock Signal //TODO signal ConfirmedBlock
	my *ProducerPluginImpl
}

type RuntimeOptions struct {
	MaxTransactionTime      *int32
	MaxIrreversibleBlockAge *int32
	ProduceTimeOffsetUs     *int32
	LastBlockTimeOffsetUs   *int32
	SubjectiveCpuLeewayUs   *int32
	IncomingDeferRadio      *float64
}

type WhitelistAndBlacklist struct {
	ActorWhitelist    *AccountNameSet
	ActorBlacklist    *AccountNameSet
	ContractWhitelist *AccountNameSet
	ContractBlacklist *AccountNameSet
	ActionBlacklist   *NamePairSet
	KeyBlacklist      *PublicKeySet
}

type WhitelistAndBlacklistParam struct {
	ActorWhitelist    []common.AccountName
	ActorBlacklist    []common.AccountName
	ContractWhitelist []common.AccountName
	ContractBlacklist []common.AccountName
	ActionBlacklist   []common.NamePair
	KeyBlacklist      []string
}

type GreylistParams struct {
	Accounts []common.AccountName
}

func NewProducerPlugin(io *asio.IoContext) *ProducerPlugin {
	plugin := &ProducerPlugin{}

	plugin.my = NewProducerPluginImpl(io)
	plugin.my.Self = plugin

	return plugin
}

//func (p *ProducerPlugin) chain() *Chain.Controller {
//	//return App().FindPlugin(chain_plugin.ChainPlug).(*chain_plugin.ChainPlugin).Chain()
//	return Chain.GetControllerInstance()
//}

func (p *ProducerPlugin) SetProgramOptions(options *[]cli.Flag) {
	*options = append(*options,
		cli.BoolFlag{
			Name:  "enable-stale-production, e",
			Usage: "Enable block production, even if the chain is stale.",
		},
		cli.BoolFlag{
			Name:  "pause-on-startup, x",
			Usage: "Start this node in a state where production is paused.",
		},
		cli.IntFlag{
			Name:  "max-transaction-age",
			Usage: "Limits the maximum time (in milliseconds) that is allowed a pushed transaction's code to execute before being considered invalid",
			Value: 30,
		},
		cli.IntFlag{
			Name:  "max-irreversible-block-age",
			Usage: "Limits the maximum age (in seconds) of the DPOS Irreversible Block for a chain this node will produce blocks on (use negative value to indicate unlimited)",
			Value: -1,
		},
		cli.StringSliceFlag{
			Name:  "producer-name, p",
			Usage: "ID of producer controlled by this node(e.g. inita; may specify multiple times)",
		},
		cli.StringSliceFlag{
			Name:  "private-key",
			Usage: "(DEPRECATED - Use signature-provider instead) Tuple of [public key, WIF private key] (may specify multiple times)",
		},
		cli.StringSliceFlag{
			Name: "signature-provider",
			Usage: "Key=Value pairs in the form <public-key>=<provider-spec>\n" +
				"Where:\n" +
				"   <public-key>    \tis a string form of a vaild EOSIO public key\n\n" +
				"   <provider-spec> \tis a string in the form <provider-type>:<data>\n\n" +
				"   <provider-type> \tis KEY, or KEOSD\n\n" +
				"   KEY:<data>      \tis a string form of a valid EOSIO private key which maps to the provided public key\n\n",
		},
		cli.IntFlag{
			Name:  "keosd-provider-timeout",
			Usage: "Limits the maximum time (in milliseconds) that is allowd for sending blocks to a keosd provider for signing",
			Value: 5,
		},
		cli.StringSliceFlag{
			Name:  "greylist-account",
			Usage: "account that can not access to extended CPU/NET virtual resources",
		},
		cli.IntFlag{
			Name:  "produce-time-offset-us",
			Usage: "offset of non last block producing time in micro second. Negative number results in blocks to go out sooner, and positive number results in blocks to go out later",
			Value: 0,
		},
		cli.IntFlag{
			Name:  "last-block-time-offset-us",
			Usage: "offset of last block producing time in micro second. Negative number results in blocks to go out sooner, and positive number results in blocks to go out later",
			Value: 0,
		},
		cli.Float64Flag{
			Name:  "incoming-defer-ratio",
			Usage: "ratio between incoming transations and deferred transactions when both are exhausted",
			Value: 1.0,
		},
	)
}

func (p *ProducerPlugin) PluginInitialize(c *cli.Context) {
	Try(func() {
		//app.Action = func(c *cli.Context) {
		for _, n := range c.StringSlice("producer-name") {
			name := common.N(n)
			p.my.Producers.Add(name)
		}

		for _, keyIdToWifPairString := range c.StringSlice("private-key") {
			Try(func() {
				var keyIdToWifPair [2]string
				json.Unmarshal([]byte(keyIdToWifPairString), &keyIdToWifPair)
				pubKey, err := ecc.NewPublicKey(keyIdToWifPair[0])
				if err != nil {
					Throw(err)
				}
				priKey, err2 := ecc.NewPrivateKey(keyIdToWifPair[1])
				if err2 != nil {
					Throw(err2)
				}
				p.my.SignatureProviders[pubKey] = makeKeySignatureProvider(priKey)
				// wlog("\"private-key\" is DEPRECATED, use \"signature-provider=${pub}=KEY:${priv}\"", ("pub",key_id_to_wif_pair.first)("priv", blanked_privkey));
			}).Catch(func(e error) {
				log.Error("Malformed private key pair: \"%s\", ignoring", keyIdToWifPairString)
			}).End()
		}

		for _, keySpecPair := range c.StringSlice("signature-provider") {
			Try(func() {
				delim := strings.Index(keySpecPair, "=")
				EosAssert(delim >= 0, &PluginConfigException{}, "Missing \"=\" in the key spec pair")
				pubKeyStr := keySpecPair[0:delim]
				specStr := keySpecPair[delim+1:]

				specDelim := strings.Index(specStr, ":")
				EosAssert(specDelim >= 0, &PluginConfigException{}, "Missing \":\" in the key spec pair")
				specTypeStr := specStr[0:specDelim]
				specData := specStr[specDelim+1:]

				pubKey, err := ecc.NewPublicKey(pubKeyStr)
				if err != nil {
					Throw(err)
				}

				if specTypeStr == "KEY" {
					priKey, e := ecc.NewPrivateKey(specData)
					if e != nil {
						Throw(e)
					}
					p.my.SignatureProviders[pubKey] = makeKeySignatureProvider(priKey)
				} else if specTypeStr == "KEOSD" {
					p.my.SignatureProviders[pubKey] = makeKeosdSignatureProvider(p.my, specData, pubKey)
				}

			}).Catch(func(interface{}) {
				log.Error("Malformed signature provider: \"%s\", ignoring!", keySpecPair)
			}).End()
		}

		p.my.ProductionEnabled = c.Bool("enable-stale-production")

		p.my.ProductionPaused = c.Bool("pause-on-startup")

		p.my.KeosdProviderTimeoutUs = common.Milliseconds(int64(c.Int("keosd-provider-timeout")))

		p.my.ProduceTimeOffsetUs = int32(c.Int("produce-time-offset-us"))

		p.my.MaxTransactionTimeMs = int32(c.Int("max-transaction-age"))

		p.my.MaxIrreversibleBlockAgeUs = common.Seconds(int64(c.Int("max-irreversible-block-age")))

		p.my.IncomingDeferRadio = c.Float64("incoming-defer-ratio")

		if greylist := c.StringSlice("greylist-account"); len(greylist) > 0 {
			param := GreylistParams{}
			for _, a := range greylist {
				param.Accounts = append(param.Accounts, common.N(a))
			}
			p.AddGreylistAccounts(param)
		}

		App().GetMethod(BlockSync).Register(&BlockSyncCaller{Caller: p.my.OnIncomingBlock})
		App().GetMethod(TransactionAsync).Register(&TransactionAsyncCaller{Caller: p.my.OnIncomingTransactionAsync})

	}).FcLogAndRethrow().End()
}

func (p *ProducerPlugin) PluginStartup() {
	Try(func() {
		log.Info("producer plugin:  plugin_startup() begin")

		p.my.Chain = App().GetPlugin(chain_plugin.ChainPlug).(*chain_plugin.ChainPlugin).Chain()
		chain := p.my.Chain

		EosAssert(p.my.Producers.Size() == 0 || chain.GetReadMode() == Chain.DBReadMode(Chain.SPECULATIVE), &PluginConfigException{},
			"node cannot have any producer-name configured because block production is impossible when read_mode is not \"speculative\"")

		EosAssert(p.my.Producers.Size() == 0 || chain.GetValidationMode() == Chain.ValidationMode(Chain.FULL), &PluginConfigException{},
			"node cannot have any producer-name configured because block production is not safe when validation_mode is not \"full\"")

		chain.AcceptedBlock.Connect(&AcceptedBlockCaller{Caller: func(bsp *types.BlockState) { p.my.OnBlock(bsp) }})
		chain.IrreversibleBlock.Connect(&IrreversibleBlockCaller{Caller: func(bsp *types.BlockState) { p.my.OnIrreversibleBlock(bsp.SignedBlock) }})

		libNum := chain.LastIrreversibleBlockNum()
		lib := chain.FetchBlockByNumber(libNum)
		if lib != nil {
			p.my.OnIrreversibleBlock(lib)
		} else {
			p.my.IrreversibleBlockTime = common.MaxTimePoint()
		}

		if p.my.Producers.Size() > 0 {
			log.Info("Launching block production for %d producers at %s.", p.my.Producers.Size(), common.Now())

			if p.my.ProductionEnabled {
				if chain.HeadBlockNum() == 0 {
					newChainBanner(chain)
				}
			}
		}

		p.my.ScheduleProductionLoop()

		log.Info("producer plugin:  plugin_startup() end")

	}).FcCaptureAndRethrow().End()
}

func (p *ProducerPlugin) PluginShutdown() {
	p.my.Timer.Cancel()
	log.Info("producer plugin shutdown")
}

func (p *ProducerPlugin) IsProducerKey(key ecc.PublicKey) bool {
	privateKey := p.my.SignatureProviders[key]
	if privateKey != nil {
		return true
	}
	return false
}

func (p *ProducerPlugin) SignCompact(key *ecc.PublicKey, digest crypto.Sha256) *ecc.Signature {
	if !key.Compare(*ecc.NewPublicKeyNil()) {
		privateKeyFunc := p.my.SignatureProviders[*key]
		EosAssert(privateKeyFunc != nil, &ProducerPrivKeyNotFound{}, "Local producer has no private key in config.ini corresponding to public key %s", key)

		return privateKeyFunc(digest)
	}
	return ecc.NewSigNil()
}

func (p *ProducerPlugin) Pause() {
	p.my.ProductionPaused = true
}

func (p *ProducerPlugin) Resume() {
	p.my.ProductionPaused = false
	// it is possible that we are only speculating because of this policy which we have now changed
	// re-evaluate that now
	//
	if p.my.PendingBlockMode == PendingBlockMode(speculating) {
		chain := p.my.Chain
		chain.AbortBlock()
		p.my.ScheduleProductionLoop()
	}
}

func (p *ProducerPlugin) Paused() bool {
	return p.my.ProductionPaused
}

func (p *ProducerPlugin) UpdateRuntimeOptions(options RuntimeOptions) {
	checkSpeculation := false

	if options.MaxTransactionTime != nil {
		p.my.MaxTransactionTimeMs = *options.MaxTransactionTime
	}

	if options.MaxIrreversibleBlockAge != nil {
		p.my.MaxIrreversibleBlockAgeUs = common.Seconds(int64(*options.MaxIrreversibleBlockAge))
		checkSpeculation = true
	}

	if options.ProduceTimeOffsetUs != nil {
		p.my.ProduceTimeOffsetUs = *options.ProduceTimeOffsetUs
	}

	if options.LastBlockTimeOffsetUs != nil {
		p.my.LastBlockTimeOffsetUs = *options.LastBlockTimeOffsetUs
	}

	if options.IncomingDeferRadio != nil {
		p.my.IncomingDeferRadio = *options.IncomingDeferRadio
	}

	if checkSpeculation && p.my.PendingBlockMode == PendingBlockMode(speculating) {
		//chain := p.my.Chain
		p.my.Chain.AbortBlock()
		p.my.ScheduleProductionLoop()
	}

	if options.SubjectiveCpuLeewayUs != nil {
		//chain := p.my.Chain
		p.my.Chain.SetSubjectiveCpuLeeway(common.Microseconds(*options.SubjectiveCpuLeewayUs))
	}
}

func (p *ProducerPlugin) GetRuntimeOptions() RuntimeOptions {
	var maxIrreversibleBlockAge int32 = -1
	if p.my.MaxIrreversibleBlockAgeUs.Count() >= 0 {
		maxIrreversibleBlockAge = int32(p.my.MaxIrreversibleBlockAgeUs.Count() / 1e6)
	}
	return RuntimeOptions{
		&p.my.MaxTransactionTimeMs,
		&maxIrreversibleBlockAge,
		&p.my.ProduceTimeOffsetUs,
		&p.my.LastBlockTimeOffsetUs,
		nil, nil,
	}
}

func (p *ProducerPlugin) AddGreylistAccounts(params GreylistParams) {
	for _, acc := range params.Accounts {
		p.my.Chain.AddResourceGreyList(acc)
	}
}

func (p *ProducerPlugin) RemoveGreylistAccounts(params GreylistParams) {
	for _, acc := range params.Accounts {
		p.my.Chain.RemoveResourceGreyList(acc)
	}
}

func (p *ProducerPlugin) GetGreylist() GreylistParams {
	result := GreylistParams{}
	list := p.my.Chain.GetResourceGreyList()
	result.Accounts = make([]common.AccountName, list.Size())

	list.Each(func(acc common.AccountName) {
		result.Accounts = append(result.Accounts, acc)
	})
	return result
}

func (p *ProducerPlugin) GetWhitelistBlacklist() WhitelistAndBlacklist {
	chain := p.my.Chain
	return WhitelistAndBlacklist{
		chain.GetActorWhiteList(),
		chain.GetActorBlackList(),
		chain.GetContractWhiteList(),
		chain.GetContractBlackList(),
		chain.GetActionBlackList(),
		chain.GetKeyBlackList(),
	}
}

func (p *ProducerPlugin) SetWhitelistBlacklist(params WhitelistAndBlacklist) {
	chain := p.my.Chain
	if params.ActorWhitelist != nil {
		chain.SetActorWhiteList(params.ActorWhitelist)
	}
	if params.ActorBlacklist != nil {
		chain.SetActorBlackList(params.ActorBlacklist)
	}
	if params.ContractWhitelist != nil {
		chain.SetContractWhiteList(params.ContractWhitelist)
	}
	if params.ContractBlacklist != nil {
		chain.SetContractBlackList(params.ContractBlacklist)
	}
	if params.ActionBlacklist != nil {
		chain.SetActionBlackList(params.ActionBlacklist)
	}
	if params.KeyBlacklist != nil {
		chain.SetKeyBlackList(params.KeyBlacklist)
	}
}

func failureIsSubjective(e Exception, deadlineIsSubjective bool) bool {
	code := e.Code()
	return (code == BlockCpuUsageExceeded{}.Code()) ||
		(code == BlockCpuUsageExceeded{}.Code()) ||
		(code == DeadlineException{}.Code() && deadlineIsSubjective)
}

func makeDebugTimeLogger() func() {
	start := time.Now()
	return func() {
		fmt.Println(time.Now().Sub(start))
	}
}

func makeKeySignatureProvider(key *ecc.PrivateKey) signatureProviderType {
	signFunc := func(digest crypto.Sha256) *ecc.Signature {
		sign, err := key.Sign(digest.Bytes())
		if err != nil {
			Throw(err)
		}
		return &sign
	}
	return signFunc
}

func makeKeosdSignatureProvider(produce *ProducerPluginImpl, url string, publicKey ecc.PublicKey) signatureProviderType {
	signFunc := func(digest crypto.Sha256) *ecc.Signature {
		if produce != nil {
			//TODO
			return ecc.NewSigNil()
		} else {
			return ecc.NewSigNil()
		}
	}
	return signFunc
}

func newChainBanner(db *Chain.Controller) {
	fmt.Print("\n" +
		"*******************************\n" +
		"*                             *\n" +
		"*   ------ NEW CHAIN ------   *\n" +
		"*   -  Welcome to EOSIO!  -   *\n" +
		"*   -----------------------   *\n" +
		"*                             *\n" +
		"*******************************\n" +
		"\n")

	if db.HeadBlockState().Header.Timestamp.ToTimePoint() < common.Now().SubUs(common.Microseconds(200*common.DefaultConfig.BlockIntervalMs)) {
		fmt.Print("Your genesis seems to have an old timestamp\n" +
			"Please consider using the --genesis-timestamp option to give your genesis a recent timestamp\n\n")
	}
}
