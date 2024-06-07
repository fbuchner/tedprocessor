package transform

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"tedprocessor/convert"

	"github.com/rs/zerolog/log"
)

// Processes a JSON file, appends result to target file
func ProcessJSON(jsonFilepath, targetFilepath string, csvSeparator string) error {
	log.Debug().Str("JSON file", jsonFilepath).Msg("Processing JSON file")

	procurementProcedure, err := ReadJSON(jsonFilepath)
	if err != nil {
		return fmt.Errorf("failed to convert json to object: %v", err)
	}

	var dataRow DataRow

	// Core attributes of the contract notice
	dataRow.NoticeID = procurementProcedure.NoticeID
	dataRow.IssueDate = procurementProcedure.IssueDate
	dataRow.NoticeTypeCode = procurementProcedure.NoticeTypeCode
	dataRow.NoticeLanguageCode = procurementProcedure.NoticeLanguageCode
	dataRow.SubTypeCode = procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.NoticeSubType.SubTypeCode
	dataRow.NoticePublicationID = procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Publication.NoticePublicationID
	dataRow.ContractingActivityTypeCode = procurementProcedure.ContractingParty.ContractingActivity.ActivityTypeCode

	//TODO check if array empty in UBLExtensions
	dataRow.MainOrgGroupLeadIndicator = procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Organizations.Organization[0].GroupLeadIndicator
	dataRow.MainOrgAcquiringCPBIndicator = procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Organizations.Organization[0].AcquiringCPBIndicator
	dataRow.MainOrgAwardingCPBIndicator = procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Organizations.Organization[0].AwardingCPBIndicator
	dataRow.MainOrgWebsiteURI = procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Organizations.Organization[0].Company.WebsiteURI
	dataRow.MainOrgPartyIdentificationId = procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Organizations.Organization[0].Company.PartyIdentification.ID
	dataRow.MainOrgPartyName = procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Organizations.Organization[0].Company.PartyName.Name
	dataRow.MainOrgPostalAddressStreetName = procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Organizations.Organization[0].Company.PostalAddress.StreetName
	dataRow.MainOrgPostalAddressCityName = procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Organizations.Organization[0].Company.PostalAddress.CityName
	dataRow.MainOrgPostalAddressPostalZone = procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Organizations.Organization[0].Company.PostalAddress.PostalZone
	dataRow.MainOrgPostalAddressCountrySubentityCode = procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Organizations.Organization[0].Company.PostalAddress.CountrySubentityCode
	dataRow.MainOrgPostalAddressCountryIdentificationCode = procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Organizations.Organization[0].Company.PostalAddress.Country.IdentificationCode
	dataRow.MainOrgPartyLegalEntityCompanyID = procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Organizations.Organization[0].Company.PartyIdentification.ID
	dataRow.MainOrgContactName = procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Organizations.Organization[0].Company.Contact.Name
	dataRow.MainOrgContactTelephone = procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Organizations.Organization[0].Company.Contact.Telephone
	dataRow.MainOrgContactTelefax = procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Organizations.Organization[0].Company.Contact.Telefax
	dataRow.MainOrgContactElectronicMail = procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Organizations.Organization[0].Company.Contact.ElectronicMail

	dataRow.TenderingProcessProcedureCode = procurementProcedure.TenderingProcess.ProcedureCode

	// Set default procurement project data
	dataRow.ProcurementProjectID = procurementProcedure.ProcurementProject.ID
	dataRow.ProcurementProjectName = procurementProcedure.ProcurementProject.Name
	dataRow.ProcurementProjectDescription = procurementProcedure.ProcurementProject.Description
	dataRow.ProcurementProjectProcurementTypeCode = procurementProcedure.ProcurementProject.ProcurementTypeCode
	dataRow.ProcurementProjectNote = procurementProcedure.ProcurementProject.Note
	dataRow.ProcurementProjectMainCommodityClassification = procurementProcedure.ProcurementProject.MainCommodityClassification.ItemClassificationCode

	// Now flatten the lot structure by creating several copies of the structure if necessary
	for _, lot := range procurementProcedure.ProcurementProjectLot {
		dataRow.LotID = lot.ID
		dataRow.LotName = lot.ProcurementProject.Name
		dataRow.LotDescription = lot.ProcurementProject.Description
		dataRow.LotProcurementTypeCode = lot.ProcurementProject.ProcurementTypeCode
		dataRow.LotProjectNote = lot.ProcurementProject.Note
		dataRow.LotMainCommodityClassification = lot.ProcurementProject.MainCommodityClassification.ItemClassificationCode

		// TODO there is also a RealizedLocation based on the Main Procurement Project. One could compare the info and take the Main Procurement one as default and overwrite with lot info as long as that is not empty
		dataRow.RealizedLocationAddressStreetName = lot.ProcurementProject.RealizedLocation.Address.StreetName
		dataRow.RealizedLocationAddressCityName = lot.ProcurementProject.RealizedLocation.Address.CityName
		dataRow.RealizedLocationAddressPostalZone = lot.ProcurementProject.RealizedLocation.Address.PostalZone
		dataRow.RealizedLocationAddressCountrySubentityCode = lot.ProcurementProject.RealizedLocation.Address.CountrySubentityCode
		dataRow.RealizedLocationAddressCountryIdentificationCode = lot.ProcurementProject.RealizedLocation.Address.Country.IdentificationCode

		// Use conditional assignment to not assign only a space as value if empty
		if lot.ProcurementProject.PlannedPeriod.DurationMeasure.Value != "" {
			dataRow.DurationMeasure = lot.ProcurementProject.PlannedPeriod.DurationMeasure.Value + " " + lot.ProcurementProject.PlannedPeriod.DurationMeasure.UnitCode
		}

		// Add the dataRow item as a line in the CSV file
		err = appendStructToCSV(targetFilepath, dataRow, csvSeparator)
		if err != nil {
			return fmt.Errorf("failed to append data row: %v", err)
		}
	}

	return nil
}

