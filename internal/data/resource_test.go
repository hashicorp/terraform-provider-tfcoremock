package data

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestResource_symmetry(t *testing.T) {
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
						Object: &map[string]Value{},
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
						List: &[]Value{},
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
			checkSymmetry(t, testCase.Resource)
		})
	}
}

func toJson(t *testing.T, obj Resource) string {
	data, err := json.Marshal(obj)
	if err != nil {
		t.Fatalf("found unexpected error when marshalling json: %v", err)
	}
	return string(data)
}

func checkResourceEqual(t *testing.T, expected, actual Resource) {
	expectedString := toJson(t, expected)
	actualString := toJson(t, actual)
	if expectedString != actualString {
		t.Fatalf("expected did not match actual\nexpected:\n%s\nactual:\n%s", expectedString, actualString)
	}
}

func checkSymmetry(t *testing.T, resource Resource) {
	raw, err := resource.ToTerraform5Value()
	if err != nil {
		t.Fatalf("found unexpected error in ToTerraform5Value(): %v", err)
	}

	value := tftypes.NewValue(resource.objectType, raw)
	actual := Resource{}
	err = actual.FromTerraform5Value(value)
	if err != nil {
		t.Fatalf("found unexpected error in FromTerraform5Value(): %v", err)
	}

	checkResourceEqual(t, resource, actual)
}
