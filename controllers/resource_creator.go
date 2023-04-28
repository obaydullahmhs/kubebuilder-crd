package controllers

import (
	aadeeappsv1 "github.com/obaydullahmhs/kubebuilder-crd/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func newDeployment(aadee *aadeeappsv1.Aadee, deploymentName string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: aadee.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(aadee, aadeeappsv1.GroupVersion.WithKind("Aadee")),
			},
		},
		Spec: appsv1.DeploymentSpec{

			Replicas: aadee.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":  aadee.Name,
					"kind": "Aadee",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  aadee.Name,
						"kind": "Aadee",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  aadee.Name,
							Image: aadee.Spec.Container.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: aadee.Spec.Container.Port,
								},
							},
						},
					},
				},
			},
		},
	}
}
func newService(aadee *aadeeappsv1.Aadee, serviceName string) *corev1.Service {
	labels := map[string]string{
		"app":  aadee.Name,
		"kind": "Aadee",
	}
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind: "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: aadee.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(aadee, aadeeappsv1.GroupVersion.WithKind("Aadee")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Protocol:   "TCP",
					Port:       aadee.Spec.Service.ServicePort,
					TargetPort: intstr.FromInt(int(aadee.Spec.Container.Port)),
					NodePort:   aadee.Spec.Service.ServiceNodePort,
				},
			},
			Type: func() corev1.ServiceType {
				if aadee.Spec.Service.ServiceType == "NodePort" {
					return corev1.ServiceTypeNodePort
				} else {
					return corev1.ServiceTypeClusterIP
				}
			}(),
		},
	}
}
