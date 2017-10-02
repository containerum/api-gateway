package version

//Version keeps curent app version
var Version string

//String return app version
func String() string {
	if Version == "" {
		return "1.0.0-dev"
	}
	return Version
}
