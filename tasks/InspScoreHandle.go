package tasks

import (
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
	"horizon/model"
	"strconv"
)

type Score struct {
	inspection model.Inspection
	instance   model.Instance
}

func (s *Score) ScoreCPU(metric model.Metric, result []byte) {
	// result预处理
	summary := resultSummary(result)
	avgNum, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", summary["avg"]), 64)
	// 计算评分
	if avgNum >= 30 {
		deduction := int8(1 + (avgNum/100-0.3)*20)
		message := fmt.Sprintf("%s%.2f%s过高 (>=30%%)", metric.Name, avgNum, metric.Unit)
		if avgNum >= 50 {
			deduction = 12
			message = fmt.Sprintf("%s%.2f%s过载 (>=50%%)", metric.Name, avgNum, metric.Unit)
		}
		model.Db.Debug().Create(&model.Score{
			InspId:    s.inspection.InspId,
			Metric:    metric.Key,
			Value:     decimal.NewFromFloat(avgNum),
			Deduction: deduction,
			Message:   message,
		})
	}
	// 写入评分
	// TODO 生成异常代办
}

func resultSummary(result []byte) map[string]float32 {
	json.Get(result).ToString()
	values := json.Get(result, 0).Get("values")
	var arr []float32
	for i := 0; i < values.Size(); i++ {
		arr = append(arr, values.Get(i, 1).ToFloat32())
	}
	maxNum := getMaxNum(arr)
	avgNum := getAvgNum(arr)
	minNum := getMinNum(arr)
	summaryMap := make(map[string]float32)
	summaryMap["max"] = maxNum
	summaryMap["avg"] = avgNum
	summaryMap["min"] = minNum
	return summaryMap
}

func getMaxNum(ary []float32) float32 {
	if len(ary) == 0 {
		return 0
	}
	maxVal := ary[0]
	for i := 1; i < len(ary); i++ {
		if maxVal < ary[i] {
			maxVal = ary[i]
		}
	}
	return maxVal
}

func getAvgNum(ary []float32) float32 {
	if len(ary) == 0 {
		return 0
	}
	sumVal := ary[0]
	for i := 1; i < len(ary); i++ {
		sumVal += ary[i]
	}
	return sumVal / float32(len(ary))
}

func getMinNum(ary []float32) float32 {
	if len(ary) == 0 {
		return 0
	}
	minVal := ary[0]
	for i := 1; i < len(ary); i++ {
		if minVal > ary[i] {
			minVal = ary[i]
		}
	}
	return minVal
}
