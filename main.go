package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Response struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PushedAt  time.Time `json:"pushed_at"`
}

func main() {

	// Open File of Apps
	f, err := os.Open("apps.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Create File of Apps
	csvFile, err := os.Create("./data.csv")

	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)

	for i, line := range data {

		if i > 0 { // omit header line
			project := "fury_" + line[0]
			fmt.Println("Processing project>", project)

			///////////////////////////////////////////
			url := "https://api.github.com/repos/mercadolibre/" + project

			// Create a Bearer string by appending string access token
			var bearer = "Bearer " + " Token "

			// Create a new request using http
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Println("Error on response.\n[ERROR] -", err)
			}
			// add authorization header to the req
			req.Header.Add("Authorization", bearer)

			// Send req using http Client
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Println("Error on response.\n[ERROR] -", err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("Error while reading the response bytes:", err)
			}
			//log.Println(string([]byte(body)))

			var r Response
			err = json.Unmarshal(body, &r)
			if err != nil {
				panic(err)
			}

			var row []string
			row = append(row, strconv.Itoa(r.ID))
			row = append(row, r.Name)
			row = append(row, string(r.CreatedAt.Format("2006-01-02 15:04:05")))
			row = append(row, string(r.UpdatedAt.Format("2006-01-02 15:04:05")))
			row = append(row, string(r.PushedAt.Format("2006-01-02 15:04:05")))
			writer.Write(row)
			///////////////////////////////////////////
			defer resp.Body.Close()
		}
	}
}
