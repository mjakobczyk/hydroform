package overrides

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"gopkg.in/yaml.v3"
)

func Test_ReadOverridesFromCluster(t *testing.T) {
	component := "monitoring"

	// fake k8s with override ConfigMaps
	k8sMock := fake.NewSimpleClientset(
		&v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "monitoring-overrides",
				Namespace: "kyma-installer",
				Labels:    map[string]string{"installer": "overrides", "component": component},
			},
			Data: map[string]string{
				"componentOverride1": "test1",
				"componentOverride2": "test2",
			},
		},
		&v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "global-overrides",
				Namespace: "kyma-installer",
				Labels:    map[string]string{"installer": "overrides"},
			},
			Data: map[string]string{
				"global.globalOverride1": "test1",
				"global.globalOverride2": "test2",
			},
		},
	)

	// Read additional overrides file
	data, err := ioutil.ReadFile("../test/data/overrides.yaml")
	var overrides map[string]interface{}
	err = yaml.Unmarshal(data, &overrides)
	require.NoError(t, err)

	data, err = ioutil.ReadFile("../test/data/overrides-colliding.yaml")
	var overridesColliding map[string]interface{}
	err = yaml.Unmarshal(data, &overridesColliding)
	require.NoError(t, err)

	t.Run("Should properly read overrides with no colliding data", func(t *testing.T) {
		testProvider, err := New(k8sMock, overrides, log.Printf)
		require.NoError(t, err)

		err = testProvider.ReadOverridesFromCluster()
		require.NoError(t, err)

		// Monitoring should have two sample component overrides + one from additional overrides
		// + global overrides under one "global" key
		res := testProvider.OverridesGetterFunctionFor(component)()
		require.Equal(t, 4, len(res), "Number of component overrides not as expected")

		// Another component without any component override should only have two global overrides + one from additional overrides
		res2 := testProvider.OverridesGetterFunctionFor("anotherComponent")()
		require.Contains(t, res2, "global")
		require.Equal(t, 1, len(res2))
		require.Equal(t, 3, len(res2["global"].(map[string]interface{})), "Number of global overrides not as expected")
	})

	t.Run("Should not duplicate additional overrides when reading overrides many times", func(t *testing.T) {
		testProvider, err := New(k8sMock, overrides, log.Printf)
		require.NoError(t, err)

		err = testProvider.ReadOverridesFromCluster()
		require.NoError(t, err)

		err = testProvider.ReadOverridesFromCluster()
		require.NoError(t, err)

		err = testProvider.ReadOverridesFromCluster()
		require.NoError(t, err)

		// Monitoring should have two sample component overrides + one from additional overrides
		// + global overrides under one "global" key
		res := testProvider.OverridesGetterFunctionFor(component)()
		require.Equal(t, 4, len(res), "Number of component overrides not as expected")

		// Another component without any component override should only have two global overrides + one from additional overrides
		res2 := testProvider.OverridesGetterFunctionFor("anotherComponent")()
		require.Contains(t, res2, "global")
		require.Equal(t, 1, len(res2))
		require.Equal(t, 3, len(res2["global"].(map[string]interface{})), "Number of global overrides not as expected")
	})

	t.Run("Should always put additionalOverrides on top of other overrides", func(t *testing.T) {
		testProvider, err := New(k8sMock, overridesColliding, log.Printf)
		require.NoError(t, err)

		err = testProvider.ReadOverridesFromCluster()
		require.NoError(t, err)

		// Monitoring should have two sample component overrides with one overridden by additional overrides
		// + global overrides under one "global" key - two global overrides from this file with one overridden by additional overrides
		res := testProvider.OverridesGetterFunctionFor(component)()
		require.Equal(t, 3, len(res), "Number of component overrides not as expected")
		require.Equal(t, "changed", res["componentOverride1"], "Override from additional overrides not on top of regular overrides")

		// Another component without any component override should only have two global overrides with one overridden by additional overrides
		res2 := testProvider.OverridesGetterFunctionFor("anotherComponent")()
		require.Contains(t, res2, "global")
		require.Equal(t, 1, len(res2))
		globalOverrides, ok := res2["global"].(map[string]interface{})
		require.True(t, ok)
		require.Equal(t, 2, len(res2["global"].(map[string]interface{})), "Number of global overrides not as expected")
		require.Equal(t, "changed", globalOverrides["globalOverride1"], "Override from additional overrides not on top of regular overrides")
	})

	t.Run("Test non-nested overrides", func(t *testing.T) {
		illegalOverrides := make(map[string]interface{})
		illegalOverrides["componentX"] = "value1"
		illegalOverrides["componentY"] = "value2"
		_, err := New(k8sMock, illegalOverrides, log.Printf)
		require.Error(t, err)
	})
}
