package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/asankov/k8s-chaos-monkey/config"
)

func main() {
	// Create an in-cluster kubeConfig for Kubernetes.
	// This will fail if our workload is not running inside a Kubernetes cluster
	// but for now it should be fine, since we expect that we are always running inside Kubernetes.
	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] error while getting in cluster config: %v\n", err)
		return
	}
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] error while building Kubernetes client: %v\n", err)
		return
	}

	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] error while getting application config: %v\n", err)
		return
	}

	for {
		if err := deletePod(clientset, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR]: %v\n", err)
		}
		// }

		time.Sleep(time.Duration(cfg.PeriodInSeconds) * time.Second)
	}
}

func deletePod(clientset *kubernetes.Clientset, cfg *config.Config) error {
	// List all pods in the given namespace.
	// If we want to have some sort of timeout for this API call we can create a new TimeoutContext via context.WithTimeout.
	pods, err := clientset.CoreV1().Pods(cfg.Namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error while listing Pods: %w", err)
	}
	fmt.Fprintf(os.Stderr, "There are %d pods in the [%v] namespace\n", len(pods.Items), cfg.Namespace)

	// Choose a random pod to be deleted.
	// Use crypto/rand, because math/rand produces deterministic results and should not be used in production code.
	i, err := rand.Int(rand.Reader, big.NewInt(int64(len(pods.Items))))
	if err != nil {
		return fmt.Errorf("error while generating random number: %w", err)
	}
	podToBeDeleted := pods.Items[i.Int64()]

	fmt.Fprintf(os.Stderr, "Chose to delete pod [%v]\n", podToBeDeleted.Name)

	// Delete the pod.
	// If we want to have some sort of timeout for this API call we can create a new TimeoutContext via context.WithTimeout.
	if err := clientset.CoreV1().Pods(cfg.Namespace).Delete(context.Background(), podToBeDeleted.Name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("error while deleting pod [%v]: %w", podToBeDeleted.Name, err)
	}

	fmt.Fprintf(os.Stderr, "Succesfully deleted the pod [%v]\n", podToBeDeleted.Name)

	return nil
}
