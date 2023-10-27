package help

import (
	"time"

	"github.com/gimtwi/go-library-project/types"
)

func CheckAvailability(hold *types.Hold, book *types.Book, holdRepo types.HoldRepository, loanRepo types.LoanRepository) error {
	holds, err := holdRepo.GetByBookID(book.ID)
	if err != nil {
		return err
	}

	loans, err := loanRepo.GetByBookID(book.ID)
	if err != nil {
		return err
	}

	if book.Quantity > uint(len(loans)) && book.Quantity > uint(len(holds)) {
		hold.IsAvailable = true
	} else {
		hold.IsAvailable = false
	}

	return nil
}

func CalculateDaysToWait(hold *types.Hold, book *types.Book) error {
	// 14 days wait per user per book in line
	fullCycles := (hold.InLinePosition - 1) / book.Quantity

	waitingDays := fullCycles * 14

	hold.EstimatedWeeksToWait = daysToWeeks(waitingDays)
	return nil
}

func daysToWeeks(days uint) uint {
	weeks := days / 7
	return uint(weeks)
}

func RearrangeHolds(bookID uint, holdRepo types.HoldRepository, loanRepo types.LoanRepository, bookRepo types.BookRepository) error {
	holds, err := holdRepo.GetByBookID(uint(bookID))
	if err != nil {
		return err
	}

	book, err := bookRepo.GetByID(uint(bookID))
	if err != nil {
		return err
	}

	for i, hold := range holds {
		if err := CheckAvailability(&hold, book, holdRepo, loanRepo); err != nil {
			return err
		}

		if hold.IsAvailable {
			hold.InLinePosition = 0
			hold.ExpiryDate = time.Now().Add(3 * 24 * time.Hour) // if book is available for loner than 3 days the hold will expire automatically
			hold.EstimatedWeeksToWait = 0
		} else {
			hold.InLinePosition = uint(i + 1)

			if err := CalculateDaysToWait(&hold, book); err != nil {
				return err
			}
		}

		if err := holdRepo.Update(&hold); err != nil {
			return err
		}
	}
	return nil
}

func CalculateHoldsRatio(ownedCopies, holdsCount uint) int {
	if ownedCopies == 0 {
		return 0
	} else if ownedCopies < holdsCount {
		holdsRatio := int(holdsCount / ownedCopies)
		return holdsRatio
	} else {
		return 0
	}
}
