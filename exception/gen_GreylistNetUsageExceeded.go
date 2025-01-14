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

var GreylistNetUsageExceededName = reflect.TypeOf(GreylistNetUsageExceeded{}).Name()

type GreylistNetUsageExceeded struct {
	_ResourceExhaustedException
	Elog log.Messages
}

func NewGreylistNetUsageExceeded(parent _ResourceExhaustedException, message log.Message) *GreylistNetUsageExceeded {
	return &GreylistNetUsageExceeded{parent, log.Messages{message}}
}

func (e GreylistNetUsageExceeded) Code() int64 {
	return 3080007
}

func (e GreylistNetUsageExceeded) Name() string {
	return GreylistNetUsageExceededName
}

func (e GreylistNetUsageExceeded) What() string {
	return "Transaction exceeded the current greylisted account network usage limit"
}

func (e *GreylistNetUsageExceeded) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e GreylistNetUsageExceeded) GetLog() log.Messages {
	return e.Elog
}

func (e GreylistNetUsageExceeded) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e GreylistNetUsageExceeded) DetailMessage() string {
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

func (e GreylistNetUsageExceeded) String() string {
	return e.DetailMessage()
}

func (e GreylistNetUsageExceeded) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3080007,
		Name: GreylistNetUsageExceededName,
		What: "Transaction exceeded the current greylisted account network usage limit",
	}

	return json.Marshal(except)
}

func (e GreylistNetUsageExceeded) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*GreylistNetUsageExceeded):
		callback(&e)
		return true
	case func(GreylistNetUsageExceeded):
		callback(e)
		return true
	default:
		return false
	}
}
