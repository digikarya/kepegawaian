package auth

import (
	"github.com/gorilla/mux"
    "gitlab.com/quick-count/go/auth/services"
	"gitlab.com/quick-count/go/config"
	"gitlab.com/quick-count/go/helper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

// App has router and db instances
type Auth struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Initializes the app with predefined configuration
func (a *Auth) Initialize(config *config.Config,route *mux.Router) {
	db, err := gorm.Open(mysql.Open(config.GetDSN()), &gorm.Config{
		PrepareStmt: true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatal("Could not connect database")
	}
	sqlDB, err := db.DB()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5*time.Minute)
	a.DB = db
	a.Router = route
	a.setRouters()
	//a.Router.Use(mux.CORSMethodMiddleware(a.Router))
	log.Println("Auth server is running")
}

// setRouters sets the all required routers
func (a *Auth) setRouters() {
	//// Routing for handling the projects
	a.Post("/auth", a.guest(services.AuthToken))
	//a.Delete("/logout", a.guest(controllers.AuthToken))
}


// Post wraps the router for POST method
func (a *Auth) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods(http.MethodPost,http.MethodOptions)
}


// Delete wraps the router for DELETE method
func (a *Auth) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// Run the app on it's router
func (a *Auth) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

type RequestHandlerFunction func(db *gorm.DB, w http.ResponseWriter, r *http.Request)

func (a *Auth) guest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		helper.CorsHelper(w,r)
		if r.Method == http.MethodOptions {
			return
		}
		handler(a.DB, w, r)
	}
}

func (a *Auth) guard(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		helper.CorsHelper(w,r)
		if r.Method == http.MethodOptions {
			return
		}
		if err := services.Authorization(a.DB,r,"admin");err != nil {
			helper.RespondJSONError(w,http.StatusUnauthorized,err)
			return
		}
		handler(a.DB, w, r)
	}
}




