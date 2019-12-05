package main

import (
	"strings"
)

type uniqueItem struct {
	Id string
	Name string
	Corrupted bool
	OriginalPrice string
	CalculatedPrice float64
}

type currencyItem struct {
	Id string
	Type string
	OriginalPrice string
	CalculatedPrice float64
}

type divinationCardItem struct {
	Id string
	Name string
	Mods string
	MaxStackSize int
	OriginalPrice string
	OriginalQuantity int
	CalculatedPrice float64
}

func getUniqueFromStashItem(si stashItem) uniqueItem {
	u := uniqueItem{ Name: si.Name}
	u.Id = si.Id
	u.Corrupted = si.Corrupted
	u.OriginalPrice = si.Note
	u.CalculatedPrice = formatPrice(si.Note)
	return u
}

func getCurrencyFromStashItem(si stashItem) currencyItem {
	c := currencyItem{Type: si.TypeLine}
	c.Id = si.Id
	c.OriginalPrice = si.TypeLine
	c.CalculatedPrice = formatPrice(si.TypeLine)
	return c
}

func getDivinationFromStashItem(si stashItem) divinationCardItem {
	d := divinationCardItem{Name: si.TypeLine}
	d.Id = si.Id
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
