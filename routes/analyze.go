package routes

import (
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	ErrInvalidOptionType       = errors.New("invalid option type")
	ErrInvalidLongShort        = errors.New("invalid long/short value")
	ErrInvalidStrikePrice      = errors.New("invalid strike price")
	ErrInvalidBidAsk           = errors.New("bid and ask prices must be positive")
	ErrExpirationDateRequired  = errors.New("expiration date is required")
	ErrorInvalidExpirationDate = errors.New("invalid expiration date")
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
		return ErrInvalidOptionType
	}
	if oc.LongShort != Long && oc.LongShort != Short {
		return ErrInvalidLongShort
	}
	if oc.StrikePrice <= 0 {
		return ErrInvalidStrikePrice
	}

	if oc.Bid < 0 || oc.Ask < 0 {
		return ErrInvalidBidAsk
	}

	if oc.ExpirationDate == "" {
		return ErrExpirationDateRequired
	}

	_, err := time.Parse(time.RFC3339, oc.ExpirationDate)
	if err != nil {
		return ErrorInvalidExpirationDate
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

func analyze(c *gin.Context) {
	var contracts []OptionsContract

	if err := c.ShouldBindJSON(&contracts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, contract := range contracts {
		if err := contract.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	result := analyzeOptionsContracts(contracts)
	c.JSON(http.StatusOK, result)
}

func analyzeOptionsContracts(contracts []OptionsContract) AnalysisResult {
	const (
		tolerance  = 0.01
		priceRange = 50.0 // Assuming a range of 50 points above and below the max strike price for calculations
	)
	maxStrike := 0.0
	count := 0
	var result AnalysisResult

	for _, contract := range contracts {
		if contract.StrikePrice > maxStrike {
			maxStrike = contract.StrikePrice
		}
	}

	priceChan := make(chan float64)
	go func() {
		for price := maxStrike - priceRange; price <= maxStrike+priceRange; price += 0.01 {
			priceChan <- price
		}
		close(priceChan)
	}()

	for price := range priceChan {
		totalProfitLoss := 0.0
		for _, contract := range contracts {
			profitLoss := calculateProfitLoss(contract, price)
			totalProfitLoss += profitLoss
		}
		if count%10 == 0 {
			result.GraphData = append(result.GraphData, GraphPoint{X: round(price, 2), Y: round(totalProfitLoss, 2)})
		}
		count++
		if totalProfitLoss > result.MaxProfit {
			result.MaxProfit = round(totalProfitLoss, 2)
		}
		if totalProfitLoss < result.MaxLoss {
			result.MaxLoss = round(totalProfitLoss, 2)
		}
		if math.Abs(totalProfitLoss) < tolerance {
			result.BreakEvenPoints = append(result.BreakEvenPoints, round(price, 2))
		}
	}

	return result
}

func calculateProfitLoss(contract OptionsContract, price float64) float64 {
	optionValue := 0.0
	switch contract.Type {
	case Call:
		if price > contract.StrikePrice {
			optionValue = price - contract.StrikePrice
		}
	case Put:
		if price < contract.StrikePrice {
			optionValue = contract.StrikePrice - price
		}
	}

	switch contract.LongShort {
	case Long:
		return optionValue - contract.Ask
	case Short:
		return contract.Bid - optionValue
	}

	return 0
}

func round(num float64, digits int) float64 {
	multiplier := math.Pow(10, float64(digits))
	return math.Round(num*multiplier) / multiplier
}
