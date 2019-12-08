package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type apiResponse struct {
	NextChangeId string  `json:"next_change_id"`
	Stashes      []stash `json:"stashes"`
}

type stash struct {
	AccountName       string      `json:"accountName"`
	LastCharacterName string      `json:"lastCharacterName"`
	Id                string      `json:"id"`
	StashName         string      `json:"stash"`
	StashType         string      `json:"stashType"`
	Items             []stashItem `json:"items"`
	PublicStash       bool        `json:"public"`
}

type stashItem struct {
	Id           string   `json:"id"`
	League       string   `json:"league"`
	FrameType    int      `json:"frameType"`      // 3=unique, 5=currency, 6=divcard
	Note         string   `json:"note,omitempty"` // Price
	Name         string   `json:"name"`
	TypeLine     string   `json:"typeLine"`
	EnchantMods  []string `json:"enchantMods,omitempty"`
	ImplicitMods []string `json:"implicitMods,omitempty"`
	ExplicitMods []string `json:"explicitMods,omitempty"`
	CraftedMods  []string `json:"craftedMods,omitempty"`
	Elder        bool     `json:"elder,omitempty"`
	Shaper       bool     `json:"shaper,omitempty"`
	Veiled       bool     `json:"veiled,omitempty"`
	Corrupted    bool     `json:"corrupted,omitempty"`
	Ilvl         int      `json:"ilvl"`
	StackSize    int      `json:"stackSize"`
	MaxStackSize int      `json:"maxStackSize"`
}

func apiGet(lastChangeId string) *apiResponse {
	resp, err := http.Get("http://api.pathofexile.com/public-stash-tabs?id=" + lastChangeId)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var apiRes apiResponse
	err = json.Unmarshal(body, &apiRes)
	if err != nil {
		log.Fatal(err)
	}

	return &apiRes
}

func findItemsByCategoryAndLeague(tab stash, categories []int, league string) []stashItem {
	var stashItems []stashItem
	for _, item := range tab.Items {
		if item.League != league {
			continue
		}
		for _, category := range categories {
			if category == item.FrameType {
				if item.Note == "" {
					item.Note = tab.StashName
				}
				stashItems = append(stashItems, item)
			}
		}
	}
	return stashItems
}
