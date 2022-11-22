package ka50outputparser_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	mocks "github.com/dimchansky/dcs-hmd/internal/mocks/ka50outputparser"
	"github.com/dimchansky/dcs-hmd/ka50outputparser"
)

func TestOutputParser_HandleMessage(t *testing.T) {
	testCases := []struct {
		message            string
		expectedRotorPitch *float64
		expectedRotorRPM   *float64
	}{
		{"637beb27*53=0.9362:52=0.7792\n", pFloat64(14.1068), pFloat64(85.712)},
		{"637beb27*52=0.7792\n", nil, pFloat64(85.712)},
		{"637beb27*53=0.9362\n", pFloat64(14.1068), nil},
	}

	for _, tt := range testCases {
		message := tt.message
		expectedRotorPitch := tt.expectedRotorPitch
		expectedRotorRPM := tt.expectedRotorRPM

		t.Run(message, func(t *testing.T) {
			testObj := &mocks.ValuesSetter{}

			if expectedRotorPitch != nil {
				testObj.On("SetRotorPitch", mock.AnythingOfType("float64"))
			}
			if expectedRotorRPM != nil {
				testObj.On("SetRotorRPM", mock.AnythingOfType("float64"))
			}

			p := ka50outputparser.New(testObj)
			p.HandleMessage([]byte(message))

			if expectedRotorPitch != nil {
				testObj.AssertCalled(t, "SetRotorPitch", *expectedRotorPitch)
			}
			if expectedRotorRPM != nil {
				testObj.AssertCalled(t, "SetRotorRPM", *expectedRotorRPM)
			}
		})
	}
}

func pFloat64(v float64) *float64 {
	return &v
}
