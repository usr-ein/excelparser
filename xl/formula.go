package xl

import (
	"errors"
)

type Formula string

func getFormula(s string) (Formula, error) {
	if len(s) > 1 && s[0] == '=' {
		return Formula(s), nil
	}
	return "", errors.New("not a formula")
}

func (f Formula) String() string {
	return string(f)
}
