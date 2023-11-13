# Copyright 2023 Shaun Mitchell

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#  	http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#!/bin/bash

APIGATEWAYURL='<API-GATEWAY-URL>'
#APIGATEWAYURL='localhost:8080'
APIKEY='<API-KEY>'

curl -H "Content-Type: application/json" \
-d @mock_payload.json \
-X POST \
"$APIGATEWAYURL/v1/currentConditions:lookup?key=$APIKEY"
