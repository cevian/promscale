// This file and its contents are licensed under the Apache License 2.0.
// Please see the included NOTICE for copyright information and
// LICENSE for a copy of the license

package testhelpers

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

var (
	useDocker    = flag.Bool("use-docker", true, "start database using a docker container")
	testDatabase = flag.String("database", "tmp_db_timescale_migrate_test", "database to run integration tests on")
)

const extensionState = timescaleBit | promscaleBit | postgres12Bit

func TestPGConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	db, err := pgx.Connect(context.Background(), PgConnectURL(defaultDB, Superuser))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close(context.Background())
	var res int
	err = db.QueryRow(context.Background(), "SELECT 1").Scan(&res)
	if err != nil {
		t.Fatal(err)
	}
	if res != 1 {
		t.Errorf("Res is not 1 but %d", res)
	}
}

func TestOtelCollectorConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping checking Jaeger connection")
	}
	jContainer, err := StartJaegerContainer(true)
	require.NoError(t, err)

	// Start Otel collector.
	otelContainer, _, _, err := StartOtelCollectorContainer(fmt.Sprintf("%s:%s", jContainer.ContainerIp, jContainer.GrpcReceivingPort.Port()), true)
	require.NoError(t, err)
	require.NoError(t, jContainer.Container.Terminate(context.Background()))
	require.NoError(t, otelContainer.Terminate(context.Background()))
}

func TestJaegerConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping checking Jaeger connection")
	}
	jContainer, err := StartJaegerContainer(true)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, jContainer.Container.Terminate(context.Background()))
	}()

	const servicesEndpoint = "api/services"

	resp, err := http.Get(fmt.Sprintf("http://%s:%s/%s", jContainer.Host, jContainer.UIPort.Port(), servicesEndpoint))
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)

	bSlice, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Greater(t, len(bSlice), 0)
}

func TestWithDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	WithDB(t, *testDatabase, Superuser, false, extensionState, func(db *pgxpool.Pool, t testing.TB, connectURL string) {
		var res int
		err := db.QueryRow(context.Background(), "SELECT 1").Scan(&res)
		if err != nil {
			t.Fatal(err)
		}
		if res != 1 {
			t.Errorf("Res is not 1 but %d", res)
		}
	})
}

func runMain(m *testing.M) int {
	flag.Parse()
	ctx := context.Background()
	if !testing.Short() && *useDocker {
		_, closer, err := StartPGContainer(ctx, extensionState, "", false)
		if err != nil {
			fmt.Println("Error setting up container", err)
			os.Exit(1)
		}

		tmpDir := ""
		if runtime.GOOS == "darwin" {
			// Docker on Mac lacks access to default os tmp dir - "/var/folders/random_number"
			// so switch to cross-user tmp dir
			tmpDir = "/tmp"
		}
		path, err := ioutil.TempDir(tmpDir, "prom_test")
		if err != nil {
			fmt.Println("Error getting temp dir for Prometheus storage", err)
			os.Exit(1)
		}
		err = os.Mkdir(filepath.Join(path, "wal"), 0777)
		if err != nil {
			fmt.Println("Error getting temp dir for Prometheus storage", err)
			os.Exit(1)
		}
		promContainer, _, _, err := StartPromContainer(path, ctx)
		if err != nil {
			fmt.Println("Error setting up container", err)
			os.Exit(1)
		}
		defer func() {
			_ = closer.Close()
			err = promContainer.Terminate(ctx)
			if err != nil {
				panic(err)
			}
		}()
	}
	return m.Run()
}
func TestMain(m *testing.M) {
	os.Exit(runMain(m))
}
