all: prepare dependencies build_key_getter generate_keys clean push_image


prepare:
	mkdir -p build

dependencies:
	GO111MODULE=on go mod vendor -v

build_key_getter:
	go build -o ./build/get_rsa_key main.go

generate_keys:
	ssh-keyscan github.com > ./build/known_hosts

push_image:

clean:
	rm -rf vendor
	rm -rf build
