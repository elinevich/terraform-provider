package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/boxboat/terraform-provider-youtrack/youtrack"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: youtrack.Provider})
}
