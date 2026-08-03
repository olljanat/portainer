package cli

import (
	"context"

	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getPortainerUserDefaultPolicies() []rbacv1.PolicyRule {
	return []rbacv1.PolicyRule{
		{
			Verbs:     []string{"get", "list", "watch"},
			Resources: []string{"*"},
			APIGroups: []string{"*"},
		},
		{
			Verbs:     []string{"create"},
			Resources: []string{"pods/exec"},
			APIGroups: []string{""},
		},
		{
			Verbs:     []string{"delete"},
			Resources: []string{"jobs", "pods"},
			APIGroups: []string{""},
		},
	}
}

func (kcl *KubeClient) upsertPortainerK8sClusterRoles() error {
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: portainerUserCRName,
		},
		Rules: getPortainerUserDefaultPolicies(),
	}

	_, err := kcl.cli.RbacV1().ClusterRoles().Create(context.TODO(), clusterRole, metav1.CreateOptions{})
	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
			_, err = kcl.cli.RbacV1().ClusterRoles().Update(context.TODO(), clusterRole, metav1.UpdateOptions{})
		}
		if err != nil {
			return err
		}
	}

	return nil
}
