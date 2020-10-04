package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"ese_spiel_scoreboard/internal/pkg/database"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type config struct {
	ListenAddress string
	AdminKey      string
	DBConn        string
	DBType        string
}

var conf config

//GET /api/scoreboard
func apiScoreboard(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "bla")
}

//POST /api/task/create
func apiTaskCreate(w http.ResponseWriter, r *http.Request) {
	points, err := strconv.Atoi(r.FormValue("points"))
	if err != nil {
		log.Fatalf("invalid points argument in create task call: %s", err)
	}
	t := database.Task{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Key:         r.FormValue("key"),
		Storyline:   r.FormValue("storyline"),
		Points:      points,
	}
	database.Insert_task(&t)
}

//GET /
func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/template/login.html")
	if err != nil {
		log.Fatalf("could not parse template: %s", err)
	}
	t.Execute(w, database.Scoreboard())
}

//POST /login
func login(w http.ResponseWriter, r *http.Request) {
	if database.VerifyPassword(r.FormValue("username"), r.FormValue("password")) {
		h := sha256.New()
		h.Write([]byte(strconv.Itoa(rand.Intn(100000000))))
		code := h.Sum(nil)
		codestr := hex.EncodeToString(code)
		cookieValue := r.FormValue("username") + ":" + codestr
		database.InsertSession(database.Session{User: r.FormValue("username"), Key: codestr})
		expire := time.Now().AddDate(0, 0, 5)
		http.SetCookie(w, &http.Cookie{Name: "SessionID", Value: cookieValue, Expires: expire, HttpOnly: true})
		http.Redirect(w, r, "/board", http.StatusFound)
	}
}

//GET /board
func userBoard(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("SessionID")
	if err != nil {
		http.Error(w, "session cookie not set", http.StatusUnauthorized)
		return
	}
	if database.Authenticate(c.Value) {
		t, err := template.ParseFiles("web/template/board.html")
		if err != nil {
			log.Fatalf("failed to open board template: %s", err)
			http.Error(w, "templating error", http.StatusInternalServerError)
			return
		}
		s := strings.Split(c.Value, ":")
		user := s[0]
		t.Execute(w, database.BuildBoard(user))

	} else {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}
}

//POST /api/user/create
func userCreate(w http.ResponseWriter, r *http.Request) {
	u := database.User{
		Name: r.FormValue("name"),
		//TODO hash password
		Password:    r.FormValue("password"),
		Description: r.FormValue("description"),
	}
	database.Insert_user(&u)
}

//GET /user/modify

func parseConfig(path string) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("could not read config file: %s", err)
		os.Exit(1)
	}

	err = json.Unmarshal(file, &conf)
}

//GET /verify
func verifyTask(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("SessionID")
	if err != nil {
		http.Error(w, "session cookie not set", http.StatusUnauthorized)
		return
	}
	if database.Authenticate(c.Value) {
		t, err := template.ParseFiles("web/template/verify.html")
		if err != nil {
			log.Fatalf("failed to open board template: %s", err)
			http.Error(w, "templating error", http.StatusInternalServerError)
			return
		}
		task := r.URL.Query().Get("task")

		t.Execute(w, task)

	} else {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}
}

func apiVerifyTask(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("SessionID")
	if err != nil {
		http.Error(w, "session cookie not set", http.StatusUnauthorized)
		return
	}
	if database.Authenticate(c.Value) {
		if err != nil {
			log.Fatalf("failed to open board template: %s", err)
			http.Error(w, "templating error", http.StatusInternalServerError)
			return
		}
		task := r.FormValue("task")
		key := r.FormValue("key")
		user := strings.Split(c.Value, ":")[0]

		database.VerifyTask(user, task, key)

		http.Redirect(w, r, "/board", http.StatusFound)

	} else {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}
}

func main() {
	parseConfig(os.Args[1])

	database.Initialize(conf.DBType, conf.DBConn)
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/board", userBoard)
	http.HandleFunc("/verify", verifyTask)
	http.HandleFunc("/api/verify", apiVerifyTask)
	http.HandleFunc("/api/task/create", apiTaskCreate)
	http.HandleFunc("/api/user/create", userCreate)
	http.HandleFunc("/api/scoreboard", apiScoreboard)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	go log.Printf("listening on %s", conf.ListenAddress)
	http.ListenAndServe(conf.ListenAddress, nil)
}
