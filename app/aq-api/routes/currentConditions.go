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

package routes

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gitlab.com/google-cloud-ce/googlers/shaunmitchell/air-quality-mock-endpoint/app/api/data"
)

// NOTE: Some response objects have been commented out as there is no data in the current mock mock
// They are here for response completeness and can be added / uncommented if added to the database

// currentConditions JSON Request structure
type CurrentConditions struct {
	Location 			*latLng 				`json:"location,omitempty"`
	ExtraComputations	[]*extraComputations 	`json:"extraComputations,omitempty"` // see enums.go
	UaqiColorPalette 	*colorPalette 			`json:"uaqiColorPalette,omitempty"` // see enums.go
	CustomLocalAqi 		[]*customLocalAqi 		`json:"customLocalAqis,omitempty"`
	UniversalAqi		*bool 					`json:"universalAqi,omitempty"`
	LanguageCode		string 					`json:"languageCode,omitempty"`
}

type latLng struct {
	Latitude float64 	`json:"latitude,omitempty"`
	Longitude float64 	`json:"longitude,omitempty"`
}

type customLocalAqi struct {
	RegionCode string 	`json:"regionCode,omitempty"`
	Aqi string 			`json:"aqi,omitempty"`
}

// Removed Health Recommendations since they are not populated in DB
type Response struct {
	DateTime				string					`json:"dateTime,omitempty"`
	RegionCode				string					`json:"regionCode,omitempty"`
	Indexes					[]*airQualityIndex		`json:"indexes,omitempty"`
	Pollutants				[]*pollutant			`json:"pollutants,omitempty"`
//	HealthRecommendations	*healthRecommendations	`json:"healthRecommendations"`
}

type airQualityIndex struct {
	Code				string		`json:"code,omitempty"`
	DisplayName			string		`json:"displayName,omitempty"`
	AqiDisplay			string		`json:"aqiDisplay,omitempty"`
	Color				*color		`json:"color,omitempty"`
	Category			string		`json:"category,omitempty"`
	DominantPollutant	string		`json:"dominantPollutant,omitempty"`
	Aqi					int64		`json:"aqi,omitempty"`
}

type color struct {
	Red		float64	`json:"red,omitempty"`
	Green	float64	`json:"green,omitempty"`
	Blue	float64	`json:"blue,omitempty"`
	Alpha	float64	`json:"alpha,omitempty"`
}

type concentration struct {
	Units	*units	`json:"units,omitempty"` //see enums.go
	Value	int32	`json:"value,omitempty"`
}

type additionalInfo struct {
	Sources		string	`json:"sources,omitempty"`
	Effects		string	`json:"effects,omitempty"`
}

type pollutant struct {
	Code			string			`json:"code,omitempty"`
	DisplayName		string			`json:"displayName,omitempty"`
	FullName		string			`json:"fullName,omitempty"`
	Concentration	*concentration	`json:"concentration,omitempty"`
	AdditionalInfo	*additionalInfo	`json:"additionalInfo,omitempty"`
}

// type healthRecommendations struct {
// 	GeneralPopulation		string	`json:"generalPopulation"`
// 	Elderly					string	`json:"ederly"`
// 	LungDiseasePopulation	string	`json:"lungDiseasePopulation"`
// 	HeartDiseasePopulation	string	`json:"heartDiseasePopulation"`
// 	Athletes				string	`json:"athletes"`
// 	PregnantWomen			string	`json:"pregnantWomen"`
// 	Children				string	`json:"children"`
// }

// Routes is the sub-router for the /v1/currentConditions endpoint
// The endpoint has a :{action} attached to it so we have to read that in order to
// direct to the proper handler
func (cr CurrentConditions) Routes() http.Handler {
	mountHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		action := chi.URLParam(r, "action")

		switch action {
		case "lookup":
			cr.GetCurrentConditions(w, r)
		}
	})

	return mountHandler
}

//GetCurrentConditions is the controlling function for the lookup action
func (cr CurrentConditions) GetCurrentConditions(w http.ResponseWriter, r *http.Request) {
	cc := &CurrentConditions{}

	if err := render.Bind(r, cc); err != nil {
		log.Printf("render.bind error: %v", err)
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	results, err := data.FindLocation(cc.Location.Latitude, cc.Location.Longitude)

	if err != nil {
		log.Printf("unable to read data from firestore: %v", err)
		render.Render(w, r, ErrNotFound())
		return
	}

	var response *Response
	var category string
	var arQI []*airQualityIndex
//	var pollutants []*pollutant

	if results == (data.AQ{}) {
		log.Printf("no data found for lat/lng (%v,%v)", cc.Location.Latitude, cc.Location.Longitude)
		render.Render(w, r, ErrNotFound())
		return
	}

	c := color{
		Red: 	results.Red,
		Green: 	results.Green,
		Blue:	results.Blue,
		Alpha:	results.Alpha,
	}

	if cc.LanguageCode == "en" {
		category = results.Category
	} else {
		category = results.Categorie
	}

	arQI = append(arQI, &airQualityIndex{
		Code:				results.Code,
		DisplayName: 		results.DisplayName,
		AqiDisplay: 		results.AqiDisplay,
		Color:				&c,
		Category:			category,
		DominantPollutant: results.DominantPollutant,
//		Aqi:				int64(results.Aqi),
	})

	// units := UNIT_UNSPECIFIED

	// concentration := &concentration{
	// 	Units: &units,
	// 	Value: 0,
	// }

	// additionalInfo := &additionalInfo{
	// 	Sources: "",
	// 	Effects: "",
	// }

	// healthRecommendations := &healthRecommendations{
	// 	GeneralPopulation: "",
	// 	Elderly: "",
	// 	LungDiseasePopulation: "",
	// 	HeartDiseasePopulation: "",
	// 	Athletes: "",
	// 	PregnantWomen: "",
	// 	Children: "",
	// }

	// pollutants = append(pollutants, &pollutant{
	// 	Code:			"",
	// 	DisplayName: 	"",
	// 	FullName:		"",
	// 	Concentration:  concentration,
	// 	AdditionalInfo:	additionalInfo,
	// })

	response = new(Response)
	response.DateTime = results.DateTime
	response.RegionCode = results.RegionCode
	response.Indexes = arQI
//	response.Pollutants = pollutants
//	response.HealthRecommendations = healthRecommendations

	render.DefaultResponder(w, r, response)
}

// Bind will unmarshal the JSON request payload into the CurrentConditions struct
// Some default data / error processing is done here for the mock
func (cr *CurrentConditions) Bind(r *http.Request) error {
	// Default to english if languagecode is not set or not set to either en or fr
	if cr.LanguageCode != "en" && cr.LanguageCode != "fr" {
		cr.LanguageCode = "en"
	}

	// Default the universalA to true if not set
	var universalAqDefault = func(b bool) *bool {
		return &b
	}

	if cr.UniversalAqi == nil {
		cr.UniversalAqi = universalAqDefault(true)
	}


	return nil
}