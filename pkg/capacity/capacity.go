// Copyright 2019 Kube Capacity Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package capacity

import (
	"context"
	"fmt"
	"os"

	"github.com/robscott/kube-capacity/pkg/kube"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

// FetchAndPrint gathers cluster resource data and outputs it
func FetchAndPrint(opts Options) {
	clientset, err := kube.NewClientSet(opts.KubeContext, opts.KubeConfig, opts.InsecureSkipTLSVerify, opts.ImpersonateUser, opts.ImpersonateGroup)
	if err != nil {
		fmt.Printf("Error connecting to Kubernetes: %v\n", err)
		os.Exit(1)
	}

	podList := getPods(clientset, opts.PodLabels, opts.Namespaces)

	var pmList *v1beta1.PodMetricsList

	if opts.ShowUtil {
		mClientset, err := kube.NewMetricsClientSet(opts.KubeContext, opts.KubeConfig, opts.InsecureSkipTLSVerify)
		if err != nil {
			fmt.Printf("Error connecting to Metrics API: %v\n", err)
			os.Exit(4)
		}

		pmList = getPodMetrics(mClientset, opts.Namespaces)
	}

	cm := buildClusterMetric(podList, pmList)
	printList(&cm, opts)
}

func getPods(clientset kubernetes.Interface, podLabels string, namespaces []string) *corev1.PodList {
	var podList *corev1.PodList
	for _, namespace := range namespaces {
		pl, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: podLabels,
		})
		if podList == nil {
			podList = pl
		} else {
			podList.Items = append(podList.Items, pl.Items...)
		}
		if err != nil {
			fmt.Printf("Error listing Pods in namespace %s: %v\n", namespace, err)
			// os.Exit(3)
		}
	}
	if podList == nil {
		podList = new(corev1.PodList)
	}

	return podList
}

func getPodMetrics(mClientset *metrics.Clientset, namespaces []string) *v1beta1.PodMetricsList {
	var metricsList *v1beta1.PodMetricsList

	for _, namespace := range namespaces {
		pmList, err := mClientset.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), metav1.ListOptions{})
		if metricsList == nil {
			metricsList = pmList
		} else {
			metricsList.Items = append(metricsList.Items, pmList.Items...)
		}
		if err != nil {
			fmt.Printf("Error getting Pod Metrics: %v\n", err)
			fmt.Println("For this to work, metrics-server needs to be running in your cluster")
			os.Exit(6)
		}
	}

	return metricsList
}

func getNodeMetrics(mClientset *metrics.Clientset, nodeList *corev1.NodeList, nodeLabels string) *v1beta1.NodeMetricsList {
	nmList, err := mClientset.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{
		LabelSelector: nodeLabels,
	})

	if err != nil {
		fmt.Printf("Error getting Node Metrics: %v\n", err)
		fmt.Println("For this to work, metrics-server needs to be running in your cluster")
		os.Exit(7)
	}

	return nmList
}
