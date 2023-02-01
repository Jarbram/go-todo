package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thedevsaddam/renderer"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var rnd *renderer.Render
var db *mgo.Database

const (
	hostName       string = "localhost:27017"
	dbName         string = "demo_todo"
	collectionName string = "todo"
	port           string = ":9000"
)

type (
	todoModel struct {
		ID        bson.ObjectId `bson:"_id,omitempty"`
		Title     string        `bson:"title"`
		Completed bool          `bson:"completed"`
		CreateAt  time.Time     `bson:"createAt"`
	}
	todo struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Completed string    `json:"completed"`
		CreateAt  time.Time `json:"createAt"`
	}
)

func init() {
	rnd = renderer.New()
	session, err := mgo.Dial(hostName)
	CheckErr(err)
	db = session.DB(dbName)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Mount("/todo", todoHandler())

	svr := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("listening on port", port)
		if err := svr.ListenAndServe(); err != nil {
			log.Printf("listen:%s\n", err)
		}
	}()
}

func todoHandler() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", getTodos)
		r.Post("/", createTodoHandler)
		r.Put("/{id}", updateTodoHandler)
		r.Delete("/{id}", deleteTodoHandler)
	})
	return rg
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
