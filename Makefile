all: prepare dependencies build_key_getter generate_keys clean push_image

registry = registry.strsqr.cloud
namespace = cattle-pipeline
secretKey = jenkins
secretValue = jenkins-id-rsa

prepare:
	mkdir -p build

dependencies:
	GO111MODULE=on go mod vendor -v

build_key_getter:
	go build -o ./build/get_rsa_key main.go

generate_keys:
	ssh-keyscan github.com > ./build/known_hosts
	./build/get_rsa_key -imageName $(registry)/golang-all-in:latest -namespace $(namespace) -secretKey $(secretKey) -secretValue $(secretValue)

push_image:
	docker push $(registry)/golang-all-in:latest

clean:
	rm -rf vendor
	rm -rf build
