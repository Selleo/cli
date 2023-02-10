package generators

// errors rerturn first non empty error
func errors(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
