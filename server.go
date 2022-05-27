package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"tawesoft.co.uk/go/dialog"
)

type Login struct {
	Name  string
	Email string
	Image string
	Post  int
	Id    int
	Sub   int
}

type Post struct {
	Title          string
	Content        string
	Date           string
	Id             int
	Like           int
	Dislike        int
	Image          string
	Author         string
	Filter         int
	AuthorComment  string
	ContentComment string
	DateComment    string
	// CountCom       int
}

type PostData struct {
	Title   string
	Content string
	Date    string
	Id      int
	Like    int
	Dislike int
	Image   string
	Author  string
	Filter  int
}

var user Login
var allUser []Login
var allResult []Post
var allData []PostData

type Register struct {
	Id       int
	Pseudo   string
	Email    string
	Password string
	Log      int
}

//Define new CookiesSessions
var store = sessions.NewCookieStore([]byte("mysession"))

// Initialise DataBase, and create it with his tables
func initDatabase(database string) *sql.DB {
	db, err := sql.Open("sqlite3", "assets/Database/db.db")
	if err != nil {
		log.Fatal(err)
	}
	sqlStmt := `
				CREATE TABLE IF NOT EXISTS register (
					id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					pseudo TEXT NOT NULL, 
					email TEXT NOT NULL, 
					password TEXT NOT NULL,
					image TEXT NOT NULL,
					post INT NOT NULL,
					subscribers INT NOT NULL
				);

				CREATE TABLE IF NOT EXISTS post (
					id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					author TEXT NOT NULL,
					date TEXT NOT NULL,
					title TEXT NOT NULL,
					content TEXT NOT NULL,
					like INT NOT NULL,
					dislike INT NOT NULL,
					filter INT NOT NULL
					
				);

				CREATE TABLE IF NOT EXISTS comment (
					id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					postid INT NOT NULL,
					date TEXT NOT NULL,
					author TEXT NOT NULL,
					content TEXT NOT NULL
				);

				CREATE TABLE IF NOT EXISTS like (
					postid INTEGER NOT NULL,
					author TEXT NOT NULL,
					like INT NOT NULL,
					dislike INT NOT NULL,
					PRIMARY KEY (postid, author)
				);
				`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

//Get the last id from the table post
func getBookLastID() int {
	db := initDatabase("assets/Database/db.db/")
	var id int

	err := db.QueryRow("select ifnull(max(id), 0) as id from post").Scan(&id)
	if err != nil {
		panic(err)
	}
	return id + 1
}

//Function use for Insert elements in the table "register"
func insertIntoRegister(db *sql.DB, pseudo string, email string, password string, image string) (int64, error) {
	result, _ := db.Exec(`INSERT INTO register (pseudo, email, password, image, post, subscribers) values (?, ?, ?, ?, 0, 0)`, pseudo, email, password, image)
	return result.LastInsertId()
}

//Function use for Insert elements in the table "post"
func insertIntoPost(db *sql.DB, title string, content string, author string) (int64, error) {
	result, _ := db.Exec(`INSERT INTO post (author, date, title, content, like, dislike, filter) values (?, ?, ?, ?, 0, 0, 0)`, author, time.Now(), title, content)
	return result.LastInsertId()
}

//Function use for Insert elements in the table "comment"
func insertIntoComment(db *sql.DB, postid int, author string, content string) (int64, error) {
	result, _ := db.Exec(`INSERT INTO comment (postid, date, author, content) values (?, ?, ?, ?)`, postid, "0", author, content)
	return result.LastInsertId()
}

//Function use for Insert elements in the table "like"
func insertIntoLike(db *sql.DB, postid string, author string) (int64, error) {
	result, _ := db.Exec(`INSERT INTO like (postid, author, like, dislike) values (?, ?, 1, 1)`, postid, author)
	return result.LastInsertId()
}

//Get all data from the struct Post (for the post data). Useful for send data in the html page with templates
func getPostData() {
	db := initDatabase("assets/Database/db.db")
	var temp Post

	rows, _ :=
		db.Query(`SELECT * FROM post`)
	allResult = nil
	for rows.Next() {
		rows.Scan(&temp.Id, &temp.Author, &temp.Date, &temp.Title, &temp.Content, &temp.Like, &temp.Dislike, &temp.Filter)
		allResult = append([]Post{temp}, allResult...)
	}
}

//Get all data from the struct Post (for the comment data). Useful for send data in the html page with templates
func getCommentData(idInfo int) {
	db := initDatabase("assets/Database/db.db")
	var temp Post

	rows, _ :=
		db.Query("SELECT author, content, date FROM comment WHERE postid = ?", idInfo)
	allResult = nil
	for rows.Next() {
		rows.Scan(&temp.AuthorComment, &temp.ContentComment, &temp.DateComment)
		allResult = append(allResult, temp)
	}
}

//Get all data from the struct Post (for the post data). Useful for send data in the html page with templates
func getPostDataById(idInfo int) {
	db := initDatabase("assets/Database/db.db")
	var temp PostData

	rows, _ :=
		db.Query("SELECT id, author, date, title, content, like, dislike, filter FROM post WHERE id = ?", idInfo)
	allData = nil
	for rows.Next() {
		rows.Scan(&temp.Id, &temp.Author, &temp.Date, &temp.Title, &temp.Content, &temp.Like, &temp.Dislike, &temp.Filter)
		allData = append(allData, temp)
	}
}

//Get all data from the struct Post and filters the data
func getPostDataByFilter(filter int) {
	db := initDatabase("assets/Database/db.db")
	var temp Post

	rows, _ :=
		db.Query(`SELECT * FROM post WHERE filter = ?`, filter)
	allResult = nil
	for rows.Next() {
		rows.Scan(&temp.Id, &temp.Author, &temp.Date, &temp.Title, &temp.Content, &temp.Like, &temp.Dislike, &temp.Filter)
		allResult = append([]Post{temp}, allResult...)
	}
}

func getUserInfo(userInfo string) {
	db := initDatabase("assets/Database/db.db")
	var temp Login

	rows, _ :=
		db.Query(`SELECT pseudo, email, image, post, subscribers FROM register WHERE pseudo = ?`, userInfo)
	allUser = nil
	for rows.Next() {
		rows.Scan(&temp.Name, &temp.Email, &temp.Image, &temp.Post, &temp.Sub)
		allUser = append(allUser, temp)
	}
}

func getUserInfoByCookie(username string) {
	db := initDatabase("assets/Database/db.db")
	var temp Login

	rows, _ :=
		db.Query(`SELECT pseudo, email, image, post, subscribers FROM register WHERE pseudo = ?`, username)
	allUser = nil
	for rows.Next() {
		rows.Scan(&temp.Name, &temp.Email, &temp.Image, &temp.Post, &temp.Sub)
		allUser = append(allUser, temp)
	}
}

//This function compare the Pseudo and the password send by the user.
func login(LogPseudo string, LogPassword string) bool {
	db := initDatabase("assets/Database/db.db")
	var pseudo string
	var password string
	var result = false
	rows, _ :=
		db.Query("SELECT pseudo, password FROM register")
	for rows.Next() {
		rows.Scan(&pseudo, &password)
		if LogPseudo == pseudo && CheckPasswordHash(LogPassword, password) { //If Pseudo and password match, the function return true
			result = true
		}
	}
	return result
}

//This function check if you have already vote or no (like)
func checkLike(username string, likeId string) {
	db := initDatabase("assets/Database/db.db")
	var author string
	var postid int
	var like int
	var dislike int
	rows, _ := db.Query("SELECT postid, author, like, dislike FROM like WHERE author = ? and postid = ?", username, likeId)
	for rows.Next() {
		rows.Scan(&postid, &author, &like, &dislike)
	}
	if author != "" && postid != 0 && like != 0 {
		db.Query("SELECT postid, author, like, dislike FROM like WHERE author = ? and postid = ?", username, likeId)
		if like == 1 && dislike == 1 {
			db.Exec("UPDATE post SET like = like + 1 WHERE id = ?", likeId)
			db.Exec("UPDATE like SET like = 2 WHERE author = ?", username)
		}
		if like == 2 {
			db.Exec("UPDATE post SET like = like - 1 WHERE id = ?", likeId)
			db.Exec("UPDATE like SET like = 1 WHERE author = ?", username)
		}
	} else {
		insertIntoLike(db, likeId, username)
		db.Exec("UPDATE post SET like = like + 1 WHERE id = ?", likeId)
		db.Exec("UPDATE like SET like = 2 WHERE author = ?", username)
	}
}

//This function check if you have already vote or no (dislike)
func checkDislike(username string, likeId string) {
	db := initDatabase("assets/Database/db.db")
	var author string
	var postid int
	var like int
	var dislike int
	rows, _ := db.Query("SELECT postid, author, like, dislike FROM like WHERE author = ? and postid = ?", username, likeId)
	for rows.Next() {
		rows.Scan(&postid, &author, &like, &dislike)
	}
	if author != "" && postid != 0 && dislike != 0 {
		db.Query("SELECT postid, author, like, dislike FROM like WHERE author = ? and postid = ?", username, likeId)
		if like == 1 && dislike == 1 {
			db.Exec("UPDATE post SET dislike = dislike + 1 WHERE id = ?", likeId)
			db.Exec("UPDATE like SET dislike = 2 WHERE author = ?", username)
		}
		if dislike == 2 {
			db.Exec("UPDATE post SET dislike = dislike - 1 WHERE id = ?", likeId)
			db.Exec("UPDATE like SET dislike = 1 WHERE author = ?", username)
		}
	} else {
		insertIntoLike(db, likeId, username)
		db.Exec("UPDATE post SET dislike = dislike + 1 WHERE id = ?", likeId)
		db.Exec("UPDATE like SET dislike = 2 WHERE author = ?", username)
	}
}

//Data encryption
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//Check if the Data encryption is the same of the password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//The function considers whether the nickname / mail is already in use
func register(RegisterPseudo string, RegisterEmail string) bool {
	db := initDatabase("assets/Database/db.db")

	var pseudo string
	var email string
	var result = true
	rows, _ :=
		db.Query("SELECT  pseudo, email FROM register")
	for rows.Next() {
		rows.Scan(&pseudo, &email)
		if RegisterPseudo == pseudo || RegisterEmail == email { //if there are same Pseudo or Email, the function return False
			result = false
		}
	}
	return result
}

//HandleFunc for index.html (Get and post data)
func indexHandler(w http.ResponseWriter, r *http.Request) {
	badgeInformatique := r.FormValue("badgeInformatique")
	badgeSport := r.FormValue("badgeSport")
	badgeMusique := r.FormValue("badgeMusique")
	badgeGame := r.FormValue("badgeGame")
	badgeFood := r.FormValue("badgeFood")
	if badgeInformatique == "1" {
		getPostDataByFilter(1)
	} else if badgeSport == "2" {
		getPostDataByFilter(2)
	} else if badgeMusique == "3" {
		getPostDataByFilter(3)
	} else if badgeGame == "4" {
		getPostDataByFilter(4)
	} else if badgeFood == "5" {
		getPostDataByFilter(5)
	} else {
		getPostData()
	}
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, allResult) //Execute the value of all post
}

