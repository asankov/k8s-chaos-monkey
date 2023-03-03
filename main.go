package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"github.com/asankov/k8s-chaos-monkey/config"
)

func main() {
	// Create an in-cluster kubeConfig for Kubernetes.
	// This will fail if our workload is not running inside a Kubernetes cluster
	// but for now it should be fine, since we expect that we are always running inside Kubernetes.
	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("[ERROR] error while getting in cluster config: %v", err)
		return
	}
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		fmt.Printf("[ERROR] error while building Kubernetes client: %v", err)
		return
	}

	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		fmt.Printf("[ERROR] error while getting application config: %v", err)
		return
	}

	// Listen for SIGINT and SIGKILL OS channel
	// and try to exit gracefully once received.
	cancelSignal := make(chan os.Signal)
	signal.Notify(cancelSignal, os.Interrupt, os.Kill)

	for {
		select {
		case <-time.After(time.Duration(cfg.PeriodInSeconds) * time.Second):
			if err := deletePod(clientset, cfg); err != nil {
				fmt.Printf("[ERROR]: %v", err)
			}
		case <-cancelSignal:
			fmt.Printf("Detected cancel signal, stopping application")
			return
		}
	}
}

func deletePod(clientset *kubernetes.Clientset, cfg *config.Config) error {
	// List all pods in the given namespace.
	// If we want to have some sort of timeout for this API call we can create a new TimeoutContext via context.WithTimeout.
	pods, err := clientset.CoreV1().Pods(cfg.Namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error while listing Pods: %w", err)
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	// Choose a random pod to be deleted.
	// Use crypto/rand, because math/rand produces deterministic results and should not be used in production code.
	i, err := rand.Int(rand.Reader, big.NewInt(int64(len(pods.Items))))
	if err != nil {
		return fmt.Errorf("error while generating random number: %w", err)
	}
	podToBeDeleted := pods.Items[i.Int64()]

	fmt.Printf("Chose to delete pod [%v]", podToBeDeleted.Name)

	// Delete the pod.
	// If we want to have some sort of timeout for this API call we can create a new TimeoutContext via context.WithTimeout.
	if err := clientset.CoreV1().Pods(cfg.Namespace).Delete(context.Background(), podToBeDeleted.Name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("error while deleting pod [%v]: %w", podToBeDeleted.Name, err)
	}

	fmt.Printf("Succesfully deleted the pod [%v]", podToBeDeleted.Name)

	return nil
}
