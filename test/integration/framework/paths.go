/*
Copyright 2020 The Jetstack cert-manager contributors.

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

package framework

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// setEnvTestEnv configures environment variables for controller-runtime's
// 'envtest' package.
func setUpEnvTestEnv() {
	maybeSetEnv("TEST_ASSET_ETCD", "etcd", "hack", "bin", "etcd")
	maybeSetEnv("TEST_ASSET_KUBE_APISERVER", "kube-apiserver", "hack", "bin", "kube-apiserver")
	maybeSetEnv("TEST_ASSET_KUBECTL", "kubectl", "hack", "bin", "kubectl")
}

func maybeSetEnv(key, bin string, path ...string) {
	if os.Getenv(key) != "" {
		return
	}
	p, err := getPath(bin, path...)
	if err != nil {
		panic(fmt.Sprintf(`Failed to find integration test dependency %q.
Either re-run this test using "bazel test //test/integration/{name}" or set the %s environment variable.`, bin, key))
	}
	os.Setenv(key, p)
}

// getCRDsPath returns a path to a directory containing cert-manager CRDs.
func getCRDsPath() (string, error) {
	bazelPath := filepath.Join(os.Getenv("RUNFILES_DIR"), "com_github_jetstack_cert_manager", "deploy", "charts", "cert-manager", "crds")
	d, err := os.Stat(bazelPath)
	if err == os.ErrNotExist {
		return filepath.Join(filepath.Dir(os.Getenv("GOMOD")), "deploy", "charts", "cert-manager", "crds"), nil
	}
	if err != nil {
		return "", err
	}
	if m := d.Mode(); !m.IsDir() {
		return "", fmt.Errorf("directory containing CRDs is not a directory")
	}
	return bazelPath, nil
}

func getPath(name string, path ...string) (string, error) {
	bazelPath := filepath.Join(append([]string{os.Getenv("RUNFILES_DIR"), "com_github_jetstack_cert_manager"}, path...)...)
	p, err := exec.LookPath(bazelPath)
	if err == nil {
		return p, nil
	}

	return exec.LookPath(name)
}
