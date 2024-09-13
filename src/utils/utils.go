package utils

import "fmt"

func AppendError(existErr, newErr error) error {
	if existErr == nil {
		return newErr
	} else {
		return fmt.Errorf("%v, %w", existErr, newErr)
	}
}
