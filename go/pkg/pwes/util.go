package pwes

const TimeLayout = "2006-01-02 15:04:05"

// TODO: Apply a better function: compute score : 1 / (x/10 + 1) * 95 + 5
func computeScore(nbValidations int64) int64 {
	return 1/(nbValidations/10+1)*95 + 5
}
