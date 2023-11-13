// Copyright 2023 Shaun Mitchell

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"net/http"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/spf13/viper"
	"cloud.google.com/go/profiler"

	"gitlab.com/google-cloud-ce/googlers/shaunmitchell/air-quality-mock-endpoint/app/api/routes"
	"gitlab.com/google-cloud-ce/googlers/shaunmitchell/air-quality-mock-endpoint/app/api/data"
)

func main() {
	// Read config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.SetEnvPrefix("aq")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config file found")
		}
	}

	cfg := profiler.Config{
		Service:        "aq-mock-service",
		ServiceVersion: "0.0.4",
		// ProjectID must be set if not running on GCP.
		ProjectID: viper.GetString("project_id"),

		// For OpenCensus users:
		// To see Profiler agent spans in APM backend,
		// set EnableOCTelemetry to true
		// EnableOCTelemetry: true,
	}

	if err := profiler.Start(cfg); err != nil {
		checkErr(err)
	}

	err := data.CheckForDefaultData()
	if err != nil {
		checkErr(err)
	}

	port := viper.GetString("port")
	address := viper.GetString("address")

	r := chi.NewRouter()

	// Register middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Invalid Request"))
		checkErr(err)
	})

	r.Mount("/v1/currentConditions:{action}", routes.CurrentConditions{}.Routes())

	log.Printf("Starting server at %s:%s", address, port)
	checkErr(http.ListenAndServe(address + ":" + port, r))
}

// checkErr will check if the error object is set and then
// log the error.
// @todo Add in ability to respond to writeback to the API response writer (if needed?)
func checkErr(err error,) {
	if err != nil {
		log.Println(err)
	}
}