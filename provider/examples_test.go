package provider_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: fix this test
// func TestExamples(t *testing.T) {
// 	t.Parallel()
//
// 	t.Run("gigo_parameter", func(t *testing.T) {
// 		resource.Test(t, resource.TestCase{
// 			Providers: map[string]*schema.Provider{
// 				"gigo": provider.New(),
// 			},
// 			IsUnitTest: true,
// 			Steps: []resource.TestStep{
// 				{
// 					Config: mustReadFile(t, "../examples/provider/provider.tf"),
// 				},
// 			},
// 		})
// 	})
// }

func mustReadFile(t *testing.T, path string) string {
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	return string(content)
}
