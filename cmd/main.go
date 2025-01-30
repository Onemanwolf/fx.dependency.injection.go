package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"fx.dependency.injection/models"
	"fx.dependency.injection/repositories"
	"fx.dependency.injection/services"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			repositories.NewAzureSQLUserRepository,
			func(repo *repositories.AzureSQLUserRepository) repositories.UserRepository {
				return repo
			},
			services.NewUserService,
			NewRouter,
		),
		fx.Invoke(StartServer),
	)

	app.Run()
}

func NewRouter(userService *services.UserService) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/users", createUserHandler(userService)).Methods("POST")
	r.HandleFunc("/users/{id}", getUserHandler(userService)).Methods("GET")
	r.HandleFunc("/users/{id}", updateUserHandler(userService)).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUserHandler(userService)).Methods("DELETE")
	return r
}

func StartServer(lc fx.Lifecycle, router *mux.Router) {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go srv.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}

func createUserHandler(userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := userService.CreateUser(r.Context(), &user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

func getUserHandler(userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		user, err := userService.GetUserByID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if user == nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

func updateUserHandler(userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user.ID = id

		if err := userService.UpdateUser(r.Context(), &user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

func deleteUserHandler(userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		if err := userService.DeleteUser(r.Context(), id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
