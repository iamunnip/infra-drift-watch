package adapters

import "github.com/iamunnip/infra-drift-watch/internal/config"

// awsEnvironment scopes selectors to each subprocess without changing the host.
// Both region variables are set because the invoked tools use different AWS SDKs.
func awsEnvironment(settings config.AWS) []string {
	env := []string{}
	if settings.Profile != nil {
		env = append(env, "AWS_PROFILE="+*settings.Profile, "AWS_DEFAULT_PROFILE="+*settings.Profile)
	}
	if settings.Region != nil {
		env = append(env, "AWS_REGION="+*settings.Region, "AWS_DEFAULT_REGION="+*settings.Region)
	}
	return env
}
