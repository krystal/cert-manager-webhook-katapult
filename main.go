package main

import (
	"log"
	"os"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	"github.com/krystal/cert-manager-webhook-katapult/pkg/dns"
)

// GroupName is the name under which the webhook will be available
var GroupName = os.Getenv("GROUP_NAME")

func main() {
	if GroupName == "" {
		panic("GROUP_NAME must be specified")
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)

	cmd.RunWebhookServer(GroupName,
		&dns.ProviderSolver{Logger: logger},
	)
}
