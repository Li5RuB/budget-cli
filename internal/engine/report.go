package engine

import (
	"sort"
	"strconv"
)

type Report struct {
	cc *CrossCourse
}

func NewReport() *Report {
	return &Report{
		cc: NewCrossCourse(),
	}
}

func (r *Report) CalculateTagReport(transactions []Transaction) (map[string]float64, error) {
	report := make(map[string]float64)

	for _, tx := range transactions {
		ccAmount, err := GetAdjCrossCourseAmount(tx, r.cc)
		if err != nil {
			return nil, err
		}

		if len(tx.Tags) == 0 {
			report["other"] += ccAmount
			continue
		}

		for _, tag := range tx.Tags {
			report[tag] += ccAmount
		}
	}

	return report, nil
}

func (r *Report) CalculateWalletReport(transactions []Transaction) (map[string]float64, error) {
	report := make(map[string]float64)

	for _, tx := range transactions {
		ccAmount, err := GetAdjCrossCourseAmount(tx, r.cc)
		if err != nil {
			return nil, err
		}

		report[tx.Wallet] += ccAmount
	}

	return report, nil
}

func (r *Report) GenerateMonthlyReport(transactions []Transaction) (map[string]float64, []string, error) {
	reportMap := make(map[string]float64)

	for _, tx := range transactions {
		ccAmount, err := GetAdjCrossCourseAmount(tx, r.cc)

		if err != nil {
			return nil, nil, err
		}

		monthKey := tx.CreateAt.Format("2006-Jan")
		reportMap[monthKey] += ccAmount
	}

	months := make([]string, 0, len(reportMap))
	for month := range reportMap {
		months = append(months, month)
	}

	sort.Strings(months)

	return reportMap, months, nil
}

func (r *Report) GetTotal(transactions []Transaction) (*TotalReport, error) {
	var total TotalReport
	for _, tx := range transactions {
		ccAmount, err := GetAdjCrossCourseAmount(tx, r.cc)
		if err != nil {
			return nil, err
		}

		total.Total += ccAmount

		if ccAmount > 0 {
			total.TotalGain += ccAmount
		} else {
			total.TotalLost += ccAmount
		}
	}

	return &total, nil
}

func GetAdjCrossCourseAmount(tx Transaction, cc *CrossCourse) (float64, error) {
	var ccAmount float64

	if tx.CurId != 0 {
		rate, err := cc.Get(CrossCourseParams{CurrencyID: strconv.FormatInt(tx.CurId, 10), ParamMode: 0})
		if err != nil {
			return 0, err
		}
		ccAmount = (tx.Amount / float64(rate.CurScale)) * rate.CurOfficialRate
	} else {
		ccAmount = tx.Amount
	}

	return ccAmount, nil
}
