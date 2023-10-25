package help

import (
	"time"

	"github.com/gimtwi/go-library-project/types"
)

func AssignAndUpdateHoldsCount(hold *types.Hold, holds []types.Hold, holdRepo types.HoldRepository) error {
	newHoldCount := uint(len(holds) + 1)

	hold.HoldsCount = newHoldCount

	// update HoldsCount for other holds
	if err := UpdateHoldsCountForHolds(newHoldCount, holds, holdRepo); err != nil {
		return err
	}

	return nil
}

func UpdateHoldsCountForHolds(newHoldCount uint, holds []types.Hold, holdRepo types.HoldRepository) error {
	for i := range holds {
		holds[i].HoldsCount = newHoldCount
		if err := holdRepo.Update(&holds[i]); err != nil {
			return err
		}
	}
	return nil
}

func MoveHoldsUpOnePlace(holds []types.Hold, holdRepo types.HoldRepository) error {

	for i := range holds {
		// update the HoldListPosition for the current hold
		holds[i].HoldListPosition = uint(i) - 1

		// update the hold's position in the database
		if err := holdRepo.Update(&holds[i]); err != nil {
			return err
		}
	}

	return nil
}

func HandleOwnedCopies(hold *types.Hold, loans []types.Loan) error {
	hold.OwnedCopies = uint(len(loans))
	return nil
}

func CalculateHoldsRatio(ownedCopies, holdsCount uint) int {
	if ownedCopies < holdsCount {
		holdsRatio := int(holdsCount / ownedCopies)
		return holdsRatio
	} else {
		return 0
	}
}

func UpdateHoldsRatio(holds []types.Hold, holdsRatio int, holdRepo types.HoldRepository) error {
	for i := range holds {
		holds[i].HoldsRatio = holdsRatio
		if err := holdRepo.Update(&holds[i]); err != nil {
			return err
		}
	}
	return nil
}

func UpdateEstimatedWaitDays(holds []types.Hold, holdRepo types.HoldRepository) error {
	holdPointers := make([]*types.Hold, len(holds))

	for i := range holds {
		holdPointers[i] = &holds[i]
	}

	for _, hold := range holdPointers {
		if err := CalculateEstimatedWaitDays(hold); err != nil {
			return err
		}

		if err := holdRepo.Update(hold); err != nil {
			return err
		}
	}
	return nil
}

func UpdateHoldsRatioAndEstimatedWaitDays(holdsRatio int, holds []types.Hold, holdRepo types.HoldRepository) error {

	for i := range holds {
		holds[i].HoldsRatio = holdsRatio

		if err := CalculateEstimatedWaitDays(&holds[i]); err != nil {
			return err
		}

		if err := holdRepo.Update(&holds[i]); err != nil {
			return err
		}
	}

	return nil
}

func CalculateEstimatedWaitDays(hold *types.Hold) error {
	if hold.HoldListPosition == 1 {
		hold.EstimatedWaitDays = 0 // user is first in line, no waiting
		return nil
	}
	hold.EstimatedWaitDays = (hold.HoldListPosition - 1) * 14 // 14 days per user in line
	return nil
}

func CalculateAvailableCopies(book *types.Book, hold *types.Hold, loans []types.Loan) error {
	availableCopies := book.Quantity - uint(len(loans))
	hold.AvailableCopies = availableCopies
	if availableCopies > 0 {
		hold.IsAvailable = true
	} else {
		hold.IsAvailable = false
	}
	if hold.IsAvailable {
		hold.ExpiryDate = time.Now().Add(3 * 24 * time.Hour)
	}
	return nil
}

func UpdateHoldsAfterDelete(bookID uint, holds []types.Hold, holdRepo types.HoldRepository, loanRepo types.LoanRepository) error {
	newHoldCount := uint(len(holds))

	if err := UpdateHoldsCountForHolds(newHoldCount, holds, holdRepo); err != nil {
		return err
	}

	if err := MoveHoldsUpOnePlace(holds, holdRepo); err != nil {
		return err
	}

	loans, err := loanRepo.GetByBookID(bookID)
	if err != nil {
		return err
	}

	holdsRatio := CalculateHoldsRatio(uint(len(loans)), newHoldCount)

	if err := UpdateHoldsRatioAndEstimatedWaitDays(holdsRatio, holds, holdRepo); err != nil {
		return err
	}

	return nil
}
