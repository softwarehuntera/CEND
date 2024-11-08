package datamanagment

// The Name objects reflect rows in the Database, where the Name, Taxonomy pairs are unique.
type Name struct {
	Name string
	Taxonomy	string
	PrimaryName *string
}