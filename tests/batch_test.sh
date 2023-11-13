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
APIKEY='<API-KEY>'

batchLocations=(1.00 1.00 2.00 2.00)
start=$(date +%s)

for l in ${!batchLocations[@]}; do
    if [ $((l % 2)) != 0 ];
    then
        curl -H "Content-Type: application/json" \
        -d '{
                "location": {
                    "latitude": '"${batchLocations[($l-1)]}"',
                    "longitude": '"${batchLocations[$l]}"'
                },
                "extraComputations": [
                    "LOCAL_AQI"
                ],
                "customLocalAqis": [{
                    "regionCode": "fr",
                "aqi": "eaqi"
                }],
                "languageCode": "en"
            }' \
        -X POST \
        "$APIGATEWAYURL/v1/currentConditions:lookup?key=$APIKEY"
    fi
done
end=$(date +%s)
echo "Elapsed Time: $(($end-$start)) seconds"