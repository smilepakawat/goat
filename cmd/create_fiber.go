package cmd

var createFiberCmd = createProject(
	"create-fiber",
	"Create a new Go Fiber project",
	"Creates a new Go Fiber project with a basic structure and specified options.",
	[]string{
		"templates/base/gitignore.tmpl",
		"templates/fiber/main.go.tmpl",
		"templates/fiber/go.mod.tmpl",
	},
)

func init() {
	rootCmd.AddCommand(createFiberCmd)
}
