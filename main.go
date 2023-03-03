package main

import (
	"context"
	"fmt"
	"math/rand"
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
)

func run() error {
	// Create an in-cluster config for Kubernetes.
	// This will fail if our workload is not running inside a Kubernetes cluster
	// but for now it should be fine, since we expect that we are always running inside Kubernetes.
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	for {
		pods, err := clientset.CoreV1().Pods("TODO").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("[ERROR] Error while listing Pods: %v", err)
			continue
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// choose a random pod to be deleted
		i := rand.Intn(len(pods.Items))
		podToBeDeleted := pods.Items[i]

		fmt.Printf("Chose to delete pod [%v]", podToBeDeleted.Name)

		if err := clientset.CoreV1().Pods("TODO").Delete(context.Background(), podToBeDeleted.Name, metav1.DeleteOptions{}); err != nil {
			fmt.Printf("[ERROR] Error while deleting pod [%v]: %v", podToBeDeleted.Name, err)
		}

		time.Sleep(10 * time.Second)
	}
}
