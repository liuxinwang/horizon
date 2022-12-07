package tasks

import (
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
	"horizon/model"
	"math"
	"strconv"
)

type Score struct {
	inspection model.Inspection
	instance   model.Instance
}

type Result struct {
	Total int8
}

type HealthLevel struct {
	level     string
	levelName string
}

func (s *Score) ScoreCPU(metric model.Metric, result []byte) {
	// result预处理
	summary := resultSummary(result)
	avgNum, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", summary["avg"]), 64)
	// 计算评分
	if avgNum >= 30 {
		deduction := int8(1 + (avgNum/100.00-0.30)*20)
		message := fmt.Sprintf("%s%.2f%s过高 (>=30%%)", metric.Name, avgNum, metric.Unit)
		if avgNum >= 50 {
			deduction = 12
			message = fmt.Sprintf("%s%.2f%s过载 (>=50%%)", metric.Name, avgNum, metric.Unit)
		}
		// 写入评分
		model.Db.Create(&model.Score{
			InspId:    s.inspection.InspId,
			Metric:    metric.Key,
			Value:     decimal.NewFromFloat(avgNum),
			Deduction: deduction,
			Message:   message,
		})
	}
	// TODO 生成异常代办
}

func (s *Score) ScoreMemory(metric model.Metric, result []byte, availableResult []byte) {
	// result预处理
	summary := resultSummary(result)
	avgNum, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", summary["avg"]), 64)
	availableSummary := resultSummary(availableResult)
	availableAvgNum, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", availableSummary["avg"]), 64)

	// 计算评分
	if avgNum >= 80 {
		if availableAvgNum < 4*1024*1024*1024 {
			deduction := 3
			message := fmt.Sprintf("%s%.2f%s过高 (>=80%%)，当前可用内存%.2fG(<4G)", metric.Name, avgNum, metric.Unit, availableAvgNum)

			if availableAvgNum < 2*1024*1024*1024 {
				deduction = 12
				message = fmt.Sprintf("%s%.2f%s过载 (>=80%%)，当前可用内存%.2fG(<2G)",
					metric.Name, avgNum, metric.Unit, availableAvgNum/1024/1024/1024)
			}

			// 写入评分
			model.Db.Create(&model.Score{
				InspId:    s.inspection.InspId,
				Metric:    metric.Key,
				Value:     decimal.NewFromFloat(avgNum),
				Deduction: int8(deduction),
				Message:   message,
			})
		}
	}
	// TODO 生成异常代办
}

func (s *Score) ScoreSwap(metric model.Metric, result []byte) {
	// result预处理
	summary := resultSummary(result)
	avgNum, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", summary["avg"]), 64)
	// 计算评分
	if avgNum > 0 {
		deduction := int8(4)
		message := fmt.Sprintf("%s%.2f%s (>0%s)", metric.Name, avgNum, metric.Unit, metric.Unit)
		// 写入评分
		model.Db.Create(&model.Score{
			InspId:    s.inspection.InspId,
			Metric:    metric.Key,
			Value:     decimal.NewFromFloat(avgNum),
			Deduction: deduction,
			Message:   message,
		})
	}
	// TODO 生成异常代办
}

func (s *Score) ScoreDisk(metric model.Metric, resultAvailable []byte, resultGrowth []byte) {
	// result预处理
	availableValue := json.Get(resultAvailable, 0).Get("value", 1).ToInt64()
	avgGrowthValue := json.Get(resultGrowth, 0).Get("value", 1).ToInt64()
	availableDays := 90
	if avgGrowthValue > 0 {
		availableDays = int(availableValue / avgGrowthValue)
	}
	if availableDays <= 30 {
		deduction := int8(12 - availableDays/3)
		message := fmt.Sprintf("磁盘可用空间剩余%d天 (<=30天)", availableDays)
		// 写入评分
		model.Db.Create(&model.Score{
			InspId:    s.inspection.InspId,
			Metric:    metric.Key,
			Value:     decimal.NewFromFloat(float64(availableDays)),
			Deduction: deduction,
			Message:   message,
		})
	}
}

