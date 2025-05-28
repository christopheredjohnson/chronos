package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

type ExportProject struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Elapsed string `json:"elapsed"`
}

func exportToCSV(projects []Project, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"ID", "Name", "Elapsed"})
	for _, p := range projects {
		writer.Write([]string{

			fmt.Sprintf("%d", p.ID),
			p.Name,
			formatDuration(p.Elapsed),
		})
	}
	return nil
}

func exportToJSON(projects []Project, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var exportedData []ExportProject

	for _, p := range projects {
		exportedData = append(exportedData, ExportProject{
			ID:      p.ID,
			Name:    p.Name,
			Elapsed: formatDuration(p.Elapsed),
		})
	}
	return json.NewEncoder(file).Encode(exportedData)
}
