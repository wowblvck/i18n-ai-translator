package cli

func Execute(binaryName string, args []string) int {
	if len(args) > 0 && args[0] == "init" {
		return runInitCommand(binaryName, args[1:])
	}
	return runTranslateCommand(binaryName, args)
}
