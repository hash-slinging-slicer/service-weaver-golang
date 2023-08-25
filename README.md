# service-weaver-golang
Latihan Service Weaver Golang CRUD
- Required weaver, https://serviceweaver.dev/docs.html#installation

# Command
## UNIX
go mod tidy

weaver generate .

go build .

weaver multi deploy config_unix.toml

## Windows
go mod tidy

weaver generate .

go build .

weaver multi deploy config.toml

# Database
download database here: https://drive.google.com/file/d/1wxJZFpqgppamYMyHYH3H3Efkx6jgKtnV/view?usp=sharing
