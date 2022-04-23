package util

import (
	"testing"
	"time"
)

func TestWeibull_Draw(t *testing.T) {
	type fields struct {
		shape    float64
		scale    float64
		location float64
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{name: "testWeibull 1.1", fields: fields{shape: 1, scale: 1, location: 0}},
		{name: "testWeibull 1.1", fields: fields{shape: 1, scale: 1, location: 0}},
		{name: "testWeibull 1.1", fields: fields{shape: 1, scale: 1, location: 0}},
		{name: "testWeibull 1.1", fields: fields{shape: 1, scale: 1, location: 0}},
		{name: "testWeibull 1.2", fields: fields{shape: 1, scale: 2, location: 0}},
		{name: "testWeibull 1.2", fields: fields{shape: 1, scale: 2, location: 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time.Sleep(time.Second)
			w := &Weibull{
				shape:    tt.fields.shape,
				scale:    tt.fields.scale,
				location: tt.fields.location,
			}
			t.Log(w.Draw())
		})
	}
}
