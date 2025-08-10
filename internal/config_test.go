package internal

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	type want struct {
		cfg Config
	}

	type test struct {
		name     string
		flags    []string
		envSetup func()
		want     want
	}

	tests := []test{
		{
			name:     "Case 1 - default values",
			flags:    []string{"test"},
			envSetup: nil,
			want: want{
				cfg: Config{
					Addr:        defaultAddr,
					Port:        defaultPort,
					DBDSN:       defaultDBStr,
					MigratePath: "migrations",
				},
			},
		},
		{
			name: "Case 2 - custom falgs",
			flags: []string{"test",
				"-addr", "1.1.1.1",
				"-port", "7777",
				"-db", "db_dsn_test",
				"-migrate", "ne-mig",
			},
			envSetup: nil,
			want: want{
				cfg: Config{
					Addr:        "1.1.1.1",
					Port:        7777,
					DBDSN:       "db_dsn_test",
					MigratePath: "ne-mig",
				},
			},
		},
		{
			name: "Case 3 - custom falgs and defaults",
			flags: []string{"test",
				"-addr", "1.1.1.1",
				"-migrate", "ne-mig",
			},
			envSetup: nil,
			want: want{
				cfg: Config{
					Addr:        "1.1.1.1",
					Port:        defaultPort,
					DBDSN:       defaultDBStr,
					MigratePath: "ne-mig",
				},
			},
		},
		{
			name:  "Case 4 - envs variables",
			flags: []string{"test"},
			envSetup: func() {
				t.Setenv("DB_DSN", "db_dsn_test")
				t.Setenv("MIGRATE_PATH", "ne-mig")
			},
			want: want{
				cfg: Config{
					Addr:        defaultAddr,
					Port:        defaultPort,
					DBDSN:       "db_dsn_test",
					MigratePath: "ne-mig",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			os.Args = tc.flags
			flag.CommandLine = flag.NewFlagSet(tc.flags[0], flag.ExitOnError)
			if tc.envSetup != nil {
				tc.envSetup()
			}

			testCfg := ReadConfig()

			assert.Equal(t, tc.want.cfg, *testCfg)
		})
	}
}
