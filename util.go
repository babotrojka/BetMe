package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var jsonContent = []byte(`{"success":true,"data":[{"key":"americanfootball_ncaaf","active":true,"group":"American Football","details":"US College Football","title":"NCAAF","has_outrights":false},{"key":"americanfootball_nfl","active":true,"group":"American Football","details":"US Football","title":"NFL","has_outrights":false},{"key":"aussierules_afl","active":true,"group":"Aussie Rules","details":"Aussie Football","title":"AFL","has_outrights":false},{"key":"baseball_mlb","active":true,"group":"Baseball","details":"Major League Baseball \ud83c\uddfa\ud83c\uddf8","title":"MLB","has_outrights":false},{"key":"basketball_nba","active":true,"group":"Basketball","details":"US Basketball","title":"NBA","has_outrights":false},{"key":"cricket_odi","active":true,"group":"Cricket","details":"One Day Internationals","title":"One Day Internationals","has_outrights":false},{"key":"cricket_test_match","active":true,"group":"Cricket","details":"International Test Matches","title":"Test Matches","has_outrights":false},{"key":"esports_lol","active":true,"group":"Esports","details":"LCS, LEC, LCK, LPL","title":"League of Legends","has_outrights":false},{"key":"icehockey_nhl","active":true,"group":"Ice Hockey","details":"US Ice Hockey","title":"NHL","has_outrights":false},{"key":"mma_mixed_martial_arts","active":true,"group":"Mixed Martial Arts","details":"Mixed Martial Arts","title":"MMA","has_outrights":false}]}`)

func prva() {
	var f interface{}
	err := json.Unmarshal(jsonContent, &f)
	if err != nil {
		fmt.Println("Error in Unmarshal")
	}
	m := f.(map[string]interface{})
	for k, v := range m {
		switch vv := v.(type) {
		case []interface{}:
			fmt.Println(k, "is an array:")
			for _, u := range vv {
				um := u.(map[string]interface{})
				fmt.Printf("type of um is %T\n", um)
				for kljuc, vrij := range um {
					switch uu := vrij.(type) {
					case []interface{}:
						fmt.Println(kljuc, "is an array:")
					case string:
						fmt.Println(kljuc, "is string", uu)
					case int:
						fmt.Println(kljuc, "is int", u)
					case bool:
						fmt.Println(kljuc, "is bool", uu)
					default:
						fmt.Println(kljuc, "is have no idea what")
					}
				}
			}
		case string:
			fmt.Println(k, "is string", vv)
		case int:
			fmt.Println(k, "is int", v)
		case bool:
			fmt.Println(k, "is bool", vv)
		default:
			fmt.Println(k, "is have no idea what")
		}
	}
}

type SportLeague struct {
	Key           string
	Active        bool
	Group         string
	Details       string
	Title         string
	Has_outrights bool
}
type Response struct {
	Success bool          `json:"success"`
	Data    []interface{} `json:"data"`
}

func druga() map[string][]SportLeague {
	var response Response
	err := json.Unmarshal(jsonContent, &response)
	if err != nil {
		fmt.Println("greska kod Unmarshal\n")
		panic(err)
	} else {
		data := response.Data
		var sportLeague []SportLeague
		marshJson, _ := json.Marshal(data)
		err = json.Unmarshal(marshJson, &sportLeague)
		if err != nil {
			panic(err)
		}
		mapa := make(map[string][]SportLeague)
		for _, league := range sportLeague {
			if x, found := mapa[league.Group]; found == false {
				s := make([]SportLeague, 0, 10)
				s = append(s, league)
				mapa[league.Group] = s
			} else {
				s := append(x, league)
				mapa[league.Group] = s
			}
		}
		return mapa
	}
	return nil
}

type Response2 struct {
	Success bool          `json:"success"`
	Data    []SportLeague `json:"data"`
}

func treca() map[string][]SportLeague {
	var response2 Response2
	err := json.Unmarshal(jsonContent, &response2)
	if err != nil {
		panic(err)
	} else {
		mapa := make(map[string][]SportLeague)
		for _, league := range response2.Data {
			if x, found := mapa[league.Group]; found == false {
				s := make([]SportLeague, 0, 10)
				s = append(s, league)
				mapa[league.Group] = s
			} else {
				s := append(x, league)
				mapa[league.Group] = s
			}
		}
		fmt.Println(mapa)
		return mapa
	}
	return nil
}

type Response3 struct {
	Success bool   `json:"success"`
	Games   []Game `json:"data"`
}
type Game struct {
	Sport_key     string   `json:"sport_key"`
	Sport_nice    string   `json:"sport_nice"`
	Teams         []string `json:"teams"`
	Commence_time int64    `json:"commence_time"`
	Home_Team     string   `json:"home_team"`
	Sites         []Site   `json:"sites"`
	Sites_count   int      `json:"sites_count"`
}
type Site struct {
	Site_key    string      `json:"site_key"`
	Site_nice   string      `json:"site_nice"`
	Last_update int64       `json:"last_update"`
	Odds        interface{} `json:"odds"`
}
type Odd struct {
	H2h []float64 `json:"h2h"`
}

func sportPrva() []Game {
	bodyFromFile, _ := ioutil.ReadFile("tests/cricket.txt")
	var responser3 Response3
	err := json.Unmarshal(bodyFromFile, &responser3)
	if err != nil {
		panic(err)
	} else {
		return responser3.Games
	}
	return nil
}
func main() {
	games := sportPrva()
	for _, u := range games {
		fmt.Print(u.Teams[0] + " " + "vs" + " " + u.Teams[1] + "\t")
		mapa := u.Sites[0].Odds.(map[string]interface{})
		for _, oneInterface := range mapa["h2h"].([]interface{}) {
			fmt.Printf("Value is %f, type is %T\n", oneInterface, oneInterface)
		}
	}
}
