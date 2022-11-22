package ka50outputparser

import (
	"fmt"
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

// HandleMessage implements udplistener.MessageHandler interface.
func (p *OutputParser) HandleMessage(msg []byte) {
	res := r.FindAllSubmatch(msg, -1)
	for _, rg := range res {
		argBs := rg[2]
		arg, err := strconv.Atoi(string(argBs))
		if err != nil {
			fmt.Println("failed to parse argument", err)
			continue
		}
		valBs := rg[3]
		switch arg {
		case 52: // rotor RPM
			val, err := strconv.ParseFloat(string(valBs), 64)
			if err != nil {
				fmt.Println("failed to parse rotor rpm", err)
				continue
			}
			p.s.SetRotorRPM(val * 110.0)
		case 53: // rotor pitch
			val, err := strconv.ParseFloat(string(valBs), 64)
			if err != nil {
				fmt.Println("failed to parse rotor pitch", err)
				continue
			}
			p.s.SetRotorPitch(val*(15.0-1.0) + 1.0)
		}
	}
}
