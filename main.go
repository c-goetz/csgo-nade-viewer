package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Map string

const (
	MapInferno  Map = "inferno"
	MapDust2    Map = "dust2"
	MapAncient  Map = "ancient"
	MapMirage   Map = "mirage"
	MapTrain    Map = "train"
	MapNuke     Map = "nuke"
	MapOverpass Map = "overpass"
)

var allMaps = []Map{
	MapInferno,
	MapDust2,
	MapAncient,
	MapMirage,
	MapTrain,
	MapNuke,
	MapOverpass,
}

func toMap(s string) Map {
	switch s {
	case "inferno":
		return MapInferno
	case "dust2":
		return MapDust2
	case "ancient":
		return MapAncient
	case "mirage":
		return MapMirage
	case "train":
		return MapTrain
	case "nuke":
		return MapNuke
	case "overpass":
		return MapOverpass
	default:
		panic("unknown map: " + s)
	}
}

type Side string

const (
	SideT  Side = "t"
	SideCT Side = "ct"
)

var allSides = []Side{
	SideT,
	SideCT,
}

func toSide(s string) Side {
	switch s {
	case "t":
		return SideT
	case "ct":
		return SideCT
	default:
		panic("unknown side: " + s)
	}
}

type ThrowMod string

const (
	ThrowModRegular    ThrowMod = "regular"
	ThrowModLeftClick  ThrowMod = "lc"
	ThrowModRightClick ThrowMod = "rc"
	ThrowModW          ThrowMod = "w"
	ThrowModA          ThrowMod = "a"
	ThrowModD          ThrowMod = "d"
	ThrowModShift      ThrowMod = "shift"
	ThrowModJump       ThrowMod = "jump"
)

func toThrowMod(s string) ThrowMod {
	switch s {
	case "regular":
		return ThrowModRegular
	case "lc":
		return ThrowModLeftClick
	case "rc":
		return ThrowModRightClick
	case "w":
		return ThrowModW
	case "a":
		return ThrowModA
	case "d":
		return ThrowModD
	case "shift":
		return ThrowModShift
	case "jump":
		return ThrowModJump
	default:
		panic("unknown throwMod: " + s)
	}
}

type Nade struct {
	Map       Map
	Side      Side
	ThrowMods []ThrowMod
	Name      string
	Images    []string
}

func main() {
	entries, err := os.ReadDir("./img")
	if err != nil {
		panic(err)
	}
	nades := make(map[Map]map[Side]map[string]*Nade)
	for _, m := range allMaps {
		nades[m] = make(map[Side]map[string]*Nade)
		for _, s := range allSides {
			nades[m][s] = make(map[string]*Nade)
		}
	}
	for i := range entries {
		p := parseImage(entries[i].Name())
		if n, ok := nades[p.Map][p.Side][p.Name]; ok {
			n.Images = append(n.Images, p.Image)
		} else {
			nades[p.Map][p.Side][p.Name] = &Nade{
				Name:      p.Name,
				Map:       p.Map,
				Side:      p.Side,
				ThrowMods: p.ThrowMods,
				Images:    []string{p.Image},
			}
		}
	}
	tmpl := template.Must(template.ParseFiles("index-template.html"))
	f, err := os.Create("index.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = tmpl.Execute(f, T{
		allMaps,
		allSides,
		nades,
	})
	if err != nil {
		panic(err)
	}
}

type T struct {
	AllMaps  []Map
	AllSides []Side
	Nades    map[Map]map[Side]map[string]*Nade
}

type ParsedNade struct {
	Map       Map
	Side      Side
	ThrowMods []ThrowMod
	Name      string
	Image     string
}

func parseImage(n string) ParsedNade {
	ext := filepath.Ext(n)
	base, _ := strings.CutSuffix(n, ext)
	// <map>_<side>_<mod>-<mod>_<name>_<index>
	parts := strings.Split(base, "_")
	if len(parts) != 5 {
		log.Fatalf("invalid filename (not enough parts): %s", n)
	}
	fmt.Println(n)
	return ParsedNade{
		Map:       toMap(parts[0]),
		Side:      toSide(parts[1]),
		ThrowMods: parseThrowMods(parts[2]),
		Name:      strings.ReplaceAll(parts[3], "-", " "),
		Image:     n,
	}
}

func parseThrowMods(s string) []ThrowMod {
	parts := strings.Split(s, "-")
	r := make([]ThrowMod, len(parts))
	for i := range parts {
		p := parts[i]
		r[i] = toThrowMod(p)
	}
	return r
}
