package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createKubeConfignetwork() (*rest.Config, error) {
	// Check if running inside a Kubernetes cluster
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" && os.Getenv("KUBERNETES_SERVICE_PORT") != "" {
		// Use in-cluster config
		return rest.InClusterConfig()
	}

	// Create out-of-cluster config (e.g., for local development)
	userHomeDir, _ := os.UserHomeDir()
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func createNetworkPolicy(clientset *kubernetes.Clientset, policy *networkingv1.NetworkPolicy) error {
	_, err := clientset.NetworkingV1().NetworkPolicies(policy.Namespace).Create(context.TODO(), policy, metav1.CreateOptions{})
	return err
}

func deleteNetworkPolicy(clientset *kubernetes.Clientset, namespace, policyName string) error {
	return clientset.NetworkingV1().NetworkPolicies(namespace).Delete(context.TODO(), policyName, metav1.DeleteOptions{})
}

func TestNetworkPoliciesss(t *testing.T) {
	// Create a Kubernetes client using the default configuration
	config, err := createKubeConfignetwork()
	if err != nil {
		t.Fatalf("Error creating Kubernetes client config: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Error creating Kubernetes client: %v", err)
	}

	Convey("Given a Kubernetes client", t, func() {
		namespace := "parent-namespace"
		policyName := "allow-nginx"

		// Define a NetworkPolicy for testing (allow traffic to pods with label "app=nginx")
		policy := &networkingv1.NetworkPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      policyName,
				Namespace: namespace,
			},
			Spec: networkingv1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "nginx"},
				},
				Ingress: []networkingv1.NetworkPolicyIngressRule{
					{
						From: []networkingv1.NetworkPolicyPeer{},
					},
				},
				PolicyTypes: []networkingv1.PolicyType{networkingv1.PolicyTypeIngress},
			},
		}

		Convey("When we create a network policy for testing", func() {
			err := createNetworkPolicy(clientset, policy)
			So(err, ShouldBeNil)

			Convey("Then the network policy should exist", func() {
				_, err := clientset.NetworkingV1().NetworkPolicies(namespace).Get(context.TODO(), policyName, metav1.GetOptions{})
				So(err, ShouldBeNil)

				Convey("When we delete the network policy", func() {
					err := deleteNetworkPolicy(clientset, namespace, policyName)
					So(err, ShouldBeNil)

					Convey("Then the network policy should be deleted", func() {
						_, err := clientset.NetworkingV1().NetworkPolicies(namespace).Get(context.TODO(), policyName, metav1.GetOptions{})
						So(err, ShouldNotBeNil)
					})
				})
			})
		})
	})
}
