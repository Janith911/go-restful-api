package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID         uint   `gorm:"unique;primaryKey;autoIncrement:true"`
	Firstname  string `json:"Firstname"`
	Lastname   string `json:"Lastname"`
	Occupation string `json:"Occupation"`
}

type Users struct {
	Users []User `json:"Users"`
}

func initDB(dsn string) *gorm.DB {
	log.Println("Attempting to connect to MySQL Database")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DisableAutomaticPing: false})

	if err != nil {
		log.Fatal("An ERROR Occurred : ", err)
	}

	log.Println("Connected to MySQL Database")

	// Migrate the Schema
	db.AutoMigrate(&User{})

	return db
}

// Handler Functions
func getUsers(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)

		// List all the users and store the result in 'users'
		users := []User{}
		db.Find(&users)

		bs, err := json.MarshalIndent(Users{Users: users}, "", " ")
		if err != nil {
			log.Println("An ERROR Occurred : ", err)
		}
		w.Write(bs)
	}
}

func getUser(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		user := User{}

		//Convert the ID to intiger and perform a validation
		user_id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Println("An ERROR Occurred : ", err)
			http.Error(w, "Invalid User ID", http.StatusNotFound)
		} else {

			// Find the user record with given User ID
			db.Find(&user, "id = ?", user_id)

			bs, err := json.MarshalIndent(user, "", " ")

			if err != nil {
				log.Println("An ERROR Occurred : ", err)
				http.Error(w, "An ERROR Occerued", http.StatusNotFound)
			} else if user.ID == 0 {
				w.Write([]byte("{}"))
			} else {
				w.Write(bs)
			}
		}
	}
}

func updateUser(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		user := User{}

		// Parse the incomming request body and store the result in 'user'
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			log.Println("An ERROR Occurred : ", err)
			http.Error(w, "Unable to parse data", http.StatusBadRequest)
		} else {
			user_id, err := strconv.Atoi(r.PathValue("id"))
			if err != nil {
				log.Println("An ERROR Occurred : ", err)
				http.Error(w, "Invalid User ID", http.StatusNotFound)
			} else {
				user_db := User{}
				db.Find(&user_db, "id = ?", user_id)
				if user_db.ID == 0 {
					log.Println("Record ID : ", user_id, "Record does not exist")
					http.Error(w, "Record does not exist", http.StatusBadRequest)
				} else {
					// Update the record with the given User ID
					db.Model(&User{}).Where("id = ?", user_id).Updates(user)
				}
				w.Write([]byte("Record updated successfuly"))
			}
		}
	}
}

func createUser(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		user := User{}
		// Parse the request body into JSON object
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			log.Println("An ERROR Occurred : ", err)
			http.Error(w, "Unable to parse data", http.StatusBadRequest)
		} else {
			// Create the record
			db.Create(&user)
			w.Write([]byte("{}"))
		}
	}
}

func deleteUser(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		user_id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Println("An ERROR Occurred : ", err)
			http.Error(w, "Invalid User ID", http.StatusNotFound)
		} else {
			// Delete the record
			db.Delete(&User{}, user_id)
			w.Write([]byte("Record Deleted"))
		}

	}
}

func main() {
	mysql_dsn := "goapiuser:goapipasswd@tcp(127.0.0.1:3306)/goapi"

	// Connect to MySQL database instance
	db := initDB(mysql_dsn)

	// Create new Router
	r := http.NewServeMux()

	// Bind Handler functions
	r.HandleFunc("GET /users", getUsers(db))
	r.HandleFunc("GET /users/{id}", getUser(db))
	r.HandleFunc("POST /users/{id}", updateUser(db))
	r.HandleFunc("POST /users", createUser(db))
	r.HandleFunc("DELETE /users/{id}", deleteUser(db))

	//Start the server and start listening
	http.ListenAndServe(":8000", r)
}
