package convert

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

func ProcessXML(xmlFilepath string, jsonFolderpath string, countryFilter string) error {
	log.Debug().Str("filepath", xmlFilepath).Msg("Processing XML file")

	// Open the XML file
	xmlFile, err := os.Open(xmlFilepath)
	if err != nil {
		return fmt.Errorf("failed to read xml file: %v", err)
	}
	defer xmlFile.Close()

	// Read the XML file
	byteValue, _ := io.ReadAll(xmlFile)
	// Unmarshal the XML into the struct
	var procurementProcedure ProcurementProcedure
	err = xml.Unmarshal(byteValue, &procurementProcedure)
	if err != nil {
		return fmt.Errorf("failed to extract xml data: %v", err)
	}

	// Make sure the ProcurementProcedure XML format is valid, there are some exotic formats mixed into TED
	if procurementProcedure.NoticeID == "" {
		log.Warn().Str("XML file", xmlFilepath).Msg("Mismatch in XML schema, skipping file")
	}

	// Print the struct (for debugging purposes)
	//fmt.Printf("%+v\n", procurementProcedure)

	// Skip processing if the country filter does not match the realized location country
	countryCode1 := procurementProcedure.ProcurementProject.RealizedLocation.Address.Country.IdentificationCode
	countryCode2 := ""
	if len(procurementProcedure.UBLExtensions.UBLExtension) > 0 && len(procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Organizations.Organization) > 0 {
		countryCode2 = procurementProcedure.UBLExtensions.UBLExtension[0].ExtensionContent.EformsExtension.Organizations.Organization[0].Company.PostalAddress.Country.IdentificationCode
	}
	if countryFilter != "" && countryCode1 != countryFilter && countryCode2 != countryFilter {
		log.Debug().Str("filepath", xmlFilepath).Str("Country code 1", countryCode1).Str("Country code 2", countryCode2).Msg("Skipping file due to country filter")
		return nil
	}

	// Set filename for JSON file
	filenameWithExt := filepath.Base(xmlFilepath)
	ext := filepath.Ext(filenameWithExt)
	filenameWithoutExt := strings.TrimSuffix(filenameWithExt, ext)
	targetPath := filepath.Join(jsonFolderpath, filenameWithoutExt+".json")
	// Write data out to JSON
	writeJSON(procurementProcedure, targetPath)
	log.Debug().Str("JSON path", targetPath).Msg("Writing JSON file")

	return nil
}

func writeJSON(procurementProcedure ProcurementProcedure, filePath string) error {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(procurementProcedure)
	if err != nil {
		return fmt.Errorf("failed to convert object to json: %v", err)
	}

	// Save the JSON data to a file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file for json data: %v", err)

	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed write json data to file: %v", err)
	}

	//fmt.Println("JSON data saved.")
	return nil
}
