package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	sitereports "github.com/KamiQuasi/site-reports-api"
	"github.com/boltdb/bolt"
	"github.com/husobee/vestigo"
)

// Time format for sorting "2006-01-02T15:04:05Z07:00"

var db *bolt.DB

func main() {
	var err error
	router := vestigo.NewRouter()

	router.SetGlobalCors(&vestigo.CorsAccessControl{
		AllowOrigin:   []string{"*", "*"},
		AllowHeaders:  []string{"Content-Type", "Access-Control-Allow-Origin"},
		ExposeHeaders: []string{"Content-Type", "Access-Control-Allow-Origin"},
	})

	db, err = bolt.Open("/data/site-reports.db", 0600, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("Error opening DB: %s", err))
	}

	// Static Component Routes
	bc := http.FileServer(http.Dir("./frontend/bower_components"))
	router.Get("/bower_components/*", http.StripPrefix("/bower_components/", bc).ServeHTTP)
	wc := http.FileServer(http.Dir("./frontend/src"))
	router.Get("/src/*", http.StripPrefix("/src/", wc).ServeHTTP)

	// Property Routes
	router.Get("/p", propertiesHandler(db))
	router.Post("/p/add", addPropertyHandler(db))
	router.SetCors("/p/add", &vestigo.CorsAccessControl{
		AllowMethods: []string{"POST"},
		AllowHeaders: []string{"Content-Type", "Access-Control-Allow-Origin"},
	})
	router.Get("/p/:pid", propertyHandler(db))
	router.Post("/p/:pid", editPropertyHandler(db))
	router.SetCors("/p/:pid", &vestigo.CorsAccessControl{
		AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Access-Control-Allow-Origin"},
	})

	// Site Routes
	router.Get("/s", sitesHandler(db))
	router.Post("/s/add", addSiteHandler(db))
	router.SetCors("/s/add", &vestigo.CorsAccessControl{
		AllowMethods: []string{"POST"},
		AllowHeaders: []string{"Content-Type", "Access-Control-Allow-Origin"},
	})
	router.Get("/s/:sid", siteHandler(db))
	router.Post("/s/:sid", editSiteHandler(db))
	router.SetCors("/s/:sid", &vestigo.CorsAccessControl{
		AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Access-Control-Allow-Origin"},
	})
	router.Post("/s/:sid/page", addPageHandler(db))
	router.Delete("/s/:sid/page", removePageHandler(db))
	router.SetCors("/s/:sid/page", &vestigo.CorsAccessControl{
		AllowMethods: []string{"POST"},
		AllowHeaders: []string{"Content-Type", "Access-Control-Allow-Origin"},
	})
	// Team Routes
	router.Get("/t", tempHandler(db))
	router.Post("/t/add", tempHandler(db))
	router.SetCors("/t/add", &vestigo.CorsAccessControl{
		AllowMethods: []string{"POST"},
		AllowHeaders: []string{"Content-Type", "Access-Control-Allow-Origin"},
	})
	router.Get("/t/:tid", tempHandler(db))
	router.Post("/t/:tid", tempHandler(db))
	router.SetCors("/t/:tid", &vestigo.CorsAccessControl{
		AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Access-Control-Allow-Origin"},
	})
	// Report Routes
	router.Get("/r", tempHandler(db))
	router.Post("/r/add", tempHandler(db))
	router.SetCors("/r/add", &vestigo.CorsAccessControl{
		AllowMethods: []string{"POST"},
		AllowHeaders: []string{"Content-Type", "Access-Control-Allow-Origin"},
	})
	router.Get("/r/:rid", tempHandler(db))
	router.Post("/r/:rid", tempHandler(db))
	router.SetCors("/r/:rid", &vestigo.CorsAccessControl{
		AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Access-Control-Allow-Origin"},
	})
	// User Routes
	router.Get("/u", tempHandler(db))
	router.Post("/u/add", tempHandler(db))
	router.SetCors("/u/add", &vestigo.CorsAccessControl{
		AllowMethods: []string{"POST"},
		AllowHeaders: []string{"Content-Type", "Access-Control-Allow-Origin"},
	})
	router.Get("/u/:uid", tempHandler(db))
	router.Post("/u/:uid", tempHandler(db))
	router.SetCors("/u/:uid", &vestigo.CorsAccessControl{
		AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Access-Control-Allow-Origin"},
	})
	router.Get("/*", homeHandler)

	bind := fmt.Sprintf(":%s", "8024")
	fmt.Printf("listening on %s...", bind)
	err = http.ListenAndServe(bind, router)
	//err = http.ListenAndServeTLS(bind, viper.GetString("cert"), viper.GetString("key"), nil)
	if err != nil {
		log.Fatal(fmt.Errorf("Listener Error: %s", err))
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
}

