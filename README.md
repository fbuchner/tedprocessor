# Tender Europe Daily (TED) Procurement Data

## Initial Situation and Objectives
Data on public procurement above the EU procurement threshold of approximately 221,000 EUR (for services, more in [Commission Delegated Regulation (EU) 2023/2495](https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX:32023R2495)) are published on the [Tenders Europe Daily (TED) portal](https://ted.europa.eu). The data is available license-free on the [EU OpenData Portal](https://data.europa.eu/data/datasets/ted-csv?locale=en).
The TEDprocessor establishes a processing pipeline to prepare the procurement data following an ETL (Extract-Transform-Load) schema. For this purpose, a custom data model is used, which simplifies the scope of attributes and relations.

## Data Model
The data model for TED data is complex and is described in the ["Annex to the second amendment of the eForms Implementing Regulation - 2023/288"](https://ec.europa.eu/docsroom/documents/58074). It contains more than 300 fields in total. In addition, there are various relations between the individual classes or objects. The full data model is documented as an [open-source ontology](https://docs.ted.europa.eu/EPO/latest/_attachments/html_reports/ePO/index.html). Part of the processing by TEDprocessor is a simplification of this data model. For this purpose, a significantly reduced and heavily opinionated non-hierarchical data model was developed.

### General Considerations
**Reduction of Attributes**: For the analysis of procurement data, only part of the attributes are necessary. In the present data model, a selection based on experience was made, which can be adapted in the program code itself according to one's own needs.

**Denormalization**: To enable simple evaluation, the dataset was also denormalized, i.e., brought into the first normal form (1NF) so it can easily be imported into your preferred analysis tool. For this purpose, the structure of the lots was dissolved in particular. Each lot is thus treated as a separate procurement, whose core data are identical across lots. The eForms specification describes organizations separately from their actual uses in the dataset (i.e., an organization is described centrally and can then appear as a purchaser, place of performance, etc.). This was also resolved as part of the denormalization.

**State of program code**: This program was developed out of professional curiosity. Therefore it might contain bugs (both in program code or business logic) or lack features that you might wish for. Please feel free to contribute to the code base and open a pull request or file an issue. 

### Procurement Procedure Object
The procurement procedure object describes an ongoing procurement procedure, i.e., the call for competition.

| Attribute                               | Description |
|-----------------------------------------|-------------|
| NoticeID                                | Hexadecimal TED ID. The ID is unique for each procurement but can occur several times due to different lots of the same procurement.        |
| IssueDate                               | Timestamp of the publication            |
| NoticeTypeCode                          | Indicates the type of notice, for example a standard publication or an update because of changes            |
| NoticeLanguageCode                      | Language of the publication, German is abbreviated as "DEU"            |
| SubTypeCode                             |             |
| NoticePublicationID                     | Running tally of procurements in the style of an incrementing number dash year            |
| ContractingActivityTypeCode             | Industry (like general public sector, health, educationm ... )            |
| MainOrgGroupLeadIndicator               |             |
| MainOrgAcquiringCPBIndicator            |             |
| MainOrgAwardingCPBIndicator             |             |
| MainOrgWebsiteURI                       | Main organization website         |
| MainOrgPartyIdentificationId            | Main organization internal id for the publication (like ORG-001)            |
| MainOrgPartyName                        | Main organization name, so the buying party            |
| MainOrgPostalAddressStreetName          |            |
| MainOrgPostalAddressCityName            |             |
| MainOrgPostalAddressPostalZone          |             |
| MainOrgPostalAddressCountrySubentityCode|             |
| MainOrgPostalAddressCountryIdentificationCode |             |
| MainOrgPartyLegalEntityCompanyID        |             |
| MainOrgContactName                      | Might contain personal information - be aware (GDPR)           |
| MainOrgContactTelephone                 |             |
| MainOrgContactTelefax                   |             |
| MainOrgContactElectronicMail            | Might contain personal information - be aware (GDPR)             |
| TenderingProcessProcedureCode           |             |
| ProcurementProjectID                    | Buyer internal procurement ID            |
| ProcurementProjectName                  | Short name for the procurement            |
| ProcurementProjectDescription           | Procurement description. Main field to search for keywords            |
| ProcurementProjectProcurementTypeCode   | Indicated the type of procurement, e.g. services or supplies            |
| ProcurementProjectNote                  | Additional notes about the procurement for the bidders           |
| ProcurementProjectMainCommodityClassification | Main CPV code used to classify the type of procuremnt            |
| LotID                                   | The ID of the Lot, usually starting at Lot-0001. When there is only one lot the lot data is usually equivalent to the Procurement Project data           |
| LotName                                 | See above at Procurement Project             |
| LotDescription                          | See above at Procurement Project             |
| LotProcurementTypeCode                  | See above at Procurement Project             |
| LotProjectNote                          | See above at Procurement Project             |
| LotMainCommodityClassification          | See above at Procurement Project             |
| DurationMeasure                         | Contract runtime (e.g. 48 MONTH)             |
| RealizedLocationAddressStreetName       |             |
| RealizedLocationAddressCityName         |             |
| RealizedLocationAddressPostalZone       |             |
| RealizedLocationAddressCountrySubentityCode |             |
| RealizedLocationAddressCountryIdentificationCode |             |


## Processing pipeline
The TED procurement data is processed in a three-step logic, roughly oriented at an ETL logic.

1) **First step: Download data (as XML)**
Data is downloaded on a per-month basis from the TED portal. The downloaded data is stored in a temporary file and immediately extracted into the download folder. There will be one folder for each month of the year (e.g. 05 for May indepdent of year) containing an archive file for each specific day of the year (TED does not publish notices every single day so there will be gaps). 
Within this step all the archives will also be unpacked into the xml directory. This directory will contain one folder per almost daily archive file. On any given day there will be several thousand notices published and therefore the same amount of xml files.
2) **Second step: Convert data (to JSON)**
Each XML file will be converted into a JSON file (and stored in the JSON directory). In this step not all data present in the XML file will be taken along. In case an XML file does not match the defined data model it will be discarded and this will be logged.
3) **Third step: Transform data (to CSV)**
Each of the previously created JSON files will be taken and the data transformed to fit the target data model. For each of the JSON file between one and n different data rows will be created (depending on the amount of lots in each procurement) and appended to a CSV file. This CSV file is the main output of tedprocessor and can then be imported into your database or data analytics solution of choice.

