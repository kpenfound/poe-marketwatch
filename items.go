package main

import (
	"strings"
)

type uniqueItem struct {
	Id string
	Name string
	Corrupted bool
	OriginalPrice price
}

type currencyItem struct {
	Id string
	Type string
	OriginalPrice price
	OriginalQuantity int
}

type divinationCardItem struct {
	Id string
	Name string
	Mods string
	MaxStackSize int
	OriginalPrice price
	OriginalQuantity int
}

func getUniqueFromStashItem(si stashItem) uniqueItem {
	u := uniqueItem{ Name: si.Name}
	u.Id = si.Id
	u.Corrupted = si.Corrupted
	u.OriginalPrice = interpretPrice(si.Note)
	return u
}

func getCurrencyFromStashItem(si stashItem) currencyItem {
	c := currencyItem{Type: si.TypeLine}
	c.Id = si.Id
	c.Type = si.TypeLine
	c.OriginalPrice = interpretPrice(si.Note)
	c.OriginalQuantity = si.StackSize
	return c
}

func getDivinationFromStashItem(si stashItem) divinationCardItem {
	d := divinationCardItem{Name: si.TypeLine}
	d.Id = si.Id
	d.Mods = strings.Join(si.ExplicitMods[:], ",")
	d.MaxStackSize = si.MaxStackSize
	d.OriginalPrice = interpretPrice(si.Note)
	d.OriginalQuantity = si.StackSize
	return d
}

func formatPrice(raw string) float64 {
	return 1.12
}
