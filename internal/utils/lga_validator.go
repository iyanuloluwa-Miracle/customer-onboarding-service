package utils

var StateLGAMap = map[string][]string{
	"Lagos": {"Ikeja", "Epe", "Ikorodu"},
	"Oyo":   {"Ibadan North", "Ibadan South"},
}

func IsValidLGA(state, lga string) bool {
	for _, v := range StateLGAMap[state] {
		if v == lga {
			return true
		}
	}
	return false
}
