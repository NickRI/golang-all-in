all: prepare dependencies build_key_getter generate_keys generate_image push_image

prepare:
	mkdir -p build

dependencies:
	GO111MODULE=on go mod vendor -v

build_key_getter:
	go build -o ./build/get_rsa_key main.go

generate_keys:
	./build/get_rsa_key
	ssh-keyscan github.com > ./build/known_hosts

generate_image:
	docker build . -t registry.firefly.red/golang-all-in:latest

push_image:
	docker push registry.firefly.red/golang-all-in:latest

clean:
	rm -rf vendor
	rm -rf build
