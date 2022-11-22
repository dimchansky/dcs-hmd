package ka50outputparser

import (
	"regexp"
	"strconv"
)

type ValuesSetter interface {
	SetRotorPitch(val float64)
	SetRotorRPM(val float64)
}

func New(s ValuesSetter) *OutputParser {
	return &OutputParser{s: s}
}

type OutputParser struct {
	s ValuesSetter
}

var (
	r = regexp.MustCompile(`((\d+)=('[^']*'|[^:\r\n]*))`)
)

const (
	rotorRPMArg   = 52
	rotorPitchArg = 53
)

// HandleMessage implements udplistener.MessageHandler interface.
func (p *OutputParser) HandleMessage(msg []byte) {
	res := r.FindAllSubmatch(msg, -1)
	for _, rg := range res {
		argBs := rg[2]

		arg, err := strconv.Atoi(string(argBs))
		if err != nil {
			continue
		}

		valBs := rg[3]

		switch arg {
		case rotorRPMArg: // rotor RPM
			val, err := strconv.ParseFloat(string(valBs), 64)
			if err != nil {
				continue
			}

			const maxRotorRPM = 110.0

			p.s.SetRotorRPM(val * maxRotorRPM)

		case rotorPitchArg: // rotor pitch
			val, err := strconv.ParseFloat(string(valBs), 64)
			if err != nil {
				continue
			}

			const (
				maxRotorPitch = 15.0
				minRotorPitch = 1.0
			)

			p.s.SetRotorPitch(val*(maxRotorPitch-minRotorPitch) + minRotorPitch)
		}
	}
}
