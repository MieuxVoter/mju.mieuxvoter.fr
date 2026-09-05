// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/tyler-sommer/stick"
	"github.com/tyler-sommer/stick/twig"
	"golang.org/x/text/language"
	"log/slog"
	"main/src/locales"
	"main/src/public"
	"main/src/templates"
	"main/src/version"
	"net/http"
	"os"
)

// polyglotKey MUST be defined in all available language files.
// We use it as a bit of a workaround to detect a usable language for the user.
var polyglotKey = "AppTitle"

func main() {

	logger := slog.Default()
	loadDotEnv()

	serverPort, foundServerPort := os.LookupEnv("WEB_PORT")
	if !foundServerPort {
		serverPort = "8034"
		logger.Info(fmt.Sprintf(
			"Environment variable WEB_PORT was not found.  "+
				"Defaulting to WEB_PORT=%s",
			serverPort,
		))
	}

	logger.Info("Starting web server…")

	localization := &locales.Localization{
		Logger: logger,
	}
	localization.Init(language.English)

	templateEngine := twig.New(
		&templates.EmbedFSLoader{
			FS: templates.TemplatesFS,
		},
	)
	twigErr := templateEngine.Register(locales.LocalizationExtension{
		Localization: localization,
		Localizers:   make(map[string]*locales.Localizer),
	})
	if twigErr != nil {
		panic(twigErr)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, userLanguage := localization.NewLocalizerAndLanguage(
			polyglotKey,
			r.Header.Get("Accept-Language"),
		)

		err := templateEngine.Execute(
			"index.html.twig",
			w,
			map[string]stick.Value{
				"version":  version.GetVersion(),
				"language": userLanguage.String(),
				"url":      r.URL.String(),
			},
		)
		if err != nil {
			handleServerError(err, w)
			return
		}
	})

	router.Get("/p/{encodedPoll}", func(w http.ResponseWriter, r *http.Request) {
		_, userLanguage := localization.NewLocalizerAndLanguage(
			polyglotKey,
			r.Header.Get("Accept-Language"),
		)

		err := templateEngine.Execute(
			"poll.html.twig",
			w,
			map[string]stick.Value{
				"version":  version.GetVersion(),
				"language": userLanguage.String(),
				"url":      r.URL.String(),
			},
		)
		if err != nil {
			handleServerError(err, w)
			return
		}
	})

	router.Get("/b/{encodedBallots}", func(w http.ResponseWriter, r *http.Request) {
		_, userLanguage := localization.NewLocalizerAndLanguage(
			polyglotKey,
			r.Header.Get("Accept-Language"),
		)

		err := templateEngine.Execute(
			"ballots.html.twig",
			w,
			map[string]stick.Value{
				"version":  version.GetVersion(),
				"language": userLanguage.String(),
				"url":      r.URL.String(),
			},
		)
		if err != nil {
			handleServerError(err, w)
			return
		}
	})

	// We also want to serve some static files, like CSS and the favicon
	//staticFiles := http.FileServer(http.Dir("public")) // using dir
	staticFiles := http.FileServer(http.FS(public.EmbedFS)) // or using embed
	router.Handle("/*", http.StripPrefix("/", staticFiles))

	// Finally, let's start the webserver and wait for an interrupting signal
	logger.Info(fmt.Sprintf("Visit http://localhost:%s in your favorite web browser.", serverPort))
	_ = http.ListenAndServe(":"+serverPort, router)
}

func handleServerError(err error, writer http.ResponseWriter) {
	writer.WriteHeader(500)
	_, _ = writer.Write([]byte(err.Error()))
}

//func handleUserError(err error, writer http.ResponseWriter) {
//	writer.WriteHeader(401)
//	_, _ = writer.Write([]byte(err.Error()))
//}

// loadDotEnv loads Environment variables from files, for convenience.
func loadDotEnv() {
	err := godotenv.Load(".env.local")
	if err != nil {
		// It's okay, we do not actually *require* a local env, since we have no secrets.
		//fmt.Println("No .env.local file found.  You may create one by copying .env.")
	}
	err = godotenv.Load() // .env
	if err != nil {
		// It's okay, we do not actually *require* a local env, since we have no secrets.
		//fmt.Println("No .env file found.  Odd.")
	}
}
