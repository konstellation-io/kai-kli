version="$1"
current_date=$(date +'%Y-%m-%d')

echo "Building version $version - $current_date"

go build -o ./dist/kli -ldflags="-X 'main.Version=$version' -X 'main.Date=$current_date'" cmd/main.go

chmod +x ./dist/kli
