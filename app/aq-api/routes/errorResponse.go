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

import(
	"net/http"
	"errors"

	"github.com/go-chi/render"
)


type ErrResponse struct {
	Error   error `json:"-"` // low-level runtime error
	Er		*er	`json:"error,omitempty"` // QA Error response structure
}

type er struct {
	Code		int		`json:"code,omitempty"`
	Message		string	`json:"message,omitempty"`
	Status		string	`json:"status,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.Er.Code)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Error:	err,
		Er:		&er{
			Code: 400,
			Message: err.Error(),
			Status:	"INVALID_REQUEST",
		},
	}
}

func ErrNotFound() render.Renderer {
	return &ErrResponse {
		Error:	errors.New(""),
		Er:		&er {
			Code: 400,
			Message: "information is unavailable for this location. Please try a different location.",
			Status: "INVALID_ARGUMENT",
		},
	}
}