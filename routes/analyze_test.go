package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOptionsContractValidate(t *testing.T) {
	tests := []struct {
		name    string
		oc      OptionsContract
		wantErr bool
	}{
		{
			name: "valid contract",
			oc: OptionsContract{
				Type:           Call,
				StrikePrice:    100.0,
				Bid:            1.0,
				Ask:            2.0,
				ExpirationDate: "2022-12-31T23:59:59Z",
				LongShort:      Long,
			},
			wantErr: false,
		},
		{
			name: "invalid option type",
			oc: OptionsContract{
				Type:           "InvalidType",
				StrikePrice:    100.0,
				Bid:            1.0,
				Ask:            2.0,
				ExpirationDate: "2022-12-31T23:59:59Z",
				LongShort:      Long,
			},
			wantErr: true,
		},
		{
			name: "invalid long/short value",
			oc: OptionsContract{
				Type:           Call,
				StrikePrice:    100.0,
				Bid:            1.0,
				Ask:            2.0,
				ExpirationDate: "2022-12-31T23:59:59Z",
				LongShort:      "InvalidLongShort",
			},
			wantErr: true,
		},
		{
			name: "invalid strike price",
			oc: OptionsContract{
				Type:           Call,
				StrikePrice:    -100.0,
				Bid:            1.0,
				Ask:            2.0,
				ExpirationDate: "2022-12-31T23:59:59Z",
				LongShort:      Long,
			},
			wantErr: true,
		},
		{
			name: "bid and ask prices must be positive",
			oc: OptionsContract{

				Type:           Call,
				StrikePrice:    100.0,
				Bid:            -1.0,
				Ask:            -2.0,
				ExpirationDate: "2022-12-31T23:59:59Z",
				LongShort:      Long,
			},
			wantErr: true,
		},
		{
			name: "expiration date is required",
			oc: OptionsContract{
				Type:           Call,
				StrikePrice:    100.0,
				Bid:            1.0,
				Ask:            2.0,
				ExpirationDate: "",
				LongShort:      Long,
			},
			wantErr: true,
		},
		{
			name: "invalid expiration date",
			oc: OptionsContract{
				Type:           Call,
				StrikePrice:    100.0,
				Bid:            1.0,
				Ask:            2.0,
				ExpirationDate: "invalid",
				LongShort:      Long,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.oc.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("OptionsContract.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAnalyzeEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/analyze", analyze)

	tests := []struct {
		name           string
		contracts      []OptionsContract
		expectedStatus int
	}{
		{
			name: "single call contract",
			contracts: []OptionsContract{
				{
					Type:           Call,
					StrikePrice:    100.0,
					Bid:            1.0,
					Ask:            2.0,
					ExpirationDate: "2022-12-31T23:59:59Z",
					LongShort:      Long,
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "multiple contracts",
			contracts: []OptionsContract{
				{
					Type:           Call,
					StrikePrice:    100.0,
					Bid:            10.05,
					Ask:            12.04,
					ExpirationDate: "2025-12-17T00:00:00Z",
					LongShort:      Long,
				},
				{
					Type:           Call,
					StrikePrice:    102.50,
					Bid:            12.10,
					Ask:            14.0,
					ExpirationDate: "2025-12-17T00:00:00Z",
					LongShort:      Long,
				},
				{
					Type:           Put,
					StrikePrice:    103.0,
					Bid:            14.0,
					Ask:            15.50,
					ExpirationDate: "2025-12-17T00:00:00Z",
					LongShort:      Short,
				},
				{
					Type:           Put,
					StrikePrice:    105.0,
					Bid:            16.0,
					Ask:            18.0,
					ExpirationDate: "2025-12-17T00:00:00Z",
					LongShort:      Long,
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
