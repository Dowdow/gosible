package action

// Compile-time check that every action type satisfies Args.
var (
	_ Args = (*CopyArgs)(nil)
	_ Args = (*DirArgs)(nil)
	_ Args = (*DockerArgs)(nil)
	_ Args = (*FileArgs)(nil)
	_ Args = ShellArgs("")
)
