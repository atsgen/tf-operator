package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetSecretKey returns value for a key in the specified secret
func GetSecretKey(c client.Client, namespacedName types.NamespacedName, key string) (string, error) {
	ctx := context.TODO()
	instance := &corev1.Secret{}
	err := c.Get(ctx, namespacedName, instance)
	if err != nil {
		return "", err
	}

	value, _ := instance.Data[key]
	return string(value), nil
}

