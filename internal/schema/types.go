// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

type Type string

const (
	Boolean Type = "boolean"
	Float   Type = "float"
	Integer Type = "integer"
	Number  Type = "number"
	String  Type = "string"

	List   Type = "list"
	Map    Type = "map"
	Object Type = "object"
	Set    Type = "set"
)
