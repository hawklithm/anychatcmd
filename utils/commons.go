package utils

func If(b bool, t, f interface{}) interface{} {
	if b {
		return t
	}
	return f
}
