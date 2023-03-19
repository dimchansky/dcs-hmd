package outputparser_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/dimchansky/dcs-hmd/aircraft/ka-50/outputparser"
	mocks "github.com/dimchansky/dcs-hmd/internal/mocks/aircraft/ka-50/outputparser"
)

func TestOutputParser_HandleMessage(t *testing.T) {
	testCases := []struct {
		message                  string
		expectedRotorPitch       *float64
		expectedRotorRPM         *float64
		expectedVerticalVelocity *float64
	}{
		{"637beb27*53=0.9362:52=0.7792\n", pFloat64(14.1068), pFloat64(85.712), nil},
		{"637beb27*52=0.7792\n", nil, pFloat64(85.712), nil},
		{"637beb27*53=0.9362\n", pFloat64(14.1068), nil, nil},

		{"637beb27*53=0.9362:52=0.7792\r\n", pFloat64(14.1068), pFloat64(85.712), nil},
		{"637beb27*52=0.7792\r\n", nil, pFloat64(85.712), nil},
		{"637beb27*53=0.9362\r\n", pFloat64(14.1068), nil, nil},

		{"637beb27*53=0.9362:52=0.7792", pFloat64(14.1068), pFloat64(85.712), nil},
		{"637beb27*52=0.7792", nil, pFloat64(85.712), nil},
		{"637beb27*53=0.9362", pFloat64(14.1068), nil, nil},

		{"637beb27*1000= :53=0.9362:52=0.7792\n", pFloat64(14.1068), pFloat64(85.712), nil},
		{"637beb27*1000= :52=0.7792\n", nil, pFloat64(85.712), nil},
		{"637beb27*1000= :53=0.9362\n", pFloat64(14.1068), nil, nil},

		{"637beb27*1000=:53=0.9362:52=0.7792\n", pFloat64(14.1068), pFloat64(85.712), nil},
		{"637beb27*1000=:52=0.7792\n", nil, pFloat64(85.712), nil},
		{"637beb27*1000=:53=0.9362\n", pFloat64(14.1068), nil, nil},

		{"637beb27*1000='53=0.1:52=0.2':53=0.9362:52=0.7792\n", pFloat64(14.1068), pFloat64(85.712), nil},
		{"637beb27*1000='53=0.1:52=0.2':52=0.7792\n", nil, pFloat64(85.712), nil},
		{"637beb27*1000='53=0.1:52=0.2':53=0.9362\n", pFloat64(14.1068), nil, nil},

		{"637beb27*53=0.9362:52=0.7792:1000='53=0.1:52=0.2'\n", pFloat64(14.1068), pFloat64(85.712), nil},
		{"637beb27*52=0.7792:1000='53=0.1:52=0.2'\n", nil, pFloat64(85.712), nil},
		{"637beb27*53=0.9362:1000='53=0.1:52=0.2'\n", pFloat64(14.1068), nil, nil},

		{"637beb27*53=0.9362:52=0.7792:1000='53=0.1:52=0.2'\r\n", pFloat64(14.1068), pFloat64(85.712), nil},
		{"637beb27*52=0.7792:1000='53=0.1:52=0.2'\r\n", nil, pFloat64(85.712), nil},
		{"637beb27*53=0.9362:1000='53=0.1:52=0.2'\r\n", pFloat64(14.1068), nil, nil},

		{":53=0.9362:52=0.7792\n", pFloat64(14.1068), pFloat64(85.712), nil},
		{":52=0.7792\n", nil, pFloat64(85.712), nil},
		{":53=0.9362\n", pFloat64(14.1068), nil, nil},

		{"53=0.9362:52=0.7792\n", pFloat64(14.1068), pFloat64(85.712), nil},
		{"52=0.7792\n", nil, pFloat64(85.712), nil},
		{"53=0.9362\n", pFloat64(14.1068), nil, nil},

		{"24=0.0000\n", nil, nil, pFloat64(0)},
		{"24=-1.0000\n", nil, nil, pFloat64(-30)},
		{"24=1.0000\n", nil, nil, pFloat64(30)},
	}

	for _, tt := range testCases {
		message := tt.message
		expectedRotorPitch := tt.expectedRotorPitch
		expectedRotorRPM := tt.expectedRotorRPM
		expectedVerticalVelocity := tt.expectedVerticalVelocity

		t.Run(message, func(t *testing.T) {
			testObj := &mocks.ValuesSetter{}

			if expectedRotorPitch != nil {
				testObj.On("SetRotorPitch", mock.AnythingOfType("float64"))
			}
			if expectedRotorRPM != nil {
				testObj.On("SetRotorRPM", mock.AnythingOfType("float64"))
			}
			if expectedVerticalVelocity != nil {
				testObj.On("SetVerticalVelocity", mock.AnythingOfType("float64"))
			}

			p := outputparser.New(testObj)
			p.HandleMessage([]byte(message))

			if expectedRotorPitch != nil {
				testObj.AssertNumberOfCalls(t, "SetRotorPitch", 1)
				testObj.AssertCalled(t, "SetRotorPitch", *expectedRotorPitch)
			}
			if expectedRotorRPM != nil {
				testObj.AssertNumberOfCalls(t, "SetRotorRPM", 1)
				testObj.AssertCalled(t, "SetRotorRPM", *expectedRotorRPM)
			}
			if expectedVerticalVelocity != nil {
				testObj.AssertNumberOfCalls(t, "SetVerticalVelocity", 1)
				testObj.AssertCalled(t, "SetVerticalVelocity", *expectedVerticalVelocity)
			}
		})
	}
}

func BenchmarkOutputParser_HandleMessage(b *testing.B) {
	vs := emptyValuesSetter{}
	p := outputparser.New(vs)
	msg := []byte("637beb27*53=0.9362:52=0.7792:24=-0.0100:1000=1.2345:1001=1.2345:1002=1.2345:1003=1.2345\n")

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		p.HandleMessage(msg)
	}
}

type emptyValuesSetter struct{}

func (s emptyValuesSetter) SetRotorPitch(float64)       {}
func (s emptyValuesSetter) SetRotorRPM(float64)         {}
func (s emptyValuesSetter) SetVerticalVelocity(float64) {}

func pFloat64(v float64) *float64 {
	return &v
}
