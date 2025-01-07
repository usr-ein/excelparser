package xl

import "strings"

type Relativeness struct {
	Row bool
	Col bool
}

func getRelativeness(localAddress LocalAddress) Relativeness {
	/**
	 * Determines if the address is relative or absolute.
	 * E.g. A1 is absolute, $A$1 is absolute, A$1 is mixed, $A1 is mixed.
	 */
	relativeness := Relativeness{
		Row: true,
		Col: true,
	}

	dollarCount := strings.Count(string(localAddress), "$")
	if dollarCount > 0 {
		relativeness.Col = !strings.HasPrefix(string(localAddress), "$")

		if (dollarCount == 2 && !relativeness.Col) || (dollarCount == 1 && relativeness.Col) {
			relativeness.Row = false
		}
	}

	return relativeness
}
