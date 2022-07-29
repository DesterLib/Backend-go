# Build Frontend
bash scripts/build.sh

# Build Backend
go mod vendor
go mod tidy
go build -ldflags "-s -w" -o dester