package simple

import "github.com/hashicorp/terraform-provider-mock/internal/schema"

var (
	Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bool": {
				Optional: true,
				Type:     schema.Boolean,
			},
			"number": {
				Optional: true,
				Type:     schema.Number,
			},
			"string": {
				Optional: true,
				Type:     schema.String,
			},
			"float": {
				Optional: true,
				Type:     schema.Float,
			},
			"integer": {
				Optional: true,
				Type:     schema.Integer,
			},
		},
	}
)
