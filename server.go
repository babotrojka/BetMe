package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/thedevsaddam/renderer"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

var rnd *renderer.Render
var keys []string
var tiket []string
var kof float64
var store = sessions.NewCookieStore([]byte("cookieSifra"))
var check bool
var username string

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "bazepodataka"
	dbname   = "bet_me"
)

const API_KEY = "35a53e6ac5e21dcb213cc9e3d65f146e"

var jsonContent = []byte(`{"success":true,"data":[{"key":"americanfootball_ncaaf","active":true,"group":"American Football","details":"US College Football","title":"NCAAF","has_outrights":false},{"key":"americanfootball_nfl","active":true,"group":"American Football","details":"US Football","title":"NFL","has_outrights":false},{"key":"aussierules_afl","active":true,"group":"Aussie Rules","details":"Aussie Football","title":"AFL","has_outrights":false},{"key":"baseball_mlb","active":true,"group":"Baseball","details":"Major League Baseball \ud83c\uddfa\ud83c\uddf8","title":"MLB","has_outrights":false},{"key":"basketball_nba","active":true,"group":"Basketball","details":"US Basketball","title":"NBA","has_outrights":false},{"key":"cricket_odi","active":true,"group":"Cricket","details":"One Day Internationals","title":"One Day Internationals","has_outrights":false},{"key":"cricket_test_match","active":true,"group":"Cricket","details":"International Test Matches","title":"Test Matches","has_outrights":false},{"key":"esports_lol","active":true,"group":"Esports","details":"LCS, LEC, LCK, LPL","title":"League of Legends","has_outrights":false},{"key":"icehockey_nhl","active":true,"group":"Ice Hockey","details":"US Ice Hockey","title":"NHL","has_outrights":false},{"key":"mma_mixed_martial_arts","active":true,"group":"Mixed Martial Arts","details":"Mixed Martial Arts","title":"MMA","has_outrights":false}]}`)

func init() {
	rnd = renderer.New(
		renderer.Options{
			ParseGlobPattern: "./static/template/*.html",
		},
	)
	check = false
	username = ""
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

func readSports(jsonContent []byte) map[string][]SportLeague {
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

type PageData struct {
	Check  bool
	Sports []string
	User   string
	Msg    string
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
type PageSportData struct {
	Is3   bool
	Team1 string
	Team2 string
	Odd1  float64
	OddX  float64
	Odd2  float64
}

var sportsData []PageSportData

type CompletePageSportData struct {
	Check     bool
	User      string
	SportData []PageSportData
	Sports    []string
	Tiket     []string
	Kof       string
}

var session *sessions.Session

func sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ = store.Get(r, "ivan")
		next.ServeHTTP(w, r)
	})
}
func home(w http.ResponseWriter, r *http.Request) {
	//	json, _ := http.Get("https://api.the-odds-api.com/v3/sports/?apiKey=" + API_KEY)
	//	resp, _ := ioutil.ReadAll(json.Body)
	//	sportsMap := readSports(resp)
	data := PageData{
		Check:  check,
		User:   username,
		Sports: keys,
		Msg:    "",
	}
	err := rnd.HTML(w, http.StatusOK, "home", data)
	if err != nil {
		panic(err)
	}
}

func citajParove(key string) []Game {
	//bodyFromFile, _ := ioutil.ReadFile("tests/soccer.txt")
	jsonR, _ := http.Get("https://api.the-odds-api.com/v3/odds/?apiKey=" + API_KEY + "&sport=" + key + "&region=eu")
	bodyFromFile, _ := ioutil.ReadAll(jsonR.Body)
	var responser3 Response3
	err := json.Unmarshal(bodyFromFile, &responser3)
	if err != nil {
		panic(err)
	} else {
		return responser3.Games
	}
	return nil
}
func paroviUPageSports(games []Game) []PageSportData {
	data := make([]PageSportData, 0, len(games))
	for _, u := range games {
		var oneData PageSportData
		h2hMap := u.Sites[0].Odds.(map[string]interface{})
		switch len(h2hMap["h2h"].([]interface{})) {
		case 3:
			oneData = PageSportData{
				true,
				u.Teams[0],
				u.Teams[1],
				h2hMap["h2h"].([]interface{})[0].(float64),
				h2hMap["h2h"].([]interface{})[1].(float64),
				h2hMap["h2h"].([]interface{})[2].(float64),
			}
		case 2:
			oneData = PageSportData{
				false,
				u.Teams[0],
				u.Teams[1],
				h2hMap["h2h"].([]interface{})[0].(float64),
				0,
				h2hMap["h2h"].([]interface{})[1].(float64),
			}
		default:
			fmt.Println("We have an alien with no 2 or 3 odds")
		}
		data = append(data, oneData)
	}
	return data
}

