package errs

import "fmt"

func AppendInfo(err error, msg string) error {
	return fmt.Errorf("[%s]: %w", msg, err)
}
