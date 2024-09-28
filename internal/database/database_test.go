package database

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	ctx := context.Background()

	srv, _ := New(ctx)
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	ctx := context.Background()

	srv, err := New(ctx)

	if err != nil {
		t.Errorf("unable to connect to db: %s", err.Error())
	}

	if ok, err := srv.Health(ctx); ok == false {
		t.Errorf("service is dead with error: %s", err.Error())
	}
}
