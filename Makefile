

run: build
	./banners --debug


build:
	go build -o banners ./vendor/banners 