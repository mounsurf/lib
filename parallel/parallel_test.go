package parallel

import (
	"fmt"
	"github.com/mounsurf/lib/util"
	"testing"
	"time"
)

func handler(item *Item) {
	fmt.Println("进行处理", item)
	time.Sleep(time.Second)
	fmt.Println("处理完毕", item)
}

func TestRun1(t *testing.T) {
	dataSet := []int{1, 2, 3, 4, 5, 6}
	err := Run(handler, 2, dataSet)
	util.CheckError(err)
}

func TestRun2(t *testing.T) {
	dataSet := []string{"1", "2", "3", "4", "5", "6"}
	go func() {
		for {
			dataSet = append(dataSet, util.RandStr(2, "123456789"))
			fmt.Println(dataSet)
			time.Sleep(time.Second / 2)
		}
	}()
	err := Run(handler, 2, dataSet)
	util.CheckError(err)
}
func TestRun3(t *testing.T) {
	dataSet := []string{"1", "2", "3", "4", "5", "6"}
	go func() {
		// 这里存在如果写的不及时，就存在导致认为已经读完了的问题。不推荐使用。
		for i := 0; i < 10; i++ {
			dataSet = append(dataSet, util.RandStr(2, "123456789"))
			fmt.Println(dataSet)
			time.Sleep(time.Second / 2)
		}
	}()
	err := Run(handler, 10, &dataSet)
	util.CheckError(err)
}

func TestRun4(t *testing.T) {
	dataSet := []int{1, 2, 3, 4, 5, 6}
	go func() {
		for i := 0; i < 10; i++ {
			dataSet = append(dataSet, 111222)
			fmt.Println(dataSet)
			time.Sleep(time.Second / 2)
		}
	}()
	err := Run(handler, 2, &dataSet)
	util.CheckError(err)
}

func TestRun5(t *testing.T) {
	dataSet := make(chan interface{})
	go func() {
		for i := 0; i < 10; i++ {
			dataSet <- i
		}
		close(dataSet)
	}()
	err := Run(handler, 2, dataSet)
	util.CheckError(err)
}

func TestRun6(t *testing.T) {
	dataSet := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			dataSet <- i
		}
		close(dataSet)

	}()
	err := Run(handler, 2, dataSet)
	util.CheckError(err)
}
