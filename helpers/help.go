package help

import (
	"fmt"

	"github.com/gimtwi/go-library-project/types"
)

func CheckAuthorsGenresKinds(item *types.Item, ar types.AuthorRepository, gr types.GenreRepository, kr types.KindRepository) error {
	associatedAuthors := make([]types.Author, 0)
	associatedGenres := make([]types.Genre, 0)
	associatedKinds := make([]types.Kind, 0)

	for _, author := range item.Authors {
		a, err := ar.GetByID(author.ID)
		if err != nil {
			return fmt.Errorf("author not found")
		}
		associatedAuthors = append(associatedAuthors, *a)
	}

	for _, genre := range item.Genres {
		g, err := gr.GetByID(genre.ID)
		if err != nil {
			return fmt.Errorf("genre not found")
		}
		associatedGenres = append(associatedGenres, *g)
	}

	for _, kind := range item.Kinds {
		k, err := kr.GetByID(kind.ID)
		if err != nil {
			return fmt.Errorf("kind not found")
		}
		associatedKinds = append(associatedKinds, *k)
	}

	item.Authors = associatedAuthors
	item.Genres = associatedGenres
	item.Kinds = associatedKinds

	return nil
}

func DisassociateAuthorsGenresKinds(item *types.Item, ir types.ItemRepository) error {
	var existingItem types.Item

	i, err := ir.GetByID(item.ID)

	if err != nil {
		return fmt.Errorf("item not found")
	}
	existingItem = *i

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
			if err := ir.DisassociateAuthor(&existingItem, &existingAuthor); err != nil {
				return err
			}
		}
	}

	for _, existingGenre := range existingItem.Genres {
		if !associatedGenreIDs[existingGenre.ID] {
			if err := ir.DisassociateGenre(&existingItem, &existingGenre); err != nil {
				return err
			}
		}
	}

	for _, existingKind := range existingItem.Kinds {
		if !associatedKindIDs[existingKind.ID] {
			if err := ir.DisassociateKind(&existingItem, &existingKind); err != nil {
				return err
			}
		}
	}

	return nil
}
