build:
	go1.18beta1 install ./cmd/migrate/.

run:
	cd ./usage && migrate run