func ReadJSON(filePath string) (convert.ProcurementProcedure, error) {
	// Create an instance of ProcurementProcedure from contracting package
	var procurementProcedure convert.ProcurementProcedure

	// Open the JSON file
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return procurementProcedure, fmt.Errorf("failed to open json file: %v", err)
	}
	defer jsonFile.Close()

	// Unmarshal the JSON data into the struct
	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&procurementProcedure)
	if err != nil {
		return procurementProcedure, fmt.Errorf("failed to decode json file: %v", err)
	}

	return procurementProcedure, nil
}

// Generic CSV writer that appends a struct to a CSV file
func appendStructToCSV(filePath string, data interface{}, csvSeparator string) error {
	if len(csvSeparator) == 0 {
		csvSeparator = ";"
	}
	var csvSeparatorRune rune = rune(csvSeparator[0])

	var writeHeaders bool

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		writeHeaders = true
	}

	// Open the file in append mode, create if it doesn't exist
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	writer.Comma = csvSeparatorRune
	defer writer.Flush()

	// Use reflection to get the struct's field names and values
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Struct {
		return &reflect.ValueError{Method: "appendStructToCSV", Kind: v.Kind()}
	}

	if writeHeaders {
		// Get the field names (headers) only if the file is new
		headers := []string{}
		for i := 0; i < v.NumField(); i++ {
			headers = append(headers, v.Type().Field(i).Name)
		}
		if err := writer.Write(headers); err != nil {
			return err
		}
	}

	// Get the field values
	values := []string{}
	for i := 0; i < v.NumField(); i++ {
		values = append(values, fmt.Sprintf("%v", v.Field(i).Interface()))
	}

	// Write the values as a row in the CSV file
	if err := writer.Write(values); err != nil {
		return err
	}

	return nil
}
