package routes

import (
	"math"
	"net/http"

	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-theluckiestsoul/models"
	"github.com/gin-gonic/gin"
)

func analyze(c *gin.Context) {
	var contracts []models.OptionsContract

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

func analyzeOptionsContracts(contracts []models.OptionsContract) models.AnalysisResult {
	const (
		tolerance  = 0.01
		priceRange = 50.0 // Assuming a range of 50 points above and below the max strike price for calculations
	)
	maxStrike := 0.0
	count := 0
	var result models.AnalysisResult

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
			result.GraphData = append(result.GraphData, models.GraphPoint{X: round(price, 2), Y: round(totalProfitLoss, 2)})
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

func calculateProfitLoss(contract models.OptionsContract, price float64) float64 {
	optionValue := 0.0
	switch contract.Type {
	case models.Call:
		if price > contract.StrikePrice {
			optionValue = price - contract.StrikePrice
		}
	case models.Put:
		if price < contract.StrikePrice {
			optionValue = contract.StrikePrice - price
		}
	}

	switch contract.LongShort {
	case models.Long:
		return optionValue - contract.Ask
	case models.Short:
		return contract.Bid - optionValue
	}

	return 0
}

func round(num float64, digits int) float64 {
	multiplier := math.Pow(10, float64(digits))
	return math.Round(num*multiplier) / multiplier
}
