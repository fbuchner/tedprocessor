package transform

// Target data model is a flat file model, no hierarchy (data table)
type DataRow struct {
	NoticeID                    string
	IssueDate                   string
	NoticeTypeCode              string
	NoticeLanguageCode          string
	SubTypeCode                 string
	NoticePublicationID         string
	ContractingActivityTypeCode string

	MainOrgGroupLeadIndicator                     string
	MainOrgAcquiringCPBIndicator                  string
	MainOrgAwardingCPBIndicator                   string
	MainOrgWebsiteURI                             string
	MainOrgPartyIdentificationId                  string
	MainOrgPartyName                              string
	MainOrgPostalAddressStreetName                string
	MainOrgPostalAddressCityName                  string
	MainOrgPostalAddressPostalZone                string
	MainOrgPostalAddressCountrySubentityCode      string
	MainOrgPostalAddressCountryIdentificationCode string
	MainOrgPartyLegalEntityCompanyID              string
	MainOrgContactName                            string
	MainOrgContactTelephone                       string
	MainOrgContactTelefax                         string
	MainOrgContactElectronicMail                  string

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
}
