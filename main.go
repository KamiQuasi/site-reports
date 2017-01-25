package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/quad"
	"github.com/spf13/viper"
)

// func loadSite(site string) (*Site, error) {
// 	store, err := cayley.NewMemoryGraph()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	p := cayley.StartPath(store, quad.String(site))

// 	err = p.Iterate(nil).EachValue(nil, func(value quad.Value) {
// 		nativeValue := quad.NativeOf(value)

// 	})

// 	store.AddQuad(quad.Make("phrase of the day", "is of course", "Hello World!", nil))

// 	// Now we iterate over results. Arguments:
// 	// 1. Optional context used for cancellation.
// 	// 2. Flag to optimize query before execution.
// 	// 3. Quad store, but we can omit it because we have already built path with it.
// 	err = p.Iterate(nil).EachValue(nil, func(value quad.Value) {
// 		nativeValue := quad.NativeOf(value) // this converts RDF values to normal Go types
// 		fmt.Println(nativeValue)
// 	})
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// }

// Site struct
type Site struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

func main() {
	viper.SetConfigName("app") // no need to include file extension
	viper.AddConfigPath("config")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.SetDefault("ip", "127.0.0.1")
	viper.SetDefault("port", "8080")
	// viper.SetDefault("db", "/data/sites.boltdb")
	viper.SetEnvPrefix("OPENSHIFT")
	viper.BindEnv("ip", "GO_IP")
	viper.BindEnv("port", "GO_PORT")

	// graph.InitQuadStore("bolt", viper.GetString("db"), nil)

	// store, err := cayley.NewGraph("bolt", viper.GetString("db"), nil)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	store, err := cayley.NewMemoryGraph()
	if err != nil {
		log.Fatalln(err)
	}

	rhd := &Site{ID: "rhd", Name: "Red Hat Developers", URL: "http://developers.redhat.com"}
	store.AddQuad(quad.Make(quad.StringToValue(rhd.ID), "name", quad.StringToValue(rhd.Name), nil))
	store.AddQuad(quad.Make(quad.StringToValue(rhd.ID), "url", quad.StringToValue(rhd.URL), nil))
	store.Close()

	bc := http.FileServer(http.Dir("./bower_components"))
	nm := http.FileServer(http.Dir("./node_modules"))
	assets := http.FileServer(http.Dir("./assets"))
	http.Handle("/bower_components/", http.StripPrefix("/bower_components/", bc))
	http.Handle("/node_modules/", http.StripPrefix("/node_modules/", nm))
	http.Handle("/assets/", http.StripPrefix("/assets/", assets))

	http.HandleFunc("/setup/", setupHandler)
	http.HandleFunc("/svc/", serviceHandler)
	http.HandleFunc("/site/", siteHandler)
	http.HandleFunc("/", home)

	bind := fmt.Sprintf("%s:%s", viper.Get("ip"), viper.Get("port"))
	fmt.Printf("listening on %s...", bind)
	//err = http.ListenAndServeTLS(bind, "cert.pem", "key.pem", nil)
	http.ListenAndServe(bind, nil)
	if err != nil {
		panic(err)
	}
}

func setupHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/setup.html")
	t.Execute(w, nil)
}

func serviceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s - %s\n", r.Method, r.URL)
	fmt.Fprint(w, "Service!")
}

func siteHandler(w http.ResponseWriter, r *http.Request) {
	//store, err := cayley.NewGraph("bolt", viper.GetString("db"), nil)
	store, err := cayley.NewMemoryGraph()
	if err != nil {
		log.Fatalln(err)
	}

	// rhd := &Site{ID: "rhd", Name: "Red Hat Developers", URL: "http://developers.redhat.com"}
	// store.AddQuad(quad.Make(rhd.ID, "name", rhd.Name, nil))
	site := &Site{ID: r.URL.Path[len("/site/"):]}

	p := cayley.StartPath(store).Is(quad.String(site.ID)).Out()
	vals, err := p.Iterate(nil).AllValues(nil)
	if err != nil {
		log.Fatalln(err)
	} else if len(vals) == 0 {
		site.Name = "Unnamed Site"
		site.URL = "http://127.0.0.1"
		store.AddQuad(quad.Make(site.ID, "name", site.Name, nil))
		store.AddQuad(quad.Make(site.ID, "url", site.URL, nil))
	} else {
		site.Name = vals[0].Native().(string)
		site.URL = vals[1].Native().(string)
	}
	// rhd := &Site{ID: "rhd", Name: "Red Hat Developers", URL: "http://developers.redhat.com"}
	// store.AddQuad(quad.Make(rhd.ID, "name", rhd.Name, nil))

	t, _ := template.ParseFiles("templates/site.html")
	t.Execute(w, site)
	store.Close()
}
func home(w http.ResponseWriter, r *http.Request) {
	var sites []Site
	// Create a brand new graph
	//store, err := cayley.NewGraph("bolt", viper.GetString("db"), nil)
	store, err := cayley.NewMemoryGraph()
	if err != nil {
		log.Fatalln(err)
	}

	// store.AddQuad(quad.Make(group, "type", "group", nil))
	// store.AddQuad(quad.Make(siteURL, "type", "property", nil))
	// store.AddQuad(quad.Make(siteURL, "name", "Red Hat Developers", nil))
	// store.AddQuad(quad.Make(siteURL, "allows protocol", "http", nil))
	// store.AddQuad(quad.Make(siteURL, "allows protocol", "https", nil))
	// store.AddQuad(quad.Make(siteURL, "scores", 72, nil))
	// Now we create the path, to get to our data
	p := cayley.StartPath(store).Has("name").Save("name", "name")

	// Now we iterate over results. Arguments:
	// 1. Optional context used for cancellation.
	// 2. Flag to optimize query before execution.
	// 3. Quad store, but we can omit it because we have already built path with it.
	err = p.Iterate(nil).EachValue(nil, func(value quad.Value) {
		nativeValue := quad.StringOf(value) // this converts RDF values to normal Go types
		fmt.Println(json.Marshal(quad.StringOf(value)))
		site := Site{ID: nativeValue}
		sites = append(sites, site)
	})
	if err != nil {
		log.Fatalln(err)
	}

	t, _ := template.ParseFiles("templates/home.html")
	pusher, ok := w.(http.Pusher)
	if ok { // Push is supported. Try pushing rather than waiting for the browser.
		if err := pusher.Push("/assets/raw-element.html", nil); err != nil {
			log.Printf("Failed to push: %v", err)
		}
		if err := pusher.Push("/assets/styles.css", nil); err != nil {
			log.Printf("Failed to push: %v", err)
		}
	}
	t.Execute(w, sites)
	store.Close()
	// client := github.NewClient(nil)

	// orgs, _, err := client.Organizations.List("KamiQuasi", nil)

	// developers := Site{URL: "https://github.com/redhat-developer/developers.redhat.com"}
	// developers.analyze()

	// b, err := json.Marshal(developers.Scores[0])
	// if err != nil {

	// }
	//fmt.Fprintf(res, "%s", b)
}
