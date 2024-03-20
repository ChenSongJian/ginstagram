package utils

func IsInIntSlice(target int, slice []int) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}
