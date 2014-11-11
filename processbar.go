package quest

type ProgressBar struct {
	Qreq     *Qrequest
	Total    int64
	Current  int64
	Progress func(current, total, expected int64)
}

func (p *ProgressBar) Write(b []byte) (int, error) {
	current := len(b)
	// Response ContentLength
	if p.Total != -1 {
		p.Current += int64(current)
		p.Progress(p.Current, p.Total, p.Total-p.Current)
	}
	return current, nil
}
