package k8s

import (
	"bytes"
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	clientset *kubernetes.Clientset
}

func NewClient() (*Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	return &Client{clientset: clientset}, nil
}

func (c *Client) UpsertSecret(ctx context.Context, name, namespace string, envVars map[string]string, replace bool) (int, error) {
	data := make(map[string][]byte)
	for key, value := range envVars {
		data[key] = []byte(value)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Type: corev1.SecretTypeOpaque,
		Data: data,
	}

	existingSecret, err := c.clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		_, createErr := c.clientset.CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
		if createErr != nil {
			return 0, fmt.Errorf("failed to create secret: %w", createErr)
		}
		return len(envVars), nil
	}

	var keysModified int

	if replace {
		existingSecret.Data = data
		keysModified = len(envVars)
	} else {
		if existingSecret.Data == nil {
			existingSecret.Data = make(map[string][]byte)
		}

		for key, value := range data {
			existingValue, exists := existingSecret.Data[key]
			if !exists || !bytes.Equal(existingValue, value) {
				keysModified++
			}
			existingSecret.Data[key] = value
		}
	}
	_, err = c.clientset.CoreV1().Secrets(namespace).Update(ctx, existingSecret, metav1.UpdateOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to update secret: %w", err)
	}

	return keysModified, nil
}
