package validate

import (
	"strings"
	"strconv"
)

func ValidateCPF(document string) string {
	doc := strings.Split(document, "")
	d1, _ := strconv.Atoi(doc[9])
	d2, _ := strconv.Atoi(doc[10])
	testD1, testD2 := 0, d1*2
	eq := true
	for i := 0; i <= 9; i++ {
		eq = eq && doc[i] == doc[i+1]
	}
	for i := 0; i <= 8; i++ {
		r, _ := strconv.Atoi(doc[i])
		testD1 += r * (10 - i)
		testD2 += r * (11 - i)
	}
	if (testD1*10)%11 == d1 && (testD2*10)%11 == d2 && !eq {
		return ""
	} else {
		return  "CPF invalid"
	}
}

func ValidateCNPJ(document string) string {
	doc := document
	eq := true
	for i := 0; i <= 9; i++ {
		eq = eq && doc[i] == doc[i+1]
	}
	d := doc[:12]
	digit := calculateDigit(d, 5)
	d = d + digit
	digit = calculateDigit(d, 6)

	if doc == d+digit && !eq {
		return ""
	} else {
		return "CNPJ invalid"
	}
}

func calculateDigit(doc string, positions int) string {
	sum := 0
	for i := 0; i < len(doc); i++ {
		digit, _ := strconv.ParseInt(string(doc[i]), 10, 0)
		sum += int(digit) * positions
		positions--
		if positions < 2 {
			positions = 9
		}
	}
	sum %= 11
	if sum < 2 {
		return "0"
	}
	return strconv.FormatInt(int64(11-sum), 10)
}
