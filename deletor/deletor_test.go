package deletor

import (
	"context"
	"errors"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const namespace = "test"

var (
	podNames = []string{"pod-1", "pod-2", "pod-3"}
)

func TestDeletor(t *testing.T) {
	var calledWithName string

	pods := make([]v1.Pod, 0, len(podNames))
	for _, name := range podNames {
		pods = append(pods, v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace}})
	}
	d := NewDeletor(namespace, func(ctx context.Context) (*v1.PodList, error) {
		return &v1.PodList{Items: pods}, nil
	}, func(ctx context.Context, name string) error {
		calledWithName = name
		return nil
	})

	err := d.DeletePod()
	if err != nil {
		t.Errorf("Got error from DeletePod, expected no error - %v", err)
	}

	assertCalledWith(t, calledWithName, podNames)
}

func TestDeletorErrors(t *testing.T) {
	t.Run("Error while listing pods", func(t *testing.T) {
		errListPods := errors.New("error while listing pods")
		d := NewDeletor(namespace, func(ctx context.Context) (*v1.PodList, error) {
			return nil, errListPods
		}, func(ctx context.Context, name string) error {
			return nil
		})

		err := d.DeletePod()
		if !errors.Is(err, errListPods) {
			t.Errorf("Expected error to be errListPods, was %v", err)
		}
	})

	t.Run("Error while deleting pod", func(t *testing.T) {
		errDeletePod := errors.New("error while deleting pod")
		d := NewDeletor(namespace, func(ctx context.Context) (*v1.PodList, error) {
			return &v1.PodList{Items: []v1.Pod{{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: namespace}}}}, nil
		}, func(ctx context.Context, name string) error {
			return errDeletePod
		})

		err := d.DeletePod()
		if !errors.Is(err, errDeletePod) {
			t.Errorf("Expected error to be errListPods, was %v", err)
		}
	})

}

func assertCalledWith(t *testing.T, actual string, expected []string) {
	for _, name := range expected {
		if name == actual {
			return
		}
	}

	t.Errorf("Delete method was called with name [%v], expected one of [%v]", actual, expected)
}
