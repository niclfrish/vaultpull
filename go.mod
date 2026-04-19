module github.com/vaultpull/vaultpull

go 1.21

require (
	github.com/hashicorp/vault/api v1.10.0
	github.com/spf13/cobra v1.8.0
	github.com/joho/godotenv v1.5.1
)

require (
	// indirect dependencies required by vault/api
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.4 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.7 // indirect
)
