package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	// _ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

// TMovie - movie struct
type TMovie struct {
	ID     string
	Name   string
	URL    string
	Poster string
	User   string
	Rents  int
}

// TMovieList - movie struct
type TMovieList struct {
	Title string
	User  string
	W     string
	List  []TMovie
}

// TUser - user struct
type TUser struct {
	ID    string
	Name  string
	Email string
}

var myList = TMovieList{
	Title: "Список фильмов",
	User:  "",
	W:     "0",
	List: []TMovie{
		{ID: "1", Name: "Pulp fiction", URL: "https://www.youtube.com/embed/qMvb5jUMR80", Poster: "https://m.media-amazon.com/images/M/MV5BNGNhMDIzZTUtNTBlZi00MTRlLWFjM2ItYzViMjE3YzI5MjljXkEyXkFqcGdeQXVyNzkwMjQ5NzM@._V1_UY268_CR1,0,182,268_AL_.jpg"},
		{ID: "2", Name: "Snatch", URL: "https://www.youtube.com/embed/_c1kBwsv8PM", Poster: "https://m.media-amazon.com/images/M/MV5BMTA2NDYxOGYtYjU1Mi00Y2QzLTgxMTQtMWI1MGI0ZGQ5MmU4XkEyXkFqcGdeQXVyNDk3NzU2MTQ@._V1_UY268_CR0,0,182,268_AL_.jpg"},
		{ID: "3", Name: "Generation P", URL: "https://www.youtube.com/embed/Rts2oc1oilI", Poster: "https://m.media-amazon.com/images/M/MV5BMTAxNzM1MDI1OTheQTJeQWpwZ15BbWU4MDk1NjQ1ODAx._V1_UY268_CR4,0,182,268_AL_.jpg"},
	},
}

// Server - server struct
type Server struct {
	// db       *sql.DB
	db       *pgx.Conn
	currUser string
	currRole string
	email    string
	movieID  string
}

