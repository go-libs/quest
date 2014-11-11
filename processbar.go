package quest

type ProgressBar struct {
	Qreq                     *Qrequest
	Progress                 func(current, total, expected int64)
	Total, Current, Expected int64
}

func (p *ProgressBar) Write(b []byte) (int, error) {
	current := len(b)
	// Response ContentLength
	if p.Total != -1 {
		p.calculate(int64(current))
		p.Progress(p.Current, p.Total, p.Expected)
	}
	return current, nil
}

func (p *ProgressBar) calculate(i int64) {
	p.Current += i
	p.Expected = p.Total - p.Current
}
