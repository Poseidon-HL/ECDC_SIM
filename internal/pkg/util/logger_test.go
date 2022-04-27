package util

import "testing"

func TestClearLogDir(t *testing.T) {
	type args struct {
		dirPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "ClearLogDir", args: args{dirPath: `D:\Code\Projects\Go\ECDC_SIM\output\event_log\`}},
		{name: "ClearLogDir", args: args{dirPath: `D:\Code\Projects\Go\ECDC_SIM\output\data_center_log\`}},
		{name: "ClearLogDir", args: args{dirPath: `D:\Code\Projects\Go\ECDC_SIM\output\default_log\`}},
		{name: "ClearLogDir", args: args{dirPath: `D:\Code\Projects\Go\ECDC_SIM\output\event_log\`}},
		{name: "ClearLogDir", args: args{dirPath: `D:\Code\Projects\Go\ECDC_SIM\output\logrus_test_log\`}},
		{name: "ClearLogDir", args: args{dirPath: `D:\Code\Projects\Go\ECDC_SIM\output\test_log\`}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ClearLogDir(tt.args.dirPath); (err != nil) != tt.wantErr {
				t.Errorf("ClearLogDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
