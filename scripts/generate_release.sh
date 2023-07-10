version="$1"
current_date=$(date +'%Y-%m-%d')

echo "Building version $version - $current_date"

go build -o ./dist/kli -ldflags="-X 'main.version=$version' -X 'main.date=$current_date'" cmd/main.go

chmod +x ./dist/kli
