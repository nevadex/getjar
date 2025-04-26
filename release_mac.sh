# little script that builds 4 binaries for release

rm getjar_amd64; rm getjar_arm64; rm getjar_amd64.exe; rm getjar_darwin;

echo "deleted old builds, press enter to compile new builds"
# shellcheck disable=SC2162
read -s

echo "building linux (amd64 linux)"
GOOS=linux GOARCH=amd64 go build -o getjar_amd64 .
echo "building raspbian (arm64 linux)"
GOOS=linux GOARCH=arm64 go build -o getjar_arm64 .
echo "building windows (amd64 windows)"
GOOS=windows GOARCH=amd64 go build -o getjar_amd64.exe .
echo "building mac (amd64 darwin)"
GOOS=darwin GOARCH=amd64 go build -o getjar_darwin .