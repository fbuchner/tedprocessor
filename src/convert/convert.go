package convert

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

// Structs to represent XML elements, ContractNotice being the encompassing root element
type ContractNotice struct {
	IssueDate             string                `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 IssueDate"`
	NoticeTypeCode        string                `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 NoticeTypeCode"`
	NoticeLanguageCode    string                `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 NoticeLanguageCode"`
	UBLExtensions         UBLExtensions         `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonExtensionComponents-2 UBLExtensions"`
	ContractingParty      ContractingParty      `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 ContractingParty"`
	TenderingProcess      TenderingProcess      `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 TenderingProcess"`
	ProcurementProject    ProcurementProject    `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 ProcurementProject"`
	ProcurementProjectLot ProcurementProjectLot `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 ProcurementProjectLot"`
}

type UBLExtensions struct {
	UBLExtension []UBLExtension `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonExtensionComponents-2 UBLExtension"`
}

type UBLExtension struct {
	ExtensionContent ExtensionContent `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonExtensionComponents-2 ExtensionContent"`
}

type ExtensionContent struct {
	EformsExtension EformsExtension `xml:"http://data.europa.eu/p27/eforms-ubl-extensions/1 EformsExtension"`
}

type EformsExtension struct {
	NoticeSubType NoticeSubType `xml:"http://data.europa.eu/p27/eforms-ubl-extension-aggregate-components/1 NoticeSubType"`
	Organizations Organizations `xml:"http://data.europa.eu/p27/eforms-ubl-extension-aggregate-components/1 Organizations"`
	Publication   Publication   `xml:"http://data.europa.eu/p27/eforms-ubl-extension-aggregate-components/1 Publication"`
}

type NoticeSubType struct {
	SubTypeCode string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 SubTypeCode"`
}

type Organizations struct {
	Organization []Organization `xml:"http://data.europa.eu/p27/eforms-ubl-extension-aggregate-components/1 Organization"`
}

type Organization struct {
	GroupLeadIndicator    string  `xml:"http://data.europa.eu/p27/eforms-ubl-extension-basic-components/1 GroupLeadIndicator"`
	AcquiringCPBIndicator string  `xml:"http://data.europa.eu/p27/eforms-ubl-extension-basic-components/1 AcquiringCPBIndicator"`
	AwardingCPBIndicator  string  `xml:"http://data.europa.eu/p27/eforms-ubl-extension-basic-components/1 AwardingCPBIndicator"`
	Company               Company `xml:"http://data.europa.eu/p27/eforms-ubl-extension-aggregate-components/1 Company"`
}

type Company struct {
	WebsiteURI          string              `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 WebsiteURI"`
	PartyIdentification PartyIdentification `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 PartyIdentification"`
	PartyName           PartyName           `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 PartyName"`
	PostalAddress       PostalAddress       `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 PostalAddress"`
	PartyLegalEntity    PartyLegalEntity    `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 PartyLegalEntity"`
	Contact             Contact             `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 Contact"`
}

type PartyIdentification struct {
	ID string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ID"`
}

type PartyName struct {
	Name string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Name"`
}

type PostalAddress struct {
	StreetName           string  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 StreetName"`
	CityName             string  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 CityName"`
	PostalZone           string  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 PostalZone"`
	CountrySubentityCode string  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 CountrySubentityCode"`
	Country              Country `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 Country"`
}

type Country struct {
	IdentificationCode string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 IdentificationCode"`
}

type PartyLegalEntity struct {
	CompanyID string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 CompanyID"`
}

type Contact struct {
	Name           string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Name"`
	Telephone      string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Telephone"`
	Telefax        string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Telefax"`
	ElectronicMail string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ElectronicMail"`
}

type Publication struct {
	NoticePublicationID string `xml:"http://data.europa.eu/p27/eforms-ubl-extension-basic-components/1 NoticePublicationID"`
}

type ContractingParty struct {
	BuyerProfileURI      string               `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 BuyerProfileURI"`
	ContractingPartyType ContractingPartyType `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 ContractingPartyType"`
	ContractingActivity  ContractingActivity  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 ContractingActivity"`
	Party                Party                `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 Party"`
}

type ContractingPartyType struct {
	PartyTypeCode string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 PartyTypeCode"`
}

type ContractingActivity struct {
	ActivityTypeCode string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ActivityTypeCode"`
}

type Party struct {
	PartyIdentification PartyIdentification `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 PartyIdentification"`
}

type TenderingProcess struct {
	ProcedureCode string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ProcedureCode"`
}

type ProcurementProject struct {
	ID                          string                      `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ID"`
	Name                        string                      `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Name"`
	Description                 string                      `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Description"`
	ProcurementTypeCode         string                      `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ProcurementTypeCode"`
	Note                        string                      `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Note"`
	MainCommodityClassification MainCommodityClassification `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 MainCommodityClassification"`
	RealizedLocation            RealizedLocation            `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 RealizedLocation"`
}

type MainCommodityClassification struct {
	ItemClassificationCode string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ItemClassificationCode"`
}

type RealizedLocation struct {
	Address Address `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 Address"`
}

type Address struct {
	StreetName           string  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 StreetName"`
	CityName             string  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 CityName"`
	PostalZone           string  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 PostalZone"`
	CountrySubentityCode string  `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 CountrySubentityCode"`
	Country              Country `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 Country"`
}

type ProcurementProjectLot struct {
	ID string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ID"`
}

func ReadXML(countryFilter string) error {
	// Open the XML file
	xmlFile, err := os.Open("/Users/frederic/Downloads/00262102_2024.xml")
	if err != nil {
		return fmt.Errorf("failed to read xml file: %v", err)
	}
	defer xmlFile.Close()

	// Read the XML file
	byteValue, _ := io.ReadAll(xmlFile)
	// Unmarshal the XML into the struct
	var contractNotice ContractNotice
	err = xml.Unmarshal(byteValue, &contractNotice)
	if err != nil {
		return fmt.Errorf("failed to extract xml data: %v", err)
	}

	// Print the struct (for debugging purposes)
	//fmt.Printf("%+v\n", contractNotice)

	// Skip processing if the country filter does not match the realized location country
	if contractNotice.ProcurementProject.RealizedLocation.Address.Country.IdentificationCode != countryFilter {
		return nil
	}

	//write data out to JSON
	writeJSON(contractNotice, "/Users/frederic/Downloads/contractnotice.json")

	return nil
}

func writeJSON(contractNotice ContractNotice, filePath string) error {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(contractNotice)
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
