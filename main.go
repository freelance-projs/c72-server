package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type User struct {
	Name string
	Age  int
}

type TransactionLog struct {
	Entity    string `json:"entity"`
	TagName   string `json:"tag_name"`
	Action    string `json:"action"`
	Count     int    `json:"count"`
	CreatedAt string `json:"created_at"`
}

func main() {
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile("./account_credentials.json"))
	if err != nil {
		panic(err)
	}

	spreadsheetId := "1XIMAojHp1g-SMt8aOY-IGaPB0hT-KnW1HsBWB4VcV64"
	createSheet(srv, spreadsheetId)
}

func createSheet(srv *sheets.Service, spreadsheetID string) {
	// Retrieve the spreadsheet to check existing sheet names
	spreadsheet, err := srv.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve spreadsheet: %v", err)
	}
	sheetName := time.Now().Format("2006-01-02") // Name of the new sheet

	// Check if the sheet exists
	sheetExists := false
	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.Title == sheetName {
			sheetExists = true
			break
		}
	}

	// If the sheet doesn't exist, create it
	if sheetExists {
		return
	}

	// Create a new sheet request
	// addSheetRequest := &sheets.Request{
	// 	AddSheet: &sheets.AddSheetRequest{
	// 		Properties: &sheets.SheetProperties{
	// 			Title: sheetName,
	// 		},
	// 	},
	// }

	// Prepare the batch update request
	// batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{
	// 	Requests: []*sheets.Request{addSheetRequest},
	// }

	// Call BatchUpdate to add the new sheet
	// _, err = srv.Spreadsheets.BatchUpdate(spreadsheetID, batchUpdateRequest).Do()
	// if err != nil {
	// 	log.Fatalf("Unable to add new sheet: %v", err)
	// }
	//
	// fmt.Printf("New sheet '%s' created successfully!\n", sheetName)

	// Now, insert headers into the first row of the new sheet
	// header := []interface{}{"Phòng ban", "Hành vi", "Loại đồ", "Số lượng", "Thời gian thực hiện"}
	//
	// // Prepare the request to insert the header into the first row
	// valueRange := &sheets.ValueRange{
	// 	Values: [][]interface{}{header},
	// }
	//
	// // Insert headers into the new sheet
	// _, err = srv.Spreadsheets.Values.Update(spreadsheetID, sheetName+"!A1", valueRange).
	// 	ValueInputOption("RAW").Do()
	// if err != nil {
	// 	log.Fatalf("Unable to insert header row: %v", err)
	// }

	// columnRange := sheetName + "!D2:D" // Assuming you want to apply the validation to Column D (from row 2 downwards)

	// Define the allowed values for the column (brown, return)
}

func insertSheet(srv *sheets.Service, id string) {
	readRange := "Sheet1!A:B" // Range where data will be inserted (starting from A2, B2)

	// List of users to insert
	users := []User{
		{"John", 30},
		{"Alice", 25},
		{"Bob", 40},
	}

	// Prepare the data in the format expected by the API
	var values [][]interface{}
	for _, user := range users {
		values = append(values, []interface{}{user.Name, user.Age})
	}

	// Create the request body
	valueRange := &sheets.ValueRange{
		Values: values,
	}

	// Append the data to the sheet
	_, err := srv.Spreadsheets.Values.Append(id, readRange, valueRange).
		ValueInputOption("RAW"). // Use "RAW" to insert data as is
		Do()

	if err != nil {
		log.Fatalf("Unable to append data: %v", err)
	}

	fmt.Println("Data inserted successfully")
}
