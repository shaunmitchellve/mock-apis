# Mock AQ API

This program will read in the JSON request, query the firebase DB for an exact match to the location lat/lng and return the JSON results.

The endpoint is written to mimic the `https://developers.google.com/maps/documentation/air-quality/reference/rest/v1/currentConditions/lookup` request and response specs. HOWEVER, not all fields will be used in the request and not all fields will be populated in the response.

Some of the response  fields have been commented out as they are not currently used / supplied in the Firestore DB but are here for completeness and any future requirements.