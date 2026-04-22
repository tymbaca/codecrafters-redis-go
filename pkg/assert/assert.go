package assert

func True(cond bool) {
	if !cond {
		panic("assertion failure")
	}
}
