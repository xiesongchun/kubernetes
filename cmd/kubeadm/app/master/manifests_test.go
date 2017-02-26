/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package master

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"k8s.io/apimachinery/pkg/util/intstr"
	api "k8s.io/client-go/pkg/api/v1"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
)

func TestWriteStaticPodManifests(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Couldn't create tmpdir")
	}
	defer os.RemoveAll(tmpdir)

	// set up tmp GlobalEnvParams values for testing
	oldEnv := kubeadmapi.GlobalEnvParams
	kubeadmapi.GlobalEnvParams.KubernetesDir = fmt.Sprintf("%s/etc/kubernetes", tmpdir)
	defer func() { kubeadmapi.GlobalEnvParams = oldEnv }()

	var tests = []struct {
		cfg      *kubeadmapi.MasterConfiguration
		expected bool
	}{
		{
			cfg:      &kubeadmapi.MasterConfiguration{},
			expected: true,
		},
	}
	for _, rt := range tests {
		actual := WriteStaticPodManifests(rt.cfg)
		if (actual == nil) != rt.expected {
			t.Errorf(
				"failed WriteStaticPodManifests with an error:\n\texpected: %t\n\t  actual: %t",
				rt.expected,
				(actual == nil),
			)
		}
	}
}

func TestEtcdVolume(t *testing.T) {
	var tests = []struct {
		cfg      *kubeadmapi.MasterConfiguration
		expected api.Volume
	}{
		{
			cfg: &kubeadmapi.MasterConfiguration{},
			expected: api.Volume{
				Name: "etcd",
				VolumeSource: api.VolumeSource{
					HostPath: &api.HostPathVolumeSource{
						Path: kubeadmapi.GlobalEnvParams.HostEtcdPath},
				}},
		},
	}

	for _, rt := range tests {
		actual := etcdVolume(rt.cfg)
		if actual.Name != rt.expected.Name {
			t.Errorf(
				"failed etcdVolume:\n\texpected: %s\n\t  actual: %s",
				rt.expected.Name,
				actual.Name,
			)
		}
		if actual.VolumeSource.HostPath.Path != rt.expected.VolumeSource.HostPath.Path {
			t.Errorf(
				"failed etcdVolume:\n\texpected: %s\n\t  actual: %s",
				rt.expected.VolumeSource.HostPath.Path,
				actual.VolumeSource.HostPath.Path,
			)
		}
	}
}

func TestEtcdVolumeMount(t *testing.T) {
	var tests = []struct {
		expected api.VolumeMount
	}{
		{
			expected: api.VolumeMount{
				Name:      "etcd",
				MountPath: "/var/lib/etcd",
			},
		},
	}

	for _, rt := range tests {
		actual := etcdVolumeMount()
		if actual.Name != rt.expected.Name {
			t.Errorf(
				"failed etcdVolumeMount:\n\texpected: %s\n\t  actual: %s",
				rt.expected.Name,
				actual.Name,
			)
		}
		if actual.MountPath != rt.expected.MountPath {
			t.Errorf(
				"failed etcdVolumeMount:\n\texpected: %s\n\t  actual: %s",
				rt.expected.MountPath,
				actual.MountPath,
			)
		}
	}
}

func TestCertsVolume(t *testing.T) {
	var tests = []struct {
		cfg      *kubeadmapi.MasterConfiguration
		expected api.Volume
	}{
		{
			cfg: &kubeadmapi.MasterConfiguration{},
			expected: api.Volume{
				Name: "certs",
				VolumeSource: api.VolumeSource{
					HostPath: &api.HostPathVolumeSource{
						Path: "/etc/ssl/certs"},
				}},
		},
	}

	for _, rt := range tests {
		actual := certsVolume(rt.cfg)
		if actual.Name != rt.expected.Name {
			t.Errorf(
				"failed certsVolume:\n\texpected: %s\n\t  actual: %s",
				rt.expected.Name,
				actual.Name,
			)
		}
		if actual.VolumeSource.HostPath.Path != rt.expected.VolumeSource.HostPath.Path {
			t.Errorf(
				"failed certsVolume:\n\texpected: %s\n\t  actual: %s",
				rt.expected.VolumeSource.HostPath.Path,
				actual.VolumeSource.HostPath.Path,
			)
		}
	}
}