func (s *Score) ScoreIOPS(metric model.Metric, result []byte) {
	// result预处理
	summary := resultSummary(result)
	avgNum, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", summary["avg"]), 64)
	// 计算评分
	if avgNum >= 4000 {
		deduction := 4
		message := fmt.Sprintf("IOPS使用%.2f过高 (>=4000)", avgNum)
		if avgNum >= 4500 {
			deduction = 12
			message = fmt.Sprintf("IOPS使用%.2f过载 (>=4500)", avgNum)
		}
		// 写入评分
		model.Db.Create(&model.Score{
			InspId:    s.inspection.InspId,
			Metric:    metric.Key,
			Value:     decimal.NewFromFloat(avgNum),
			Deduction: int8(deduction),
			Message:   message,
		})
	}
	// TODO 生成异常代办
}

func (s *Score) ScoreDeadlock(metric model.Metric, result []byte) {
	// result预处理
	summary := resultSummary(result)
	maxNum, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", summary["max"]), 64)
	minNum, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", summary["min"]), 64)
	dlNum := maxNum - minNum
	// 计算评分
	if dlNum > 0 {
		deduction := 4
		message := fmt.Sprintf("发现%s%.2f%s (>=1%s)", metric.Name, dlNum, metric.Unit, metric.Unit)
		// 写入评分
		model.Db.Create(&model.Score{
			InspId:    s.inspection.InspId,
			Metric:    metric.Key,
			Value:     decimal.NewFromFloat(dlNum),
			Deduction: int8(deduction),
			Message:   message,
		})
	}
	// TODO 生成异常代办
}

func (s *Score) ScoreSlowSQLNum(metric model.Metric, slowTotalNum int) {
	// 计算评分
	if slowTotalNum >= 100 {
		deduction := 2 + (slowTotalNum-100)/50
		message := fmt.Sprintf("%s%d%s (>=100%s)", metric.Name, slowTotalNum, metric.Unit, metric.Unit)
		if slowTotalNum >= 500 {
			deduction = 18
			message = fmt.Sprintf("%s%d%s (>=500%s)", metric.Name, slowTotalNum, metric.Unit, metric.Unit)
		}
		// 写入评分
		model.Db.Create(&model.Score{
			InspId:    s.inspection.InspId,
			Metric:    metric.Key,
			Value:     decimal.NewFromFloat(float64(slowTotalNum)),
			Deduction: int8(deduction),
			Message:   message,
		})
	}
	// TODO 生成异常代办
}

func (s *Score) ScoreIncrementIdOverflow(metric model.Metric, result []byte) {
	// result预处理
	if string(result) == "null" {
		return
	}
	tableNum := json.Get(result).Size()
	if tableNum == 0 {
		return
	}
	// 计算评分
	deduction := int8(4)
	message := fmt.Sprintf("发现%s风险%d%s (<20%%)", metric.Name, tableNum, metric.Unit)
	// 写入评分
	model.Db.Create(&model.Score{
		InspId:    s.inspection.InspId,
		Metric:    metric.Key,
		Value:     decimal.NewFromFloat(float64(tableNum)),
		Deduction: deduction,
		Message:   message,
	})
	// TODO 生成异常代办
}

func (s *Score) ScoreLockWait(metric model.Metric, result []byte) {
	// result预处理
	summary := resultSummary(result)
	maxNum, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", summary["max"]), 64)
	// 计算评分
	if maxNum > 3 {
		deduction := int8(4)
		message := fmt.Sprintf("%s%.2f%s (>=3%s)", metric.Name, maxNum, metric.Unit, metric.Unit)
		// 写入评分
		model.Db.Create(&model.Score{
			InspId:    s.inspection.InspId,
			Metric:    metric.Key,
			Value:     decimal.NewFromFloat(maxNum),
			Deduction: deduction,
			Message:   message,
		})
	}
	// TODO 生成异常代办
}

