package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-theluckiestsoul/models"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)



func TestAnalyzeEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/analyze", analyze)

	tests := []struct {
		name           string
		contracts      []models.OptionsContract
		expectedStatus int
	}{
		{
			name: "single call contract",
			contracts: []models.OptionsContract{
				{
					Type:           models.Call,
					StrikePrice:    100.0,
					Bid:            1.0,
					Ask:            2.0,
					ExpirationDate: "2022-12-31T23:59:59Z",
					LongShort:      models.Long,
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "multiple contracts",
			contracts: []models.OptionsContract{
				{
					Type:          models.Call,
					StrikePrice:    100.0,
					Bid:            10.05,
					Ask:            12.04,
					ExpirationDate: "2025-12-17T00:00:00Z",
					LongShort:      models.Long,
				},
				{
					Type:           models.Call,
					StrikePrice:    102.50,
					Bid:            12.10,
					Ask:            14.0,
					ExpirationDate: "2025-12-17T00:00:00Z",
					LongShort:     models.Long,
				},
				{
					Type:          models.Put,
					StrikePrice:    103.0,
					Bid:            14.0,
					Ask:            15.50,
					ExpirationDate: "2025-12-17T00:00:00Z",
					LongShort:      models.Short,
				},
				{
					Type:           models.Put,
					StrikePrice:    105.0,
					Bid:            16.0,
					Ask:            18.0,
					ExpirationDate: "2025-12-17T00:00:00Z",
					LongShort:      models.Long,
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contractsJSON, _ := json.Marshal(tt.contracts)

			req, err := http.NewRequest(http.MethodPost, "/analyze", bytes.NewBuffer(contractsJSON))
			assert.NoError(t, err)

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			if status := resp.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
			res := resp.Result()
			defer res.Body.Close()

			cupaloy.SnapshotT(t, dumpResponse(t, res))
		})
	}
}

func dumpResponse(t *testing.T, r *http.Response) string {
	t.Helper()
	body, err := httputil.DumpResponse(r, true)
	if err != nil {
		t.Fatalf("failed to dump response: %v", err)
	}
	return string(body)
}
