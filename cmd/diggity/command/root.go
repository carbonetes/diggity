package command

func Run() error {
	err := scan.Execute()
	if err != nil {
		return err
	}
	return nil
}
