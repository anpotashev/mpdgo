package mpdapi

// Заворачивет ошибку полученную при вызове функции из internal
func wrapPkgError(err error) error {
	if err == nil {
		return nil
	}
	// todo implementation
	return err
}

func wrapPkgErrorIgnoringAnswer(_ []string, err error) error {
	return wrapPkgError(err)
}
