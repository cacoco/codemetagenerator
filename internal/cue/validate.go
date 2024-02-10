package cue

import (
	"fmt"
	"sort"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	cuetoken "cuelang.org/go/cue/token"
	"cuelang.org/go/encoding/json"
)

// see: https://ijmacd.github.io/rfc3339-iso8601/
const (
	schema = `
import "time"

#Context: "https://w3id.org/codemeta/3.0" | "https://raw.githubusercontent.com/codemeta/codemeta/master/codemeta.json" | "https://doi.org/10.5063/schema/codemeta-2.0" | "https://raw.githubusercontent.com/codemeta/codemeta/master/codemeta.jsonld"
#DevelopmentStatus: =~ {"(?i)^Abandoned$" | "(?i)^Active$" | "(?i)^Concept$" | "(?i)^Inactive$" | "(?i)^Moved$" | "(?i)^Suspended$" | "(?i)^Unsupported$" | "(?i)^WIP$"}

#ValidEmail: =~ "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
#ValidURL: =~ "^(http:\/\/www\\.|https:\/\/www\\.|http:\/\/|https:\/\/|\/|\/\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\\-\\.]{1}[a-z0-9]+)*\\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$"
#ValidDate: time.Format(time.RFC3339) | time.Format("2006-01-02")
#ValidTime: time.Format(time.Kitchen24) | time.Format("15:04:05Z") | time.Format("15:04:05-07:00") | time.Format("15:04:05+07:00")

#Action: {
	#Thing & {
		'@type': "Action"
	}
}

#Event: {
	#Thing & {
		'@type': "Event"
	}
}

#StructuredValue: {
	#Thing & {
		'@type': "StructuredValue"
	}
}

#PropertyValue: {
	#Thing & {
		'@type': "PropertyValue"
	}
	maxValue?: int | float
	measurementMethod?: string | #ValidURL
	measurementTechnique?: string | #ValidURL
	minValue?: int | float
	propertyID?: string | #ValidURL
	unitCode?: string | #ValidURL
	value?: string | int | float | bool | #StructuredValue
	valueReference?: string | #ValidURL | #PropertyValue | #StructuredValue
}

#Thing: {
    '@type': string
	additionalType?: string | #ValidURL
	alternateName?: string
	description?: string
	disambiguatingDescription?: string
	identifier?: string | #ValidURL | #PropertyValue
	image?: string | #ValidURL
	mainEntityOfPage?: string | #CreativeWork | #ValidURL
	name?: string
	potentialAction?: #Action | [...#Action]
	sameAs?: #ValidURL
	subjectOf?: string | #CreativeWork | #Event
	url?: #ValidURL
}

#Person: {
	#Thing & {
		'@type': "Person"
	}
	'@id'?: string
	affiliation?: string | #Organization
	description?: string
	email?: #ValidEmail
	familyName?: string
	givenName?: string
}

#Organization: {
	#Thing & {
		'@type': "Organization"
	}
	'@id'?: string
	address?: string
	description?: string
	email?: #ValidEmail
}

#ComputerLanguage: {
	#Thing & {
		'@type': "ComputerLanguage"
	}
	version?: string
}

#ListItem: {
	#Thing & {
		'@type': "ListItem"
	}
	item?: string | #Thing
	nextItem?: string | #ListItem
	position?: int | string
	previousItem?: string | #ListItem
}

#ItemList: {
	#Thing & {
		'@type': "ItemList"
	}
	itemListElement?: string | #Thing | #ListItem | [...(string | #Thing | #ListItem)]
	itemListOrder?: string
	numberOfItems?: int
}

#AggregateRating: {
	#Thing & {
		'@type': "AggregateRating"
	}
	itemReviewed?: string | #Thing
	ratingCount?: int
	reviewCount?: int
}

#DefinedTermSet: {
	#Thing & {
		'@type': "DefinedTermSet"
	}
	hasDefinedTerm?: string | #DefinedTerm | [...(string | #DefinedTerm)]
}

#DefinedTerm: {
	#Thing & {
		'@type': "DefinedTerm"
	}
	inDefinedTermSet?: string | #DefinedTermSet
	termCode?: string
}

#CreativeWork: {
	#Thing & {
	}
	about?: string | #Thing
	abstract?: string
	accessMode?: string
	accessModeSufficient?: string | #ItemList
	accessibilityAPI?: string
	accessibilityControl?: string
	accessibilityFeature?: string
	accessibilityHazard?: string
	accessibilitySummary?: string
	accountablePerson?: #Person
	acquireLicensePage?: #CreativeWork | #ValidURL
	aggregateRating?: #AggregateRating
	alternativeHeadline?: string
	archivedAt?: string | #ValidURL
	assesses?: string | #DefinedTerm
	associatedMedia?: string | #ValidURL
	audience?: string
	audio?: string
	author?: #Organization | #Person | [...(#Organization | #Person)]
	award?: string
	character?: string | #Person
	citation?: #CreativeWork | #ValidURL | string
	comment?: string
	commentCount?: int
	conditionsOfAccess?: string
	contentLocation?: string
	contentRating?: string
	contentReferenceTime?: #ValidDate
	contributor?: #Organization | #Person | [...(#Organization | #Person)]
	copyrightHolder?: #Organization | #Person | [...(#Organization | #Person)]
	copyrightNotice?: string
	copyrightYear?: int | float
	correction?: string | #ValidURL
	countryOfOrigin?: string
	creativeWorkStatus?: string | #DefinedTerm
	creator?: #Organization | #Person | [...(#Organization | #Person)]
	creditText?: string
	dateCreated?: #ValidDate
	dateModified?: #ValidDate
	datePublished?: #ValidDate
	digitalSourceType?: string
	discussionUrl?: #ValidURL
	editEIDR?: string | #ValidURL
	editor?: #Person | [...#Person]
	educationalAlignment?: string | [...string]
	educationalLevel?: string | #DefinedTerm | #ValidURL
	educationalUse?: string | #DefinedTerm
	encoding?: string
	encodingFormat?: string | #ValidURL
	exampleOfWork?: string | #CreativeWork
	expires?: #ValidDate
	funder?: string | #Organization | #Person
	funding?: string
	genre?: string | #ValidURL
	hasPart?: string | #CreativeWork
	headline?: string
	inLanguage?: string
	interactionStatistic?: string
	interactivityType?: string
	interpretedAsClaim?: string
	isAccessibleForFree?: bool
	isBasedOn?: string | #CreativeWork | #ValidURL
	isFamilyFriendly?: bool
	isPartOf?: string | #CreativeWork | #ValidURL
	keywords?: string | [...string]
	learningResourceType?: string | #DefinedTerm
	license?: #CreativeWork | #ValidURL
	locationCreated?: string
	mainEntity?: string | #Thing
	maintainer?: #Person | #Organization | [...(#Person | #Organization)]
	material?: string | #ValidURL
	materialExtent?: string
	mentions?: #Thing
	offers?: string
	pattern?: string | #DefinedTerm
	position?: int | string
	producer?: #Organization | #Person | [...(#Organization | #Person)]
	provider?: #Organization | #Person | [...(#Organization | #Person)]
	publication?: string
	publisher?: #Organization | #Person | [...(#Organization | #Person)]
	publisherInprint?: string | #Organization
	publishingPrinciples?: string | #CreativeWork | #ValidURL
	recordedAt?: string
	releasedEvent?: string
	review?: string
	schemaVersion?: string | #ValidURL
	sdDatePublished?: #ValidDate
	sdLicense?: string | #CreativeWork | #ValidURL
	sdPublisher?: string | #Organization | #Person
	size?: string | #DefinedTerm
	sourceOrganization?: string | #Organization
	spatial?: string
	spatialCoverage?: string
	sponsor?: #Organization | #Person | [...(#Organization | #Person)]
	teaches?: string | #DefinedTerm
	temporal?: string | #ValidDate
	temporalCoverage?: string | #ValidDate | #ValidURL
	text?: string
	thumbnail?: string
	thumbnailUrl?: #ValidURL
	timeRequired?: string
	translatonOfWork?: string | #CreativeWork
	translator?: #Organization | #Person | [...(#Organization | #Person)]
	typicalAgeRange?: string
	usageInfo?: string | #CreativeWork | #ValidURL
	version?: string | int | float
	video?: string
	workExample?: string | #CreativeWork
	workTranslation?: string | #CreativeWork
}

#SoftwareApplication: {
	#Thing & {
		'@type': "SoftwareApplication"
	}
	#CreativeWork & {
	}
	name: string
	version?: string
	provider?: (#Organization | #Person) | [...(#Organization | #Person)]
}

#SoftwareSourceCode: {
	#Thing & {
		'@type': "SoftwareSourceCode"
	}
	#CreativeWork & {
	}
	'@context': #Context
	// schema.org terms
	applicationCategory?: string | #ValidURL
	applicationSubCategory?: string | #ValidURL
	codeRepository?: #ValidURL
	downloadUrl?: #ValidURL
	fileFormat?: string | #ValidURL
	fileSize?: string
	installUrl?: #ValidURL
	memoryRequirements?: string | #ValidURL
	operatingSystem?: string
	permissions?: string
	processorRequirements?: string
	programmingLanguage?: string | #ComputerLanguage   
	relatedLink?: #ValidURL
	releaseNotes?: string | #ValidURL
	runtimePlatform?: string
	softwareHelp?: string | #SoftwareSourceCode
	softwareRequirements?: string | [...(#SoftwareApplication | #SoftwareSourceCode)]
	storageRequirements?: string | #ValidURL
	supportingData?: string
	targetProduct?: string
	// codemeta terms
	buildInstructions?: #ValidURL
	continuousIntegration?: #ValidURL | [...#ValidURL]
	developmentStatus?: #DevelopmentStatus
	embargoEndDate?: #ValidDate
	funding?: string
	hasSourceCode?: #SoftwareSourceCode
	isSourceOf?: #SoftwareSourceCode
	issueTracker?: #ValidURL
	maintainer?: (#Organization | #Person) | [...(#Organization | #Person)]
	readme?: #ValidURL 
	referencePublication?: string
	softwareSuggestions?: string | [...(#SoftwareApplication | #SoftwareSourceCode)]
}

{#SoftwareSourceCode}
`
)

