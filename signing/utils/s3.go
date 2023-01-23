package utils

func URIEncode(input string) string {
	const upperhex = "0123456789ABCDEF"

	hexCount := 0

	for i := 0; i < len(input); i++ {
		if shouldEscape(input[i]) {
			hexCount++
		}
	}

	output := make([]byte, len(input)+2*hexCount)

	j := 0
	for i := 0; i < len(input); i++ {
		switch c := input[i]; {
		case c == ' ':
			output[j] = '+'
			j++
		case shouldEscape(c):
			output[j] = '%'
			output[j+1] = upperhex[c>>4]
			output[j+2] = upperhex[c&15]
			j += 3
		default:
			output[j] = input[i]
			j++
		}
	}

	return string(output)
}

func shouldEscape(c byte) bool {
	switch {
	case 'a' <= c && c <= 'z':
	case 'A' <= c && c <= 'Z':
	case '0' <= c && c <= '9':
	case c == '-':
	case c == '_':
	case c == '.':
	case c == '~':
	case c == '/':
	default:
		return true
	}

	return false
}