func TestCertsVolumeMount(t *testing.T) {
	var tests = []struct {
		expected api.VolumeMount
	}{
		{
			expected: api.VolumeMount{
				Name:      "certs",
				MountPath: "/etc/ssl/certs",
			},
		},
	}

	for _, rt := range tests {
		actual := certsVolumeMount()
		if actual.Name != rt.expected.Name {
			t.Errorf(
				"failed certsVolumeMount:\n\texpected: %s\n\t  actual: %s",
				rt.expected.Name,
				actual.Name,
			)
		}
		if actual.MountPath != rt.expected.MountPath {
			t.Errorf(
				"failed certsVolumeMount:\n\texpected: %s\n\t  actual: %s",
				rt.expected.MountPath,
				actual.MountPath,
			)
		}
	}
}

func TestK8sVolume(t *testing.T) {
	var tests = []struct {
		cfg      *kubeadmapi.MasterConfiguration
		expected api.Volume
	}{
		{
			cfg: &kubeadmapi.MasterConfiguration{},
			expected: api.Volume{
				Name: "k8s",
				VolumeSource: api.VolumeSource{
					HostPath: &api.HostPathVolumeSource{
						Path: kubeadmapi.GlobalEnvParams.KubernetesDir},
				}},
		},
	}

	for _, rt := range tests {
		actual := k8sVolume(rt.cfg)
		if actual.Name != rt.expected.Name {
			t.Errorf(
				"failed k8sVolume:\n\texpected: %s\n\t  actual: %s",
				rt.expected.Name,
				actual.Name,
			)
		}
		if actual.VolumeSource.HostPath.Path != rt.expected.VolumeSource.HostPath.Path {
			t.Errorf(
				"failed k8sVolume:\n\texpected: %s\n\t  actual: %s",
				rt.expected.VolumeSource.HostPath.Path,
				actual.VolumeSource.HostPath.Path,
			)
		}
	}
}

func TestK8sVolumeMount(t *testing.T) {
	var tests = []struct {
		expected api.VolumeMount
	}{
		{
			expected: api.VolumeMount{
				Name:      "k8s",
				MountPath: "/etc/kubernetes/",
				ReadOnly:  true,
			},
		},
	}

	for _, rt := range tests {
		actual := k8sVolumeMount()
		if actual.Name != rt.expected.Name {
			t.Errorf(
				"failed k8sVolumeMount:\n\texpected: %s\n\t  actual: %s",
				rt.expected.Name,
				actual.Name,
			)
		}
		if actual.MountPath != rt.expected.MountPath {
			t.Errorf(
				"failed k8sVolumeMount:\n\texpected: %s\n\t  actual: %s",
				rt.expected.MountPath,
				actual.MountPath,
			)
		}
		if actual.ReadOnly != rt.expected.ReadOnly {
			t.Errorf(
				"failed k8sVolumeMount:\n\texpected: %t\n\t  actual: %t",
				rt.expected.ReadOnly,
				actual.ReadOnly,
			)
		}
	}
}

func TestComponentResources(t *testing.T) {
	a := componentResources("250m")
	if a.Requests == nil {
		t.Errorf(
			"failed componentResources, return value was nil",
		)
	}
}

