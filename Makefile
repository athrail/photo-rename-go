photo-rename:
	go build -o bin/photo-rename cmd/photo-rename/main.go

all: photo-rename

clean:
	rm -rf bin
