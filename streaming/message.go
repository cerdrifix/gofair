package streaming

import (
	"bytes"
	"encoding/json"
)

func addCRLF(data []byte) []byte {
	return append(data, []byte{'\r', '\n'}...)
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func scanCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{'\r', '\n'}); i >= 0 {
		// We have a full newline-terminated line.
		return i + 2, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

func getOp(buf []byte) (string, error) {

	// Unmarshal raw bytes to JSON
	tmp := make(map[string]json.RawMessage)
	var op string
	err := json.Unmarshal(buf, &tmp)
	if err != nil {
		return op, err
	}

	// Peek to see the op code
	err = json.Unmarshal(tmp["op"], &op)
	if err != nil {
		return op, err
	}

	return op, nil
}
