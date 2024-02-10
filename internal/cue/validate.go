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

#Thing: {
    "@type": string
	"@id"?: string
	additionalType?: string | #ValidURL
	alternateName?: string
	description?: string | #Thing
	disambiguatingDescription?: string
	identifier?:  #Thing | string | #ValidURL
	image?: #Thing | #ValidURL
	mainEntityOfPage?: #CreativeWork | #ValidURL
	name?: string
	potentialAction?: #Thing
	sameAs?: #ValidURL
	subjectOf?: #CreativeWork | #Thing
	url?: #ValidURL
}

#Person: {
	#Thing & {
		"@type": "Person"
	}
	affiliation?: string | #Organization
	description?: string
	email?: #ValidEmail
	familyName?: string
	givenName?: string
}

#Organization: {
	#Thing & {
		"@type": "Organization"
	}
	address?: string
	description?: string
	email?: #ValidEmail
}

#ComputerLanguage: {
	#Thing & {
		"@type": "ComputerLanguage"
	}
	version?: string
}

#ListItem: {
	#Thing & {
		"@type": "ListItem"
	}
	item?: string | #Thing
	nextItem?: string | #ListItem
	position?: int | string
	previousItem?: string | #ListItem
}

#ItemList: {
	#Thing & {
		"@type": "ItemList"
	}
	itemListElement?: string | #Thing | #ListItem | [...(string | #Thing | #ListItem)]
	itemListOrder?: string
	numberOfItems?: int
}

#AggregateRating: {
	#Thing & {
		"@type": "AggregateRating"
	}
	itemReviewed?: string | #Thing
	ratingCount?: int
	reviewCount?: int
}

#DefinedTermSet: {
	#Thing & {
		"@type": "DefinedTermSet"
	}
	hasDefinedTerm?: string | #DefinedTerm | [...(string | #DefinedTerm)]
}

#DefinedTerm: {
	#Thing & {
		"@type": "DefinedTerm"
	}
	inDefinedTermSet?: string | #DefinedTermSet
	termCode?: string
}

