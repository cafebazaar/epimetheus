package epimetheus

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/cactus/go-statsd-client/statsd/statsdtest"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

func TestTimerRunWithErrorRecordsOKLabel(t *testing.T) {
	workError := errors.New("work failed")
	tests := []struct {
		name       string
		labelNames []string
		values     []string
		workError  error
		wantValues []string
	}{
		{"only ok success", []string{"ok"}, []string{"pending"}, nil, []string{"true"}},
		{"only ok error", []string{"ok"}, []string{"pending"}, workError, []string{"false"}},
		{"first ok success", []string{"ok", "service"}, []string{"pending", "api"}, nil, []string{"true", "api"}},
		{"first ok error", []string{"ok", "service"}, []string{"pending", "api"}, workError, []string{"false", "api"}},
		{"later ok success", []string{"service", "ok"}, []string{"api", "pending"}, nil, []string{"api", "true"}},
		{"later ok error", []string{"service", "ok"}, []string{"api", "pending"}, workError, []string{"api", "false"}},
		{"without ok success", []string{"service"}, []string{"api"}, nil, []string{"api"}},
		{"without ok error", []string{"service"}, []string{"api"}, workError, []string{"api"}},
		{"no labels success", nil, nil, nil, nil},
		{"no labels error", nil, nil, workError, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := viper.New()
			config.Set("stats.namespace", "test")
			config.Set("stats.system-name", "timer")
			timer := NewEpimetheusServer(config).NewTimer("work", tt.labelNames)

			recorder := statsdtest.NewRecordingSender()
			client, err := statsd.NewClientWithSender(recorder, "")
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close()
			timer.client = &client

			registry := prometheus.NewRegistry()
			if err := registry.Register(timer.watcher); err != nil {
				t.Fatal(err)
			}

			calls := 0
			values := append([]string(nil), tt.values...)
			gotError := timer.RunWithError(func() error {
				calls++
				return tt.workError
			}, values...)
			if gotError != tt.workError {
				t.Errorf("RunWithError returned %v, want %v", gotError, tt.workError)
			}
			if calls != 1 {
				t.Errorf("work called %d times, want 1", calls)
			}

			families, err := registry.Gather()
			if err != nil {
				t.Fatal(err)
			}
			if len(families) != 1 || len(families[0].GetMetric()) != 1 {
				t.Fatalf("expected one Prometheus metric, got %v", families)
			}
			metric := families[0].GetMetric()[0]
			if metric.GetHistogram().GetSampleCount() != 1 {
				t.Errorf("histogram sample count = %d, want 1", metric.GetHistogram().GetSampleCount())
			}
			gotLabels := make(map[string]string)
			for _, label := range metric.GetLabel() {
				gotLabels[label.GetName()] = label.GetValue()
			}
			wantLabels := make(map[string]string)
			for i, name := range tt.labelNames {
				wantLabels[name] = tt.wantValues[i]
			}
			if !reflect.DeepEqual(gotLabels, wantLabels) {
				t.Errorf("Prometheus labels = %v, want %v", gotLabels, wantLabels)
			}

			sent := recorder.GetSent()
			if len(sent) != 2 {
				t.Fatalf("expected Statsd timing and count, got %v", sent)
			}
			wantName := "test.timer.work." + strings.Join(tt.wantValues, ".")
			for _, stat := range sent {
				if !stat.Parsed || stat.Stat != wantName {
					t.Errorf("Statsd metric name = %q (parsed: %t), want %q", stat.Stat, stat.Parsed, wantName)
				}
			}
			if sent[0].Tag != "ms" || sent[1].Tag != "c" || sent[1].Value != "1" {
				t.Errorf("expected Statsd timing and one count, got %v", sent)
			}
		})
	}
}
