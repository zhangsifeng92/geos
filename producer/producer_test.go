package producer_plugin

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	var apply = false
	var timer = new(scheduleTimer)
	var blockNum = 1

	//sigs := make(chan os.Signal, 1)
	//signal.Notify(sigs, syscall.SIGINT)

	var scheduleProductionLoop func()

	scheduleProductionLoop = func() {
		timer.cancel()
		base := time.Now()
		minTimeToNextBlock := int64(common.BlockIntervalUs) - base.UnixNano()/1e3%int64(common.BlockIntervalUs)
		wakeTime := base.Add(time.Microsecond * time.Duration(minTimeToNextBlock))

		timer.expiresUntil(wakeTime)

		// test after 12 block need to apply new block to continue
		if blockNum%12 == 0 {
			apply = true
			return
		}

		timerCorelationId++
		cid := timerCorelationId
		timer.asyncWait(func() bool { return cid == timerCorelationId }, func() {
			fmt.Println("exec async1...", time.Now())
			fmt.Println("add.blockNum", blockNum)
			blockNum++

			scheduleProductionLoop()
		})
	}

	applyBlock := func() {
		for {
			if apply {
				apply = false
				blockNum++
				fmt.Println("exec apply...", time.Now(), "\n-----------add.", blockNum)
				scheduleProductionLoop()
			}
		}
	}

	naughty := func() {
		for {
			time.Sleep(666 * time.Millisecond)
			scheduleProductionLoop()
		}
	}

	//go func() {
	//	sig := <-sigs
	//	fmt.Println("sig: ", sig)
	//}()

	scheduleProductionLoop()
	applyBlock()
	naughty() //try to break the schedule timer
}
