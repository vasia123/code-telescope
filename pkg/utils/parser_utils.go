package utils

// FindMatchingCloseBracket находит позицию закрывающей скобки, соответствующей открывающей скобке
// Эта функция может быть использована различными парсерами для обработки синтаксиса
func FindMatchingCloseBracket(str string, openPos int) int {
	if openPos >= len(str) || str[openPos] != '(' {
		return -1
	}

	depth := 1
	for i := openPos + 1; i < len(str); i++ {
		if str[i] == '(' {
			depth++
		} else if str[i] == ')' {
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}