#CreativeWork: {
	#Thing & {
		"@type": string | *"CreativeWork"
	}
	about?: #Thing
	abstract?: string
	accessMode?: string
	accessModeSufficient?: #ItemList
	accessibilityAPI?: string
	accessibilityControl?: string
	accessibilityFeature?: string
	accessibilityHazard?: string
	accessibilitySummary?: string
	accountablePerson?: #Person | [...#Person]
	acquireLicensePage?: #CreativeWork | #ValidURL
	aggregateRating?: #AggregateRating
	alternativeHeadline?: string
	archivedAt?: #ValidURL | #Thing
	assesses?: #DefinedTerm | string
	associatedMedia?: #Thing
	audience?: #Thing
	audio?: #Thing
	author?: #Organization | #Person | [...(#Organization | #Person)]
	award?: string
	character?: #Person
	citation?: #CreativeWork | string
	comment?: #Thing
	commentCount?: int
	conditionsOfAccess?: string
	contentLocation?: #Thing
	contentRating?: #Thing | string
	contentReferenceTime?: #ValidDate
	contributor?: #Organization | #Person | [...(#Organization | #Person)]
	copyrightHolder?: #Organization | #Person | [...(#Organization | #Person)]
	copyrightNotice?: string
	copyrightYear?: int | float
	correction?: #Thing | string | #ValidURL
	countryOfOrigin?: #Thing
	creativeWorkStatus?: #DefinedTerm | string
	creator?: #Organization | #Person | [...(#Organization | #Person)]
	creditText?: string
	dateCreated?: #ValidDate
	dateModified?: #ValidDate
	datePublished?: #ValidDate
	digitalSourceType?: #Thing
	discussionUrl?: #ValidURL
	editEIDR?: #ValidURL | string
	editor?: #Person | [...#Person]
	educationalAlignment?: #Thing
	educationalLevel?: #DefinedTerm | string | #ValidURL
	educationalUse?: #DefinedTerm | string
	encoding?: #Thing
	encodingFormat?: #ValidURL | string
	exampleOfWork?: #CreativeWork
	expires?: #ValidDate
	funder?: #Organization | #Person | [...(#Organization | #Person)]
	funding?: string | #Thing
	genre?: string | #ValidURL
	hasPart?: #CreativeWork
	headline?: string
	inLanguage?: #Thing | string | [...(#Thing | string)]
	interactionStatistic?: #Thing
	interactivityType?: string
	interpretedAsClaim?: #Thing
	isAccessibleForFree?: bool
	isBasedOn?: #CreativeWork | #Thing | #ValidURL
	isFamilyFriendly?: bool
	isPartOf?:  #CreativeWork | #ValidURL
	keywords?: #DefinedTerm | string | #ValidURL | [...(#DefinedTerm | string | #ValidURL)]
	learningResourceType?: #DefinedTerm | string
	license?: #CreativeWork | #ValidURL
	locationCreated?: #Thing
	mainEntity?: #Thing
	maintainer?: #Person | #Organization | [...(#Person | #Organization)]
	material?: #Thing | string | #ValidURL
	materialExtent?: string
	mentions?: #Thing
	offers?: #Thing
	pattern?: #DefinedTerm | string
	position?: int | string
	producer?: #Organization | #Person | [...(#Organization | #Person)]
	provider?: #Organization | #Person | [...(#Organization | #Person)]
	publication?: #Thing
	publisher?: #Organization | #Person | [...(#Organization | #Person)]
	publisherInprint?: #Organization
	publishingPrinciples?: #CreativeWork | #ValidURL
	recordedAt?: #Thing
	releasedEvent?: #Thing
	review?: #Thing
	schemaVersion?: string | #ValidURL
	sdDatePublished?: #ValidDate
	sdLicense?: #CreativeWork | #ValidURL
	sdPublisher?: #Organization | #Person
	size?: #DefinedTerm | #Thing | string
	sourceOrganization?: #Organization
	spatial?: #Thing
	spatialCoverage?: #Thing
	sponsor?: #Organization | #Person | [...(#Organization | #Person)]
	teaches?: #DefinedTerm | string
	temporal?: #ValidDate | string
	temporalCoverage?: #ValidDate | string | #ValidURL
	text?: string
	thumbnail?: #Thing
	thumbnailUrl?: #ValidURL
	timeRequired?: #Thing
	translatonOfWork?: #CreativeWork
	translator?: #Organization | #Person | [...(#Organization | #Person)]
	typicalAgeRange?: string
	usageInfo?: #CreativeWork | #ValidURL
	version?: int | float | string
	video?: #Thing
	workExample?: #CreativeWork
	workTranslation?: #CreativeWork
}

#SoftwareApplication: {
	#CreativeWork & {
		"@type": "SoftwareApplication"
	}
	name: string
	provider?: (#Organization | #Person) | [...(#Organization | #Person)]
}

#SoftwareSourceCode: {
	"@context": #Context
	#CreativeWork & {
		"@type": "SoftwareSourceCode"
	}
	// schema.org terms
	applicationCategory?: string | #ValidURL
	applicationSubCategory?: string | #ValidURL
	codeRepository?: #ValidURL
	codeSampleType?: string
	downloadUrl?: #ValidURL
	fileFormat?: string | #ValidURL
	fileSize?: string
	installUrl?: #ValidURL
	memoryRequirements?: string | #ValidURL
	operatingSystem?: string
	permissions?: string
	processorRequirements?: string
	programmingLanguage?: #ComputerLanguage | string
	relatedLink?: #ValidURL
	releaseNotes?: string | #ValidURL
	runtimePlatform?: string
	softwareHelp?: #CreativeWork | #SoftwareSourceCode
	softwareRequirements?: #SoftwareApplication | #SoftwareSourceCode | [...(#SoftwareApplication | #SoftwareSourceCode)]
	storageRequirements?: string | #ValidURL
	supportingData?: #Thing
	targetProduct?: #SoftwareApplication
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
