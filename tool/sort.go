package tool

type StringLenSort []string

func (p StringLenSort) Less(i, j int) bool { return len(p[i]) < len(p[j]) }
func (p StringLenSort) Len() int           { return len(p) }
func (p StringLenSort) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
