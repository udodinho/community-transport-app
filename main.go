package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type TrafficResponse struct {
	Rows []struct {
		Elements []struct {
			Distance struct {
				Text  string
				Value int
			}
			Duration struct {
				Text  string
				Value int
			}
			Duration_In_Traffic struct {
				Text  string
				Value int
			}
		}
	}
}

func main() {
	// Please input your longitude and latitude coordinates after typing go run main.go
	if len(os.Args) != 5 {
		fmt.Println("Usage: go run main.go <latitudeOrigin> <longitudeOrigin> <destinationLatitude> <destinationLongitude>")
		os.Exit(1)
	}

	latitudeOrigin := os.Args[1]
	longitudeOrigin := os.Args[2]
	destinationLatitude := os.Args[3]
	destinationLongitude := os.Args[4]

	trafficCondition, err := fetchTrafficInformation(latitudeOrigin, longitudeOrigin, destinationLatitude, destinationLongitude)
	if err != nil {
		exit(fmt.Sprintf("error getting traffic information: %v\n", err.Error()))
	}

	fmt.Printf("Traffic condition from %s, %s to %s %s is: %s\n", latitudeOrigin, longitudeOrigin, destinationLatitude, destinationLongitude, trafficCondition)

}

// fetchTrafficInformation retrieves real-time traffic information for a specific location.
func fetchTrafficInformation(latitudeOrigin, longitudeOrigin, destinationLatitude, destinationLongitude string) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		exit(fmt.Sprintf("unable to load env: %s", err.Error()))
	}

	// Retrieve the API KEY from environment varaible
	apiKey := os.Getenv("API_KEY")

	URL := fmt.Sprintf("https://maps.googleapis.com/maps/api/distancematrix/json?origins=%s,%s&destinations=%s,%s&departure_time=now&traffic_model=best_guess&key=%s", latitudeOrigin, longitudeOrigin, destinationLatitude, destinationLongitude, apiKey)

	client := &http.Client{}

	res, err := client.Get(URL)
	if err != nil {
		exit(fmt.Sprintf("error while making api call to the url: %s", err.Error()))
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		exit(fmt.Sprintf("unable to read body: %s", err.Error()))
	}

	var data TrafficResponse
	err = json.Unmarshal(body, &data)

	if err != nil {
		exit(fmt.Sprintf("error while unmarshaling body: %s", err.Error()))
	}

	if len(data.Rows) == 0 || len(data.Rows[0].Elements) == 0 {
		exit(fmt.Sprintf("no information found for this location: %s", err.Error()))
	}

	// Extract the traffic condition
	trafficCondition := data.Rows[0].Elements[0].Duration_In_Traffic.Text
	return trafficCondition, nil
}

// exit returns a message to the console and exit gracefully.
func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