func TestComponentProbe(t *testing.T) {
	var tests = []struct {
		port int
		path string
	}{
		{
			port: 1,
			path: "foo",
		},
	}
	for _, rt := range tests {
		actual := componentProbe(rt.port, rt.path)
		if actual.Handler.HTTPGet.Port != intstr.FromInt(rt.port) {
			t.Errorf(
				"failed componentProbe:\n\texpected: %v\n\t  actual: %v",
				rt.port,
				actual.Handler.HTTPGet.Port,
			)
		}
		if actual.Handler.HTTPGet.Path != rt.path {
			t.Errorf(
				"failed componentProbe:\n\texpected: %s\n\t  actual: %s",
				rt.path,
				actual.Handler.HTTPGet.Path,
			)
		}
	}
}

func TestComponentPod(t *testing.T) {
	var tests = []struct {
		n string
	}{
		{
			n: "foo",
		},
	}

	for _, rt := range tests {
		c := api.Container{Name: rt.n}
		v := api.Volume{}
		actual := componentPod(c, v)
		if actual.ObjectMeta.Name != rt.n {
			t.Errorf(
				"failed componentPod:\n\texpected: %s\n\t  actual: %s",
				rt.n,
				actual.ObjectMeta.Name,
			)
		}
	}
}

func TestGetComponentBaseCommand(t *testing.T) {
	var tests = []struct {
		c        string
		expected []string
	}{
		{
			c:        "foo",
			expected: []string{"kube-foo", "--v=2"},
		},
		{
			c:        "bar",
			expected: []string{"kube-bar", "--v=2"},
		},
	}

	for _, rt := range tests {
		actual := getComponentBaseCommand(rt.c)
		for i := range actual {
			if actual[i] != rt.expected[i] {
				t.Errorf(
					"failed getComponentBaseCommand:\n\texpected: %s\n\t  actual: %s",
					rt.expected[i],
					actual[i],
				)
			}
		}
	}
}

