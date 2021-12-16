package main

import (
	"context"
	"flag"
	"fmt"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
	"go-open-api/db"
	"go-open-api/handlers"
	"go-open-api/middlewares"
	"go-open-api/repository"
	"go-open-api/services"
	"go-open-api/utils"
	"log"
	"net/http"
)

var port int

func init() {
	flag.IntVar(&port, "port", 8080, "set api port")
	db.LoadConfigsFromFlags(flag.CommandLine)
	flag.Parse()
}

func main() {
	cfg := db.NewConfig()

	log.Println(cfg.Source())

	dbConn := db.New(cfg)
	defer utils.HandleClose(dbConn)

	err := dbConn.PingContext(context.Background())
	if err != nil {
		log.Panicln(err)
	}
	log.Println("database connected.")

	usersRepository := repository.NewUsersRepository(dbConn)
	usersService := services.NewUsersService(usersRepository)

	ws := new(restful.WebService)
	restful.Add(ws)

	restful.Filter(middlewares.NCSACommonLogFormatLogger())
	middlewares.CORS(restful.DefaultContainer)
	restful.Filter(restful.DefaultContainer.OPTIONSFilter)

	handlers.RegisterUsersHandlers(usersService, ws)

	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(), // you control what services are visible
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject,
	}

	restful.Add(restfulspec.NewOpenAPIService(config))

	fileServer := http.FileServer(http.Dir("./public"))

	http.Handle("/apidocs/", http.StripPrefix("/apidocs/", fileServer))

	log.Printf("start listening on localhost:%d\n", port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: restful.DefaultContainer,
	}

	log.Fatalln(server.ListenAndServe())
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "Users Service",
			Description: "Resource for managing Users",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "Eric Lau",
					Email: "eric.devtt@gmail.com",
					URL:   "https://github.com/EricLau1",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "MIT",
					URL:  "http://mit.org",
				},
			},
			Version: "1.0.0",
		},
	}
	swo.Tags = []spec.Tag{
		{
			TagProps: spec.TagProps{
				Name:        "users",
				Description: "Managing users",
			},
		},
	}
}
