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

type AQ struct {
	Latitude			float64		`firestore:"latitude"`
	Longitude			float64		`firestore:"longitude"`
	DateTime			string		`firestore:"dateTime,omitempty"`
	RegionCode			string		`firestore:"regionCode,omitempty"`
	Code				string		`firestore:"code,omitempty"`
	DisplayName			string		`firestore:"displayName,omitempty"`
	Aqi					int32		`firestore:"aqi,omitempty"`
	AqiDisplay			string		`firestore:"aqiDisplay,omitempty"`
	Red					float64		`firestore:"red,omitempty"`
	Green				float64		`firestore:"green,omitempty"`
	Blue				float64		`firestore:"blue,omitempty"`
	Alpha				float64		`firestore:"alpha,omitempty"`
	DominantPollutant	string		`firestore:"dominantPollutant,omitempty"`
	Category			string		`firestore:"category,omitempty"`
	Categorie			string		`firestore:"categorie,omitempty"`
}