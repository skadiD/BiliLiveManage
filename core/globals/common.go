package globals

type Common struct {
	DB string `yaml:"db"`
}

var (
	Setting *Common
)