## Implementation
The program was implemented using the [Go (Golang)](https://go.dev/) programming language. This allows compilation into a single binary that can then be executed via cron jobs or - if preferred - containerized and managed by the preferred platform. 

### Preparing the program
First build the program from source or [download one of the releases](https://github.com/fbuchner/tedprocessor/releases). Make sure to chose the right executable for your platform.  

Then create a config file (suggested name: config.json), you can [find a template here](https://github.com/fbuchner/tedprocessor/blob/main/src/config.json). The configuration allows you to customize the following parameters.

All folders and files will be created as needed automatically by tedprocessor.

| Config item   | Description | Sample value |
|---------------|-------------|--------------|
| bulkdata_url | URL to download the monthly collection of procurement notices from | https://ted.europa.eu/packages/monthly
| download_dir  | Path to the folder where the procurement notices archives should be stored at  | /home/user/ted/downloads  |
| xml_dir  | Path to the folder where the XML files from the downloaded archives should be unpacked into  | /home/user/ted/xml  |
| json_dir  | Path to the folder where the to JSON converted procurement notices should be stored at  |  /home/user/ted/json |
| extracted_file  | Path to the CSV-file that will contain all procurement notice data. The file will contain headers | /home/user/ted/output.csv |
| filter_for_country | Discards all procurement notices that do not fit the given language code. Both the place of buyer and place of performance will be taken into account. | DEU  |
| run_steps { run_step_download } | Sets whether the download step should be run | true |
| run_steps { run_step_processxml } | Sets whether the process xml files to JSON step should be run  | true  |
| run_steps { run_step_transform } | Sets whether the transform to target data model and save as CSV step should be run  |  false |
| delete_after_processing | tells whether to delete files that have been successfully processed  | false |
| download_period: { from_year } | Tells the start year from which to download the notices. The start year is included. | 2024  |
| download_period: { from_month } | Tells the start month from which to download the notices. The start month is included. | 5  |
| download_period: { to_year } | Tells the end year until which to download the notices. The end year is included. | 2024  |
| download_period: { to_month } | Tells the end month  until which to download the notices. The end month is included. | 5  |
| log_level | Sets the log level of tedprocessor. Supported levels are debug, info, warn, error, fatal, panic. For production the level info is recommended. | info |
| csv_separator_char | Separator character for the CSV output file, for example ";" or "\t". In case more than one character is provided only the first one will be taken. |  ; |

### Running the program

**Linux / Mac**
1) open a shell/terminal and navigate to the executable
2) make tedprocessor executable: `chmod +x tedprocessor`
3) run tedprocessor: `./tedprocessor -config=config.json`

**Windows**
1) open the command line (Windows + R, cmd)
2) run tedprocessor:  `tedprocessor.exe -config=config.json`

The -config parameter can be ommitted on all platforms. In this case tedprocessor looks for a "config.json" in the same directory.

### See Also
[OP-TED](https://github.com/OP-TED) - TED and EU Github Repository of the Public Procurement Unit of the Publications Office of the European Union.
[TED Developer Docs](https://docs.ted.europa.eu/home/index.html) - Developer portal with documentation of the eForms SDK, TED API, ontology, and other resources.
[Tendex](https://github.com/jesperkonincks/Tendex) - Python script for downloading TED data by calendar years.