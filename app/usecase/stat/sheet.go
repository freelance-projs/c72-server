package stat

import (
	"context"
	"reflect"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type sheetService struct {
	srv *sheets.Service
}

func newSheetService() *sheetService {
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile("./account_credentials.json"))
	if err != nil {
		panic(err)
	}

	return &sheetService{
		srv: srv,
	}
}

func (s *sheetService) createDailySheet(spreadsheetID string, sheetName string, header []any) error {
	srv := s.srv
	// Retrieve the spreadsheet to check existing sheet names
	spreadsheet, err := srv.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		return err
	}

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
		return nil
	}

	// Create a new sheet request
	addSheetRequest := &sheets.Request{
		AddSheet: &sheets.AddSheetRequest{
			Properties: &sheets.SheetProperties{
				Title: sheetName,
			},
		},
	}

	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{addSheetRequest},
	}

	if _, err := srv.Spreadsheets.BatchUpdate(spreadsheetID, batchUpdateRequest).Do(); err != nil {
		return err
	}

	// Prepare the request to insert the header into the first row
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{header},
	}

	// Insert headers into the new sheet
	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, sheetName+"!A1", valueRange).
		ValueInputOption("RAW").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *sheetService) insert(spreadsheetID string, sheetName string, data []any) error {
	found, err := s.sheetExists(spreadsheetID, sheetName)
	if err != nil {
		return err
	}
	if !found {
		header := getStructHeader(data[0])
		if err := s.createDailySheet(spreadsheetID, sheetName, header); err != nil {
			return err
		}
	}

	values := make([][]interface{}, len(data))
	for i, row := range data {
		rv := reflect.ValueOf(row)
		for j := range rv.NumField() {
			values[i] = append(values[i], rv.Field(j).Interface())
		}
	}

	// Prepare the request to insert the data into the sheet
	valueRange := &sheets.ValueRange{
		Values: values,
	}

	// Insert data into the sheet
	if _, err := s.srv.Spreadsheets.Values.Clear(spreadsheetID, sheetName+"!A2:Z1000", &sheets.ClearValuesRequest{}).Do(); err != nil {
		return err
	}
	if _, err := s.srv.Spreadsheets.Values.Append(spreadsheetID, sheetName+"!A2", valueRange).
		ValueInputOption("RAW").InsertDataOption("INSERT_ROWS").Do(); err != nil {
		return err
	}

	return nil
}

func (s *sheetService) sheetExists(spreadsheetID string, sheetTitle string) (bool, error) {
	// Retrieve the spreadsheet to get the sheet details
	spreadsheet, err := s.srv.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		return false, err
	}

	// Loop through the sheets to find if the sheet with the given title exists
	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.Title == sheetTitle {
			return true, nil // Sheet exists
		}
	}

	// Sheet does not exist
	return false, nil
}

func getStructHeader(s any) []any {
	rt := reflect.TypeOf(s)
	if rt.Kind() != reflect.Struct {
		return nil
	}
	header := make([]any, rt.NumField())
	for i := range rt.NumField() {
		header[i] = rt.Field(i).Tag.Get("header")
	}
	return header
}
