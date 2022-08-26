package data

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestResource_Symmetry(t *testing.T) {
	testCases := []struct {
		TestCase string
		Resource Resource
	}{
		{
			TestCase: "basic",
			Resource: Resource{
				objectType: tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"number": tftypes.Number,
					},
				},
				Values: map[string]Value{
					"number": {Number: big.NewFloat(0)},
				},
			},
		},
		{
			TestCase: "missing_object",
			Resource: Resource{
				objectType: tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"object": tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"number": tftypes.Number,
							},
						},
					},
				},
				Values: map[string]Value{},
			},
		},
		{
			TestCase: "missing_object_attribute",
			Resource: Resource{
				objectType: tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"object": tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"number": tftypes.Number,
							},
						},
					},
				},
				Values: map[string]Value{
					"object": {
						Object: map[string]Value{},
					},
				},
			},
		},
		{
			TestCase: "missing_list",
			Resource: Resource{
				objectType: tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"list": tftypes.List{
							ElementType: tftypes.Number,
						},
					},
				},
				Values: map[string]Value{},
			},
		},
		{
			TestCase: "empty_list",
			Resource: Resource{
				objectType: tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"list": tftypes.List{
							ElementType: tftypes.Number,
						},
					},
				},
				Values: map[string]Value{
					"list": {
						List: []Value{},
					},
				},
			},
		},
		{
			TestCase: "missing_map",
			Resource: Resource{
				objectType: tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"map": tftypes.Map{
							ElementType: tftypes.Number,
						},
					},
				},
				Values: map[string]Value{},
			},
		},
		{
			TestCase: "missing_set",
			Resource: Resource{
				objectType: tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"set": tftypes.Set{
							ElementType: tftypes.Number,
						},
					},
				},
				Values: map[string]Value{},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.TestCase, func(t *testing.T) {
			CheckSymmetry(t, testCase.Resource)
		})
	}
}

func toJson(t *testing.T, obj Resource) string {
	data, err := json.Marshal(obj)
	require.NoError(t, err)
	return string(data)
}

func CheckResourceEqual(t *testing.T, expected, actual Resource) {
	expectedString := toJson(t, expected)
	actualString := toJson(t, actual)
	require.Equal(t, expectedString, actualString)
}

func CheckSymmetry(t *testing.T, resource Resource) {
	raw, err := resource.ToTerraform5Value()
	require.NoError(t, err)

	value := tftypes.NewValue(resource.objectType, raw)
	actual := Resource{}
	require.NoError(t, actual.FromTerraform5Value(value))

	CheckResourceEqual(t, resource, actual)
}
