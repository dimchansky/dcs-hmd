package ka50outputparser

import (
	"bytes"
	"fmt"
	"strconv"
	"unsafe"
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

const (
	rotorRPMArg   = 52
	rotorPitchArg = 53
)

// HandleMessage implements udplistener.MessageHandler interface.
func (p *OutputParser) HandleMessage(msg []byte) {
	pSimPrefix := parseSimPrefix(msg)
	msg = pSimPrefix.Rest

	// alternatively try parse semicolon prefix
	if !pSimPrefix.Ok {
		// skip ':'
		if len(msg) >= 1 && msg[0] == ':' {
			msg = msg[1:]
		}
	}

	for {
		pArg := parseArg(msg)
		msg = pArg.Rest

		if !pArg.Ok {
			return
		}

		arg := pArg.Result
		pVal := parseVal(msg)
		msg = pVal.Rest

		if !pVal.Ok {
			return
		}

		valBs := pVal.Result

		switch arg {
		case rotorRPMArg: // rotor RPM
			handleRotorRPM(p.s, valBs)

		case rotorPitchArg: // rotor pitch
			handleRotorPitch(p.s, valBs)
		}
	}
}

func handleRotorRPM(s ValuesSetter, valBs []byte) {
	val, err := strconv.ParseFloat(*(*string)(unsafe.Pointer(&valBs)), 64)
	if err != nil {
		return
	}

	const maxRotorRPM = 110.0

	s.SetRotorRPM(val * maxRotorRPM)
}

func handleRotorPitch(s ValuesSetter, valBs []byte) {
	val, err := strconv.ParseFloat(*(*string)(unsafe.Pointer(&valBs)), 64)
	if err != nil {
		return
	}

	const (
		maxRotorPitch = 15.0
		minRotorPitch = 1.0
	)

	s.SetRotorPitch(val*(maxRotorPitch-minRotorPitch) + minRotorPitch)
}

func parseSimPrefix(msg []byte) parserResult[uint64] {
	rest := msg

	var tmp []byte

	// Take until '*'
	pos := bytes.IndexByte(rest, '*')
	if pos >= 0 {
		tmp = rest[:pos]
		rest = rest[pos+1:]
	} else {
		return parseErr[uint64](msg, nil)
	}

	if tmpUint, err := strconv.ParseUint(*(*string)(unsafe.Pointer(&tmp)), 16, 64); err != nil {
		return parseErr[uint64](msg, fmt.Errorf("parsing `%s` into field Sim(hex): %w", *(*string)(unsafe.Pointer(&tmp)), err))
	} else {
		return parseOk[uint64](rest, tmpUint)
	}
}

// parse: Arg(uint) '='
func parseArg(msg []byte) parserResult[uint64] {
	rest := msg

	var tmp []byte

	// Take until '=' as Arg
	pos := bytes.IndexByte(rest, '=')
	if pos >= 0 {
		tmp = rest[:pos]
		rest = rest[pos+1:]
	} else {
		return parseErr[uint64](msg, nil)
	}

	if tmpUint, err := strconv.ParseUint(*(*string)(unsafe.Pointer(&tmp)), 10, 64); err != nil {
		return parseErr[uint64](msg, fmt.Errorf("parsing `%s` into field Arg(uint): %w", *(*string)(unsafe.Pointer(&tmp)), err))
	} else {
		return parseOk[uint64](rest, tmpUint)
	}
}

func parseVal(msg []byte) parserResult[[]byte] {
	rest := msg

	var val []byte

	// Checks if the rest starts with '\' and pass it
	if len(rest) >= 1 && rest[0] == '\'' {
		pTextVal := parseTextVal(msg)
		if pTextVal.Ok {
			return pTextVal
		}
	}

	// Take until ':' or '\r' or '\n' (or all the rest if not found) as Val
	if pos := bytes.IndexByte(rest, ':'); pos >= 0 {
		val = rest[:pos]
		rest = rest[pos+1:]
	} else {
		if pos = bytes.IndexByte(rest, '\r'); pos >= 0 {
			val = rest[:pos]
		} else if pos = bytes.IndexByte(rest, '\n'); pos >= 0 {
			val = rest[:pos]
		} else {
			val = rest
		}
		rest = rest[len(rest):]
	}

	return parseOk(rest, val)
}

func parseTextVal(msg []byte) parserResult[[]byte] {
	rest := msg

	var val []byte

	// Checks if the rest starts with '\' and pass it
	if len(rest) < 1 || rest[0] != '\'' {
		return parseErr[[]byte](msg, nil)
	}

	rest = rest[1:]

	// Take until '\'' as Val(string)
	pos := bytes.IndexByte(rest, '\'')
	if pos >= 0 {
		val = rest[:pos]
		rest = rest[pos+1:]
	} else {
		return parseErr[[]byte](msg, nil)
	}

	if len(rest) >= 1 {
		nextCh := rest[0]
		if nextCh == ':' {
			rest = rest[1:]
		} else if nextCh == '\r' || nextCh == '\n' {
			rest = rest[len(rest):]
		}
	}

	return parseOk(rest, val)
}

type parserResult[T any] struct {
	Result T
	Rest   []byte
	Ok     bool
	Err    error
}

func parseErr[T any](rest []byte, err error) parserResult[T] {
	return parserResult[T]{Rest: rest, Err: err}
}

func parseOk[T any](rest []byte, v T) parserResult[T] {
	return parserResult[T]{Result: v, Rest: rest, Ok: true}
}
