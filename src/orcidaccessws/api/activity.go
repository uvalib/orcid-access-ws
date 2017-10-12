package api

import (
	"sort"
)

//
// ActivityUpdate -- the activity update request
//
type ActivityUpdate struct {
	UpdateCode string     `json:"update_code,omitempty"`
	Work       WorkSchema `json:"work,omitempty"`
}

//
// WorkSchema -- the work schema
//
type WorkSchema struct {
	Title           string   `json:"title,omitempty"`
	Abstract        string   `json:"abstract,omitempty"`
	PublicationDate string   `json:"publication_date,omitempty"`
	URL             string   `json:"url,omitempty"`
	Authors         []Person `json:"authors,omitempty"`
	ResourceType    string   `json:"resource_type,omitempty"`
}

//
// Person -- the basic person details used for authors
//
type Person struct {
	Index     int    `json:"index"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

//
// SortPeople -- helpers to sort the people lists
//
func SortPeople(people []Person) []Person {
	sortedPeople := make([]Person, len(people))
	copy(sortedPeople, people)
	sort.Sort(PeopleSorter(sortedPeople))
	return sortedPeople
}

// PeopleSorter sorts people by index
type PeopleSorter []Person

func (people PeopleSorter) Len() int           { return len(people) }
func (people PeopleSorter) Swap(i, j int)      { people[i], people[j] = people[j], people[i] }
func (people PeopleSorter) Less(i, j int) bool { return people[i].Index < people[j].Index }

//
// end of file
//
