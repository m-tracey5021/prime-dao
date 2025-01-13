package root

type TSFError string

func (err TSFError) Error() string {

	return string(err)
}

const (
	CorruptedInstallation TSFError = TSFError("tsf installation is corrupted")

	ComponentDoesNotExist TSFError = TSFError("tsf component does not exist")
)
