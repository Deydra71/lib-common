/*
Copyright 2025 Red Hat

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

package keystone

import (
	"context"
	"fmt"
	"strings"

	"github.com/openstack-k8s-operators/lib-common/modules/common/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetKeystoneOverridesByLabel - Get keystone override configuration from a secret by label.
// Can override keystone configuration to point to a central keystone instance.
// Only keys present in the secret are returned.
func GetKeystoneOverridesByLabel(
	ctx context.Context,
	h *helper.Helper,
	namespace string,
	labelKey string,
) (map[string]string, error) {
	// Use direct label selector string for key existence check
	secrets, err := h.GetKClient().CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelKey,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting secrets with label %s: %w", labelKey, err)
	}

	if len(secrets.Items) == 0 {
		h.GetLogger().V(1).Info(fmt.Sprintf("No secrets found with label %s", labelKey))
		return map[string]string{}, nil
	}

	if len(secrets.Items) > 1 {
		return nil, fmt.Errorf("multiple secrets found with label %s, only one is allowed", labelKey)
	}

	// Extract keystone override data from the secret
	secretObj := &secrets.Items[0]
	overrides := make(map[string]string)

	if v, ok := secretObj.Data["region"]; ok {
		overrides["region"] = strings.TrimSpace(string(v))
	} else {
		h.GetLogger().Info("region key not found in keystone overrides secret", "secretName", secretObj.Name)
	}

	if v, ok := secretObj.Data["auth_url"]; ok {
		overrides["auth_url"] = strings.TrimSpace(string(v))
	} else {
		h.GetLogger().Info("auth_url key not found in keystone overrides secret", "secretName", secretObj.Name)
	}

	if v, ok := secretObj.Data["www_authenticate_uri"]; ok {
		overrides["www_authenticate_uri"] = strings.TrimSpace(string(v))
	} else {
		h.GetLogger().Info("www_authenticate_uri key not found in keystone overrides secret", "secretName", secretObj.Name)
	}

	return overrides, nil
}