func dodajUTiket(w http.ResponseWriter, r *http.Request) {
	reqValue := r.PostFormValue("bet")
	if reqValue == "clear" {
		tiket = nil
		kof = 0
	} else {
		if tiket == nil {
			tiket = make([]string, 0, 1)
		}
		tiket = append(tiket, reqValue)
		betSlice := strings.Split(reqValue, " ")
		kofa, _ := strconv.ParseFloat(betSlice[len(betSlice)-1], 64)
		if kof == 0 {
			kof = kofa
		} else {
			kof *= kofa
		}
	}
	session.Values["tiket"] = tiket
	data := CompletePageSportData{
		Check:     check,
		User:      username,
		SportData: sportsData,
		Sports:    keys,
		Tiket:     tiket,
		Kof:       fmt.Sprintf("%.2f", kof),
	}
	rnd.HTML(w, http.StatusOK, "sport", data)

}
func registracija(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data := struct {
			Sports []string
			Error  string
		}{
			keys,
			"",
		}
		err := rnd.HTML(w, http.StatusOK, "registracija", data)
		if err != nil {
			panic(err)
		}
	} else if r.Method == "POST" {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			fmt.Println("Greska kod pristupa bazi")
			return
		}
		defer db.Close()
		err = db.Ping()
		if err != nil {
			fmt.Println("Greska kod pinga")
		}
		fmt.Println("Spojio se na bazu")
		_, err = db.Exec(`INSERT INTO korisnik (username, sifra, firstname, lastname) VALUES($1, $2, $3, $4);`, r.PostFormValue("username"), r.PostFormValue("password"), r.PostFormValue("firstName"), r.PostFormValue("lastName"))
		//_, err = db.Exec("select from test where valuen = ?", 3)
		if err != nil {
			fmt.Println(err)
			data := struct {
				Sports []string
				Error  string
			}{
				keys,
				"true",
			}
			err = rnd.HTML(w, http.StatusOK, "registracija", data)
			return
		}
		username = r.PostFormValue("username")
		check = true
		data := PageData{
			Check:  check,
			Sports: keys,
			User:   username,
			Msg:    "",
		}
		err = rnd.HTML(w, http.StatusOK, "home", data)
		if err != nil {
			panic(err)
		}
	}
}

func prijava(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data := struct {
			Sports []string
			Error  string
		}{
			keys,
			"",
		}
		rnd.HTML(w, http.StatusOK, "prijava", data)
	} else if r.Method == "POST" {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		err = db.Ping()
		if err != nil {
			panic(err)
		}
		var pwd string
		err = db.QueryRow(`SELECT sifra from korisnik where username = $1`, r.PostFormValue("username")).Scan(&pwd)
		if err != nil {
			data := struct {
				Sports []string
				Error  string
			}{
				keys,
				"Ne postoji username, molim registrirajte se",
			}
			rnd.HTML(w, http.StatusOK, "prijava", data)
			return
		}
		if pwd != r.PostFormValue("password") {
			data := struct {
				Sports []string
				Error  string
			}{
				keys,
				"Lozinka je pogrešna",
			}
			rnd.HTML(w, http.StatusOK, "prijava", data)
			return
		}
		username = r.PostFormValue("username")
		check = true
		data := PageData{
			Sports: keys,
			Check:  check,
			User:   username,
			Msg:    "",
		}
		rnd.HTML(w, http.StatusOK, "home", data)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	//radi sa sessionima
	username = ""
	check = false
	data := PageData{
		Sports: keys,
		Check:  check,
	}
	rnd.HTML(w, http.StatusOK, "home", data)
}
func uplati(w http.ResponseWriter, r *http.Request) {
	if check == false {
		data := struct {
			Sports []string
			Error  string
		}{
			keys,
			"Molimo registrirajte se ili prijavite",
		}
		rnd.HTML(w, http.StatusOK, "prijava", data)
		return
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	tiketInput := ""
	for _, u := range tiket {
		tiketInput += u + ","
	}
	_, err = db.Exec("INSERT INTO listic VALUES($1, $2)", username, tiketInput)
	if err != nil {
		fmt.Println("Error kod uplacivanja listica")
		panic(err)
	}
	kof = 0
	tiket = nil
	data := PageData{
		Sports: keys,
		Check:  check,
		User:   username,
		Msg:    "Vaš listić je uplaćen!",
	}
	rnd.HTML(w, http.StatusOK, "home", data)

}

func mojiListici(w http.ResponseWriter, r *http.Request) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	type List struct {
		Parovi []string
	}
	listici := make([]List, 0, 2)
	rows, err := db.Query("SELECT sadrzajlistica FROM listic WHERE username = $1", username)
	if err != nil {
		fmt.Println("Greska kod dohvacanja listica")
		panic(err)
	}
	defer rows.Close()
	var redak string
	for rows.Next() {
		err := rows.Scan(&redak)
		if err != nil {
			fmt.Println("Greska kod citanja iz baze")
			panic(err)
		}
		jedanListic := List{strings.Split(redak, ",")}
		listici = append(listici, jedanListic)
	}
	data := struct {
		Sports  []string
		Listici []List
	}{
		keys,
		listici,
	}
	rnd.HTML(w, http.StatusOK, "mojiListici", data)
}
func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	//	sportsMap := readSports(jsonContent)//read hardcoded string
	byteFromFile, _ := ioutil.ReadFile("tests/sports.txt")
	sportsMap := readSports(byteFromFile) //Read from file
	keys = make([]string, 0, len(sportsMap))
	for k, _ := range sportsMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, u := range keys {
		http.Handle("/"+u, sessionMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sport, _ := url.QueryUnescape(r.URL.String())
			fmt.Println(sportsMap[sport[1:]][0].Key)
			key := sportsMap[sport[1:]][0].Key
			sportsData = paroviUPageSports(citajParove(key))
			if session.Values["tiket"] != nil {
				tiket = session.Values["tiket"].([]string)
			}
			data := CompletePageSportData{
				User:      username,
				Check:     check,
				SportData: sportsData,
				Sports:    keys,
				Tiket:     tiket,
				Kof:       fmt.Sprintf("%.2f", kof),
			}
			rnd.HTML(w, http.StatusOK, "sport", data)
		})))
	}
	http.Handle("/dodajUTiket", sessionMiddleware(http.HandlerFunc(dodajUTiket)))
	http.HandleFunc("/registracija", registracija)
	http.HandleFunc("/prijava", prijava)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/uplati", uplati)
	http.HandleFunc("/mojiListici", mojiListici)
	http.HandleFunc("/", home)
	fmt.Println(keys)

	http.ListenAndServe(":8080", nil)
}
