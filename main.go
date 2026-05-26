package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/notwinterdust/otp-server/handlers"
	"github.com/notwinterdust/otp-server/middleware"
	"github.com/notwinterdust/otp-server/storage"
)

func main() {
	//add user command
	addUser := flag.NewFlagSet("add-user", flag.ExitOnError)
	addEmail := addUser.String("email", "", "User email")
	addPassword := addUser.String("password", "", "User password")

	if len(os.Args) > 1 && os.Args[1] == "add-user" {
		addUser.Parse(os.Args[2:])
		if *addEmail == "" || *addPassword == "" {
			log.Fatal("--email and --password are required")
		}
		db, err := openDB()
		if err != nil {
			log.Fatal(err)
		}
		if err := db.CreateUser(*addEmail, *addPassword); err != nil {
			log.Fatalf("create user: %v", err)
		}
		fmt.Printf("User %s created.\n", *addEmail)
		return
	}

	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}

	// if no user, create a user using the env file
	if exists, _ := db.UserExists(); !exists {
		email := os.Getenv("INITIAL_EMAIL")
		password := os.Getenv("INITIAL_PASSWORD")
		if email != "" && password != "" {
			if err := db.CreateUser(email, password); err != nil {
				log.Printf("warning: could not create initial user: %v", err)
			} else {
				log.Printf("Initial user %s created.", email)
			}
		}
	}

	jwtSecret := []byte(requireEnv("JWT_SECRET"))
	port := envOr("PORT", "8080")

	authH := &handlers.AuthHandler{DB: db, JWTSecret: jwtSecret}
	accountsH := &handlers.AccountsHandler{DB: db}
	authMW := middleware.Auth(jwtSecret)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/auth/login", authH.Login)
	mux.HandleFunc("GET /api/v1/health", handlers.Health)
	mux.Handle("POST /api/v1/accounts/sync", authMW(http.HandlerFunc(accountsH.Sync)))
	mux.Handle("GET /api/v1/accounts", authMW(http.HandlerFunc(accountsH.Pull)))

	addr := ":" + port
	log.Printf("OTP Sync Server v%s listening on %s", handlers.ServerVersion, addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func openDB() (*storage.DB, error) {
	path := envOr("DB_PATH", "/data/otp.db")
	return storage.Open(path)
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return v
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
