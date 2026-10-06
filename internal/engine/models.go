package engine

import (
	"time"
)

type Rate struct {
	CurID           int     `json:"Cur_ID"`
	Date            string  `json:"Date"`
	CurCode         string  `json:"Cur_Code"`
	CurAbbreviation string  `json:"Cur_Abbreviation"`
	CurScale        int     `json:"Cur_Scale"`
	CurName         string  `json:"Cur_Name"`
	CurOfficialRate float64 `json:"Cur_OfficialRate"`
}

type CrossCourseParams struct {
	CurrencyID  string
	ParamMode   int
	Periodicity int
	OnDate      time.Time
}

type Transaction struct {
	TransactionId int64
	CurId         int64
	Amount        float64
	Tags          []string
	Comment       string
	CreateAt      time.Time
	Wallet        string
}

type FilterOptions struct {
	Tag      string
	FromDate time.Time
	ToDate   time.Time
}

type Currency struct {
	CurrencyId int64
	Symbol     string
	NbrbCurId  int64
}

type TotalReport struct {
	Total     float64
	TotalLost float64
	TotalGain float64
}