func tempHandler(db *bolt.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Coming Soon!")
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Coming Soon!")
}

func addPropertyHandler(db *bolt.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var property sitereports.Property
		err := json.NewDecoder(r.Body).Decode(&property)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = property.Add(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(property)
		log.Printf("Property added: %s", property)
	}
}

func editPropertyHandler(db *bolt.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var property sitereports.Property
		pid := vestigo.Param(r, "pid")
		err := json.NewDecoder(r.Body).Decode(&property)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		property.ID = pid
		property.Update(db)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(property)
		log.Printf("Property edited: %s", property)
	}
}

func propertiesHandler(db *bolt.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		properties := sitereports.GetProperties(db)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(properties)

		log.Printf("Properties returned: %s", properties)
	}
}

func propertyHandler(db *bolt.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pid := vestigo.Param(r, "pid")
		property, err := sitereports.GetProperty(pid, db)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			log.Println(err)
			json.NewEncoder(w).Encode(sitereports.Property{})
		} else {
			json.NewEncoder(w).Encode(property)
		}
		log.Printf("Property returned: %s", property)
	}
}

func teamHandler(db *bolt.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func reportHandler(db *bolt.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func userHandler(db *bolt.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func svcHandler(db *bolt.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&sitereports.Site{ID: "rhd", URL: "https://developers.redhat.com", Name: "Red Hat Developers"})
	}
}

func addPageHandler(db *bolt.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sid := vestigo.Param(r, "sid")
		site, err := sitereports.GetSite(sid, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var page sitereports.Page
		err = json.NewDecoder(r.Body).Decode(&page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		site.AddPage(page, db)
		log.Printf("Page added: %s", page.Name)
		w.WriteHeader(200)
	}
}

func removePageHandler(db *bolt.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sid := vestigo.Param(r, "sid")
		site, err := sitereports.GetSite(sid, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var page sitereports.Page
		err = json.NewDecoder(r.Body).Decode(&page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		site.RemovePage(page, db)
		log.Printf("Page removed: %s", page.Name)
		w.WriteHeader(200)
	}
}

func sitesHandler(db *bolt.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sites []sitereports.Site
		err := db.View(func(tx *bolt.Tx) error {
			p := tx.Bucket([]byte("Sites"))
			if p != nil {
				c := p.Cursor()

				for k, _ := c.First(); k != nil; k, _ = c.Next() {
					site, _ := sitereports.GetSite(string(k), db)
					sites = append(sites, site)
				}
			}
			return nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode([]sitereports.Site{})
		} else {
			json.NewEncoder(w).Encode(sites)
		}
		j, err := json.Marshal(sites)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Sites listed: %s", string(j))
	}
}

func siteHandler(db *bolt.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var site sitereports.Site
		sid := vestigo.Param(r, "sid")
		site, err := sitereports.GetSite(sid, db)
		if err != nil {
			log.Println(err)
			json.NewEncoder(w).Encode(sitereports.Site{})
		} else {
			json.NewEncoder(w).Encode(site)
		}
	}
}

func editSiteHandler(db *bolt.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var site sitereports.Site
		sid := vestigo.Param(r, "sid")
		err := json.NewDecoder(r.Body).Decode(&site)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		site.ID = sid
		site.Update(db)

		json.NewEncoder(w).Encode(site)
	}
}

func addSiteHandler(db *bolt.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var site sitereports.Site
		err := json.NewDecoder(r.Body).Decode(&site)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = site.Add(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sitereports.Site{})
		} else {
			json.NewEncoder(w).Encode(site)
		}

	}
}
