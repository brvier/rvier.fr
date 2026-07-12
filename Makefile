.PHONY: build serve clean

build:
	go run ./generator

serve: build
	cd public && python3 -m http.server 8000

clean:
	rm -rf public
