package controller

import (
	controllerv1 "github.com/PritamDas17021999/kubebuilder-crd/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strings"
)

func newDeployment(pritam *controllerv1.Pritam, deploymentName string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: pritam.Namespace,
			OwnerReferences: []metav1.OwnerReference{

				*metav1.NewControllerRef(pritam, controllerv1.GroupVersion.WithKind("Pritam")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pritam.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": strings.Join(buildSlice(controllerv1.Myapp, deploymentName), "-"),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": strings.Join(buildSlice(controllerv1.Myapp, deploymentName), "-"),
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "my-app",
							Image: pritam.Spec.Container.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: pritam.Spec.Container.Port,
								},
							},
						},
					},
				},
			},
		},
	}
}

func newService(pritam *controllerv1.Pritam, name string, dep_name string) *corev1.Service {
	labels := map[string]string{
		"app": strings.Join(buildSlice(controllerv1.Myapp, dep_name), "-"),
	}
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind: "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: pritam.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(pritam, controllerv1.GroupVersion.WithKind("Pritam")),
			},
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeNodePort,
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Protocol:   "TCP",
					Port:       pritam.Spec.Container.Port,
					TargetPort: intstr.FromInt(int(pritam.Spec.Container.Port)),
				},
			},
		},
	}
}
