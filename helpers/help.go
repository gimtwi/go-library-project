package help

import (
	"fmt"

	"github.com/gimtwi/go-library-project/types"
)

func CheckAuthorsAndGenres(book *types.Book, authorRepo types.AuthorRepository, genreRepo types.GenreRepository) error {
	associatedAuthors := make([]types.Author, 0)
	associatedGenres := make([]types.Genre, 0)

	for _, author := range book.Authors {
		a, err := authorRepo.GetByID(author.ID)
		if err != nil {
			return fmt.Errorf("author not found")
		}
		associatedAuthors = append(associatedAuthors, *a)
	}

	for _, genre := range book.Genres {
		g, err := genreRepo.GetByID(genre.ID)
		if err != nil {
			return fmt.Errorf("genre not found")
		}
		associatedGenres = append(associatedGenres, *g)
	}

	book.Authors = associatedAuthors
	book.Genres = associatedGenres

	return nil
}

func DisassociateAuthorsAndGenres(book *types.Book, bookRepo types.BookRepository) error {
	var existingBook types.Book

	b, err := bookRepo.GetByID(book.ID)

	if err != nil {
		return fmt.Errorf("book not found")
	}
	existingBook = *b

	associatedAuthorIDs := make(map[uint]bool)
	associatedGenreIDs := make(map[uint]bool)

	for _, author := range book.Authors {
		associatedAuthorIDs[author.ID] = true
	}

	for _, genre := range book.Genres {
		associatedGenreIDs[genre.ID] = true
	}

	for _, existingGenre := range existingBook.Genres {
		if !associatedGenreIDs[existingGenre.ID] {
			if err := bookRepo.DisassociateGenre(&existingBook, &existingGenre); err != nil {
				return err
			}
		}
	}

	for _, existingAuthor := range existingBook.Authors {
		if !associatedAuthorIDs[existingAuthor.ID] {
			if err := bookRepo.DisassociateAuthor(&existingBook, &existingAuthor); err != nil {
				return err
			}
		}
	}

	return nil
}
