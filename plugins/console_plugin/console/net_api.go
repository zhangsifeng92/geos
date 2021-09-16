package console

import (
	"github.com/robertkrimen/otto"
	"github.com/zhangsifeng92/geos/common"
	"github.com/zhangsifeng92/geos/plugins/net_plugin"
)

//NetAPI interacts with local p2p network connections
type NetAPI struct {
	c       *Console
	baseUrl string
}

func newNetAPI(c *Console) *NetAPI {
	n := &NetAPI{c: c}
	return n
}

//Connect starts a new connection to a peer
func (n *NetAPI) Connect(call otto.FunctionCall) (response otto.Value) {
	host, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	var connectInfo string
	if err := DoHttpCall(&connectInfo, common.NetConnect, host); err != nil {
		return getJsResult(call, err.Error())
	}
	return getJsResult(call, connectInfo)
}

//Disconnect closes an existing connection
func (n *NetAPI) Disconnect(call otto.FunctionCall) (response otto.Value) {
	host, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	var result string
	if err = DoHttpCall(&result, common.NetDisconnect, host); err != nil {
		return getJsResult(call, err.Error())
	}
	return getJsResult(call, result)
}

//Status status of existing connection
func (n *NetAPI) Status(call otto.FunctionCall) (response otto.Value) {
	host, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	var result net_plugin.PeerStatus
	if err = DoHttpCall(&result, common.NetStatus, host); err != nil {
		return getJsResult(call, err.Error())
	}
	return getJsResult(call, result)
}

//Peers status of exiting connection
func (n *NetAPI) Peers(call otto.FunctionCall) (response otto.Value) {
	var result []net_plugin.PeerStatus
	if err := DoHttpCall(&result, common.NetConnections, nil); err != nil {
		return getJsResult(call, err.Error())
	}
	return getJsResult(call, result)
}
