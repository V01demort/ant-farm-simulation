package main

import (
	"encoding/json"
	"os"
)
//Aruzhan
type jsonRoom struct {
	Name string `json:"name"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

type jsonLink struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type jsonExport struct {
	NumAnts int        `json:"num_ants"`
	Start   string     `json:"start"`
	End     string     `json:"end"`
	Rooms   []jsonRoom `json:"rooms"`
	Links   []jsonLink `json:"links"`
	Turns   [][]string `json:"turns"`
}

func WriteJSON(g *Graph, moves []string, path string) error {
	export := jsonExport{
		NumAnts: g.NumAnts,
		Start:   g.Start,
		End:     g.End,
	}

	for _, name := range g.RoomOrder {
		r := g.Rooms[name]
		export.Rooms = append(export.Rooms, jsonRoom{Name: r.Name, X: r.X, Y: r.Y})
	}

	seen := make(map[[2]string]bool)
	for _, a := range g.RoomOrder {
		for _, b := range g.Links[a] {
			k1 := [2]string{a, b}
			k2 := [2]string{b, a}
			if seen[k1] || seen[k2] {
				continue
			}
			seen[k1] = true
			export.Links = append(export.Links, jsonLink{From: a, To: b})
		}
	}

	for _, turn := range moves {
		export.Turns = append(export.Turns, splitTurnTokens(turn))
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(export)
}

func splitTurnTokens(turn string) []string {
	var tokens []string
	cur := ""
	for _, c := range turn {
		if c == ' ' {
			if cur != "" {
				tokens = append(tokens, cur)
				cur = ""
			}
		} else {
			cur += string(c)
		}
	}
	if cur != "" {
		tokens = append(tokens, cur)
	}
	return tokens
}
