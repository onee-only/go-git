package idxfile_test

import (
	"bytes"
	"io"

	. "github.com/go-git/go-git/v6/plumbing/format/idxfile"

	fixtures "github.com/go-git/go-git-fixtures/v5"
)

func (s *IdxfileSuite) TestDecodeEncode() {
	for _, f := range fixtures.ByTag("packfile") {
		expected, err := io.ReadAll(f.Idx())
		s.NoError(err)

		idx := new(MemoryIndex)
		d := NewDecoder(bytes.NewBuffer(expected))
		err = d.Decode(idx)
		s.NoError(err)

		result := bytes.NewBuffer(nil)
		e := NewEncoder(result)
		size, err := e.Encode(idx)
		s.NoError(err)

		s.Len(expected, size)
		s.Equal(expected, result.Bytes())
	}
}
