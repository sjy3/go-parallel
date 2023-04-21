/**
 * @Author: Shi Jinyu
 * @Description:
 * @File:  parallel_test
 * @Version: 1.0.0
 * @Date: 2023/4/21 21:58
 */
package go_parallel

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

type process struct {
	Name int
}

func (p process) Do() interface{} {
	time.Sleep(1000 * time.Millisecond)
	return p.Name
}

func TestParallelProcess_Run(t *testing.T) {
	po := NewParallelObject()

	for i := 0; i < 100; i++ {
		j := i
		p := process{
			j,
		}
		po.AppendProcess(strconv.FormatInt(int64(j), 10), p)
	}

	po.SetTimeout(5000 * time.Millisecond)

	fmt.Println(po)

	result, isTimeout := po.Run()
	fmt.Println("isTimeout = ", isTimeout)
	fmt.Println(result)
	fmt.Println(len(result))

}

func TestParallelFunc_Run(t *testing.T) {
	po := NewParallelObject()

	for i := 0; i < 100; i++ {
		j := i
		po.AppendFunc(strconv.FormatInt(int64(j), 10), func() interface{} {
			time.Sleep(1000 * time.Millisecond)
			return j
		})
	}

	po.SetTimeout(5000 * time.Millisecond)

	fmt.Println(po)

	result, isTimeout := po.Run()
	fmt.Println("isTimeout = ", isTimeout)
	fmt.Println(result)
	fmt.Println(len(result))

}

func TestParallelResult(t *testing.T) {
	po := NewParallelObject()

	for i := 0; i < 100; i++ {
		j := i
		po.AppendFunc(strconv.FormatInt(int64(j), 10), func() interface{} {
			time.Sleep(1000 * time.Millisecond)
			return j
		})
	}

	po.SetTimeout(5000 * time.Millisecond)

	fmt.Println(po)

	result, isTimeout := po.Run()
	fmt.Println("isTimeout = ", isTimeout)
	fmt.Println(result)
	fmt.Println(len(result))

}
