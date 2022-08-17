package vectorclock

type VectorClock map[string]int

func NewVectorClock(views []string) VectorClock {
	var vc = make(VectorClock)
	for _, view := range views {
		vc[view] = 0
	}
	return vc
}
