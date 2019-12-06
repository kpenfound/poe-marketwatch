package main

import (
	"regexp"
	"strings"
	"strconv"
	"log"
)

type price struct {
	Price float64
	Currency string
}

func init_currency() *regexp.Regexp {
	re := regexp.MustCompile(`(~price|~b\/o) (\d+|\d+\/\d+|\d+\.\d+) \w+`)

	return re
}

func isPrice(rawPrice string, re *regexp.Regexp) bool {
	return re.MatchString(rawPrice)
}

func interpretPrice(rawPrice string) price {
	rawParts := strings.Split(rawPrice, " ")
	pr := float64(0)
	if strings.ContainsAny(rawParts[1], "/") {
		fractionParts := strings.Split(rawParts[1], "/")

		fractionNum, err := strconv.ParseFloat(fractionParts[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		fractionDen, err := strconv.ParseFloat(fractionParts[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		pr = fractionNum / fractionDen
	} else {
		parsed, err := strconv.ParseFloat(rawParts[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		pr = parsed
	}

	p := price{Price: pr, Currency: rawParts[2]}

	return p
}
