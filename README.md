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