func (s *Score) ScoreBigTableNum(metric model.Metric, result []byte) {
	maxNum := float64(json.Get(result).Size())
	// 计算评分
	if maxNum > 0 {
		deduction := math.Min(maxNum, 4)
		message := fmt.Sprintf("%s%.2f%s (>1千万行或>10GB)", metric.Name, maxNum, metric.Unit)
		// 写入评分
		model.Db.Create(&model.Score{
			InspId:    s.inspection.InspId,
			Metric:    metric.Key,
			Value:     decimal.NewFromFloat(maxNum),
			Deduction: int8(deduction),
			Message:   message,
		})
	}
	// TODO 生成异常代办
}

func (s *Score) ScoreThreadsRunningNum(metric model.Metric, result []byte) {
	// result预处理
	summary := resultSummary(result)
	maxNum, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", summary["max"]), 64)
	// 计算评分
	if maxNum >= (2*4 + 2) {
		deduction := 2
		message := fmt.Sprintf("%s%.2f%s (>=10%s)", metric.Name, maxNum, metric.Unit, metric.Unit)
		if maxNum >= (4*4 + 4) {
			deduction = 4
			message = fmt.Sprintf("%s%.2f%s (>=20%s)", metric.Name, maxNum, metric.Unit, metric.Unit)
		}
		// 写入评分
		model.Db.Create(&model.Score{
			InspId:    s.inspection.InspId,
			Metric:    metric.Key,
			Value:     decimal.NewFromFloat(maxNum),
			Deduction: int8(deduction),
			Message:   message,
		})
	}
	// TODO 生成异常代办
}

func (s *Score) ScoreThreadsConnected(metric model.Metric, result []byte) {
	// result预处理
	summary := resultSummary(result)
	maxNum, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", summary["max"]), 64)
	// 计算评分
	if maxNum >= 70 {
		deduction := 2
		message := fmt.Sprintf("%s%.2f%s (>=70%s)", metric.Name, maxNum, metric.Unit, metric.Unit)
		if maxNum >= 80 {
			deduction = 4
			message = fmt.Sprintf("%s%.2f%s (>=80%s)", metric.Name, maxNum, metric.Unit, metric.Unit)
		}
		// 写入评分
		model.Db.Create(&model.Score{
			InspId:    s.inspection.InspId,
			Metric:    metric.Key,
			Value:     decimal.NewFromFloat(maxNum),
			Deduction: int8(deduction),
			Message:   message,
		})
	}
	// TODO 生成异常代办
}

func (s *Score) ScoreIBPCacheHitsRate(metric model.Metric, result []byte) {
	// result预处理
	summary := resultSummary(result)
	avgNum, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", summary["avg"]), 64)
	// 计算评分
	if avgNum < 99 {
		deduction := 4
		message := fmt.Sprintf("%s%.2f%s (<99%s)", metric.Name, avgNum, metric.Unit, metric.Unit)
		// 写入评分
		model.Db.Create(&model.Score{
			InspId:    s.inspection.InspId,
			Metric:    metric.Key,
			Value:     decimal.NewFromFloat(avgNum),
			Deduction: int8(deduction),
			Message:   message,
		})
	}
	// TODO 生成异常代办
}

func (s *Score) ScoreHighRiskAccount(metric model.Metric, result []byte) {
	maxNum := float64(json.Get(result).Size())
	// 计算评分
	if maxNum >= 1 {
		deduction := 4
		message := fmt.Sprintf("%s%.2f%s (>=1%s)", metric.Name, maxNum, metric.Unit, metric.Unit)
		// 写入评分
		model.Db.Create(&model.Score{
			InspId:    s.inspection.InspId,
			Metric:    metric.Key,
			Value:     decimal.NewFromFloat(maxNum),
			Deduction: int8(deduction),
			Message:   message,
		})
	}
	// TODO 生成异常代办
}

func calculateScoreLevel(totalScore int8) string {
	var level string
	switch {
	case totalScore < 60:
		level = "CRITICAL"
	case totalScore < 80:
		level = "RISKY"
	case totalScore < 95:
		level = "SUBOPTIMAL"
	default:
		level = "HEALTHY"
	}
	return level
}

func calculateTotalScore(inspId string) int8 {
	var result Result
	model.Db.Model(&Score{}).Select("100 - sum(deduction) as total").Where("insp_id = ?", inspId).Find(&result)
	return result.Total
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
