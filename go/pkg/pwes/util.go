package pwes

const TimeLayout = "2006-01-02 15:04:05"

// TODO: Apply a better function: compute score : 1 / (x/10 + 1) * 105 + 5
func computeScore(nbValidations int64) int64 {
	return int64(1/(float64(nbValidations)/10+1)*105 + 5)
}
