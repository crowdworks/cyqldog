package cyqldog

import (
	"os"
	"testing"
)

func TestRenderEnv(t *testing.T) {
	cases := []struct {
		in  []byte
		env map[string]string
		out []byte
	}{
		{
			in: []byte(`host: {{ .TEST_RENDER_ENV_HOST }}\nport: {{ .TEST_RENDER_ENV_PORT }}`),
			env: map[string]string{
				"TEST_RENDER_ENV_HOST": "db.example.com",
				"TEST_RENDER_ENV_PORT": "1234",
			},
			out: []byte(`host: db.example.com\nport: 1234`),
		},
	}
	for _, tc := range cases {
		// Setup environmental variables for testinng.
		for k, v := range tc.env {
			if err := os.Setenv(k, v); err != nil {
				t.Fatalf("failed to set environmental vairables: %s=%s", k, v)
			}
		}

		got, err := renderEnv(tc.in)
		if err != nil {
			t.Errorf("env = %v, renderEnv(%s) returns err = %+v\n", tc.env, tc.in, err)
		}

		if string(got) != string(tc.out) {
			t.Errorf("env = %v, renderEnv(%s) = %s, want = %s", tc.env, tc.in, got, tc.out)
		}

		// Reset for each test case.
		for k := range tc.env {
			os.Unsetenv(k)
		}
	}
}

func TestRenderEnvError(t *testing.T) {
	cases := []struct {
		in []byte
	}{
		{
			in: []byte(`{{ broken template!!`),
		},
	}

	for _, tc := range cases {
		_, err := renderEnv(tc.in)
		if err == nil {
			t.Errorf("expected renderEnv(%s) returns err, but err = nil\n", tc.in)
		}
	}
}
