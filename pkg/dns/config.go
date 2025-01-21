package dns

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	katapultsolver "github.com/krystal/go-katapult-dns-acme-solver/solver"
	v1 "k8s.io/api/core/v1"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ProviderConfig represents the config used for Katapult DNS
type ProviderConfig struct {
	APIToken *v1.SecretKeySelector `json:"apiToken,omitempty"`
	Host     string                `json:"host,omitempty"`
}

// Load the configuration
func loadConfig(cfgJSON *extapi.JSON) (ProviderConfig, error) {
	providerConfig := ProviderConfig{}

	if cfgJSON == nil {
		return providerConfig, nil
	}

	err := json.Unmarshal(cfgJSON.Raw, &providerConfig)
	if err != nil {
		return providerConfig, fmt.Errorf("error decoding solver config: %w", err)
	}

	return providerConfig, nil
}

// Create a solver for katapult with the API token as configured
// within the issuer.
func (p *ProviderSolver) createKatapultSolver(ch *v1alpha1.ChallengeRequest) (*katapultsolver.Solver, error) {
	config, err := loadConfig(ch.Config)
	if err != nil {
		return nil, fmt.Errorf("could not load config: %w", err)
	}

	secret, err := p.client.CoreV1().Secrets(ch.ResourceNamespace).Get(
		context.Background(),
		config.APIToken.Name,
		metav1.GetOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("could not secret %s: %w", config.APIToken.Name, err)
	}

	token, ok := secret.Data[config.APIToken.Key]
	if !ok {
		return nil, fmt.Errorf("could not find key %s in secret %s", config.APIToken.Key, config.APIToken.Name)
	}

	solver := katapultsolver.NewSolverWithHost(config.Host, string(token), p.Logger)
	return solver, nil
}
