package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {
	configFile := flag.String("config", "kube.config", "Set path to local configuration file")
	namespace := flag.String("namespace", "cattle-pipeline", "Set right namespace in cluster")
	secretKey := flag.String("secretKey", "jenkins", "Set right secret key in cluster")
	secretValue := flag.String("secretValue", "jenkins-id-rsa", "Set right secret value in cluster")
	imageName := flag.String("imageName", "registry.strsqr.cloud/golang-all-in:latest", "Set final image name")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *configFile)
	if err != nil {
		log.Panicf("Build configuration error: %+v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicf("Clietnset get error: %+v", err)
	}

	secret, err := clientset.Core().Secrets(*namespace).Get(*secretKey, meta_v1.GetOptions{})
	if err != nil {
		log.Panicf("Can't get %s from %s namespace with error: %+v", *secretKey, *namespace, err)
	}

	idRsaBody := bytes.NewBuffer(secret.Data[*secretValue])
	log.Println("id_rsa file gathered")

	dockerPath, err := exec.LookPath("docker")
	if err != nil {
		log.Fatal("Docker is not installed on your computer")
	}

	idRSAPrepared := strings.NewReplacer("\n", "\\n").Replace(idRsaBody.String())

	cmd := exec.Command(dockerPath, "build", ".", "--build-arg", "ID_RSA="+idRSAPrepared, "-t", *imageName)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("Start command error: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatalf("Wait command error: %v", err)
	}
}
