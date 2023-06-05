// package main

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"

// 	_ "github.com/go-sql-driver/mysql"
// )
// type User struct {
// 	ID             int    `json:"id"`
// 	FirstName      string `json:"first_name"`
// 	Country        string `json:"country"`
// 	ProfilePicture string `json:"profile_picture"`
// }


// var db *sql.DB

// func main() {
// 	// MySQL database connection parameters
// 	dbUser := "root"
// 	dbPass := ""
// 	dbHost := "localhost"
// 	dbPort := "3306"
// 	dbName := "user_activities"

// 	// Create the data source name (DSN) string
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

// 	// Open a connection to the MySQL database
// 	var err error
// 	db, err = sql.Open("mysql", dsn)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	// Test the database connection
// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Define API routes
// 	http.HandleFunc("/users", createUserHandler)
// 	http.HandleFunc("/users/{id}", updateUserHandler)
// 	// http.HandleFunc("/users/delete/{id}", deleteUserHandler)

// 	// Start the server
// 	log.Println("Server is running on http://localhost:8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// func createUserHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Parse request body
// 	var user User
// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	// Insert the user into the database
// 	result, err := db.Exec("INSERT INTO users (first_name, country, profile_picture) VALUES (?, ?, ?)", user.FirstName, user.Country, user.ProfilePicture)
// 	if err != nil {
// 		http.Error(w, "Failed to create user", http.StatusInternalServerError)
// 		return
// 	}

// 	// Get the ID of the created user
// 	userID, _ := result.LastInsertId()

// 	// Return the created user
// 	user.ID = int(userID)
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(user)
// }

// func updateUserHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPatch {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Get the user ID from the URL path
// 	// You can use a routing library to handle URL parameters more elegantly
// 	// Here, we assume the ID is passed as a path parameter
// 	userID := r.URL.Path[len("/users/"):]
// 	fmt.Println(userID)

// 	// Parse request body
// 	var user User
// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	// Update the user in the database
// 	_, err = db.Exec("UPDATE users SET first_name=?, country=?, profile_picture=? WHERE id=?", user.FirstName, user.Country, user.ProfilePicture, userID)
// 	if err != nil {
// 		http.Error(w, "Failed to update user", http.StatusInternalServerError)
// 		return
// 	}

// 	// Return the updated user
// 	w.WriteHeader(http.StatusOK)
// }

// func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodDelete {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Get the user ID from the URL path
// 	// You can use a routing library to handle URL parameters more elegantly
// 	// Here, we assume the ID is passed as a path parameter
// 	userID := r.URL.Path[len("/users/delete"):]
// 	fmt.Println(userID)


// 	// Delete the user from the database
// 	_, err := db.Exec("DELETE FROM users WHERE id=?", userID)
// 	if err != nil {
// 		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
// 		return
// 	}

// 	// Delete associated activity logs for the user
// 	_, err = db.Exec("DELETE FROM activity_logs WHERE user_id=?", userID)
// 	if err != nil {
// 		http.Error(w, "Failed to delete activity logs", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// }


package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID             int    `json:"id"`
	FirstName      string `json:"first_name"`
	Country        string `json:"country"`
	ProfilePicture string `json:"profile_picture"`
}

var db *sql.DB

func main() {
	// MySQL database connection parameters
	dbUser := "root"
	dbPass := ""
	dbHost := "localhost"
	dbPort := "3306"
	dbName := "user_activities"

	// Create the database connection
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	var err error
	db, err = sql.Open("mysql", dbURI)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize the router
	router := mux.NewRouter()

	// Define the API endpoints
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{id}", updateUser).Methods("PATCH")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	// Start the server
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Insert the user into the database
	result, err := db.Exec("INSERT INTO users (first_name, country, profile_picture) VALUES (?, ?, ?)", user.FirstName, user.Country, user.ProfilePicture)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Get the auto-generated user ID
	userID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get user ID", http.StatusInternalServerError)
		return
	}

	// Set the user ID in the response
	user.ID = int(userID)

	// Return the created user as the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update the user in the database
	_, err = db.Exec("UPDATE users SET first_name=?, country=?, profile_picture=? WHERE id=?", user.FirstName, user.Country, user.ProfilePicture, id)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Delete the user from the database
	_, err = db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
