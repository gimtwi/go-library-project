package help

import (
	"math"
	"time"

	"github.com/gimtwi/go-library-project/types"
)

func CheckInitialAvailability(hold *types.Hold, item *types.Item, hr types.HoldRepository, lr types.LoanRepository) error {
	holds, err := hr.GetByItemID(item.ID)
	if err != nil {
		return err
	}

	loans, err := lr.GetByItemID(item.ID)
	if err != nil {
		return err
	}

	if item.Quantity > uint(len(loans)) && item.Quantity > uint(len(holds)) {
		hold.IsAvailable = true
	} else {
		hold.IsAvailable = false
	}

	return nil
}

func ReCheckInitialAvailability(hold *types.Hold, item *types.Item, lr types.LoanRepository) error {
	loans, err := lr.GetByItemID(item.ID)
	if err != nil {
		return err
	}

	availableItems := item.Quantity - uint(len(loans))

	if availableItems >= hold.InLinePosition {
		hold.IsAvailable = true
	} else {
		hold.IsAvailable = false
	}

	return nil
}

func CalculateDaysToWait(hold *types.Hold, item *types.Item) error {
	var fullCycles float64 = float64(hold.InLinePosition-1) / float64(item.Quantity)
	waitingDays := uint(math.Round(fullCycles)) * 14

	hold.EstimatedWeeksToWait = daysToWeeks(waitingDays)

	return nil
}

func ReCalculateDaysToWait(hold *types.Hold, item *types.Item, lr types.LoanRepository) error {
	loans, err := lr.GetByItemID(item.ID)
	if err != nil {
		return err
	}

	availableItems := item.Quantity - uint(len(loans))

	if availableItems >= hold.InLinePosition {
		hold.EstimatedWeeksToWait = 0
	} else {
		var fullCycles float64 = float64(hold.InLinePosition) / float64(item.Quantity)
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

func RearrangeHolds(itemID uint, hr types.HoldRepository, lr types.LoanRepository, ir types.ItemRepository) error {
	holds, err := hr.GetByItemID(uint(itemID))
	if err != nil {
		return err
	}

	item, err := ir.GetByID(uint(itemID))
	if err != nil {
		return err
	}

	for i, hold := range holds {
		hold.InLinePosition = uint(i + 1)
		if err := ReCheckInitialAvailability(&hold, item, lr); err != nil {
			return err
		}

		if hold.IsAvailable {
			hold.ExpiryDate = time.Now().Add(3 * 24 * time.Hour) // if item is available for loner than 3 days the hold will expire automatically
			hold.EstimatedWeeksToWait = 0
		} else {

			if err := ReCalculateDaysToWait(&hold, item, lr); err != nil {
				return err
			}
		}

		if err := hr.Update(&hold); err != nil {
			return err
		}
	}
	return nil
}
