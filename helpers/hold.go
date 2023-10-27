package help

import (
	"math"
	"time"

	"github.com/gimtwi/go-library-project/types"
)

func CheckInitialAvailability(hold *types.Hold, book *types.Book, holdRepo types.HoldRepository, loanRepo types.LoanRepository) error {
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

func ReCheckInitialAvailability(hold *types.Hold, book *types.Book, loanRepo types.LoanRepository) error {

	loans, err := loanRepo.GetByBookID(book.ID)
	if err != nil {
		return err
	}

	availableBooks := book.Quantity - uint(len(loans))

	if availableBooks >= hold.InLinePosition {
		hold.IsAvailable = true
	} else {
		hold.IsAvailable = false
	}

	return nil
}

func CalculateDaysToWait(hold *types.Hold, book *types.Book) error {
	var fullCycles float64 = float64(hold.InLinePosition-1) / float64(book.Quantity)
	waitingDays := uint(math.Round(fullCycles)) * 14

	hold.EstimatedWeeksToWait = daysToWeeks(waitingDays)

	return nil
}

func ReCalculateDaysToWait(hold *types.Hold, book *types.Book, loanRepo types.LoanRepository) error {
	loans, err := loanRepo.GetByBookID(book.ID)
	if err != nil {
		return err
	}

	availableBooks := book.Quantity - uint(len(loans))

	if availableBooks >= hold.InLinePosition {
		hold.EstimatedWeeksToWait = 0
	} else {
		var fullCycles float64 = float64(hold.InLinePosition) / float64(book.Quantity)
		if uint(fullCycles) == 0 {
			waitingDays := uint(math.Ceil(fullCycles)) * 14
			hold.EstimatedWeeksToWait = daysToWeeks(waitingDays)
		} else {
			waitingDays := uint(math.Round(fullCycles)) * 14
			hold.EstimatedWeeksToWait = daysToWeeks(waitingDays)
		}

	}

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
		hold.InLinePosition = uint(i + 1)
		if err := ReCheckInitialAvailability(&hold, book, loanRepo); err != nil {
			return err
		}

		if hold.IsAvailable {
			hold.ExpiryDate = time.Now().Add(3 * 24 * time.Hour) // if book is available for loner than 3 days the hold will expire automatically
			hold.EstimatedWeeksToWait = 0
		} else {

			if err := ReCalculateDaysToWait(&hold, book, loanRepo); err != nil {
				return err
			}
		}

		if err := holdRepo.Update(&hold); err != nil {
			return err
		}
	}
	return nil
}
