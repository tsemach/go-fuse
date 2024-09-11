// from: https://dev.to/ilyakaznacheev/a-clean-way-to-pass-configs-in-a-go-application-1g64
package config

type FuseFSMountPoint struct {
	Mountpoint string `yaml:"mountpoint"`
	Targetpath string `yaml:"targetpath"`
}

type Config struct {
	Local bool `envconfig:"FUSEFS_LOCAL",yaml:"local"`
	Filesystems []FuseFSMountPoint `yaml:"filesystems"`
	Stam string `yaml:"stam"`
}
