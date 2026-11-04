package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/iamunnip/infra-drift-watch/internal/state"
)

func History(dir, target, format string, limit int, out io.Writer) error {
	if format != "terminal" && format != "json" {
		return errors.New("format must be terminal or json")
	}
	records, err := state.ReadHistory(dir, target, limit)
	if err != nil {
		return err
	}
	if format == "json" {
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(struct {
			SchemaVersion string         `json:"schema_version"`
			Records       []state.Record `json:"records"`
		}{"1", records}); err != nil {
			return errors.New("cannot write history output")
		}
		return nil
	}
	for _, record := range records {
		if _, err := fmt.Fprintf(out, "%s %s [%s]: %s (%d findings; %d delivery errors)\n", record.At.UTC().Format(time.RFC3339), record.Result.Target, record.Result.Adapter, record.Result.Status, len(record.Result.Findings), len(record.DeliveryErrors)); err != nil {
			return errors.New("cannot write history output")
		}
	}
	return nil
}
