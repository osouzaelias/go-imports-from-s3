package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
)

type Person struct {
	ID                 uuid.UUID
	FirstName          string
	LastName           string
	City               string
	StreetName         string
	Zip                string
	State              string
	Country            string
	DataExpirationDate string
}

func main() {
	file, err := os.Create("testdata.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	header := []string{"ID", "FirstName", "LastName", "City", "StreetName", "Zip", "State", "Country", "DataExpirationDate"}
	err = writer.Write(header)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 50000; i++ {
		p := Person{
			ID:                 uuid.New(),
			FirstName:          gofakeit.FirstName(),
			LastName:           gofakeit.LastName(),
			City:               gofakeit.City(),
			StreetName:         gofakeit.StreetName(),
			Zip:                gofakeit.Zip(),
			State:              gofakeit.State(),
			Country:            gofakeit.Country(),
			DataExpirationDate: addDaysToCurrentDate(1),
		}

		row := []string{
			p.ID.String(),
			p.FirstName,
			p.LastName,
			p.City,
			p.StreetName,
			p.Zip,
			p.State,
			p.Country,
			p.DataExpirationDate,
		}

		err = writer.Write(row)
		if err != nil {
			panic(err)
		}
	}

	writer.Flush()

	fmt.Println("Arquivo gerado com sucesso!")
}

func addDaysToCurrentDate(days int) string {
	t := time.Now()
	t = t.Add(time.Duration(days) * 24 * time.Hour)
	return strconv.FormatInt(t.Unix(), 10)
}