func TestGetAPIServerCommand(t *testing.T) {
	var tests = []struct {
		cfg      *kubeadmapi.MasterConfiguration
		expected []string
	}{
		{
			cfg: &kubeadmapi.MasterConfiguration{
				API:        kubeadm.API{Port: 123},
				Networking: kubeadm.Networking{ServiceSubnet: "bar"},
			},
			expected: []string{
				"kube-apiserver",
				"--insecure-bind-address=127.0.0.1",
				"--admission-control=NamespaceLifecycle,LimitRanger,ServiceAccount,PersistentVolumeLabel,DefaultStorageClass,ResourceQuota,DefaultTolerationSeconds",
				"--service-cluster-ip-range=bar",
				"--service-account-key-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/sa.pub",
				"--client-ca-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/ca.crt",
				"--tls-cert-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/apiserver.crt",
				"--tls-private-key-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/apiserver.key",
				"--kubelet-client-certificate=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/apiserver-kubelet-client.crt",
				"--kubelet-client-key=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/apiserver-kubelet-client.key",
				"--token-auth-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/tokens.csv",
				fmt.Sprintf("--secure-port=%d", 123),
				"--allow-privileged",
				"--storage-backend=etcd3",
				"--kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname",
				"--requestheader-username-headers=X-Remote-User",
				"--requestheader-group-headers=X-Remote-Group",
				"--requestheader-extra-headers-prefix=X-Remote-Extra-",
				"--requestheader-client-ca-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/front-proxy-ca.crt",
				"--requestheader-allowed-names=front-proxy-client",
				"--authorization-mode=RBAC",
				"--etcd-servers=http://127.0.0.1:2379",
			},
		},
		{
			cfg: &kubeadmapi.MasterConfiguration{
				API:        kubeadm.API{Port: 123, AdvertiseAddresses: []string{"foo"}},
				Networking: kubeadm.Networking{ServiceSubnet: "bar"},
			},
			expected: []string{
				"kube-apiserver",
				"--insecure-bind-address=127.0.0.1",
				"--admission-control=NamespaceLifecycle,LimitRanger,ServiceAccount,PersistentVolumeLabel,DefaultStorageClass,ResourceQuota,DefaultTolerationSeconds",
				"--service-cluster-ip-range=bar",
				"--service-account-key-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/sa.pub",
				"--client-ca-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/ca.crt",
				"--tls-cert-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/apiserver.crt",
				"--tls-private-key-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/apiserver.key",
				"--kubelet-client-certificate=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/apiserver-kubelet-client.crt",
				"--kubelet-client-key=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/apiserver-kubelet-client.key",
				"--token-auth-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/tokens.csv",
				fmt.Sprintf("--secure-port=%d", 123),
				"--allow-privileged",
				"--storage-backend=etcd3",
				"--kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname",
				"--requestheader-username-headers=X-Remote-User",
				"--requestheader-group-headers=X-Remote-Group",
				"--requestheader-extra-headers-prefix=X-Remote-Extra-",
				"--requestheader-client-ca-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/front-proxy-ca.crt",
				"--requestheader-allowed-names=front-proxy-client",
				"--authorization-mode=RBAC",
				"--advertise-address=foo",
				"--etcd-servers=http://127.0.0.1:2379",
			},
		},
		{
			cfg: &kubeadmapi.MasterConfiguration{
				API:        kubeadm.API{Port: 123},
				Networking: kubeadm.Networking{ServiceSubnet: "bar"},
				Etcd:       kubeadm.Etcd{CertFile: "fiz", KeyFile: "faz"},
			},
			expected: []string{
				"kube-apiserver",
				"--insecure-bind-address=127.0.0.1",
				"--admission-control=NamespaceLifecycle,LimitRanger,ServiceAccount,PersistentVolumeLabel,DefaultStorageClass,ResourceQuota,DefaultTolerationSeconds",
				"--service-cluster-ip-range=bar",
				"--service-account-key-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/sa.pub",
				"--client-ca-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/ca.crt",
				"--tls-cert-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/apiserver.crt",
				"--tls-private-key-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/apiserver.key",
				"--kubelet-client-certificate=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/apiserver-kubelet-client.crt",
				"--kubelet-client-key=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/apiserver-kubelet-client.key",
				"--token-auth-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/tokens.csv",
				fmt.Sprintf("--secure-port=%d", 123),
				"--allow-privileged",
				"--storage-backend=etcd3",
				"--kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname",
				"--requestheader-username-headers=X-Remote-User",
				"--requestheader-group-headers=X-Remote-Group",
				"--requestheader-extra-headers-prefix=X-Remote-Extra-",
				"--requestheader-client-ca-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/front-proxy-ca.crt",
				"--requestheader-allowed-names=front-proxy-client",
				"--authorization-mode=RBAC",
				"--etcd-servers=http://127.0.0.1:2379",
				"--etcd-certfile=fiz",
				"--etcd-keyfile=faz",
			},
		},
	}

	for _, rt := range tests {
		actual := getAPIServerCommand(rt.cfg, false)
		for i := range actual {
			if actual[i] != rt.expected[i] {
				t.Errorf(
					"failed getAPIServerCommand:\n\texpected: %s\n\t  actual: %s",
					rt.expected[i],
					actual[i],
				)
			}
		}
	}
}

