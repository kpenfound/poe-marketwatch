package main

import (
	"strings"
)

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
	OriginalQuantity int
	CalculatedPrice float64
}

func getUnique(si stashItem) uniqueItem {
	u := uniqueItem{ Name: si.Name}
	u.Corrupted = si.Corrupted
	u.OriginalPrice = si.Note
	u.CalculatedPrice = formatPrice(si.Note)
	return u
}

func getCurrency(si stashItem) currencyItem {
	c := currencyItem{Type: si.TypeLine}
	c.OriginalPrice = si.TypeLine
	c.CalculatedPrice = formatPrice(si.TypeLine)
	return c
}

func getDivination(si stashItem) divinationCardItem {
	d := divinationCardItem{Name: si.TypeLine}
	d.Mods = strings.Join(si.ExplicitMods[:], ",")
	d.MaxStackSize = si.MaxStackSize
	d.OriginalPrice = si.Note
	d.OriginalQuantity = si.StackSize
	d.CalculatedPrice = formatPrice(si.Note) / float64(si.StackSize)
	return d
}

func formatPrice(raw string) float64 {
	return 1.12
}
