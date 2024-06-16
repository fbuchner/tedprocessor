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
// TODO add more logging
func ProcessJSON(jsonFilepath, targetFilepath, csvSeparator string, deleteAfterProcessing bool) error {
	log.Debug().Str("JSON file", jsonFilepath).Msg("Processing JSON file")

	procurementProcedure, err := ReadJSON(jsonFilepath)
	if err != nil {
		return fmt.Errorf("failed to convert json to object: %v", err)
	}

	var dataRow DataRow

	//TODO check if array empty in UBLExtensions

	// Core attributes of the contract notice
	dataRow.NoticeID = procurementProcedure.NoticeID
	dataRow.IssueDate = procurementProcedure.IssueDate
	dataRow.NoticeTypeCode = procurementProcedure.NoticeTypeCode
	dataRow.NoticeLanguageCode = procurementProcedure.NoticeLanguageCode
	dataRow.SubTypeCode = procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.NoticeSubType.SubTypeCode
	dataRow.NoticePublicationID = procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Publication.NoticePublicationID
	dataRow.ContractingActivityTypeCode = procurementProcedure.ContractingParty.ContractingActivity.ActivityTypeCode

	//Match the buyer main organization to the separately found organization data
	var buyerOrgID string
	for _, org := range procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Organizations.Organization {
		buyerOrgID = procurementProcedure.ContractingParty.Party.PartyIdentification.ID

		if buyerOrgID != org.Company.PartyIdentification.ID {
			continue
		}

		//dataRow.MainOrgGroupLeadIndicator = org.GroupLeadIndicator
		//dataRow.MainOrgAcquiringCPBIndicator = org.AcquiringCPBIndicator
		//dataRow.MainOrgAwardingCPBIndicator = org.AwardingCPBIndicator
		dataRow.MainOrgPartyName = org.Company.PartyName.Name
		dataRow.MainOrgWebsiteURI = org.Company.WebsiteURI
		//dataRow.MainOrgPartyIdentificationId = org.Company.PartyIdentification.ID
		dataRow.MainOrgPostalAddressStreetName = org.Company.PostalAddress.StreetName
		dataRow.MainOrgPostalAddressCityName = org.Company.PostalAddress.CityName
		dataRow.MainOrgPostalAddressPostalZone = org.Company.PostalAddress.PostalZone
		//dataRow.MainOrgPostalAddressCountrySubentityCode = org.Company.PostalAddress.CountrySubentityCode
		dataRow.MainOrgPostalAddressCountryIdentificationCode = org.Company.PostalAddress.Country.IdentificationCode
		//dataRow.MainOrgPartyLegalEntityCompanyID = org.Company.PartyIdentification.ID
		dataRow.MainOrgContactName = org.Company.Contact.Name
		dataRow.MainOrgContactTelephone = org.Company.Contact.Telephone
		dataRow.MainOrgContactTelefax = org.Company.Contact.Telefax
		dataRow.MainOrgContactElectronicMail = org.Company.Contact.ElectronicMail

		// assignment done, no further loop iterations required
		break
	}

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

		/*  here it gets a bit tricky - we are in a loop through the lots to create DataRows and want to assign the winning organization
		1) first we check if there is a winning organization (and only then continue)
		2) then we get the correct tendering party from the LotTenders
		3) now we can get the organization ID from the according tendering party
		4) so we look through the organizations to find the matching one
		*/
		if len(procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.NoticeResult.LotTenders) > 0 {
			var targetTenderingPartyID string
			var winningOrgID string
			//get tendering party
			for _, lotTender := range procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.NoticeResult.LotTenders {
				if lotTender.TenderLot.ID == lot.ID {
					targetTenderingPartyID = lotTender.TenderingPartyID.ID
					break
				}
			}

			//get winning org
			for _, tenderingParty := range procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.NoticeResult.TenderingParties {
				if targetTenderingPartyID == tenderingParty.ID {
					winningOrgID = tenderingParty.Tenderer.ID
				}
			}

			if winningOrgID != "" {
				for _, org := range procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Organizations.Organization {
					if winningOrgID != org.Company.PartyIdentification.ID {
						continue
					}

					dataRow.WinningOrgWebsiteURI = org.Company.WebsiteURI
					dataRow.WinningOrgPartyIdentificationId = org.Company.PartyIdentification.ID
					dataRow.WinningOrgPartyName = org.Company.PartyName.Name
					dataRow.WinningOrgPostalAddressStreetName = org.Company.PostalAddress.StreetName
					dataRow.WinningOrgPostalAddressCityName = org.Company.PostalAddress.CityName
					dataRow.WinningOrgPostalAddressPostalZone = org.Company.PostalAddress.PostalZone
					//dataRow.WinningOrgPostalAddressCountrySubentityCode = org.Company.PostalAddress.CountrySubentityCode
					dataRow.WinningOrgPostalAddressCountryIdentificationCode = org.Company.PostalAddress.Country.IdentificationCode
					dataRow.WinningOrgPartyLegalEntityCompanyID = org.Company.PartyIdentification.ID
					dataRow.WinningOrgContactName = org.Company.Contact.Name
					dataRow.WinningOrgContactTelephone = org.Company.Contact.Telephone
					dataRow.WinningOrgContactTelefax = org.Company.Contact.Telefax
					dataRow.WinningOrgContactElectronicMail = org.Company.Contact.ElectronicMail

					break
				}
			}

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

func RemoveFile(jsonFilepath string) error {
	// Check if the file exists before trying to delete it
	if _, err := os.Stat(jsonFilepath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %v", jsonFilepath)
	}

	// Attempt to remove the file
	err := os.Remove(jsonFilepath)
	if err != nil {
		// Double-check if the file indeed doesn't exist to handle race conditions
		if os.IsNotExist(err) {
			return nil // Ignore the error if the file is already deleted
		}
		return fmt.Errorf("failed to remove json file: %v", err)
	}

	return nil
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
