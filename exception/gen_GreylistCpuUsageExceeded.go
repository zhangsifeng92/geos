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

var GreylistCpuUsageExceededName = reflect.TypeOf(GreylistCpuUsageExceeded{}).Name()

type GreylistCpuUsageExceeded struct {
	_ResourceExhaustedException
	Elog log.Messages
}

func NewGreylistCpuUsageExceeded(parent _ResourceExhaustedException, message log.Message) *GreylistCpuUsageExceeded {
	return &GreylistCpuUsageExceeded{parent, log.Messages{message}}
}

func (e GreylistCpuUsageExceeded) Code() int64 {
	return 3080008
}

func (e GreylistCpuUsageExceeded) Name() string {
	return GreylistCpuUsageExceededName
}

func (e GreylistCpuUsageExceeded) What() string {
	return "Transaction exceeded the current greylisted account CPU usage limit"
}

func (e *GreylistCpuUsageExceeded) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e GreylistCpuUsageExceeded) GetLog() log.Messages {
	return e.Elog
}

func (e GreylistCpuUsageExceeded) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e GreylistCpuUsageExceeded) DetailMessage() string {
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

func (e GreylistCpuUsageExceeded) String() string {
	return e.DetailMessage()
}

func (e GreylistCpuUsageExceeded) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3080008,
		Name: GreylistCpuUsageExceededName,
		What: "Transaction exceeded the current greylisted account CPU usage limit",
	}

	return json.Marshal(except)
}

func (e GreylistCpuUsageExceeded) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*GreylistCpuUsageExceeded):
		callback(&e)
		return true
	case func(GreylistCpuUsageExceeded):
		callback(e)
		return true
	default:
		return false
	}
}
