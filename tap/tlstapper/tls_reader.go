package tlstapper

import (
	"io"
	"time"

	"github.com/up9inc/mizu/tap/api"
)

type tlsReader struct {
	key         string
	chunks      chan *tlsChunk
	data        []byte
	doneHandler func(r *tlsReader)
	progress    *api.ReadProgress
	timer        api.SuperTimer
}

func (r *tlsReader) Read(p []byte) (int, error) {
	var chunk *tlsChunk

	for len(r.data) == 0 {
		var ok bool
		select {
		case chunk, ok = <-r.chunks:
			if !ok {
				return 0, io.EOF
			}

			r.data = chunk.getRecordedData()
		case <-time.After(time.Second * 3):
			r.doneHandler(r)
			return 0, io.EOF
		}

		if len(r.data) > 0 {
			break
		}
	}

	l := copy(p, r.data)
	r.data = r.data[l:]
	r.progress.Feed(l)

	return l, nil
}
