package models

import (
	"testing"
	"time"
	"fmt"
)

func TestStatisticOfChannelSendContract(t *testing.T) {
	start := time.Now().AddDate(0, -3, 0).Unix()
	end := time.Now().Unix()
	result, err := StatisticOfChannelSendContract(start, end)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result {
		fmt.Println(v)
	}
}

func TestStatisticOfChannelCompleteContract(t *testing.T) {
	start := time.Now().AddDate(0, -3, 0).Unix()
	end := time.Now().Unix()
	result, err := StatisticOfChannelCompleteContract(start, end)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result {
		fmt.Println(v)
	}
}

func TestStatisticOfChannelVerify(t *testing.T) {
	start := time.Now().AddDate(0, -3, 0).Unix()
	end := time.Now().Unix()
	result, err := StatisticOfChannelVerify(start, end)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result {
		fmt.Println(v)
	}
}

func TestStatisticOfCpSendContract(t *testing.T) {
	start := time.Now().AddDate(0, -3, 0).Unix()
	end := time.Now().Unix()
	result, err := StatisticOfCpSendContract(start, end)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result {
		fmt.Println(v)
	}
}

func TestStatisticOfCpCompleteContract(t *testing.T) {
	start := time.Now().AddDate(0, -3, 0).Unix()
	end := time.Now().Unix()
	result, err := StatisticOfCpCompleteContract(start, end)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result {
		fmt.Println(v)
	}
}

func TestStatisticOfCpVerify(t *testing.T) {
	start := time.Now().AddDate(0, -3, 0).Unix()
	end := time.Now().Unix()
	result, err := StatisticOfCpVerify(start, end)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result {
		fmt.Println(v)
	}
}

func TestStatisticOfChannelPaid(t *testing.T) {
	start := time.Now().AddDate(0, -3, 0).Unix()
	end := time.Now().Unix()
	result, err := StatisticOfChannelPaid(start, end)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result {
		fmt.Println(v)
	}
}

func TestStatisticOfCpPaid(t *testing.T) {
	start := time.Now().AddDate(0, -3, 0).Unix()
	end := time.Now().Unix()
	result, err := StatisticOfCpPaid(start, end)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result {
		fmt.Println(v)
	}
}

func TestStatisticOfChannelAccess(t *testing.T) {
	start := time.Now().AddDate(0, -3, 0).Unix()
	end := time.Now().Unix()
	result, err := StatisticOfChannelAccess(start, end)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result {
		fmt.Println(v)
	}
}

func TestStatisticOfChannelCompany(t *testing.T) {
	start := time.Now().AddDate(0, -3, 0).Unix()
	end := time.Now().Unix()
	result, err := StatisticOfChannelCompany(start, end)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result {
		fmt.Println(v)
	}
}

func TestStatisticOfGameEvaluate(t *testing.T) {
	start := time.Now().AddDate(0, -3, 0).Unix()
	end := time.Now().Unix()
	result, err := StatisticOfGameEvaluate(start, end)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result {
		fmt.Println(v)
	}
}

func TestStatisticOfGameAccess(t *testing.T) {
	start := time.Now().AddDate(0, -3, 0).Unix()
	end := time.Now().Unix()
	result, err := StatisticOfGameAccess(start, end)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result {
		fmt.Println(v)
	}
}

func TestStatisticOfDistributionCompany(t *testing.T) {
	start := time.Now().AddDate(0, -3, 0).Unix()
	end := time.Now().Unix()
	result, err := StatisticOfDistributionCompany(start, end)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result {
		fmt.Println(v)
	}
}

func TestStatisticOfAccounting(t *testing.T) {
	start := time.Now().AddDate(0, -3, 0).Unix()
	end := time.Now().Unix()
	data, total := StatisticOfAccounting(start, end, 100, 0)
	fmt.Println("\ntotal:", total)
	for _, v := range data {
		fmt.Println(v)
	}
}

func TestStatisticOfFinancial(t *testing.T) {
	start := time.Now().AddDate(0, -3, 0).Unix()
	end := time.Now().Unix()
	data, total := StatisticOfFinancial(start, end, 100, 0)
	fmt.Println("\ntotal:", total)
	for _, v := range data {
		fmt.Println(v)
	}
}

func TestStatisticOfChannelTrade(t *testing.T) {
	start := time.Now().AddDate(0, -3, 0).Unix()
	end := time.Now().Unix()
	data, total := StatisticOfChannelTrade(start, end, 100, 0)
	fmt.Println("\ntotal:", total)
	for _, v := range data {
		fmt.Println(v)
	}
}

func TestStatisticOfCpTrade(t *testing.T) {
	start := time.Now().AddDate(0, -3, 0).Unix()
	end := time.Now().Unix()
	data, total := StatisticOfCpTrade(start, end, 100, 0)
	fmt.Println("\ntotal:", total)
	for _, v := range data {
		fmt.Println(v)
	}
}

func TestStatisticOfOperation(t *testing.T) {
	start := time.Now().AddDate(0, -3, 0).Unix()
	end := time.Now().Unix()
	data, total := StatisticOfOperation(start, end, 100, 0)
	fmt.Println("start:", start, "end:", end,"\ntotal:", total)
	for _, v := range data {
		fmt.Println(v)
	}
}

func TestTemp(t *testing.T) {
	//start := time.Now()
	//end := start.AddDate(1, 0, 0)
	//fmt.Println(!start.Before(start))

	tmp := map[interface{}]interface{}{}

	tmp[1] = "OK"

	fmt.Println("tmp[2]:", tmp[1])

}