package main

import "context"

func main() {
	cache, err := Open("/Users/liamcervante/hashicorp/terraform-provider-mock/internal/schema/remote")
	if err != nil {
		panic(err)
	}

	if err := cache.InstallProvider("mock", "https://github.com/liamcervante/terraform-provider-mock/releases/download/v0.2.0/terraform-provider-mock_0.2.0_darwin_arm64.zip"); err != nil {
		panic(err)
	}

	if err := cache.InstallProvider("aws", "https://releases.hashicorp.com/terraform-provider-aws/4.40.0/terraform-provider-aws_4.40.0_darwin_arm64.zip"); err != nil {
		panic(err)
	}

	_, _, err = cache.GetSchema(context.Background(), "aws")
	if err != nil {
		panic(err)
	}
}
