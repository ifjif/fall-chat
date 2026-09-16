package errs

import "fmt"

type AppError struct {
	Op      string
	Err     error
	Context map[string]any
}

func (ae *AppError) Error() string {
	return fmt.Sprintf("%s faild: %v", ae.Op, ae.Err)
}

func (ae *AppError) Unwrap() error {
	return ae.Err
}

func (ae *AppError) WithContext(fields ...any) *AppError {
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			if key, ok := fields[i].(string); ok {
				ae.Context[key] = fields[i+1]
			}
		}
	}

	return ae
}

func (ae *AppError) ExtractContext() []any {
	var result []any
	for k, v := range ae.Context {
		result = append(result, k, v)
	}

	return result
}

func Wrap(op string, err error, fields ...any) error {
	if err == nil {
		return nil
	}

	ctxMap := make(map[string]any)

	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key, ok := fields[i].(string)
			if ok {
				ctxMap[key] = fields[i+1]
			}
		}
	}

	return &AppError{
		Op:      op,
		Err:     err,
		Context: ctxMap,
	}
}
