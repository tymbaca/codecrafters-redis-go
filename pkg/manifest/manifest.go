package manifest

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/pkg/iter"
)

type Manifest struct {
	Records []Record
}

type Record struct {
	File string
	Seq  int
	Type string

	StartOffsetSet bool
	StartOffset    int
}

func Decode(r io.Reader) (Manifest, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return Manifest{}, err
	}

	var records []Record

	dataStr := string(data)
	for line := range strings.SplitSeq(dataStr, "\n") {
		parts := strings.Split(line, " ")
		iter := iter.New(parts)

		fileKeyword, ok := iter.Next()
		if !ok || fileKeyword != "file" {
			return Manifest{}, fmt.Errorf("expected 'file' at the beginning of line, got '%s'", fileKeyword)
		}

		file, ok := iter.Next()
		if !ok {
			return Manifest{}, fmt.Errorf("expected file name after 'file'")
		}

		seqKeyword, ok := iter.Next()
		if !ok || seqKeyword != "seq" {
			return Manifest{}, fmt.Errorf("expected 'seq', got '%s'", seqKeyword)
		}

		seqStr, ok := iter.Next()
		if !ok {
			return Manifest{}, fmt.Errorf("expected seq number after 'seq'")
		}

		seq, err := strconv.Atoi(seqStr)
		if err != nil {
			return Manifest{}, fmt.Errorf("expected seq to be a number: %w", err)
		}

		typeKeyword, ok := iter.Next()
		if !ok || typeKeyword != "type" {
			return Manifest{}, fmt.Errorf("expected 'type', got '%s'", typeKeyword)
		}

		typ, ok := iter.Next()
		if !ok {
			return Manifest{}, fmt.Errorf("expected type name after 'type'")
		}

		rec := Record{
			File: file,
			Seq:  seq,
			Type: typ,
		}

		opt, ok := iter.Next()
		if ok {
			switch opt {
			case "startoffset":
				valStr, ok := iter.Next()
				if !ok {
					return Manifest{}, fmt.Errorf("expected number after 'startoffset'")
				}

				val, err := strconv.Atoi(valStr)
				if err != nil {
					return Manifest{}, fmt.Errorf("expected number after 'startoffset': %w", err)
				}

				rec.StartOffsetSet = true
				rec.StartOffset = val
			}
		}

		records = append(records, rec)
	}

	return Manifest{
		Records: records,
	}, nil
}

func (m Manifest) Encode(w io.Writer) error {
	buf := bytes.NewBuffer(nil)

	for _, r := range m.Records {
		fmt.Fprintf(buf, "file %s seq %d type %s", r.File, r.Seq, r.Type)

		if r.StartOffsetSet {
			fmt.Fprintf(buf, " startoffset %d", r.StartOffset)
		}

		buf.WriteString("\n")
	}

	_, err := buf.WriteTo(w)
	return err
}
