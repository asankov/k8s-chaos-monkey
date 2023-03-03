package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
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

func run() error {
	// Create an in-cluster kubeConfig for Kubernetes.
	// This will fail if our workload is not running inside a Kubernetes cluster
	// but for now it should be fine, since we expect that we are always running inside Kubernetes.
	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return err
	}

	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		return err
	}

	for {
		// List all pods in the given namespace.
		// If we want to have some sort of timeout for this API call we can create a new TimeoutContext via context.WithTimeout.
		pods, err := clientset.CoreV1().Pods(cfg.Namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("[ERROR] Error while listing Pods: %v", err)
			// Continuing here will make another API call without waiting for the period to pass.
			// We might want to check the type of the error, or implement some sort of exponential back-off
			// to prevent the Kubernetes API being overloaded by requests.
			continue
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// Choose a random pod to be deleted.
		// Use crypto/rand, because math/rand produces deterministic results and should not be used in production code.
		i, err := rand.Int(rand.Reader, big.NewInt(int64(len(pods.Items))))
		if err != nil {
			fmt.Printf("[ERROR] Error while generating random number: %v", err)
			// Continuing here will make another API call without waiting for the period to pass.
			// We might want to check the type of the error, or implement some sort of exponential back-off
			// to prevent the Kubernetes API being overloaded by requests.
			continue
		}
		podToBeDeleted := pods.Items[i.Int64()]

		fmt.Printf("Chose to delete pod [%v]", podToBeDeleted.Name)

		// Delete the pod.
		// If we want to have some sort of timeout for this API call we can create a new TimeoutContext via context.WithTimeout.
		if err := clientset.CoreV1().Pods(cfg.Namespace).Delete(context.Background(), podToBeDeleted.Name, metav1.DeleteOptions{}); err != nil {
			fmt.Printf("[ERROR] Error while deleting pod [%v]: %v", podToBeDeleted.Name, err)
			// Continuing here will make another API call without waiting for the period to pass.
			// We might want to check the type of the error, or implement some sort of exponential back-off
			// to prevent the Kubernetes API being overloaded by requests.
			continue
		}

		fmt.Printf("Succesfully deleted the pod [%v]", podToBeDeleted.Name)

		time.Sleep(time.Duration(cfg.PeriodInSeconds) * time.Second)
	}
}
