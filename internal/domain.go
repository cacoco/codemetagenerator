package internal

const (
	// JSON Keys
	Type                = "@type"
	Id                  = "@id"
	Context             = "@context"
	Identifier          = "identifier"
	GivenName           = "givenName"
	FamilyName          = "familyName"
	Email               = "email"
	Name                = "name"
	Description         = "description"
	RelatedLink         = "relatedLink"
	CodeRepository      = "codeRepository"
	IssueTracker        = "issueTracker"
	License             = "license"
	Version             = "version"
	ProgrammingLanguage = "programmingLanguage"
	RuntimePlatform     = "runtimePlatform"
	Maintainer          = "maintainer"
	Author              = "author"
	Contributor         = "contributor"
	ReleaseNotes        = "releaseNotes"
	Keywords            = "keywords"
	Readme              = "readme"
	ContIntegration     = "contIntegration"
	DevelopmentStatus   = "developmentStatus"
	URL                 = "url"
	// Implementation Values
	DefaultContext         = "https://doi.org/10.5063/schema/codemeta-2.0"
	PersonType             = "Person"
	OrganizationType       = "Organization"
	SoftwareSourceCodeType = "SoftwareSourceCode"
	ComputerLanguageType   = "ComputerLanguage"
)

type LicenseStruct struct {
	Reference             string   `json:"reference"`
	IsDeprecatedLicenseId bool     `json:"isDeprecatedLicenseId"`
	DetailsURL            string   `json:"detailsUrl"`
	ReferenceNumber       int      `json:"referenceNumber"`
	Name                  string   `json:"name"`
	LicenseId             string   `json:"licenseId"`
	SeeAlso               []string `json:"seeAlso"`
	IsOsiApproved         bool     `json:"isOsiApproved"`
	IsFsfLibre            bool     `json:"isFsfLibre"`
}

type LicensesList struct {
	LicenseListVersion string          `json:"licenseListVersion"`
	Licenses           []LicenseStruct `json:"licenses"`
}

func NewPerson(givenName *string, familyName *string, email *string, id *string) *map[string]string {
	return &map[string]string{
		Type:       PersonType,
		GivenName:  *givenName,
		FamilyName: *familyName,
		Email:      *email,
		Id:         *id,
	}
}

func NewOrganization(name *string, url *string, id *string) *map[string]string {
	return &map[string]string{
		Type: OrganizationType,
		Name: *name,
		URL:  *url,
		Id:   *id,
	}
}

func NewProgrammingLanguage(name *string, url *string) *map[string]string {
	return &map[string]string{
		Type: ComputerLanguageType,
		Name: *name,
		URL:  *url,
	}
}

func NewCodeMetaDefinition(id *string, name *string, description *string, version *string, maintainer *map[string]string, programmingLanguage *map[string]string, developmentStatus string, license *string, runtimePlatform *string, codeRepository *string, readme *string) *map[string]any {
	return &map[string]any{
		Context:             DefaultContext,
		Type:                SoftwareSourceCodeType,
		Identifier:          *id,
		Description:         *description,
		Name:                *name,
		Maintainer:          maintainer,
		CodeRepository:      *codeRepository,
		Version:             *version,
		DevelopmentStatus:   developmentStatus,
		ProgrammingLanguage: programmingLanguage,
		RuntimePlatform:     *runtimePlatform,
		License:             *license,
		Readme:              *readme,
	}
}
