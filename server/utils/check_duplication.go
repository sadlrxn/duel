package utils

func IsDuplicateInArray[T any](arr []T) bool {
	visited := make(map[interface{}]bool, 0)
	for i := 0; i < len(arr); i++ {
		if visited[arr[i]] {
			return true
		} else {
			visited[arr[i]] = true
		}
	}
	return false
}
