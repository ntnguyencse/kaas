package git

import (
	"context"

	github "github.com/google/go-github/v50/github"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const GH_Token = "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"

func Test() {
	logger := log.Log.WithName("Git modules")
	logger.Info("Using git module...")
	ctx := context.Background()
	client := github.NewTokenClient(ctx, GH_Token)
	user, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		logger.Error(err, "Error when auth with Github server")
	}
	logger.Info("Rate: ", "Ratelimit: ", resp.Rate)
	if !resp.TokenExpiration.IsZero() {
		logger.Info("Token Expiration: ", "Expiration: ", resp.TokenExpiration)
	}

	logger.Info("\n", github.Stringify(user))
}
