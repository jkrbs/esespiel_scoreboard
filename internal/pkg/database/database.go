package database

import (
	"log"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title       string
	Description string
	Key         string
	Storyline   string
	Points      int
	Finished    bool
}

type User struct {
	gorm.Model
	Name        string
	Password    string
	Description string
	Points      int
	Vorliebe    string
	Eigenschaft string
	St          uint
	Ge          uint
	Ko          uint
	In          uint
	We          uint
	Ch          uint
}

type Session struct {
	Key  string
	User string
}

var db *gorm.DB

func Initialize(mode string, dsn string) {
	var err error
	if mode == "sqlite" {
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
	}
	if mode == "postgres" {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Task{})
	db.AutoMigrate(&Session{})
	db.AutoMigrate(&Storyline{})
	db.AutoMigrate(&finished{})

}

func Insert_task(t *Task) {
	var storylines []Storyline
	db.Find(&storylines)

	f := false
	for _, s := range storylines {
		if s.Name == t.Storyline {
			f = true
		}
	}

	if !f {
		db.Create(&Storyline{Name: t.Storyline})
	}

	db.Create(t)
}

func Insert_user(u *User) {
	db.Create(u)
}

//VerifySession verifies key of user against session database table and returns true, if the session is valid
func VerifySession(key string, user int) bool {
	var sessions []Session
	db.Where("user is ?", user).Find(&sessions)
	for i := range sessions {
		if sessions[i].Key == key {
			return true
		}
	}
	return false
}

func VerifyPassword(user string, password string) bool {
	return true
}

func InsertSession(s Session) {
	db.Create(s)
}

func Scoreboard() []User {
	var users []User
	db.Order("points desc").Find(&users)
	return users
}

func Authenticate(cookie string) bool {
	s := strings.Split(cookie, ":")
	name := s[0]
	session := s[1]

	var sessions []Session
	db.Where("user is ?", name).Find(&sessions)
	for _, sess := range sessions {
		if sess.Key == session && sess.User == name {
			return true
		}
	}
	return false
}

type Board struct {
	User        string
	Vorliebe    string
	Eigenschaft string
	St          uint
	Ge          uint
	Ko          uint
	In          uint
	We          uint
	Ch          uint
	Storylines  []StorylineTask
}

type Storyline struct {
	Name string
}

type StorylineTask struct {
	Name  string
	Tasks []Task
}

func BuildBoard(user string) Board {
	var tasks []Task
	db.Find(&tasks)

	var stroylineNames []Storyline
	db.Find(&stroylineNames)

	var fin []finished
	db.Where("user is ?", user).Find(&fin)
	var user_struct User
	db.Where("name is ?", user).Find(&user_struct)

	var storylines []StorylineTask

	for _, s := range stroylineNames {
		storylines = append(storylines, StorylineTask{
			Name: s.Name,
		})
	}

	for _, t := range tasks {
		t.Finished = false
		for i := range storylines {
			if storylines[i].Name == t.Storyline {
				for _, f := range fin {
					if f.Task == t.Title {
						t.Finished = true
					}
				}

				storylines[i].Tasks = append(storylines[i].Tasks, t)
			}
		}
	}

	board := Board{
		User:        user,
		Storylines:  storylines,
		Vorliebe:    user_struct.Vorliebe,
		Eigenschaft: user_struct.Eigenschaft,
		St:          user_struct.St,
		Ge:          user_struct.Ge,
		Ko:          user_struct.Ko,
		In:          user_struct.In,
		We:          user_struct.We,
		Ch:          user_struct.Ch,
	}
	return board
}

type finished struct {
	User        string
	Vorliebe    string
	Eigenschaft string
	St          uint
	Ge          uint
	Ko          uint
	In          uint
	We          uint
	Ch          uint

	Task string
}

func VerifyTask(user string, task string, key string) {
	var tasks []Task
	db.Where("title is ?", task).Find(&tasks)
	if len(tasks) == 0 {
		return
	}
	for i, _ := range tasks {
		if tasks[i].Key == key {
			var u User
			db.Where("name is ?", user).Find(&u)
			db.Create(&finished{
				User:        user,
				Vorliebe:    u.Vorliebe,
				Eigenschaft: u.Eigenschaft,
				St:          u.St,
				Ge:          u.Ge,
				Ko:          u.Ko,
				In:          u.In,
				We:          u.We,
				Ch:          u.Ch,

				Task: task,
			})
			return
		}

		log.Printf("wrong verification key")
		return
	}
}