func main() {
	// dbHost := "movie_db:3306"
	// dbName := "movie_base"
	// dbUser := "movie_user"
	// dbPwd := "movie_user_pwd"

	// DSN := "u2:qw12345@tcp(localhost:3306)/movie_base?charset=utf8"
	DSN := "postgres://movie_user:movie_user_pwd@postgres:5432/movie_base?sslmode=disable"
	// DSN := "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable"

	// log.Println("DSN ", DSN)
	// dbs, err := sql.Open("mysql", DSN)
	ctx := context.Background()
	log.Println("Before pgx.Connect")
	conn, err := pgx.Connect(ctx, DSN)
	log.Println("Before err checking")

	if err != nil {
		for attempt := 0; attempt < 10; attempt++ {
			if conn, err := pgx.Connect(context.Background(), DSN); err == nil {
				log.Println("db connected!", conn)

				break
			}
			log.Println("Cant connect db!", attempt)
			time.Sleep(time.Duration(attempt+1) * time.Second)
		}
	}

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// db.SetConnMaxLifetime(5 * time.Minute)
	// db.SetMaxOpenConns(25)
	// db.SetMaxIdleConns(25)*

	defer func() {
		if err = conn.Close(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// db := mysql.New("tcp", "", "192.168.99.100:3306", "movie_user", "movie_user_pwd", "movie_base")

	// err := db.Connect()
	// if err != nil {
	// 	panic(err)
	// }

	// if err := dbs.Ping(); err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println("db pinged!")

	// serv := Server{db: dbs, currUser: "", currRole: "0", email: "", movieID: "0"}
	serv := Server{db: conn, currUser: "", currRole: "0", email: "", movieID: "0"}

	//router := http.NewServeMux()
	router := mux.NewRouter()

	fs := http.FileServer(http.Dir("assets"))
	// log.Println("http.StripPrefix(/assets/, fs) ", http.StripPrefix("/assets/", fs), fs)
	router.Handle("/assets/", http.StripPrefix("/assets/", fs))

	router.HandleFunc("/", serv.viewList)
	router.HandleFunc("/movie/{[0-9]+}", serv.viewMovie)
	router.HandleFunc("/login", serv.handleLogin)
	router.HandleFunc("/profile/", serv.handleProfile)
	router.HandleFunc("/payment/", serv.handlePayment)
	router.HandleFunc("/savelogin/", serv.handleSaveLogin)
	router.HandleFunc("/savepayment/", serv.handleSavePayment)
	router.HandleFunc("/saveorder/", serv.handleSaveOrder)
	router.HandleFunc("/admin/", serv.handleAdmin)
	router.HandleFunc("/users/", serv.handleUsers)
	router.HandleFunc("/user/{[0-9]+}", serv.handleUser)
	router.HandleFunc("/newuser/", serv.handleNewUser)

	log.Println("Listen at 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}

func (server *Server) viewList(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "O: %s\n", r.URL.Path)

	var tmpl = template.Must(template.New("MyTemplate").ParseFiles("./tmpl.html"))

	/*MyBlog, err := GetBlog(server.database, server.currBlog)
	if err != nil {
		log.Println(err)
		return
	}*/

	if (len(server.currUser) > 0) && len(myList.User) == 0 {
		myList.User = server.currUser
	}
	myList.W = server.currRole
	if err := tmpl.ExecuteTemplate(w, "MyTemplate", myList); err != nil {
		log.Println(err)
		return
	}
}

func (server *Server) viewMovie(w http.ResponseWriter, r *http.Request) {
	var movie = template.Must(template.New("MyMovie").ParseFiles("./movie.html"))

	mvurl := strings.Split(r.URL.Path, "/")

	myMovie := TMovie{}
	for i := range myList.List {
		if myList.List[i].ID == mvurl[len(mvurl)-1] {
			myMovie = myList.List[i]
		}
	}

	if (len(server.currUser) > 0) && len(myMovie.User) == 0 {
		myMovie.User = server.currUser
	}
	//func GetMovie(id string) (TMovie, error)

	server.movieID = myMovie.ID

	rents := 0
	row := server.db.QueryRow(context.Background(), "select count(*) from movie_base.orders where user_id = ? and movie_id = ?", server.currUser, myMovie.ID)
	log.Println("server.currUser, myMovie.ID ", server.currUser, myMovie.ID)
	err := row.Scan(&rents)

	if err == sql.ErrNoRows {
		rents = 0
	}

	myMovie.Rents = rents

	log.Println("Rents ", myMovie.Rents)

	if err := movie.ExecuteTemplate(w, "MyMovie", myMovie); err != nil {
		log.Println(err)
		return
	}
}

func (server *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var login = template.Must(template.New("MyLogin").ParseFiles("./login.html"))

	if err := login.ExecuteTemplate(w, "MyLogin", myList); err != nil {
		log.Println(err)
		return
	}
}

func (server *Server) handleProfile(w http.ResponseWriter, r *http.Request) {
	var profile = template.Must(template.New("MyProfile").ParseFiles("./profile.html"))

	type TProfileList struct {
		Name      string
		Email     string
		Phone     string
		BirthDate string
		Balance   float32
	}

	myProfileList := TProfileList{}

	tUnix := 0
	row := server.db.QueryRow(context.Background(), "select display_name, email, phone_number, birth_date from movie_base.users where id = ?", server.currUser)
	err := row.Scan(&myProfileList.Name, &myProfileList.Email, &myProfileList.Phone, &tUnix)
	if err != nil {
		return
	}

	// myProfileList.BirthDate = "01.01.2000"
	myProfileList.BirthDate = (time.Unix(int64(tUnix), 0)).Format("2 January 2006")

	pay := 0
	row = server.db.QueryRow(context.Background(), "select sum(amount) from movie_base.payments where user_id = ?", server.currUser)
	err = row.Scan(&pay)

	if err == sql.ErrNoRows {
		pay = 0
	}

	ord := 0
	row = server.db.QueryRow(context.Background(), "select sum(amount) from movie_base.orders where user_id = ?", server.currUser)
	err = row.Scan(&ord)

	if err == sql.ErrNoRows {
		ord = 0
	}

	myProfileList.Balance = float32(pay - ord)
	log.Println("myProfileList ", myProfileList)

	// var myProfileList = TProfileList{
	// 	Name:    "Петр Иванов",
	// 	Balance: 100.00,
	// }

	if err := profile.ExecuteTemplate(w, "MyProfile", myProfileList); err != nil {
		log.Println(err)
		return
	}
}

func (server *Server) handlePayment(w http.ResponseWriter, r *http.Request) {
	var payment = template.Must(template.New("MyPayment").ParseFiles("./payment.html"))
	// var payment = template.Must(template.ParseFiles("./payment.html"))

	type TPaymentList struct {
		Balance float32
	}
	var myPaymentList = TPaymentList{0}

	pay := 0
	row := server.db.QueryRow(context.Background(), "select sum(amount) from movie_base.payments where user_id = ?", server.currUser)
	err := row.Scan(&pay)

	if err == sql.ErrNoRows {
		pay = 0
	}

	ord := 0
	row = server.db.QueryRow(context.Background(), "select sum(amount) from movie_base.orders where user_id = ?", server.currUser)
	err = row.Scan(&ord)

	if err == sql.ErrNoRows {
		ord = 0
	}

	myPaymentList.Balance = float32(pay - ord)

	if err := payment.ExecuteTemplate(w, "MyPayment", myPaymentList); err != nil {
		// if err := payment.Execute(w, nil); err != nil {
		log.Println(err)
		return
	}
}

func (server *Server) handleSaveLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("err ", err)
	}

	server.email = r.FormValue("user")
	log.Println("server.email %w", server.email)

	// row := server.db.QueryRow(context.Background(), "select id, role from movie_base.users where email = ?", server.email)
	row := server.db.QueryRow(context.Background(), "select id, role from users where email = $1;", server.email)
	usr := 0
	// if err := row.Scan(&server.currUser, &server.currRole); err != nil {
	if err := row.Scan(&usr, &server.currRole); err != nil {
		log.Println("err ", err)
	}
	server.currUser = strconv.Itoa(usr)

	// log.Println("server.currRole ", server.currRole)
	// server.currUser = "1"

	// log.Println(r.Form.Get("password"))

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (server *Server) handleSavePayment(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("err ", err)
	}

	addition := r.FormValue("addition")

	// icu, _ := strconv.Atoi(server.currUser)
	res, err := server.db.Exec(context.Background(), "insert into movie_base.payments (amount, user_id, transaction_id, status, created_at, updated_at) VALUES (?,?,0,0,NULL,NULL)", addition, server.currUser)
	// res, err := server.db.Exec("insert into movie_base.payments (amount) VALUES (?)", addition)
	if err != nil {
		log.Printf("err %v, res %v", err, res)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (server *Server) handleSaveOrder(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("err ", err)
	}
	// order := r.FormValue("order")
	movieID := r.FormValue("id")
	log.Println("movieID ", movieID)

	order := 100

	res, err := server.db.Exec(context.Background(), "insert into movie_base.orders (amount, user_id, movie_id, created_at) VALUES (?,?,?,0)", order, server.currUser, server.movieID)
	if err != nil {
		log.Printf("err %v, res %v", err, res)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (server *Server) handleAdmin(w http.ResponseWriter, r *http.Request) {
	var admin = template.Must(template.New("MyAdmin").ParseFiles("./admin.html"))
	// if err := admin.ExecuteTemplate(w, "MyAdmin", myAdmin); err != nil {
	if err := admin.ExecuteTemplate(w, "MyAdmin", nil); err != nil {
		log.Println(err)
		return
	}
}

func (server *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	var users = template.Must(template.New("MyUserList").ParseFiles("./users.html"))

	type TUserList struct {
		UserList []TUser
	}

	var MyUserList TUserList

	rows, err := server.db.Query(context.Background(), "select id, display_name, email from movie_base.users")

	defer func() {
		// if err = rows.Close(); err != nil {
		rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	for rows.Next() {
		user := TUser{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			log.Println(err)
			continue
		}
		MyUserList.UserList = append(MyUserList.UserList, user)
	}

	if err := users.ExecuteTemplate(w, "MyUserList", MyUserList); err != nil {
		log.Println(err)
		return
	}
}

func (server *Server) handleUser(w http.ResponseWriter, r *http.Request) {
	var user = template.Must(template.New("MyUser").ParseFiles("./user.html"))
	url := strings.Split(r.URL.Path, "/")
	MyUser := TUser{}
	row := server.db.QueryRow(context.Background(), "select id, display_name, email from movie_base.users where id = ?", url[len(url)-1])
	err := row.Scan(&MyUser.ID, &MyUser.Name, &MyUser.Email)

	// log.Println(MyUser.Name)
	out := strings.Replace(MyUser.Name, " ", "&nbsp;", -1)
	MyUser.Name = out

	if err != nil {
		log.Fatal(err)
	}

	if err := user.ExecuteTemplate(w, "MyUser", MyUser); err != nil {
		log.Println(err)
		return
	}
}

func (server *Server) handleNewUser(w http.ResponseWriter, r *http.Request) {
	var user = template.Must(template.New("MyUser").ParseFiles("./user.html"))

	NewUser := TUser{"", "", ""}

	if err := user.ExecuteTemplate(w, "MyUser", NewUser); err != nil {
		log.Println(err)
		return
	}
}
