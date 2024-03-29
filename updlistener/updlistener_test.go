package updlistener

import (
	"bufio"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/require"
)

func Test_readLines(t *testing.T) {
	minReadBufferSize := bufio.NewReaderSize(nil, 16).Size()
	require.Equal(t, 16, minReadBufferSize, "adjust tests for the new minReadBufferSize")

	tests := []struct {
		name       string
		bufSize    int
		input      string
		wantOutput []string
	}{
		{"1", 16, "123\n456\n789", []string{"123", "456", "789"}},
		{"2", 16, "123\n456\n789\n", []string{"123", "456", "789"}},
		{"3", 16, "123\r\n456\r\n789", []string{"123", "456", "789"}},
		{"4", 16, "123\r\n456\r\n789\r\n", []string{"123", "456", "789"}},
		{"5", 16, "123\n1234567890123456\n789", []string{"123", "789"}},
		{"6", 16, "123\n1234567890123456\n789\n", []string{"123", "789"}},
		{"7", 16, "123\r\n123456789012345\r\n789", []string{"123", "789"}},
		{"8", 16, "123\r\n123456789012345\r\n789\r\n", []string{"123", "789"}},
	}
	for _, tt := range tests {
		bufSize := tt.bufSize
		input := tt.input

		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReaderSize(
				iotest.OneByteReader(strings.NewReader(input)),
				bufSize,
			)

			var msgHandler messageCollector
			readLines(reader, &msgHandler)

			require.Equal(t, tt.wantOutput, msgHandler.AsSlice())
		})
	}
}

func Benchmark_readLines(b *testing.B) {
	rowsReader := simpleRowsReader(b.N)
	reader := bufio.NewReaderSize(&rowsReader, 16)
	msgHandler := devNullMessageHandler{}

	b.ReportAllocs()
	b.ResetTimer()

	readLines(reader, msgHandler)
}

type messageCollector []string

func (m *messageCollector) HandleMessage(msg []byte) { *m = append(*m, string(msg)) }
func (m *messageCollector) AsSlice() []string        { return *m }

type devNullMessageHandler struct{}

func (d devNullMessageHandler) HandleMessage([]byte) {}

type simpleRowsReader int

func (r *simpleRowsReader) Read(b []byte) (int, error) {
	if *r < 0 {
		return 0, io.EOF
	}

	*r -= 1
	b[0] = '1'
	b[1] = '\n'

	return 2, nil
}
