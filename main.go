package main

import (
	"context"
	"fmt"
	"os"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/asankov/k8s-chaos-monkey/config"
	"github.com/asankov/k8s-chaos-monkey/deletor"
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

	deletor := deletor.NewDeletor(cfg.Namespace, func(ctx context.Context) (*v1.PodList, error) {
		return clientset.CoreV1().Pods(cfg.Namespace).List(context.Background(), metav1.ListOptions{})
	}, func(ctx context.Context, name string) error {
		return clientset.CoreV1().Pods(cfg.Namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	})

	for {
		if err := deletor.DeletePod(); err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR]: %v\n", err)
		}

		time.Sleep(time.Duration(cfg.PeriodInSeconds) * time.Second)
	}
}
