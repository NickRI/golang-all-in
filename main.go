package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"

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
	rsaKeyPath := flag.String("rsaKeyPath", "./build/id_rsa", "Set right secret value in cluster")
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

	rsaFile, err := os.Create(*rsaKeyPath)
	if err != nil {
		log.Panicf("Can't create id_rsa error: %+v", err)
	}
	defer rsaFile.Close()

	if _, err = io.Copy(rsaFile, idRsaBody); err != nil {
		log.Panicf("Can't write to file error: %+v", err)
	}

	log.Println("id_rsa file created")
}
