package main

// Remove removes s number of elements from the front of the slice and returns the removed elements
func Remove[T any](slice *[]T, s int) []T {
	var removedTiles = (*slice)[:s]
	*slice = (*slice)[s:]
	return removedTiles
}
