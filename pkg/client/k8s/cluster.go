package k8s

import (
	"context"

	"github.com/emqx/emqx-operator/api/v1beta1"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Cluster the client that knows how to interact with kubernetes to manage EMQ X Cluster
type Cluster interface {
	// UpdateCluster update the EMQ X Cluster
	UpdateCluster(namespace string, cluster v1beta1.Emqx) error
}

// ClusterOption is the EMQ X Cluster client that using API calls to kubernetes.
type ClusterOption struct {
	client client.Client
	logger logr.Logger
}

// NewCluster returns a new EMQ X Cluster client.
func NewCluster(kubeClient client.Client, logger logr.Logger) Cluster {
	logger = logger.WithValues("service", "crd.EMQXCluster")
	return &ClusterOption{
		client: kubeClient,
		logger: logger,
	}
}

// UpdateCluster implement the  Cluster.Interface
func (c *ClusterOption) UpdateCluster(namespace string, emqx v1beta1.Emqx) error {
	emqx.DescConditionsByTime()
	err := c.client.Status().Update(context.TODO(), emqx)
	if err != nil {
		c.logger.WithValues("namespace", namespace, "cluster", emqx.GetName(), "conditions", emqx.GetConditions()).
			Error(err, "emqxClusterStatus")
		return err
	}
	c.logger.WithValues("namespace", namespace, "cluster", emqx.GetName(), "conditions", emqx.GetConditions()).
		V(3).Info("emqxClusterStatus updated")
	return nil
}
