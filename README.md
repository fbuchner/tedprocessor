# Tender Europe Daily (TED) Procurement Data

## Initial Situation and Objectives
Data on public procurement above the EU procurement threshold of approximately 221,000 EUR (for services, more in [Commission Delegated Regulation (EU) 2023/2495](https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX:32023R2495)) are published on the [Tenders Europe Daily (TED) portal](https://ted.europa.eu). The data is available license-free on the [EU OpenData Portal](https://data.europa.eu/data/datasets/ted-csv?locale=en).
The TEDprocessor establishes a processing pipeline to prepare the procurement data following an ETL (Extract-Transform-Load) schema. For this purpose, a custom data model is used, which simplifies the scope of attributes and relations.

## Data Model
The data model for TED data is complex and is described in the ["Annex to the second amendment of the eForms Implementing Regulation - 2023/288"](https://ec.europa.eu/docsroom/documents/58074). It contains more than 300 fields in total. In addition, there are various relations between the individual classes or objects. The full data model is documented as an [open-source ontology](https://docs.ted.europa.eu/EPO/latest/_attachments/html_reports/ePO/index.html). Part of the processing by TEDprocessor is a simplification of this data model. For this purpose, a significantly reduced and heavily opinionated data model was developed.

### General Considerations
**Reduction of Attributes**: For the analysis of procurement data, only part of the attributes are necessary. In the present data model, a selection based on experience was made, which can be adapted in the program code itself according to one's own needs.

**Denormalization**: To enable simple evaluation, the dataset was also denormalized, i.e., brought into the first normal form (1NF). For this purpose, the structure of the lots was dissolved in particular. Each lot is thus treated as a separate procurement, whose core data are identical across lots. The eForms specification  describes organizations separately from their actual uses in the dataset (i.e., an organization is described centrally and can then appear as a purchaser, place of performance, etc.). This was also resolved as part of the denormalization.

### Procurement Procedure Object
The procurement procedure object describes an ongoing procurement procedure, i.e., the call for competition.

| Attribute | Type | Description / Comment | Code Path |
| --- | --- | --- | --- |
| TED.Token | Text |  |  |
| ProcuringEntity.Handle | Text |  |  |
| ProcuringEntity.LegalForm |  |  |  |
| ProcuringEntity.Activity |  |  |  |
| Procedure.Country |  |  |  |
| Procedure.Title |  |  |  |
| Procedure.Description |  |  |  |
| Procedure.Type |  |  |  |
| Procedure.ContractType |  |  |  |
| Procedure.ProcuringEntity |  |  |  |
| Procedure.CPV |  |  |  |
| Procedure.AllCPV |  |  |  |
| Lot.TotalNumber |  |  |  |
| Lot.Number |  |  |  |
| Lot.Title |  |  |  |
| Lot.Description |  |  |  |
| Lot.EligibilityCriteria |  |  |  |
| Lot.AwardCriteria |  |  |  |
| Lot.StartDate |  |  |  |
| Lot.EndDate |  |  |  |
| Lot.PostalCode |  |  |  |

### Award Result Object
The award result object describes a completed procurement procedure, i.e., a competition awarded to one or no economic operator. As a rule, this can be linked to a previously conducted procurement procedure (exceptions include, in particular, direct awards without a prior call for competition).

The award results objects are currently still out of scope and might be added in a future release.

## Processing pipeline
The TED procurement data is processed in an ETL step logic.

1) **Extract**
Retrieve TED notices for the selected timeframe and store the unpacked XML files (one file per procurement), then convert the XML files into JSON files for easier processing.
2) **Transform**
The most complex step, extracting the data to transform the object to fit the target data model.
3) **Load** 
Store the data in the desired target format.

As of now there are four different folders, each script moving data from one folder to another.
**01-bulkddownloads**: the downloads from the TED openData archive
**02-XMLdata**: the unpacked archives, each XML file representing one procurement object
**03-JSONdata**: the raw JSON files, directly converted from the XML files
**04-Transformed**: procurement objects converted to fit the target data model 

## Implementation
The program is implemented using the [Go (Golang)](https://go.dev/) programming language. This allows compilation into a single binary that can then be executed via cron jobs or - if preferred - containerized and managed by the preferred platform.

### See Also
[OP-TED](https://github.com/OP-TED) - TED and EU Github Repository of the Public Procurement Unit of the Publications Office of the European Union.
[TED Developer Docs](https://docs.ted.europa.eu/home/index.html) - Developer portal with documentation of the eForms SDK, TED API, ontology, and other resources.
[Tendex](https://github.com/jesperkonincks/Tendex) - Python script for downloading TED data by calendar years.