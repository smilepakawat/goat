package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunCmd(directory, name string, args ...string) {
	argsStr := strings.Join(args, " ")
	fmt.Printf("Running '%s %s'...\n", name, argsStr)
	cmd := exec.Command(name, args...)
	cmd.Dir = directory
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("failed to run '%s %s': %v\nOutput: %s", name, argsStr, err, string(output))
		os.Exit(1)
	}
	fmt.Printf("'%s %s' completed successfully.\n", name, argsStr)
	fmt.Printf("Project '%s' created successfully!\n", name)
	fmt.Printf("Next steps:\n")
	fmt.Printf("  cd %s\n", name)
	fmt.Printf("  go run main.go\n")
}