var ctx *cue.Context = cuecontext.New()

type list []errors.Error

type sortable struct {
	list
}

func (s sortable) Len() int {
	return len(s.list)
}

func (s sortable) Less(i, j int) bool {
	return s.list[i].Position().Filename() < s.list[j].Position().Filename()
}

func (s sortable) Swap(i, j int) {
	s.list[i], s.list[j] = s.list[j], s.list[i]
}

func (p list) sort() {
	sort.Sort(sortable{p})
}

func (p *list) dedupe() {
	p.sort()
	var last errors.Error
	i := 0
	for _, e := range *p {
		if last == nil || !approximateEqual(last, e) {
			last = e
			(*p)[i] = e
			i++
		}
	}
	(*p) = (*p)[0:i]
}

func approximateEqual(a, b errors.Error) bool {
	aPos := a.Position()
	bPos := b.Position()
	if aPos == cuetoken.NoPos || bPos == cuetoken.NoPos {
		return a.Error() == b.Error()
	}
	return aPos.Filename() == bPos.Filename() &&
		aPos.Line() == bPos.Line() &&
		aPos.Column() == bPos.Column() &&
		equalPath(a.Path(), b.Path())
}

func equalPath(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, x := range a {
		if x != b[i] {
			return false
		}
	}
	return true
}

func Validate(v []byte) error {
	schema := ctx.CompileString(schema, cue.Filename("codemeta.cue"))
	// ensure schema is valid cue
	if schema.Err() != nil {
		msg := errors.Details(schema.Err(), nil)
		return fmt.Errorf(msg)
	}

	var b strings.Builder
	err := json.Validate(v, schema)
	if err != nil {
		var l list = errors.Errors(err)
		l.dedupe()
		for _, e := range l {
			fmt.Fprintf(&b, "%s\n", e.Error())
		}
		return fmt.Errorf(b.String())
	}
	return nil
}
