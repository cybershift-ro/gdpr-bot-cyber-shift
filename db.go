package main

import (
	"encoding/json"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
)

func printSanction(sanction *models.Record) {
	if sanction == nil {
		fmt.Println("Invalid record")
		return
	}

	fmt.Printf("Company %s has been sanctioned with %d RON. Verified: %t\nLink: %s\n", sanction.GetStringDataValue("company_name"), sanction.GetIntDataValue("fine_amount"), sanction.GetBoolDataValue("human_verified"), sanction.GetStringDataValue("url"))
}

func saveSanction(c *models.Collection, companyName, article string, sanctionSum float64) bool {
	if c == nil || app.Dao() == nil {
		fmt.Println("Invalid db params")
		return false
	}

	newEntry := models.NewRecord(c)

	newEntry.SetDataValue("company_name", companyName)
	newEntry.SetDataValue("human_verified", false)
	newEntry.SetDataValue("url", article)
	newEntry.SetDataValue("fine_amount", sanctionSum)
	newEntry.SetDataValue("currency", "RON")
	newEntry.SetDataValue("irelevant", false)

	if err := app.Dao().SaveRecord(newEntry); err != nil {
		fmt.Println("Failed to create sanction record to db: ", err)
		return false
	}

	return true
}

func updateSanction(entry *models.Record, c *models.Collection, companyName, article string, sanctionSum float64) bool {
	if c == nil || app.Dao() == nil || entry == nil {
		fmt.Println("Invalid db params")
		return false
	}

	entry.SetDataValue("company_name", companyName)
	entry.SetDataValue("human_verified", false)
	entry.SetDataValue("url", article)
	entry.SetDataValue("fine_amount", sanctionSum)

	if err := app.Dao().SaveRecord(entry); err != nil {
		fmt.Println("Failed to update sanction record to db: ", err)
		return false
	}

	return true
}

func sanctionsToJSON() string {
	sanction_document, err := app.Dao().FindCollectionByNameOrId("sanctions")

	if err != nil {
		fmt.Println("Error while building sanction feed: ", err)
		return ""
	}

	records, _ := app.Dao().FindRecordsByExpr(sanction_document, dbx.Not(dbx.HashExp{"company_name": "unknown"}))

	all_sanctions := make([]map[string]any, len(records))

	for _, s := range records {
		all_sanctions = append(all_sanctions, s.Data())
	}

	json_data, err := json.MarshalIndent(all_sanctions, "", "    ")

	if err != nil {
		fmt.Println("Error while converting sanction feed to JSON: ", err)
		return ""
	}

	return string(json_data)
}
