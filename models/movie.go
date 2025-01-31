package models

type Movie struct {
	Id			int
	Title		string
	Description	string
	ReleaseYear	int
	Director	string
	Rating		int
	IsWatched	bool
	TrailerUrl	string
	PosterUrl	string
	Genres		[]Genre	
}

type MovieFilters struct {
	SearchTerm	string
	GenreId 	string
	IsWatched 	string
	Sort		string
}