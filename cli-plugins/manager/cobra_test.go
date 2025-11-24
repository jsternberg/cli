package manager

import (
	"testing"

	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"gotest.tools/v3/assert"
)

func init() {
	deterministicResourceAttributesEnvvar = true
}

func TestPluginResourceAttributesEnvvar(t *testing.T) {
	cmd := &cobra.Command{
		Annotations: map[string]string{
			cobra.CommandDisplayNameAnnotation: "docker",
		},
	}

	// Ensure basic usage is fine.
	env := appendPluginResourceAttributesEnvvar(nil, cmd, Plugin{Name: "compose"}, resource.Empty())
	assert.DeepEqual(t, []string{"OTEL_RESOURCE_ATTRIBUTES=docker.cli.cobra.command_path=docker%20compose"}, env)

	// Add a user-based environment variable to OTEL_RESOURCE_ATTRIBUTES.
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "a.b.c=foo")

	env = appendPluginResourceAttributesEnvvar(nil, cmd, Plugin{Name: "compose"}, resource.Empty())
	assert.DeepEqual(t, []string{"OTEL_RESOURCE_ATTRIBUTES=a.b.c=foo,docker.cli.cobra.command_path=docker%20compose"}, env)

	// Resource attributes with docker.cli prefix should be passed down as-is and other
	// values should be filtered.
	res := resource.NewSchemaless(
		attribute.Key("filtered").String("nothere"),
		attribute.Key("docker.cli.current_context").String("default"),
	)
	env = appendPluginResourceAttributesEnvvar(nil, cmd, Plugin{Name: "compose"}, res)
	assert.DeepEqual(t, []string{"OTEL_RESOURCE_ATTRIBUTES=a.b.c=foo,docker.cli.cobra.command_path=docker%20compose,docker.cli.current_context=default"}, env)
}
