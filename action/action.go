package action

// Args is implemented by every action type (copy, dir, docker, file, shell).
type Args interface {
	// Validate checks the action's own fields are structurally correct.
	// It must not depend on anything outside the action itself.
	Validate() error

	// Prepare resolves context-dependent values (relative paths, env(VAR)
	// placeholders) using the resolver functions provided by the config package.
	Prepare(resolvePath func(string) string, replaceEnv func(string) string) error

	// Run executes the action against executor and returns its output.
	Run(executor Executor) (stdout string, stderr string, err error)
}
