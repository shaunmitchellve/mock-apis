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
	"encoding/json"
)

// extraComputations and colorPallette ENUM Setup
const (
	EXTRA_COMPUTATION_UNSPECIFIED extraComputations	= iota
	LOCAL_AQI
	HEALTH_RECOMMENDATIONS
	POLLUTANT_ADDITIONAL_INFO
	DOMINANT_POLLUTANT_CONCENTRATION
	POLLUTANT_CONCENTRATION
)

const (
	COLOR_PALETTE_UNSPECIFIED colorPalette	= iota
	RED_GREEN
	INDIGO_PERSION_DARK
	INDIGO_PERSION_LIGHT
)

const (
	UNIT_UNSPECIFIED units = iota
	PARTS_PER_BILLION
	MICROGRAMS_PER_CUBIC_METER
)

type extraComputations int
type colorPalette int
type units int

func (ec extraComputations) String() string {
	return [...]string{"EXTRA_COMPUTATION_UNSPECIFIED", "LOCAL_AQI", "HEALTH_RECOMMENDATIONS", "POLLUTANT_ADDITIONAL_INFO", "DOMINANT_POLLUTANT_CONCENTRATION", "POLLUTANT_CONCENTRATION"}[ec]
}

func (ec *extraComputations) FromString(extrac string) extraComputations {
	return map[string]extraComputations {
		"EXTRA_COMPUTATION_UNSPECIFIED": EXTRA_COMPUTATION_UNSPECIFIED,
		"LOCAL_AQI": LOCAL_AQI,
		"HEALTH_RECOMMENDATIONS": HEALTH_RECOMMENDATIONS,
		"POLLUTANT_ADDITIONAL_INFO": POLLUTANT_ADDITIONAL_INFO,
		"DOMINANT_POLLUTANT_CONCENTRATION": DOMINANT_POLLUTANT_CONCENTRATION,
		"POLLUTANT_CONCENTRATION": POLLUTANT_CONCENTRATION,
	}[extrac]
}

func (ec *extraComputations) UnmarshalJSON(b []byte) error {
	var j string
	 err := json.Unmarshal(b, &j)
	 if err != nil {
		return err
	 }

	 *ec = ec.FromString(j)
	 return nil
}

func (ec extraComputations) MarshalJSON() ([]byte, error) {
	return json.Marshal(ec.String())
}

func (cp colorPalette) String() string {
	return [...]string{"COLOR_PALETTE_UNSPECIFIED", "RED_GREEN", "INDIGO_PERSIAN_DARK", "INDIGO_PERSIAN_LIGHT"}[cp]
}

func (cp *colorPalette) FromString(color string) colorPalette {
	return map[string]colorPalette {
		"COLOR_PALETTE_UNSPECIFIED": COLOR_PALETTE_UNSPECIFIED,
		"RED_GREEN": RED_GREEN,
		"INDIGO_PERSIAN_DARK": INDIGO_PERSION_DARK,
		"INDIGO_PERSIAN_LIGHT": INDIGO_PERSION_LIGHT,
	}[color]
}

func (cp *colorPalette) UnmarshalJSON(b []byte) error {
	var j string
	 err := json.Unmarshal(b, &j)
	 if err != nil {
		return err
	 }
	 *cp = cp.FromString(j)
	 return nil
}

func (cp colorPalette) MarshalJSON() ([]byte, error) {
	return json.Marshal(cp.String())
}

func (u units) String() string {
	return [...]string{"UNIT_UNSPECIFIED", "PARTS_PER_BILLION", "MICROGRAMS_PER_CUBIC_METER"}[u]
}

func (u *units) FromString(unit string) units {
	return map[string]units {
		"UNIT_UNSPECIFIED": UNIT_UNSPECIFIED,
		"PARTS_PER_BILLION": PARTS_PER_BILLION,
		"MICROGRAMS_PER_CUBIC_METER": MICROGRAMS_PER_CUBIC_METER,
	}[unit]
}

func (u *units) UnmarshalJSON(b []byte) error {
	var j string
	 err := json.Unmarshal(b, &j)
	 if err != nil {
		return err
	 }

	 *u = u.FromString(j)
	 return nil
}

func (u units) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}