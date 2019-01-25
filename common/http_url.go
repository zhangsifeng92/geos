package common

const (
	ChainUrl = "http://127.0.0.1:8888"
	// walletUrl = "http://127.0.0.1:8900"
	WalletUrl = "http://127.0.0.1:8888"
	// walletUrl = "http://127.0.0.1:8765"
	HttpEndPoint = "http://127.0.0.1:8888"
)

const (
	chainFuncBase           string = "/v1/chain"
	GetInfoFunc             string = chainFuncBase + "/get_info"
	PushTxnFunc             string = chainFuncBase + "/push_transaction"
	PushTxnsFunc            string = chainFuncBase + "/push_transactions"
	JsonToBinFunc           string = chainFuncBase + "/abi_json_to_bin"
	GetBlockFunc            string = chainFuncBase + "/get_block"
	GetBlockHeaderStateFunc string = chainFuncBase + "/get_block_header_state"
	GetAccountFunc          string = chainFuncBase + "/get_account"
	GetTableFunc            string = chainFuncBase + "/get_table_rows"
	GetTableByScopeFunc     string = chainFuncBase + "/get_table_by_scope"
	GetCodeFunc             string = chainFuncBase + "/get_code"
	GetCodeHashFunc         string = chainFuncBase + "/get_code_hash"
	GetAbiFunc              string = chainFuncBase + "/get_abi"
	GetRawAbiFunc           string = chainFuncBase + "/get_raw_abi"
	GetRawCodeAndAbiFunc    string = chainFuncBase + "/get_raw_code_and_abi"
	GetCurrencyBalanceFunc  string = chainFuncBase + "/get_currency_balance"
	GetCurrencyStatsFunc    string = chainFuncBase + "/get_currency_stats"
	GetProducersFunc        string = chainFuncBase + "/get_producers"
	GetScheduleFunc         string = chainFuncBase + "/get_producer_schedule"
	GetRequiredKeys         string = chainFuncBase + "/get_required_keys"

	historyFuncBase           string = "/v1/history"
	GetActionsFunc            string = historyFuncBase + "/get_actions"
	GetTransactionFunc        string = historyFuncBase + "/get_transaction"
	GetKeyAccountsFunc        string = historyFuncBase + "/get_key_accounts"
	GetControlledAccountsFunc string = historyFuncBase + "/get_controlled_accounts"

	accountHistoryFuncBase string = "/v1/account_history"
	GetTransactionsFunc    string = accountHistoryFuncBase + "/get_transactions"

	netFuncBase    string = "/v1/net"
	NetConnect     string = netFuncBase + "/connect"
	NetDisconnect  string = netFuncBase + "/disconnect"
	NetStatus      string = netFuncBase + "/status"
	NetConnections string = netFuncBase + "/connections"

	walletFuncBase   string = "/v1/wallet"
	WalletCreate     string = walletFuncBase + "/create"
	WalletOpen       string = walletFuncBase + "/open"
	WalletList       string = walletFuncBase + "/list_wallets"
	WalletListKeys   string = walletFuncBase + "/list_keys"
	WalletPublicKeys string = walletFuncBase + "/get_public_keys"
	WalletLock       string = walletFuncBase + "/lock"
	WalletLockAll    string = walletFuncBase + "/lock_all"
	WalletUnlock     string = walletFuncBase + "/unlock"
	WalletImportKey  string = walletFuncBase + "/import_key"
	WalletRemoveKey  string = walletFuncBase + "/remove_key"
	WalletCreateKey  string = walletFuncBase + "/create_key"
	WalletSignTrx    string = walletFuncBase + "/sign_transaction"

	// keosdStop string = "/v1/keosd/stop"

	producerFuncBase               string = "/v1/producer"
	ProducerPause                  string = producerFuncBase + "/pause"
	ProducerResume                 string = producerFuncBase + "/resume"
	ProducerPaused                 string = producerFuncBase + "/paused"
	ProducerGetRuntimeOptions      string = producerFuncBase + "/get_runtime_options"
	ProducerUpdateRuntimeOptions   string = producerFuncBase + "/update_runtime_options"
	ProducerAddGreylistAccounts    string = producerFuncBase + "/add_greylist_accounts"
	ProducerRemoveGreylistAccounts string = producerFuncBase + "/remove_greylist_accounts"
	ProducerGetGreylist            string = producerFuncBase + "/get_greylist"
	ProducerGetWhitelistBlacklist  string = producerFuncBase + "/get_whitelist_blacklist"
	ProducerSetWhitelistBlacklist  string = producerFuncBase + "/set_whitelist_blacklist"
	ProducerGetIntegrityHash       string = producerFuncBase + "/get_integrity_hash"
	ProducerCreateSnapshot         string = producerFuncBase + "/create_snapshot"
)
