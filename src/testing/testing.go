package testing

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type ContractNotice struct {
	IssueDate string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 IssueDate"`
}

func ReadXML() error {
	// Open the XML file
	xmlFile, err := os.Open("/Users/frederic/Downloads/testing.xml")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer xmlFile.Close()

	// Read the XML file
	byteValue, err := io.ReadAll(xmlFile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Initialize the Catalog struct
	var catalog ContractNotice

	// Unmarshal the XML data into the struct
	err = xml.Unmarshal(byteValue, &catalog)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Print the data

	fmt.Printf("ID: %s\n", catalog.IssueDate)

	return nil
}
