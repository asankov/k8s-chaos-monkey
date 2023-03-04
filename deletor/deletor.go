package deletor

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	v1 "k8s.io/api/core/v1"
)

// PodDeletor deletes a random post from a list of pods.
type PodDeletor struct {
	namespace string
	listPods  func(ctx context.Context) (*v1.PodList, error)
	deletePod func(ctx context.Context, name string) error
}

// NewDeletor returns a new deletor with the given properties.
func NewDeletor(namespace string, listPods func(ctx context.Context) (*v1.PodList, error), deletePod func(ctx context.Context, name string) error) *PodDeletor {
	return &PodDeletor{
		namespace: namespace,
		listPods:  listPods,
		deletePod: deletePod,
	}
}

// DeletePod calls lists the pods and deletes a random pod returned from that list.
func (p *PodDeletor) DeletePod() error {
	// List all pods in the given namespace.
	// If we want to have some sort of timeout for this API call we can create a new TimeoutContext via context.WithTimeout.
	pods, err := p.listPods(context.Background())
	if err != nil {
		return fmt.Errorf("error while listing Pods: %w", err)
	}
	fmt.Fprintf(os.Stderr, "There are %d pods in the [%v] namespace\n", len(pods.Items), p.namespace)

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
	if err := p.deletePod(context.Background(), podToBeDeleted.Name); err != nil {
		return fmt.Errorf("error while deleting pod [%v]: %w", podToBeDeleted.Name, err)
	}

	fmt.Fprintf(os.Stderr, "Succesfully deleted the pod [%v]\n", podToBeDeleted.Name)

	return nil
}
