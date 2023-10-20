package help

import (
	"github.com/gimtwi/go-library-project/types"
)

func SyncTheAuthorsAndGenres(book *types.Book, bookRepo types.BookRepository, authorRepo types.AuthorRepository, genreRepo types.GenreRepository) error {
	associatedAuthors := make([]types.Author, 0)
	associatedGenres := make([]types.Genre, 0)

	for _, author := range book.Author {
		a, err := authorRepo.GetByID(author.ID)
		if err != nil {
			return err
		}
		associatedAuthors = append(associatedAuthors, *a)
	}

	for _, genre := range book.Genre {
		g, err := genreRepo.GetByID(genre.ID)
		if err != nil {
			return err
		}
		associatedGenres = append(associatedGenres, *g)
	}

	book.Author = associatedAuthors
	book.Genre = associatedGenres

	if book.Quantity >= 1 {
		book.IsAvailable = true
	} else {
		book.IsAvailable = false
	}

	return nil
}
