package old_sys

import (
	"log"
	"testing"
	log2 "github.com/bysir-zl/bygo/log"
	"sync"
	"time"
	"runtime"
	"github.com/tealeg/xlsx"
	"fmt"
)

type ListTest struct {
	Id            int     `json:"id"`
	GameId        int     `json:"game_id"`
	Cp            string  `json:"cp"`
	StartTime     string  `json:"start_time"`
	EndTime       string  `json:"end_time"`
	Ratio         float64 `json:"ratio"`
	Rule          string  `json:"rule"`
	IsTrue        bool    `json:"isture"`
	UpdateTime    int     `json:"update_time"`
	Status        int     `json:"status"`
	IsDelete      bool    `json:"isdelete"`
	CalculateType int     `json:"calculate_type"`
	Interest      float64 `json:"interest"`
}

func TestGetList(t *testing.T) {
	s, err := GetLadder(213, "")
	if err != nil {
		t.Error(err)
	}
	log.Printf("%+v", s)
	return
}

func TestClearing(t *testing.T) {
	r, err := GetClearing(213, "", "2017-1")
	if err != nil {
		t.Error(err)
	}
	log2.Verbose("test", r)
	return
}

func TestAddRule(t *testing.T) {
	ls := []Ladder4Post{{
		StartTime:   "2017-01-01",
		EndTime:     "2018-12-31",
		Ratio:       50,
		SlottingFee: 5.8,
		Rule:        "10000<user<20000&1480212121<time<1481212121",
	}}
	err := UpdateOrAddLadderList(213, "", ls)
	if err != nil {
		t.Error(err)
	}
	log2.Verbose("test", ls)
	return
}

var ch = make(chan int, 2)
var wg *sync.WaitGroup

func TestSelect(t *testing.T) {
	//ch <- 1
	//wg = &sync.WaitGroup{}
	//wg.Add(1)
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for {
			select {
			case <-ticker.C:
				//if <-ch != 0 {
				//	defer println("hello")
				//}
				//defer wg.Done()
				println("hello")
				//runtime.Goexit()
			case <-ch:
				println("hello ch")
				runtime.Goexit()
				//default:
				//	println("world")
				//	runtime.Goexit()
			}
			//defer println("hello world")
		}
	}()
	//wg.Wait()
	<-ch
}

func TestTake(t *testing.T) {
	file, err := xlsx.OpenFile("../../docs/kpi.xlsx")
	fmt.Println(file, err)
}

func Hello() {
	println(<-ch)
}