func TestGetControllerManagerCommand(t *testing.T) {
	var tests = []struct {
		cfg      *kubeadmapi.MasterConfiguration
		expected []string
	}{
		{
			cfg: &kubeadmapi.MasterConfiguration{},
			expected: []string{
				"kube-controller-manager",
				"--address=127.0.0.1",
				"--leader-elect",
				"--master=127.0.0.1:8080",
				"--root-ca-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/ca.crt",
				"--service-account-private-key-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/sa.key",
				"--cluster-signing-cert-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/ca.crt",
				"--cluster-signing-key-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/ca.key",
				"--insecure-experimental-approve-all-kubelet-csrs-for-group=kubeadm:kubelet-bootstrap",
				"--use-service-account-credentials",
			},
		},
		{
			cfg: &kubeadmapi.MasterConfiguration{CloudProvider: "foo"},
			expected: []string{
				"kube-controller-manager",
				"--address=127.0.0.1",
				"--leader-elect",
				"--master=127.0.0.1:8080",
				"--root-ca-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/ca.crt",
				"--service-account-private-key-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/sa.key",
				"--cluster-signing-cert-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/ca.crt",
				"--cluster-signing-key-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/ca.key",
				"--insecure-experimental-approve-all-kubelet-csrs-for-group=kubeadm:kubelet-bootstrap",
				"--use-service-account-credentials",
				"--cloud-provider=foo",
			},
		},
		{
			cfg: &kubeadmapi.MasterConfiguration{Networking: kubeadm.Networking{PodSubnet: "bar"}},
			expected: []string{
				"kube-controller-manager",
				"--address=127.0.0.1",
				"--leader-elect",
				"--master=127.0.0.1:8080",
				"--root-ca-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/ca.crt",
				"--service-account-private-key-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/sa.key",
				"--cluster-signing-cert-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/ca.crt",
				"--cluster-signing-key-file=" + kubeadmapi.GlobalEnvParams.HostPKIPath + "/ca.key",
				"--insecure-experimental-approve-all-kubelet-csrs-for-group=kubeadm:kubelet-bootstrap",
				"--use-service-account-credentials",
				"--allocate-node-cidrs=true",
				"--cluster-cidr=bar",
			},
		},
	}

	for _, rt := range tests {
		actual := getControllerManagerCommand(rt.cfg, false)
		for i := range actual {
			if actual[i] != rt.expected[i] {
				t.Errorf(
					"failed getControllerManagerCommand:\n\texpected: %s\n\t  actual: %s",
					rt.expected[i],
					actual[i],
				)
			}
		}
	}
}

func TestGetSchedulerCommand(t *testing.T) {
	var tests = []struct {
		cfg      *kubeadmapi.MasterConfiguration
		expected []string
	}{
		{
			cfg: &kubeadmapi.MasterConfiguration{},
			expected: []string{
				"kube-scheduler",
				"--address=127.0.0.1",
				"--leader-elect",
				"--master=127.0.0.1:8080",
			},
		},
	}

	for _, rt := range tests {
		actual := getSchedulerCommand(rt.cfg, false)
		for i := range actual {
			if actual[i] != rt.expected[i] {
				t.Errorf(
					"failed getSchedulerCommand:\n\texpected: %s\n\t  actual: %s",
					rt.expected[i],
					actual[i],
				)
			}
		}
	}
}

func TestGetAuthzParameters(t *testing.T) {
	var tests = []struct {
		authMode string
		expected []string
	}{
		{
			authMode: "",
			expected: []string{
				"--authorization-mode=RBAC",
			},
		},
		{
			authMode: "RBAC",
			expected: []string{
				"--authorization-mode=RBAC",
			},
		},
		{
			authMode: "AlwaysAllow",
			expected: []string{
				"--authorization-mode=RBAC,AlwaysAllow",
			},
		},
		{
			authMode: "AlwaysDeny",
			expected: []string{
				"--authorization-mode=RBAC,AlwaysDeny",
			},
		},
		{
			authMode: "ABAC",
			expected: []string{
				"--authorization-mode=RBAC,ABAC",
				"--authorization-policy-file=/etc/kubernetes/abac_policy.json",
			},
		},
		{
			authMode: "Webhook",
			expected: []string{
				"--authorization-mode=RBAC,Webhook",
				"--authorization-webhook-config-file=/etc/kubernetes/webhook_authz.conf",
			},
		},
	}

	for _, rt := range tests {
		actual := getAuthzParameters(rt.authMode)
		for i := range actual {
			if actual[i] != rt.expected[i] {
				t.Errorf(
					"failed getAuthzParameters:\n\texpected: %s\n\t  actual: %s",
					rt.expected[i],
					actual[i],
				)
			}
		}
	}
}
