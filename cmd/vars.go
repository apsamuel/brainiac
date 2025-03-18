package cmd

const configurationKey = "main:configuration"
const aesKeyVariable = "BRAINIAC_AES_KEY"
const aesNonceVariable = "BRAINIAC_AES_NONCE"
const keySize = 32
const nonceSize = 12

var (
	configEngine   string
	configFile     string
	configHost     string
	configPort     int
	configDatabase string
	configTable    string
	configWrite    bool
	configRead     bool
	debug          bool
	generateSecret bool
)
