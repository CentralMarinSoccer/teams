# Teams Service
Microservice that retrieves information from TeamSnap on a scheduled basis and exposes it as JSON.

## Endpoints
 * /teams
 * /teams/{id}
 
## Configuration
Service configuration is handled through the following environment variables:
* DIVISION - TeamSnap Division ID that contains teams (Required)
* TOKEN - TeamSnap authentication token (Required)
* GOOGLE_API_KEY - API key that has prermission to geolocate and address (Required)
* REFRESHINTERVAL - How often data should be refreshed from TeamSnap
* PORT - Port service listens on. Defaults to 8080
* DOMAIN - DNS name service is available at. Defaults to localhost
