package data

// Here's a complicated way to copy an array (of type [yCt][xCt]int, in this case). It requires on Go 1.17 for
// unsafe.Slice() and teh ability to cast a slice to an array pointer.
// It might or might not be faster than: arrCopy := arr. I haven't tested it (don't have 1.17 installed, so have only
// been able to check functionality with the Go Playground, but have not done execution timing).
//arrCopy := *(*[len(arr)][len(arr[0])]int)(unsafe.Slice(&arr[0], len(arr) * len(arr[0])))
