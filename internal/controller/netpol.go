package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

// manageNetpol checks reconciles the network policy
func (r *NamespaceReconciler) manageNetpol(ctx context.Context, namespace *corev1.Namespace) error {
	netpol := &networking.NetworkPolicy{}

	err := r.Get(ctx, types.NamespacedName{Name: "gangway", Namespace: namespace.Name}, netpol)
	if errors.IsNotFound(err) {
		netpol, err := r.buildNetpol(namespace)
		if err != nil {
			return err
		}

		if err := r.Create(ctx, netpol); err != nil {
			return err
		}
	}

	return nil
}

// buildNetpo returns a NetworkPolicy object and sets the controller reference
func (r *NamespaceReconciler) buildNetpol(namespace *corev1.Namespace) (*networking.NetworkPolicy, error) {
	netpol := &networking.NetworkPolicy{
		ObjectMeta: v1.ObjectMeta{
			Name:      "gangway",
			Namespace: namespace.Name,
		},
		Spec: networking.NetworkPolicySpec{
			PodSelector: v1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "web",
				},
			},
			Ingress: []networking.NetworkPolicyIngressRule{
				networking.NetworkPolicyIngressRule{},
			},
		},
	}

	if err := ctrl.SetControllerReference(namespace, netpol, r.Scheme); err != nil {
		return nil, err
	}

	return netpol, nil
}
