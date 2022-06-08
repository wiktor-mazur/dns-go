package utils

func BoolToByte(value bool) byte {
	var result byte

	if value {
		result = 1
	}

	return result
}
