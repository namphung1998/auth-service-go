package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	goJWT "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/namphung1998/auth-service-go/internal/data"
	"github.com/namphung1998/auth-service-go/internal/httpservice"
	"github.com/namphung1998/auth-service-go/internal/jwt"
	"github.com/namphung1998/auth-service-go/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	var mongoURI string
	flag.StringVar(&mongoURI, "mongo-uri", "", "URI for MongoDB connection")
	var mongoUsername string
	flag.StringVar(&mongoUsername, "mongo-username", "", "Username for MongoDB connection")
	var mongoPassword string
	flag.StringVar(&mongoPassword, "mongo-password", "", "Password for MongoDB connection")
	var mongoDatabase string
	flag.StringVar(&mongoDatabase, "mongo-database", "", "Database name for MongoDB connection")
	var jwtSigningKey string
	flag.StringVar(&jwtSigningKey, "jwt-secret", "", "Signing key for generating JWTs")

	flag.Parse()

	clientOpts := options.Client().
		ApplyURI(mongoURI).
		SetAuth(options.Credential{
			Username: mongoUsername,
			Password: mongoPassword,
		})

	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Connect(context.Background()); err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	fmt.Println(client.Ping(context.Background(), readpref.Primary()))

	userRepo := data.NewUserRepo(client.Database("auth-app"))
	jwtService := jwt.NewService(goJWT.SigningMethodHS256, []byte(jwtSigningKey))
	userService := service.NewUser(userRepo, jwtService)
	handler := httpservice.NewHandler(userService)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Post("/greet", handler.HandleGreet())

	r.Route("/v1", func(r chi.Router) {
		r.Post("/users", handler.HandleCreateUser())
		r.Post("/users/login", handler.HandleLogin())
	})

	http.ListenAndServe(":3090", r)
}
