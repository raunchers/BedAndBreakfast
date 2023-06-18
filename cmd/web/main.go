package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/raunchers/BedAndBreakfast/pkg/config"
	"github.com/raunchers/BedAndBreakfast/pkg/handlers"
	"github.com/raunchers/BedAndBreakfast/pkg/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

// create the app variable
var app config.AppConfig
var session *scs.SessionManager

func main() {

	// change this to true when in production
	app.InProduction = false

	// create a new session
	session = scs.New()
	// how long the session will last for
	session.Lifetime = 24 * time.Hour
	// should the cookie persist when the browser window is closed
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	// should the cookie be secure
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// Get the template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	// Assign the template cache
	app.TemplateCache = tc
	// Development mode = false, since it is false the RenderTemplate func rebuilds the cache
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	// call the new template cache
	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	// _ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
