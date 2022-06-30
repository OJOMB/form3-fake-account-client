package accounts

import "time"

// Account represents an account in the form3 org section.
// See https://api-docs.form3.tech/api.html#organisation-accounts for
// more information about fields.
type AccountData struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	ID             string             `json:"id,omitempty"`
	OrganisationID string             `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Version        *int64             `json:"version,omitempty"`
	CreatedOn      *time.Time         `json:"created_on,omitempty"`
	ModifiedOn     *time.Time         `json:"modified_on,omitempty"`
}

type AccountAttributes struct {
	AcceptanceQualifier        string                      `json:"acceptance_qualifier,omitempty"`
	AccountClassification      *string                     `json:"account_classification,omitempty"`
	AccountNumber              string                      `json:"account_number,omitempty"`
	AlternativeNames           []string                    `json:"alternative_names,omitempty"`
	BankID                     string                      `json:"bank_id,omitempty"`
	BankIDCode                 string                      `json:"bank_id_code,omitempty"`
	BaseCurrency               string                      `json:"base_currency,omitempty"`
	Bic                        string                      `json:"bic,omitempty"`
	Country                    *string                     `json:"country,omitempty"`
	CustomerID                 string                      `json:"customer_id,omitempty"`
	Iban                       string                      `json:"iban,omitempty"`
	JointAccount               *bool                       `json:"joint_account,omitempty"`
	Name                       []string                    `json:"name,omitempty"`
	NameMatchingStatus         string                      `json:"name_matching_status,omitempty"`
	OrganisationIdentification *OrganisationIdentification `json:"organisation_identification,omitempty"`
	PrivateIdentification      *PrivateIdentification      `json:"private_identification,omitempty"`
	ReferenceMask              string                      `json:"reference_mask,omitempty"`
	Relationships              *Relationships              `json:"relationships,omitempty"`
	SecondaryIdentification    string                      `json:"secondary_identification,omitempty"`
	Status                     *AccountStatus              `json:"status,omitempty"`
	StatusReason               string                      `json:"status_reason,omitempty"`
	UserDefinedData            []UserDefinedData           `json:"user_defined_data,omitempty"`
	ValidationType             string                      `json:"validation_type,omitempty"`
	// Deprecated: AlternativeBankAccountNames is deprecated
	AlternativeBankAccountNames []string `json:"alternative_bank_account_names,omitempty"`
	// Deprecated: BankAccountName is deprecated
	BankAccountName string `json:"bank_account_name,omitempty"`
	// Deprecated: FirstName is deprecated
	FirstName string `json:"first_name,omitempty"`
	// Deprecated: Title is deprecated
	Title string `json:"title,omitempty"`
	// Deprecated: ProcessingService is deprecated
	ProcessingService string `json:"processing_service,omitempty"`
	// Deprecated: UserDefinedInformation is deprecated
	UserDefinedInformation string `json:"user_defined_information,omitempty"`
	// Deprecated: AccountMatchingOptOut is deprecated
	AccountMatchingOptOut bool `json:"account_matching_opt_out,omitempty"`
	// Deprecated: Switched is deprecated
	Switched *bool `json:"switched,omitempty"`
}

type PrivateIdentification struct {
	BirthDate      string   `json:"birth_date,omitempty"`
	BirthCountry   string   `json:"birth_country,omitempty"`
	Identification string   `json:"identification,omitempty"`
	Address        []string `json:"address,omitempty"`
	City           string   `json:"city,omitempty"`
	Country        string   `json:"country,omitempty"`
	// Deprecated: Title is deprecated
	Title string `json:"title,omitempty"`
	// Deprecated: FirstName is deprecated
	FirstName string `json:"first_name,omitempty"`
	// Deprecated: LastName is deprecated
	LastName string `json:"last_name,omitempty"`
	// Deprecated: DocumentNumber is deprecated
	DocumentNumber string `json:"document_number,omitempty"`
}

type OrganisationIdentification struct {
	Identification string   `json:"identification,omitempty"`
	Address        []string `json:"address,omitempty"`
	City           string   `json:"city,omitempty"`
	Country        string   `json:"country,omitempty"`
	Actors         []Actor  `json:"actors,omitempty"`
	// Deprecated: Name is deprecated
	Name string `json:"name,omitempty"`
	// Deprecated: RegistrationNumber is deprecated
	RegistrationNumber string `json:"registration_number,omitempty"`
	// Deprecated: Representative is deprecated
	Representative *OrganisationRepresentative `json:"representative,omitempty"`
}

// Deprecated: OrganisationRepresentative is deprecated
type OrganisationRepresentative struct {
	Name string `json:"name,omitempty"`
	// Deprecated: BirthDate is deprecated
	BirthDate string `json:"birth_date,omitempty"`
	// Deprecated: Residency is deprecated
	Residency string `json:"residency,omitempty"`
}

type Actor struct {
	Name      []string `json:"name,omitempty"`
	BirthDate string   `json:"birth_date,omitempty"`
	Residency string   `json:"residency,omitempty"`
}

type Relationships struct {
	AccountEvents *RelationshipAccountEvents `json:"account_events,omitempty"`
	MasterAccount *RelationshipMasterAccount `json:"master_account,omitempty"`
}

type RelationshipMasterAccount struct {
	Data []AccountReference `json:"data"`
}

type AccountReference struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type RelationshipAccountEvents struct {
	Data []AccountEventReference `json:"data"`
}

type AccountEventReference struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type UserDefinedData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
