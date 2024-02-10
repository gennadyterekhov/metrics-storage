package validators

import (
	"errors"
	"github.com/gennadyterekhov/metrics-storage/internal/exceptions"
	"github.com/gennadyterekhov/metrics-storage/internal/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_validateMetricName(t *testing.T) {
	type args struct {
		nameRaw string
	}
	possibleError := errors.New(exceptions.EmptyMetricName)

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "ok",
			args:    args{nameRaw: "name"},
			wantErr: false,
		},
		{
			name:    "empty",
			args:    args{nameRaw: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateMetricName(tt.args.nameRaw); (err != nil) != tt.wantErr {
				t.Errorf("validateMetricName() error = %v, wantErr %v", err, possibleError)
			}
		})
	}
}

func Test_validateMetricType(t *testing.T) {
	type args struct {
		metricTypeRaw string
		statusCode    int
	}

	tests := []struct {
		name         string
		args         args
		wantErr      bool
		errorMessage string
	}{
		{
			name:    "counter",
			args:    args{metricTypeRaw: "counter"},
			wantErr: false,
		},
		{
			name:    "gauge",
			args:    args{metricTypeRaw: "gauge"},
			wantErr: false,
		},
		{
			name:    "unknown",
			args:    args{metricTypeRaw: "unknown"},
			wantErr: true,
		},
		{
			name:    "wrong",
			args:    args{metricTypeRaw: "wrong"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateMetricType(tt.args.metricTypeRaw); (err != nil) != tt.wantErr {
				t.Errorf("validateMetricType() error = %v, wantErr %v", err, tt.wantErr)
				if tt.errorMessage != "" {
					require.Equal(t, tt.errorMessage, err.Error())
				}
			}
		})
	}
}

func Test_validateMetricValue(t *testing.T) {
	type args struct {
		metricTypeValidated string
		valueRaw            string
	}
	tests := []struct {
		name             string
		args             args
		wantCounterValue int64
		wantGaugeValue   float64
		wantErr          bool
	}{
		{
			name:             "counter ok",
			args:             args{metricTypeValidated: "counter", valueRaw: "10"},
			wantCounterValue: 10,
			wantGaugeValue:   0,
			wantErr:          false,
		},
		{
			name:             "gauge ok",
			args:             args{metricTypeValidated: "gauge", valueRaw: "10.45"},
			wantCounterValue: 0,
			wantGaugeValue:   10.45,
			wantErr:          false,
		},

		{
			name:             "counter float error",
			args:             args{metricTypeValidated: "counter", valueRaw: "10.45"},
			wantCounterValue: 10,
			wantGaugeValue:   0,
			wantErr:          true,
		},

		{
			name:             "counter string error",
			args:             args{metricTypeValidated: "counter", valueRaw: "hello"},
			wantCounterValue: 0,
			wantGaugeValue:   0,
			wantErr:          true,
		},
		{
			name:             "gauge string error",
			args:             args{metricTypeValidated: "gauge", valueRaw: "hello"},
			wantCounterValue: 0,
			wantGaugeValue:   0,
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCounterValue, gotGaugeValue, err := validateMetricValue(tt.args.metricTypeValidated, tt.args.valueRaw)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMetricValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCounterValue != tt.wantCounterValue {
				t.Errorf("validateMetricValue() gotCounterValue = %v, want %v", gotCounterValue, tt.wantCounterValue)
			}
			if gotGaugeValue != tt.wantGaugeValue {
				t.Errorf("validateMetricValue() gotGaugeValue = %v, want %v", gotGaugeValue, tt.wantGaugeValue)
			}
		})
	}
}

func Test_validateParameters(t *testing.T) {
	type args struct {
		metricTypeRaw string
		nameRaw       string
		valueRaw      string
	}
	tests := []struct {
		name             string
		args             args
		wantCounterValue int64
		wantGaugeValue   float64
		wantErr          bool
	}{
		{
			name:             "counter",
			args:             args{metricTypeRaw: "counter", nameRaw: "name", valueRaw: "1"},
			wantCounterValue: 1,
			wantGaugeValue:   0,
			wantErr:          false,
		},
		{
			name:             "gauge",
			args:             args{metricTypeRaw: "gauge", nameRaw: "name", valueRaw: "1.1"},
			wantCounterValue: 0,
			wantGaugeValue:   1.1,
			wantErr:          false,
		},
		{
			name:             "gauge string error",
			args:             args{metricTypeRaw: "gauge", nameRaw: "name", valueRaw: "hello"},
			wantCounterValue: 0,
			wantGaugeValue:   0,
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCounterValue, gotGaugeValue, err := validateParameters(tt.args.metricTypeRaw, tt.args.nameRaw, tt.args.valueRaw)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateParameters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCounterValue != tt.wantCounterValue {
				t.Errorf("validateParameters() gotCounterValue = %v, want %v", gotCounterValue, tt.wantCounterValue)
			}
			if gotGaugeValue != tt.wantGaugeValue {
				t.Errorf("validateParameters() gotGaugeValue = %v, want %v", gotGaugeValue, tt.wantGaugeValue)
			}
		})
	}
}

func TestGetDataToSave(t *testing.T) {
	type args struct {
		metricType  string
		metricName  string
		metricValue string
	}
	type want struct {
		metricType   string
		metricName   string
		counterValue int64
		gaugeValue   float64
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name:    "ok",
			args:    args{metricType: types.Counter, metricName: types.Counter, metricValue: "1"},
			want:    want{metricType: types.Counter, metricName: types.Counter, counterValue: int64(1), gaugeValue: float64(0)},
			wantErr: false,
		},
		{
			name:    "metricType empty",
			args:    args{metricType: "", metricName: types.Counter, metricValue: "1"},
			want:    want{metricType: "", metricName: "", counterValue: int64(0), gaugeValue: float64(0)},
			wantErr: true,
		},
		{
			name:    "metricName empty",
			args:    args{metricType: types.Counter, metricName: "", metricValue: "1"},
			want:    want{metricType: "", metricName: "", counterValue: int64(0), gaugeValue: float64(0)},
			wantErr: true,
		},
		{
			name:    "metricValue empty",
			args:    args{metricType: types.Counter, metricName: types.Counter, metricValue: ""},
			want:    want{metricType: "", metricName: "", counterValue: int64(0), gaugeValue: float64(0)},
			wantErr: true,
		},
		{
			name:    "InvalidMetricTypeChoice",
			args:    args{metricType: "hello", metricName: types.Counter, metricValue: "1"},
			want:    want{metricType: "", metricName: "", counterValue: int64(0), gaugeValue: float64(0)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filledDto, err := GetDataToSave(tt.args.metricType, tt.args.metricName, tt.args.metricValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDataToSave() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if filledDto.Type != tt.want.metricType {
				t.Errorf("GetDataToSave() gotTyp = %v, want %v", filledDto.Type, tt.want.metricType)
			}
			if filledDto.Name != tt.want.metricName {
				t.Errorf("GetDataToSave() gotName = %v, want %v", filledDto.Name, tt.want.metricName)
			}
			if filledDto.CounterValue != tt.want.counterValue {
				t.Errorf("GetDataToSave() gotCounterValue = %v, want %v", filledDto.CounterValue, tt.want.counterValue)
			}
			if filledDto.GaugeValue != tt.want.gaugeValue {
				t.Errorf("GetDataToSave() gotGaugeValue = %v, want %v", filledDto.GaugeValue, tt.want.gaugeValue)
			}
		})
	}
}
