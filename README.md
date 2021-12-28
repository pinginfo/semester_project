# Semester project

## Dependencies
 - pdftotext
 - pdfimages
 - pdfinfos
 - tesseract

## How to launch the Server
```
go run .
```

## Usage

curl

```
curl -F upload=@<filename> <url>
```
or
```
curl -F upload=@<filename> <url> | jq .
```

## Docker
### Build the image
```
docker build . -t ocr
```
### Launch a container
```
docker run -d -p 8080:8080 ocr
```
