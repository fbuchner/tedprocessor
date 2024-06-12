package transform

// Target data model is a flat file model, no hierarchy (data table)
// Uncommented lines are available but not used. Just activate if needed in here and in transform.go
type DataRow struct {
	// TODO add estimated value of procurement, also in XML extract
	NoticeID                    string
	IssueDate                   string
	NoticeTypeCode              string
	NoticeLanguageCode          string
	SubTypeCode                 string
	NoticePublicationID         string
	ContractingActivityTypeCode string

	//MainOrgGroupLeadIndicator                     string
	//MainOrgAcquiringCPBIndicator                  string
	//MainOrgAwardingCPBIndicator                   string
	MainOrgPartyName  string
	MainOrgWebsiteURI string
	//MainOrgPartyIdentificationId                  string
	MainOrgPostalAddressStreetName string
	MainOrgPostalAddressCityName   string
	MainOrgPostalAddressPostalZone string
	//MainOrgPostalAddressCountrySubentityCode      string
	MainOrgPostalAddressCountryIdentificationCode string
	//MainOrgPartyLegalEntityCompanyID              string
	MainOrgContactName           string
	MainOrgContactTelephone      string
	MainOrgContactTelefax        string
	MainOrgContactElectronicMail string

	TenderingProcessProcedureCode string

	ProcurementProjectID                          string
	ProcurementProjectName                        string
	ProcurementProjectDescription                 string
	ProcurementProjectProcurementTypeCode         string
	ProcurementProjectNote                        string
	ProcurementProjectMainCommodityClassification string

	LotID                          string
	LotName                        string
	LotDescription                 string
	LotProcurementTypeCode         string
	LotProjectNote                 string
	LotMainCommodityClassification string

	DurationMeasure string

	RealizedLocationAddressStreetName                string
	RealizedLocationAddressCityName                  string
	RealizedLocationAddressPostalZone                string
	RealizedLocationAddressCountrySubentityCode      string
	RealizedLocationAddressCountryIdentificationCode string

	WinningOrgWebsiteURI              string
	WinningOrgPartyIdentificationId   string
	WinningOrgPartyName               string
	WinningOrgPostalAddressStreetName string
	WinningOrgPostalAddressCityName   string
	WinningOrgPostalAddressPostalZone string
	//WinningOrgPostalAddressCountrySubentityCode      string
	WinningOrgPostalAddressCountryIdentificationCode string
	WinningOrgPartyLegalEntityCompanyID              string
	WinningOrgContactName                            string
	WinningOrgContactTelephone                       string
	WinningOrgContactTelefax                         string
	WinningOrgContactElectronicMail                  string
}
