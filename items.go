package main

type uniqueItem struct {
	Name string
	Corrupted bool
	OriginalPrice string
	CalculatedPrice float64
}

type currencyItem struct {
	Type string
	OriginalPrice string
	CalculatedPrice float64
}

type divinationCardItem struct {
	Name string
	Mods string
	MaxStackSize int
	OriginalPrice string
	CalculatedPrice float64
}