func nologinHandler(w http.ResponseWriter, r *http.Request) {
	badgeInformatique := r.FormValue("badgeInformatique")
	badgeSport := r.FormValue("badgeSport")
	badgeMusique := r.FormValue("badgeMusique")
	badgeGame := r.FormValue("badgeGame")
	badgeFood := r.FormValue("badgeFood")
	if badgeInformatique == "1" {
		getPostDataByFilter(1)
	} else if badgeSport == "2" {
		getPostDataByFilter(2)
	} else if badgeMusique == "3" {
		getPostDataByFilter(3)
	} else if badgeGame == "4" {
		getPostDataByFilter(4)
	} else if badgeFood == "5" {
		getPostDataByFilter(5)
	} else {
		getPostData()
	}
	t, _ := template.ParseFiles("nologin.html")
	t.Execute(w, allResult)
}

//HandleFunc for register.html (Get and post data)
func registerHandler(w http.ResponseWriter, r *http.Request) {
	pseudoForm := r.FormValue("pseudoCreate")
	emailForm := r.FormValue("emailCreate")
	passwordForm := r.FormValue("passwordCreate")
	imageForm := r.FormValue("imageCreate")
	pseudoLog := r.FormValue("pseudoLog")
	passwordLog := r.FormValue("passwordLog")

	user.Image = "http://marclimoservices.com/wp-content/uploads/2017/05/facebook-default.png"
	db := initDatabase("assets/Database/db.db")

	hash, _ := HashPassword(passwordForm)
	//register conditions
	if pseudoForm != "" && emailForm != "" && passwordForm != "" {
		if register(pseudoForm, emailForm) { //If true
			if imageForm != "" {
				insertIntoRegister(db, pseudoForm, emailForm, hash, imageForm) //insert the data send by the user in the database
				dialog.Alert("Your account has been created, Pseudo: %v | Email: %v \nPlease LogIn.", pseudoForm, emailForm)
			} else {
				insertIntoRegister(db, pseudoForm, emailForm, hash, "http://marclimoservices.com/wp-content/uploads/2017/05/facebook-default.png") //insert the data send by the user in the database
				dialog.Alert("Your account has been created, Pseudo: %v | Email: %v \nPlease LogIn.", pseudoForm, emailForm)
			}
		} else {
			dialog.Alert("Email or Pseudo already used !")
		}
	}
	//login conditions
	if login(pseudoLog, passwordLog) { //if true
		//Create and save cookie sessions
		user.Name = pseudoLog
		session, _ := store.Get(r, "mysession")
		session.Values["username"] = pseudoLog
		session.Save(r, w)
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}
	t, _ := template.ParseFiles("register.html")
	t.Execute(w, nil)

}

