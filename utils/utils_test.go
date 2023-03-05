package utils

import "testing"

func TestInterval_Sat(t *testing.T) {
	tests := []struct {
		name     string
		interval Interval
		arg      float64
		want     float64
	}{
		{
			"from -1 to 9, -2",
			Interval{
				Start: -1,
				End:   9,
			},
			-2,
			-1,
		},
		{
			"from -1 to 9, 3",
			Interval{
				Start: -1,
				End:   9,
			},
			3,
			3,
		},
		{
			"from -1 to 9, 10",
			Interval{
				Start: -1,
				End:   9,
			},
			10,
			9,
		},
		{
			"from 9 to -1, 3",
			Interval{
				Start: 9,
				End:   -1,
			},
			3,
			3,
		},
		{
			"from 9 to -1, -3",
			Interval{
				Start: 9,
				End:   -1,
			},
			-3,
			-1,
		},
		{
			"from 9 to -1, -10",
			Interval{
				Start: 9,
				End:   -1,
			},
			-10,
			-1,
		},
	}
	for _, tt := range tests {
		interval := &tt.interval
		arg := tt.arg
		want := tt.want
		t.Run(tt.name, func(t *testing.T) {
			if got := interval.Sat(arg); got != want {
				t.Errorf("Sat() = %v, want %v", got, want)
			}
		})
	}
}

func TestInterval_Length(t *testing.T) {
	tests := []struct {
		name     string
		interval Interval
		want     float64
	}{
		{
			name: "from -1 to 9",
			interval: Interval{
				Start: -1,
				End:   9,
			},
			want: 10,
		},
		{
			name: "from 9 to -1",
			interval: Interval{
				Start: 9,
				End:   -1,
			},
			want: -10,
		},
	}
	for _, tt := range tests {
		interval := &tt.interval
		want := tt.want
		t.Run(tt.name, func(t *testing.T) {
			if got := interval.Length(); got != want {
				t.Errorf("Length() = %v, want %v", got, want)
			}
		})
	}
}

func TestInterval_GetMinMax(t *testing.T) {
	tests := []struct {
		name     string
		Interval Interval
		wantMin  float64
		wantMax  float64
	}{
		{
			name: "from -1 to 9",
			Interval: Interval{
				Start: -1,
				End:   9,
			},
			wantMin: -1,
			wantMax: 9,
		},
		{
			name: "from 9 to -1",
			Interval: Interval{
				Start: 9,
				End:   -1,
			},
			wantMin: -1,
			wantMax: 9,
		},
	}
	for _, tt := range tests {
		i := &tt.Interval
		wantMin := tt.wantMin
		wantMax := tt.wantMax
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax := i.GetMinMax()
			if gotMin != wantMin {
				t.Errorf("GetMinMax() gotMin = %v, want %v", gotMin, wantMin)
			}
			if gotMax != wantMax {
				t.Errorf("GetMinMax() gotMax = %v, want %v", gotMax, wantMax)
			}
		})
	}
}

func TestIntervalTransformer_TransformForward(t1 *testing.T) {
	tests := []struct {
		name      string
		interval1 Interval
		interval2 Interval
		arg       float64
		want      float64
	}{
		{
			name: "(-10,90)->(0,10), 50->",
			interval1: Interval{
				Start: -10,
				End:   90,
			},
			interval2: Interval{
				Start: 0,
				End:   10,
			},
			arg:  50,
			want: 6,
		},
		{
			name: "(-10,90)->(0,10), 90->",
			interval1: Interval{
				Start: -10,
				End:   90,
			},
			interval2: Interval{
				Start: 0,
				End:   10,
			},
			arg:  90,
			want: 10,
		},
		{
			name: "(-10,90)->(0,10), -10->",
			interval1: Interval{
				Start: -10,
				End:   90,
			},
			interval2: Interval{
				Start: 0,
				End:   10,
			},
			arg:  -10,
			want: 0,
		},
		{
			name: "(-10,90)->(0,10), 100->",
			interval1: Interval{
				Start: -10,
				End:   90,
			},
			interval2: Interval{
				Start: 0,
				End:   10,
			},
			arg:  100,
			want: 10,
		},
	}
	for _, tt := range tests {
		t := &IntervalTransformer{
			IntervalFrom: tt.interval1,
			IntervalTo:   tt.interval2,
		}
		arg := tt.arg
		want := tt.want
		t1.Run(tt.name, func(t1 *testing.T) {
			if got := t.TransformForward(arg); got != want {
				t1.Errorf("TransformForward() = %v, want %v", got, want)
			}
		})
	}
}

func TestIntervalTransformer_TransformBackward(t1 *testing.T) {
	tests := []struct {
		name      string
		interval1 Interval
		interval2 Interval
		arg       float64
		want      float64
	}{
		{
			name: "(-10,90)->(0,10), <- 6",
			interval1: Interval{
				Start: -10,
				End:   90,
			},
			interval2: Interval{
				Start: 0,
				End:   10,
			},
			arg:  6,
			want: 50,
		},
		{
			name: "(-10,90)->(0,10), <- 10",
			interval1: Interval{
				Start: -10,
				End:   90,
			},
			interval2: Interval{
				Start: 0,
				End:   10,
			},
			arg:  10,
			want: 90,
		},
		{
			name: "(-10,90)->(0,10), <- 0",
			interval1: Interval{
				Start: -10,
				End:   90,
			},
			interval2: Interval{
				Start: 0,
				End:   10,
			},
			arg:  0,
			want: -10,
		},
		{
			name: "(-10,90)->(0,10), <- 11",
			interval1: Interval{
				Start: -10,
				End:   90,
			},
			interval2: Interval{
				Start: 0,
				End:   10,
			},
			arg:  11,
			want: 90,
		},
	}
	for _, tt := range tests {
		t := &IntervalTransformer{
			IntervalFrom: tt.interval1,
			IntervalTo:   tt.interval2,
		}
		arg := tt.arg
		want := tt.want
		t1.Run(tt.name, func(t1 *testing.T) {
			if got := t.TransformBackward(arg); got != want {
				t1.Errorf("TransformBackward() = %v, want %v", got, want)
			}
		})
	}
}
