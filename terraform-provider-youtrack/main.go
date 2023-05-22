package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/elinevich/terraform-provider/terraform-provider-youtrack"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: youtrack.Provider})
}
