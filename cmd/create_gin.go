package cmd

var createGinCmd = createProject(
	"create-gin",
	"Create a new Go Gin project",
	"Creates a new Go Gin project with a basic structure and specified options.",
	[]string{
		"templates/base/gitignore.tmpl",
		"templates/gin/main.go.tmpl",
		"templates/gin/go.mod.tmpl",
	},
)

func init() {
	rootCmd.AddCommand(createGinCmd)
}
