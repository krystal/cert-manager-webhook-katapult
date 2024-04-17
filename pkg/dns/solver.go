package dns

import (
	"fmt"
	"log"
	"strings"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	providerName = "katapult"
)

// ProviderSolver is the struct implementing the webhook.Solver interface
// for Katapult DNS
type ProviderSolver struct {
	client kubernetes.Interface
	Logger *log.Logger
}

// Name is used as the name for this DNS solver when referencing it on the ACME
// Issuer resource
func (p *ProviderSolver) Name() string {
	return providerName
}

// Present is responsible for actually presenting the DNS record with the
// DNS provider.
// This method should tolerate being called multiple times with the same value.
// cert-manager itself will later perform a self check to ensure that the
// solver has correctly configured the DNS provider.
func (p *ProviderSolver) Present(ch *v1alpha1.ChallengeRequest) error {
	sanitizedZone := strings.TrimRight(ch.ResolvedZone, ".")
	sanitizedFQDN := strings.TrimRight(ch.ResolvedFQDN, ".")

	p.Logger.Printf("presented with challenge for %s on zone %s", sanitizedFQDN, sanitizedZone)

	solver, err := p.createKatapultSolver(ch)
	if err != nil {
		return fmt.Errorf("could not create katapult solver: %w", err)
	}

	return solver.Set(sanitizedZone, sanitizedFQDN, ch.Key)
}

// CleanUp should delete the relevant TXT record from the DNS provider console.
// If multiple TXT records exist with the same record name (e.g.
// _acme-challenge.example.com) then **only** the record with the same `key`
// value provided on the ChallengeRequest should be cleaned up.
// This is in order to facilitate multiple DNS validations for the same domain
// concurrently.
func (p *ProviderSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	sanitizedZone := strings.TrimRight(ch.ResolvedZone, ".")
	sanitizedFQDN := strings.TrimRight(ch.ResolvedFQDN, ".")

	p.Logger.Printf("cleaning up %s on zone %s", sanitizedFQDN, sanitizedZone)

	solver, err := p.createKatapultSolver(ch)
	if err != nil {
		return fmt.Errorf("could not create katapult solver: %w", err)
	}

	return solver.Cleanup(sanitizedZone, sanitizedFQDN, ch.Key)
}

// Initialize will be called when the webhook first starts.
// This method can be used to instantiate the webhook, i.e. initialising
// connections or warming up caches.
// Typically, the kubeClientConfig parameter is used to build a Kubernetes
// client that can be used to fetch resources from the Kubernetes API, e.g.
// Secret resources containing credentials used to authenticate with DNS
// provider accounts.
// The stopCh can be used to handle early termination of the webhook, in cases
// where a SIGTERM or similar signal is sent to the webhook process.
func (p *ProviderSolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {

	cl, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		return fmt.Errorf("failed to get kubernetes client: %w", err)
	}

	p.client = cl

	return nil
}
