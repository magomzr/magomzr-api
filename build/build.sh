
#!/bin/bash

cd ../cmd

echo "Building Lambda function..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap

echo "Creating deployment package..."
zip function.zip bootstrap

echo "Cleaning up..."
rm bootstrap

mv function.zip ../build

echo "Build complete! Upload function.zip to AWS Lambda"