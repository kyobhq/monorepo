package broker

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func mustStartDragonflyContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "docker.dragonflydb.io/dragonflydb/dragonfly",
		ExposedPorts: []string{"6379/tcp"},
		Env: map[string]string{
			"DFLY_requirepass": os.Getenv("DRAGONFLY_DB_PASSWORD"),
		},
		WaitingFor: wait.ForLog("Starting dragonfly").WithStartupTimeout(10 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "6379/tcp")
	if err != nil {
		return container.Terminate, err
	}

	port = mappedPort.Port()
	password = os.Getenv("DRAGONFLY_DB_PASSWORD")

	return container.Terminate, nil
}

func TestMain(m *testing.M) {
	teardown, err := mustStartDragonflyContainer()
	if err != nil {
		log.Fatalf("could not start dragonfly container: %v", err)
	}

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown dragonfly container: %v", err)
	}
}

func TestNew(t *testing.T) {
	srv := New()
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	srv := New()

	stats := srv.Health()

	if stats["status"] != "up" {
		t.Fatalf("expected status to be up, got %s", stats["status"])
	}

	if _, ok := stats["error"]; ok {
		t.Fatalf("expected error not to be present")
	}

	if stats["message"] != "It's healthy" {
		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	}
}

func TestClose(t *testing.T) {
	srv := New()

	if srv.Close() != nil {
		t.Fatalf("expected Close() to return nil")
	}
}
