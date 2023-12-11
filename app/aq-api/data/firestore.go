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

package data

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"github.com/spf13/viper"
)

var ctx context.Context

// create the Firestore client
func createClient(ctx context.Context) *firestore.Client {
	client, err := firestore.NewClientWithDatabase(ctx, viper.GetString("project_id"), viper.GetString("database_id"))
	checkError("Failed to create Firestore client", err)

	return client
}

// initilize the Firestore client
func firestoreInit() (*firestore.Client, *firestore.CollectionRef) {
	ctx = context.Background()

	client := createClient(ctx)

	col := client.Collection(viper.GetString("collection_name"))

	return client, col
}

// check to see if the Firestore collection exists, if not then load in the defalt data
func CheckForDefaultData() error {
	client, col := firestoreInit()

	defer client.Close()

	q := col.Limit(1).Documents(ctx)

	defer q.Stop()

	_, err := q.Next()
	if err != iterator.Done {
		return err
	}

	log.Println("Collection empty. Setting up default data")
	defaultData(client)

	return nil

}

// find the location record
func FindLocation(lat float64, lng float64) (AQ, error) {
	_, col := firestoreInit()

	if lat != 0 && lng != 0 {
		q := col.Where("latitude", "==", lat).Where("longitude", "==", lng).Documents(ctx)

		defer q.Stop()

		d, err := q.Next()
		if err == iterator.Done {
			return AQ{}, nil
		}
		if err != nil {
			return AQ{}, err
		}

		return unMarshallAQ(d.Data()), nil
	}

	return AQ{}, errors.New("unable to find data for location")
}

// un marshall the Firestore document and put it into the AQ struct. Add missing data if needed
func unMarshallAQ(record map[string]interface{}) AQ {
	if checkNil(record["aqi"]) {
		record["aqi"] = convertInt("0")
	}

	if checkNil(record["aqiDisplay"]) {
		record["aqiDisplay"] = ""
	}

	if checkNil(record["alpha"]) {
		record["alpha"] = 0.0
	}
	
	if checkNil(record["green"]) {
		record["green"] = 0.0
	}

	if checkNil(record["blue"]) {
		record["blue"] = 0.0
	}

	if checkNil(record["red"]) {
		record["red"] = 0.0
	}
	return AQ{
		Longitude: record["longitude"].(float64),
		Latitude: record["latitude"].(float64),
		DateTime: time.Now().UTC().Format(time.RFC3339),
		RegionCode: record["regionCode"].(string),
		Code: record["code"].(string),
		DisplayName: record["displayName"].(string),
		Aqi: int32(record["aqi"].(int32)),
		AqiDisplay: record["aqiDisplay"].(string),
		Red: record["red"].(float64),
		Green: record["green"].(float64),
		Blue: record["blue"].(float64),
		Alpha: record["alpha"].(float64),
		Category: record["category"].(string),
		DominantPollutant: record["dominantPollutant"].(string),
		Categorie: record["categorie"].(string),
	}
}

// Load default data into the Firestore object store.
func defaultData(client *firestore.Client) {
	bw := client.BulkWriter(ctx)

	file, err := os.Open("./data/defaultData.csv")
	if err != nil {
		log.Fatalf("Default data CSV file not found: %v", err)
	}

	defer file.Close()

	r := csv.NewReader(file)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Unable to read default data for Firestore setup: %v", err)
		}

		if record[0] != "latitude" {
			_, err := bw.Create(client.Collection(viper.GetString("collection_name")).NewDoc(), AQ{
				Latitude: convertFloat(record[0]),
				Longitude: convertFloat(record[1]),
				DateTime:record[2],
				RegionCode: record[3],
				Code: record[4],
				DisplayName: record[5],
				Aqi: convertInt(record[6]),
				AqiDisplay: record[7],
				Red: convertFloat(record[8]),
				Green: convertFloat(record[9]),
				Blue: convertFloat(record[10]),
				Alpha: convertFloat(record[11]),
				DominantPollutant: record[12],
				Category: record[13],
				Categorie: record[14],
			})

			if err != nil {
				log.Printf("Unable to create bulk writer: %v", err)
			}
		}
	}
	bw.End()
}

// Check the error object, if not nil cause a fatal error
func checkError(msg string, e error) {
	if e != nil {
		log.Fatalf(msg + ": %v", e)
	}
}

// Check if interface is nil
func checkNil(d interface{}) bool {
	if d == nil {
		return true
	}

	return false
}

// convert a string to float64 type
func convertFloat(d string) float64 {
	if d != "" {
		r, err := strconv.ParseFloat(d, 64)
		checkError("Unable to convert string to float64", err)

		return float64(r)
	} else {
		return 0
	}
}

// convert a string to int32 type
func convertInt(d string) int32 {
	if d != "" {
		r, err := strconv.ParseInt(d, 10, 32)
		checkError("Unable to convert string to int32", err)

		return int32(r)
	} else {
		return 0
	}
}