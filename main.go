package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"neft.web/controllers"
	"neft.web/errorController"
	"neft.web/liro"
	"neft.web/middleware"
	"neft.web/rand"
	"neft.web/views"
)

var (
	isProd bool

	debug bool
)

func init() {
	flag.BoolVar(&isProd, "isProd", false, "This will ensure all pro vars are enabled")
	flag.BoolVar(&debug, "debug", false, "This will export all stats to file log.log")
}

func main() {
	views.InitTemplateController()
	go controllers.ReadInput()
	flag.Parse()
	errorController.InitLog(debug)

	errorController.InfoLogger.Println("Starting web server...")

	errorController.DebugLogger.Println("Configuring static pages...")
	staticC := controllers.NewStatic()

	errorController.DebugLogger.Println("Configuring users pages...")
	userC := controllers.NewUsers()

	errorController.DebugLogger.Println("Configuring liro modules...")

	liroC := liro.NewLiro()

	r := mux.NewRouter()
	appAssetsHandler := http.FileServer(http.Dir("./app-assets/"))
	appAssetsHandler = http.StripPrefix("/app-assets/", appAssetsHandler)
	r.PathPrefix("/app-assets/").Handler(appAssetsHandler)

	jsHandler := http.FileServer(http.Dir("./js/"))
	jsHandler = http.StripPrefix("/js/", jsHandler)
	r.PathPrefix("/js/").Handler(jsHandler)

	imgHandler := http.FileServer(http.Dir("./images/"))
	imgHandler = http.StripPrefix("/images/", imgHandler)
	r.PathPrefix("/images/").Handler(imgHandler)

	fontsHandler := http.FileServer(http.Dir("./fonts/"))
	fontsHandler = http.StripPrefix("/fonts/", fontsHandler)
	r.PathPrefix("/fonts/").Handler(fontsHandler)

	vendorHandler := http.FileServer(http.Dir("./vendor_web/"))
	vendorHandler = http.StripPrefix("/vendor/", vendorHandler)
	r.PathPrefix("/vendor/").Handler(vendorHandler)

	assetsHandler := http.FileServer(http.Dir("./assets/"))
	assetsHandler = http.StripPrefix("/assets/", assetsHandler)
	r.PathPrefix("/assets/").Handler(assetsHandler)

	assetsAppHandler := http.FileServer(http.Dir("./assets/"))
	assetsAppHandler = http.StripPrefix("/app-assets/", assetsAppHandler)
	r.PathPrefix("/app-assets/").Handler(assetsAppHandler)

	b, err := rand.Bytes(32)
	if err != nil {
		return
	}
	errorController.DebugLogger.Println("Configuring middleware...")

	errorController.IsPro = isProd

	csrfMw := csrf.Protect(b, csrf.Secure(isProd))

	userMW := middleware.User{}

	requireUseMW := middleware.RequireUser{
		User: userMW,
	}

	errorController.DebugLogger.Println("Loading routers")

	r.HandleFunc("/", staticC.NewHome).Methods("GET")

	r.NotFoundHandler = staticC.NotFound
	r.Handle("/505", staticC.Error).Methods("GET")

	// Login And Register

	r.HandleFunc("/signup", requireUseMW.Apply(userC.New)).Methods("GET")
	r.HandleFunc("/signup", requireUseMW.Apply(userC.Create)).Methods("POST")
	r.HandleFunc("/login", requireUseMW.Apply(userC.LoginNew)).Methods("GET")
	r.HandleFunc("/login", requireUseMW.Apply(userC.Login)).Methods("POST")
	r.HandleFunc("/logout", userC.Logout).Methods("POST")
	r.Handle("/forgot", userC.ForgotPwView).Methods("GET")
	r.HandleFunc("/forgot", userC.InitiateReset).Methods("POST")
	r.HandleFunc("/reset", userC.ResetPw).Methods("GET")
	r.HandleFunc("/reset", userC.CompleteReset).Methods("POST")

	r.HandleFunc("/liro", requireUseMW.RequireUser(liroC.New)).Methods("GET")
	r.HandleFunc("/liro/users", requireUseMW.RequireUser(liroC.UsersList)).Methods("GET")

	// Start server
	if debug {
		go middleware.PrintStats()
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000" // Default port if not specified
	}
	errorController.DebugLogger.Println("STARTING SERVER AT PORT: " + port)
	err = http.ListenAndServe(":"+port, csrfMw(userMW.Apply(r)))
	if err != nil {
		errorController.WarningLogger.Println("Can not start the server!" + err.Error())
	}
}
