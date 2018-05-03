package onering

import (
	"sync/atomic"
)

// my precious

type MPMC struct {
	multi
}

func (r *MPMC) init(size uint32) {
	r.multi.init(size)
	for i := range r.seq {
		r.seq[i] = -int64(i)
	}
	r.wp = r.size
	r.rp = r.size
}

func (r *MPMC) Get(i interface{}) bool {
	var (
		rp        = r.next(&r.rp)
		data, seq = r.frame(rp)
	)

	if !r.waitForEq(seq, rp) {
		return false
	}

	inject(i, *data)
	atomic.StoreInt64(seq, -rp)
	return true
}

func (r *MPMC) Consume(i interface{}) {
	var (
		fn = extractfn(i)
		it iter
	)
	for !it.stop {
		var (
			rp        = r.next(&r.rp)
			data, seq = r.frame(rp)
		)

		if !r.waitForEq(seq, rp) {
			break
		}

		fn(&it, *data)
		atomic.StoreInt64(seq, -rp)
	}
}

func (r *MPMC) Put(i interface{}) {
	var (
		wp        = r.next(&r.wp)
		data, seq = r.frame(wp)
	)
	for pread := r.size - wp; atomic.LoadInt64(seq) != pread; {
		r.wait()
	}

	*data = extractptr(i)
	atomic.StoreInt64(seq, wp)
}
