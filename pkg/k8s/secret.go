package k8s

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	clientset *kubernetes.Clientset
}

func NewClient(kubeconfigPath string) (*Client, error) {
	if kubeconfigPath == "" {
		kubeconfigPath = clientcmd.RecommendedHomeFile
	}
	
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
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

func (c *Client) GetSecretData(ctx context.Context, name, namespace string) (map[string]string, error) {
	secret, err := c.clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}

	data := make(map[string]string)
	for key, value := range secret.Data {
		// Kubernetes automatically base64 encodes secret data, so we need to decode it
		decodedValue, err := base64.StdEncoding.DecodeString(string(value))
		if err != nil {
			// If decoding fails, use the raw value (in case it wasn't base64 encoded)
			data[key] = string(value)
		} else {
			data[key] = string(decodedValue)
		}
	}

	return data, nil
}

func (c *Client) DeleteKeysFromSecret(ctx context.Context, name, namespace string, keysToDelete []string) (int, error) {
	secret, err := c.clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to get secret: %w", err)
	}

	if secret.Data == nil {
		return 0, fmt.Errorf("secret has no data")
	}

	var keysDeleted int
	for _, key := range keysToDelete {
		if _, exists := secret.Data[key]; exists {
			delete(secret.Data, key)
			keysDeleted++
		}
	}

	if keysDeleted == 0 {
		return 0, fmt.Errorf("none of the specified keys exist in the secret")
	}

	_, err = c.clientset.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to update secret: %w", err)
	}

	return keysDeleted, nil
}
