package models

import "testing"

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
