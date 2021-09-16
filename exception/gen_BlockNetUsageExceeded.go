// Code generated by gotemplate. DO NOT EDIT.

package exception

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/zhangsifeng92/geos/log"
)

// template type Exception(PARENT,CODE,WHAT)

var BlockNetUsageExceededName = reflect.TypeOf(BlockNetUsageExceeded{}).Name()

type BlockNetUsageExceeded struct {
	_ResourceExhaustedException
	Elog log.Messages
}

func NewBlockNetUsageExceeded(parent _ResourceExhaustedException, message log.Message) *BlockNetUsageExceeded {
	return &BlockNetUsageExceeded{parent, log.Messages{message}}
}

func (e BlockNetUsageExceeded) Code() int64 {
	return 3080003
}

func (e BlockNetUsageExceeded) Name() string {
	return BlockNetUsageExceededName
}

func (e BlockNetUsageExceeded) What() string {
	return "Transaction network usage is too much for the remaining allowable usage of the current block"
}

func (e *BlockNetUsageExceeded) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e BlockNetUsageExceeded) GetLog() log.Messages {
	return e.Elog
}

func (e BlockNetUsageExceeded) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e BlockNetUsageExceeded) DetailMessage() string {
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

func (e BlockNetUsageExceeded) String() string {
	return e.DetailMessage()
}

func (e BlockNetUsageExceeded) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3080003,
		Name: BlockNetUsageExceededName,
		What: "Transaction network usage is too much for the remaining allowable usage of the current block",
	}

	return json.Marshal(except)
}

func (e BlockNetUsageExceeded) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*BlockNetUsageExceeded):
		callback(&e)
		return true
	case func(BlockNetUsageExceeded):
		callback(e)
		return true
	default:
		return false
	}
}
