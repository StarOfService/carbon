package version

const DefaultVersion = "0.0.0"

var VERSION string

func GetVersion() string {
	v := DefaultVersion
	if VERSION != "" {
		v = VERSION
	}
	return v
}
