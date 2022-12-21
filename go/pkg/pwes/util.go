package pwes

// TODO: Apply a better function: compute score : 1 / (x/10 + 1) * 95 + 5
func computeScore(nbValidations int64) int64 {
	return 1/(nbValidations/10+1)*95 + 5
}
