package grafanadashboard

import (
	"testing"

	"github.com/stretchr/testify/require"

	grafanav1alpha1 "github.com/integr8ly/grafana-operator/api/integreatly/v1alpha1"
	integreatlyorgv1alpha1 "github.com/integr8ly/grafana-operator/api/integreatly/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var dashboard = &grafanav1alpha1.GrafanaDashboard{
	TypeMeta: metav1.TypeMeta{},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "dashboard2",
		Namespace: "grafana",
	},
	Spec: grafanav1alpha1.GrafanaDashboardSpec{
		Url: "url1",
	},
}

var knownDashboardsEmpty = []*integreatlyorgv1alpha1.GrafanaDashboardRef{
	{
		Name: "dashboard1",
	},
	{
		Name: "dashboard2",
	},
}

var knownDashboards = []*integreatlyorgv1alpha1.GrafanaDashboardRef{
	{
		Name:      "dashboard1",
		Namespace: "grafana",
		Hash:      "1234",
	},
	{
		Name:      "dashboard2",
		Namespace: "grafana",
		Hash:      "5678",
	},
}

func TestFindHashEmpty(t *testing.T) {
	output := findHash(knownDashboardsEmpty, dashboard)
	require.Equal(t, output, "")
}

func TestFindHash(t *testing.T) {
	output := findHash(knownDashboards, dashboard)
	require.Equal(t, output, "5678")
}
