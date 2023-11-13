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
	"context"
	"os"
	"encoding/base64"
	"io"
	"bufio"
	"strings"
	"regexp"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
)

var _ fn.Runner = &Base64File{}

type Base64File struct {
	Data map[string]string
}

// EmptyfnConfig is a workaround since kpt creates a FunctionConfig placeholder if users don't provide the functionConfig.
// `kpt fn eval` uses placeholder ConfigMap with name "function-input"
// `kpt fn render` uses placeholder "{}"
func EmptyfnConfig(o *fn.KubeObject) bool {
	if o.GetKind() == "ConfigMap" && o.GetName() == "function-input" {
		data, _, _ := o.NestedStringMap("data")
		return len(data) == 0
	}
	if o.GetKind() == "" && o.GetName() == "" {
		return true
	}

	return false
}

// Run is the core function that will read the reasourcelist passed in, check to match the target kind and filed
// then open the source file, encode it and update the targetField
// TODO: Make the matching of the field path more robust past this 1 use-case
func (r *Base64File) Run(ctx *fn.Context, fConfig *fn.KubeObject, items fn.KubeObjects, results *fn.Results) bool {
	 objects := items.WhereNot(fn.IsLocalConfig)

	if objects.Len() == 0 || objects[0] == nil {
		results.Infof("no input resources")
		return true
	}
	if EmptyfnConfig(fConfig) {
		return false
	}

	if len(r.Data["target-field"]) == 0 || len(r.Data["target-kind"]) == 0  || len(r.Data["source-file"]) == 0 {
		results.Warningf("missing one or more required arguments:`target-field`, `target-kind`, `source-file` in FunctionConfig")
		return false
	}

	for _, o := range objects {
		if o.GetKind() == r.Data["target-kind"] {
			f, err := os.Open(r.Data["source-file"])
			if err != nil {
				results.Warningf("Unable to open source-file (%s)", r.Data["source-file"])
				return false
			}

			defer f.Close()

			content, err := io.ReadAll(bufio.NewReader(f))
			if err != nil {
				results.WarningE(err)
				return false
			}

			encoded := base64.StdEncoding.EncodeToString(content)

			nestedPath := strings.Split(r.Data["target-field"], ".")

			reg, err := regexp.Compile(`\[[0-9]\]`)

			for i, e := range nestedPath {
				if reg.MatchString(e){
					se, f, err := o.NestedSlice(nestedPath[0:i]...)

					if err != nil{
						results.WarningE(err)
						return false
					}

					if f {
						nestedPath[i] = reg.ReplaceAllString(nestedPath[i], ``)
						if err = se[0].SetNestedField(&encoded, nestedPath[i:]...); err != nil {
							results.WarningE(err)
							return false
						}
					}
				}

			}

			 if err != nil {
			 	results.WarningE(err)
			 	return false
			 }
		}
	}

	results.Infof("%s (%s) updated to encoded string of %s", r.Data["target-kind"], r.Data["target-field"], r.Data["source-file"])
    return true
}

// main is th entrypoint for ktp function
func main() {
	runner := fn.WithContext(context.Background(), &Base64File{})
    if err := fn.AsMain(runner); err != nil {
        os.Exit(1)
    }
}