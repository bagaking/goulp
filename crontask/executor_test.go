package crontask

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/robfig/cron/v3"
)

var testParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

func TestExecutorJobKeysDoesNotIncludeEmptyKey(t *testing.T) {
	executor := NewExecutor(testParser, nil)
	executor.jobs["alpha"] = &Job{}
	executor.jobs["beta"] = &Job{}

	got := executor.JobKeys()
	for _, key := range got {
		if key == "" {
			t.Errorf("JobKeys() = %v, want no empty key", got)
		}
	}

	if len(got) != 2 {
		t.Fatalf("JobKeys() length = %d, want %d; keys = %v", len(got), 2, got)
	}
}

func TestExecutorRegisterRejectsInvalidCrontab(t *testing.T) {
	executor := NewExecutor(testParser, nil)
	jobMeta := JobMeta{
		Key:     "invalid",
		Crontab: "not a cron spec",
		Punch: func(context.Context) error {
			return nil
		},
	}

	err := executor.Register(context.Background(), jobMeta)
	if !errors.Is(err, ErrRegValidateFailed) {
		t.Fatalf("Register(ctx, %+v) error = %v, want %v", jobMeta, err, ErrRegValidateFailed)
	}

	if _, ok := executor.jobs[jobMeta.Key]; ok {
		t.Errorf("Register(ctx, %+v) stored job after validation failure", jobMeta)
	}
	if entries := executor.cron.Entries(); len(entries) != 0 {
		t.Errorf("Register(ctx, %+v) cron entries = %d, want %d", jobMeta, len(entries), 0)
	}
}

func TestExecutorRegisterStoresValidJob(t *testing.T) {
	executor := NewExecutor(testParser, nil)
	jobMeta := JobMeta{
		Key:     "valid",
		Crontab: "@every 1m",
		Punch: func(context.Context) error {
			return nil
		},
	}

	if err := executor.Register(context.Background(), jobMeta); err != nil {
		t.Fatalf("Register(ctx, %+v) error = %v, want nil", jobMeta, err)
	}

	gotKeys := executor.JobKeys()
	wantKeys := []string{jobMeta.Key}
	if !reflect.DeepEqual(gotKeys, wantKeys) {
		t.Errorf("JobKeys() = %v, want %v", gotKeys, wantKeys)
	}
	if entries := executor.cron.Entries(); len(entries) != 1 {
		t.Errorf("Register(ctx, %+v) cron entries = %d, want %d", jobMeta, len(entries), 1)
	}
}
