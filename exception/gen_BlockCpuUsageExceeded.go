// Code generated by gotemplate. DO NOT EDIT.

package exception

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/eosspark/eos-go/log"
)

// template type Exception(PARENT,CODE,WHAT)

var BlockCpuUsageExceededName = reflect.TypeOf(BlockCpuUsageExceeded{}).Name()

type BlockCpuUsageExceeded struct {
	_ResourceExhaustedException
	Elog log.Messages
}

func NewBlockCpuUsageExceeded(parent _ResourceExhaustedException, message log.Message) *BlockCpuUsageExceeded {
	return &BlockCpuUsageExceeded{parent, log.Messages{message}}
}

func (e BlockCpuUsageExceeded) Code() int64 {
	return 3080005
}

func (e BlockCpuUsageExceeded) Name() string {
	return BlockCpuUsageExceededName
}

func (e BlockCpuUsageExceeded) What() string {
	return "Transaction CPU usage is too much for the remaining allowable usage of the current block"
}

func (e *BlockCpuUsageExceeded) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e BlockCpuUsageExceeded) GetLog() log.Messages {
	return e.Elog
}

func (e BlockCpuUsageExceeded) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e BlockCpuUsageExceeded) DetailMessage() string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(int(e.Code())))
	buffer.WriteByte(' ')
	buffer.WriteString(e.Name())
	buffer.Write([]byte{':', ' '})
	buffer.WriteString(e.What())
	buffer.WriteByte('\n')
	for _, l := range e.Elog {
		buffer.WriteByte('[')
		buffer.WriteString(l.GetMessage())
		buffer.Write([]byte{']', ' '})
		buffer.WriteString(l.GetContext().String())
		buffer.WriteByte('\n')
	}
	return buffer.String()
}

func (e BlockCpuUsageExceeded) String() string {
	return e.DetailMessage()
}

func (e BlockCpuUsageExceeded) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3080005,
		Name: BlockCpuUsageExceededName,
		What: "Transaction CPU usage is too much for the remaining allowable usage of the current block",
	}

	return json.Marshal(except)
}

func (e BlockCpuUsageExceeded) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*BlockCpuUsageExceeded):
		callback(&e)
		return true
	case func(BlockCpuUsageExceeded):
		callback(e)
		return true
	default:
		return false
	}
}
