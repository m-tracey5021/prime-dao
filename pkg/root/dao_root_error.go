package root

type DaoError string

func (err DaoError) Error() string {

	return string(err)
}

const (
	CorruptedInstallation DaoError = DaoError("prm installation is corrupted, root path does not exist")

	ComponentDoesNotExist DaoError = DaoError("prm component does not exist")
)
