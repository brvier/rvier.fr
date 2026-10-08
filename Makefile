.PHONY: build serve clean

build:
	go run ./generator

serve:
	go run ./generator -serve localhost:8000

clean:
	rm -rf public
