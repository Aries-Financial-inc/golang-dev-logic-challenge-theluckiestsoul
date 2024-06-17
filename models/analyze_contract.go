package models

import (
	"errors"
	"time"
)

type OptionType string

const (
	Put  OptionType = "Put"
	Call OptionType = "Call"
)

type LongShort string

const (
	Long  LongShort = "long"
	Short LongShort = "short"
)

var (
	errInvalidOptionType       = errors.New("invalid option type")
	errInvalidLongShort        = errors.New("invalid long/short value")
	errInvalidStrikePrice      = errors.New("invalid strike price")
	errInvalidBidAsk           = errors.New("bid and ask prices must be positive")
	errExpirationDateRequired  = errors.New("expiration date is required")
	errorInvalidExpirationDate = errors.New("invalid expiration date")
)

type OptionsContract struct {
	Type           OptionType `json:"type"`
	StrikePrice    float64    `json:"strike_price"`
	Bid            float64    `json:"bid"`
	Ask            float64    `json:"ask"`
	ExpirationDate string     `json:"expiration_date"`
	LongShort      LongShort  `json:"long_short"`
}

func (oc OptionsContract) Validate() error {
	if oc.Type != Call && oc.Type != Put {
		return errInvalidOptionType
	}
	if oc.LongShort != Long && oc.LongShort != Short {
		return errInvalidLongShort
	}
	if oc.StrikePrice <= 0 {
		return errInvalidStrikePrice
	}

	if oc.Bid < 0 || oc.Ask < 0 {
		return errInvalidBidAsk
	}

	if oc.ExpirationDate == "" {
		return errExpirationDateRequired
	}

	_, err := time.Parse(time.RFC3339, oc.ExpirationDate)
	if err != nil {
		return errorInvalidExpirationDate
	}

	return nil
}

// AnalysisResult structure for the response body
type AnalysisResult struct {
	GraphData       []GraphPoint `json:"graph_data"`
	MaxProfit       float64      `json:"max_profit"`
	MaxLoss         float64      `json:"max_loss"`
	BreakEvenPoints []float64    `json:"break_even_points"`
}

// GraphPoint structure for X & Y values of the risk & reward graph
type GraphPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
