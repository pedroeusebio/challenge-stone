package validate

import (
	"app/model"
	"errors"
	"strconv"
	"strings"

	validator "gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

// funcao de validacao do CPF
// retorna um string de CPF invalido caso o cpf nao seja valido

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
		return "CPF invalid"
	}
}

// funcao de validacao do CNPJ

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

// funcao de validacao dos dados do invoice
// foi utilizado o validator.v9 para facilitar a validacao
// funcao feita para retornar a mensagem de erro compeensivel ao usuario

func ValidateInvoice(invoice model.Invoice) []string {
	validate = validator.New()
	error := []string{}
	Err := validate.Struct(invoice)
	if Err != nil {
		for _, err := range Err.(validator.ValidationErrors) {
			switch tag := err.Tag(); tag {
			case "required":
				error = append(error, err.Field()+": is required ")
			case "gte":
				var gte string
				switch field := err.Field(); field {
				case "Amount":
					gte = model.GteAmount
				case "Month":
					gte = model.GteMonth
				}
				error = append(error, err.Field()+": must be greater than or equals to "+gte+" ")
			case "gt":
				error = append(error, err.Field()+": must be greater than "+model.GtYear+" ")
			case "lte":
				error = append(error, err.Field()+": must be less than or equals to "+model.LteMonth+" ")
			case "numeric":
				error = append(error, err.Field()+": must be only numbers")
			case "len=11|len=14":
				error = append(error, err.Field()+": must have "+model.GteDocument+" or "+model.LteDocument+" digits")
			}
		}
	} else {
		var dErr string
		if len(invoice.Document) == 11 {
			dErr = ValidateCPF(invoice.Document)
		} else {
			dErr = ValidateCNPJ(invoice.Document)
		}
		if len(dErr) > 0 {
			error = append(error, "Document: "+dErr)
		}
	}
	return error
}

// funcao de validacao dos dados do user
// foi utilizado o validator.v9 para facilitar a validacao
// funcao feita para retornar a mensagem de erro compeensivel ao usuario

func ValidateUser(user model.User) error {
	validate = validator.New()
	var error string
	vErr := validate.Struct(user)
	if vErr != nil {
		for _, err := range vErr.(validator.ValidationErrors) {
			if len(error) > 0 {
				error += ", "
			}
			if err.Tag() == "required" {
				error += err.Field() + ": is required "
			}
			if err.Tag() == "alphanum" || err.Tag() == "excludesall" {
				error += err.Field() + ": contains invalid characters "
			}
			if err.Tag() == "gt" {
				var gt string
				if err.Field() == "Name" {
					gt = model.GtName
				} else {
					gt = model.GtPassword
				}
				error += err.Field() + ": must have more than " + gt + " characters"
			}
		}
		rErr := errors.New(error)
		return rErr
	} else {
		return nil
	}
}
