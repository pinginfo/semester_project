FROM golang:bullseye
RUN apt update -y
RUN apt install tesseract-ocr poppler-utils exiftool -y
COPY . /app
ENV TESSDATA_PREFIX=/app/tessdata
CMD cd /app && go run main.go
