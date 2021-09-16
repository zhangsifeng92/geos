wasmgo
=====

Fork from https://github.com/go-interpreter/wagon
Fork from https://github.com/perlin-network/life 

**NOTE:** `wasmgo` requires `Go >= 1.9.x`.

## examples

```
package main

import (
	"github.com/zhangsifeng92/geos/chain"
	"github.com/zhangsifeng92/geos/chain/types"
	"github.com/zhangsifeng92/geos/common"
	"github.com/zhangsifeng92/geos/crypto"
	"github.com/zhangsifeng92/geos/crypto/rlp"
	"github.com/zhangsifeng92/geos/wasmgo"
	"io/ioutil"
	"log"
)

func main() {

	name := "hello.wasm"
	code, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}

	wasmgo := wasmgo.NewWasmGo()
	param, _ := rlp.EncodeToBytes(common.N("walker")) //[]byte{0x00, 0x00, 0x00, 0x00, 0x5c, 0x05, 0xa3, 0xe1}
	applyContext := &chain.ApplyContext{
		Receiver: common.N("hello"),
		Act: &types.Action{
			Account: common.N("hello"),
			Name:    common.N("hi"),
			Data:    param,
		},
	}

	codeVersion := crypto.NewSha256Byte([]byte(code))
	wasmgo.Apply(codeVersion, code, applyContext)
	
	//"hello, walker"
	//fmt.Println(applyContext.PendingConsoleOutput)

}
```
