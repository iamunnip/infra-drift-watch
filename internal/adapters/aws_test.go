package adapters

import (
	"context"
	"os"
	"slices"
	"testing"

	"github.com/iamunnip/infra-drift-watch/internal/config"
)

func TestAWSSelectionReachesAllToolCommands(t *testing.T) {
	for _, adapter := range []string{"opentofu", "terraform", "kubernetes"} {
		t.Run(adapter, func(t *testing.T) {
			profile, region := "default", "us-east-1"
			calls := 0
			run := func(ctx context.Context, dir, exe string, args, env []string) ([]byte, int, error) {
				calls++
				for _, want := range []string{"AWS_PROFILE=default", "AWS_DEFAULT_PROFILE=default", "AWS_REGION=us-east-1", "AWS_DEFAULT_REGION=us-east-1"} {
					if !slices.Contains(env, want) {
						t.Fatalf("missing %s", want)
					}
				}
				if args[0] == "show" {
					return []byte(`{"format_version":"1.2"}`), 0, nil
				}
				return nil, 0, nil
			}
			target := config.Target{ID: "aws", Adapter: adapter, Path: ".", Context: "test", Manifests: []string{"demo.yaml"}, Settings: config.Settings{Timeout: "1m", AWS: config.AWS{Profile: &profile, Region: &region}}}
			if result := (Scanner{Run: run}).Scan(context.Background(), target); result.Status != "clean" {
				t.Fatal(result)
			}
			want := 2
			if adapter == "kubernetes" {
				want = 1
			}
			if calls != want {
				t.Fatalf("calls=%d", calls)
			}
		})
	}
}

func TestAWSSelectorsAreChildOnlyAndPreserveOIDCCredentials(t *testing.T) {
	t.Setenv("AWS_PROFILE", "ambient")
	t.Setenv("AWS_DEFAULT_PROFILE", "legacy")
	t.Setenv("AWS_REGION", "eu-west-2")
	t.Setenv("AWS_DEFAULT_REGION", "eu-north-1")
	t.Setenv("AWS_ACCESS_KEY_ID", "test-access")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test-secret")
	t.Setenv("AWS_SESSION_TOKEN", "test-session")
	profile, region := "default", "us-east-1"
	script := `test "$AWS_ACCESS_KEY_ID" = test-access && test "$AWS_SECRET_ACCESS_KEY" = test-secret && test "$AWS_SESSION_TOKEN" = test-session || exit 9
printf '%s\n' "$AWS_PROFILE" "$AWS_DEFAULT_PROFILE" "$AWS_REGION" "$AWS_DEFAULT_REGION"`
	out, code, err := Execute(context.Background(), "", "sh", []string{"-c", script}, awsEnvironment(config.AWS{Profile: &profile, Region: &region}))
	if err != nil || code != 0 || string(out) != "default\ndefault\nus-east-1\nus-east-1\n" {
		t.Fatalf("selectors=%q code=%d err=%v", out, code, err)
	}
	profile = ""
	out, code, err = Execute(context.Background(), "", "sh", []string{"-c", script}, awsEnvironment(config.AWS{Profile: &profile}))
	if err != nil || code != 0 || string(out) != "\n\neu-west-2\neu-north-1\n" {
		t.Fatalf("clear=%q code=%d err=%v", out, code, err)
	}
	if os.Getenv("AWS_PROFILE") != "ambient" || os.Getenv("AWS_REGION") != "eu-west-2" {
		t.Fatal("host selection was changed")
	}
}
