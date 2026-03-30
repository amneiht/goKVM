package util

func ClearSlice(sd []byte, remain int) []byte {

	lg := len(sd) - remain
	for i := range remain {
		sd[i] = sd[lg]
		lg++
	}
	// strip slide , reuse buffer
	return sd[:remain]
}
