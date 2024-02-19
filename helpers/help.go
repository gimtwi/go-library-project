package help

import (
	"fmt"

	"github.com/gimtwi/go-library-project/types"
)

func CheckAuthorsGenresKinds(item *types.Item, ar types.AuthorRepository, gr types.GenreRepository, kr types.KindRepository) error {
	var (
		associatedAuthors []types.Author
		associatedGenres  []types.Genre
		associatedKinds   []types.Kind
	)

	for _, author := range item.Authors {
		a, err := ar.GetByID(author.ID)
		if err != nil {
			return fmt.Errorf("author not found for ID %d: %v", author.ID, err)
		}
		associatedAuthors = append(associatedAuthors, *a)
	}

	for _, genre := range item.Genres {
		g, err := gr.GetByID(genre.ID)
		if err != nil {
			return fmt.Errorf("genre not found for ID %d: %v", genre.ID, err)
		}
		associatedGenres = append(associatedGenres, *g)
	}

	for _, kind := range item.Kinds {
		k, err := kr.GetByID(kind.ID)
		if err != nil {
			return fmt.Errorf("kind not found for ID %d: %v", kind.ID, err)
		}
		associatedKinds = append(associatedKinds, *k)
	}

	item.Authors = associatedAuthors
	item.Genres = associatedGenres
	item.Kinds = associatedKinds

	return nil
}

func DisassociateAuthorsGenresKinds(item *types.Item, ir types.ItemRepository) error {
	existingItem, err := ir.GetByID(item.ID)
	if err != nil {
		return fmt.Errorf("item not found")
	}

	associatedAuthorIDs := make(map[uint]bool)
	associatedGenreIDs := make(map[uint]bool)
	associatedKindIDs := make(map[uint]bool)

	for _, author := range item.Authors {
		associatedAuthorIDs[author.ID] = true
	}

	for _, genre := range item.Genres {
		associatedGenreIDs[genre.ID] = true
	}

	for _, kind := range item.Kinds {
		associatedKindIDs[kind.ID] = true
	}

	for _, existingAuthor := range existingItem.Authors {
		if !associatedAuthorIDs[existingAuthor.ID] {
			authorCopy := existingAuthor
			if err := ir.DisassociateAuthor(existingItem, &authorCopy); err != nil {
				return err
			}
		}
	}

	for _, existingGenre := range existingItem.Genres {
		if !associatedGenreIDs[existingGenre.ID] {
			genreCopy := existingGenre
			if err := ir.DisassociateGenre(existingItem, &genreCopy); err != nil {
				return err
			}
		}
	}

	for _, existingKind := range existingItem.Kinds {
		if !associatedKindIDs[existingKind.ID] {
			kindCopy := existingKind
			if err := ir.DisassociateKind(existingItem, &kindCopy); err != nil {
				return err
			}
		}
	}

	return nil
}