//HandleFunc for profile.html (Get and post data)
func profileHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mysession")
	username := fmt.Sprintf("%v", session.Values["username"]) //Decrypts data of the session cookies
	getUserInfoByCookie(username)
	t, _ := template.ParseFiles("profile.html")
	t.Execute(w, allUser)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	userInfo := r.URL.Path[6:]
	getUserInfo(userInfo)
	t, _ := template.ParseFiles("user.html")
	t.Execute(w, allUser)
}

//HandleFunc for register.html (Get and post data)
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mysession")
	//Delete and save coockie sessions
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/register", http.StatusSeeOther)
}

func likeHandler(w http.ResponseWriter, r *http.Request) {
	likeId := r.URL.Path[6:]
	redirect := "/info/" + likeId
	session, _ := store.Get(r, "mysession")
	username := fmt.Sprintf("%v", session.Values["username"]) //Decrypts data of the session cookies
	checkLike(username, likeId)
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

func dislikeHandler(w http.ResponseWriter, r *http.Request) {
	likeId := r.URL.Path[9:]
	redirect := "/info/" + likeId
	session, _ := store.Get(r, "mysession")
	username := fmt.Sprintf("%v", session.Values["username"]) //Decrypts data of the session cookies
	checkDislike(username, likeId)
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

//HandleFunc for post.html (Get and post data)
func postHandler(w http.ResponseWriter, r *http.Request) {
	db := initDatabase("assets/Database/db.db/")
	titleForm := r.FormValue("inputEmail")
	contentForm := r.FormValue("inputPassword")
	session, _ := store.Get(r, "mysession")
	username := fmt.Sprintf("%v", session.Values["username"]) //Decrypts data of the session cookies
	user.Post = 0

	informatique := r.FormValue("badgeInformatique")
	sport := r.FormValue("badgeSport")
	musique := r.FormValue("badgeMusique")
	jeux := r.FormValue("badgeGame")
	food := r.FormValue("badgeFood")
	lastid := getBookLastID()

	if titleForm != "" && contentForm != "" {
		insertIntoPost(db, titleForm, contentForm, username) //insert the value send by the user in the database
		db.Exec(`INSERT INTO post (date) values (?)`, time.Now())
		db.Exec(`UPDATE register SET post = post + 1 WHERE pseudo = ?`, username) //Update the value in the database
		//Get the value of checkbox, and update the database
		if informatique == "1" {
			db.Exec(`UPDATE post SET filter = ? WHERE id = ?`, 1, lastid)
		}
		if sport == "2" {
			db.Exec(`UPDATE post SET filter = ? WHERE id = ?`, 2, lastid)
		}
		if musique == "3" {
			db.Exec(`UPDATE post SET filter = ? WHERE id = ?`, 3, lastid)
		}
		if jeux == "4" {
			db.Exec(`UPDATE post SET filter = ? WHERE id = ?`, 4, lastid)
		}
		if food == "5" {
			db.Exec(`UPDATE post SET filter = ? WHERE id = ?`, 5, lastid)
		}

		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

	t, _ := template.ParseFiles("post.html")
	t.Execute(w, nil)

}

//HandleFunc for info.html (Get and post data)
func infoHandler(w http.ResponseWriter, r *http.Request) {
	db := initDatabase("assets/Database/db.db/")
	idInfo, _ := strconv.Atoi(r.URL.Path[6:]) //Get the 6th element of the url and convert it into a int
	contentComment := r.FormValue("commentArea")
	session, _ := store.Get(r, "mysession")
	username := fmt.Sprintf("%v", session.Values["username"])
	redirect := "/info/" + strconv.Itoa(idInfo)

	//Select all the lines in the database, and scan them to get the value
	getPostDataById(idInfo)

	if len(contentComment) > 0 {

		insertIntoComment(db, idInfo, username, contentComment)
		db.Exec(`UPDATE comment SET date = ? WHERE postid = ?`, time.Now(), idInfo)
		http.Redirect(w, r, redirect, http.StatusSeeOther)
	}

	getCommentData(idInfo)
	//use map for Exec 2 value
	m := map[string]interface{}{
		"Results": allResult,
		"Post":    allData,
	}
	t := template.Must(template.ParseFiles("info.html"))
	t.Execute(w, m)

}

//Initializes the server on port 8080, and manages requests
func main() {
	fs := http.FileServer(http.Dir(""))
	http.Handle("/", fs)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/user/", userHandler)
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/info/", infoHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/like/", likeHandler)
	http.HandleFunc("/dislike/", dislikeHandler)
	http.HandleFunc("/nologin", nologinHandler)
	http.ListenAndServe(":8080", nil)
}
