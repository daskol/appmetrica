package appmetrica

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestImporter(t *testing.T) {
	t.Run("Common", func(t *testing.T) {
		imp := NewEventImporter(DeviceID, "mcc", "mnc", "unexisted-col")

		if length := len(imp.header); length != 2 {
			t.Errorf("too many header elements: %d", length)
		}

		imp.Reset()

		if length := len(imp.buffer); length != 0 {
			t.Errorf("buffer wasn't reset: %d", length)
		}

		if length := len(imp.events); length != 0 {
			t.Errorf("events cache wasn't reset: %d", length)
		}
	})

	t.Run("Read", func(t *testing.T) {
		const expectedHeader = `appmetrica_device_id,application_id,event_name,event_timestamp,mcc,mnc`

		imp := NewEventImporter(DeviceID, "mcc", "mnc")
		imp.Import(&ImportEvent{
			ApplicationID:  84126,
			DeviceID:       998,
			EventName:      "SSETestEvent",
			EventTimestamp: time.Now().UTC().Unix(),
			MCC:            250,
			MNC:            25000,
		})

		buffer := new(bytes.Buffer)
		buffer.ReadFrom(imp)

		// Check final state of event importer.
		if imp.state != EndOfEvents {
			t.Errorf("wrong state of importer: %+v", imp.state)
		}

		// Check assertion on read content.
		lines := strings.Split(string(buffer.Bytes()), "\n")

		if length := len(lines); length != 3 && lines[3] == "" {
			t.Fatalf("wrong number of lines were read: %d", length)
		}

		if lines[0] != expectedHeader {
			t.Errorf("wrong header: `%s`", lines[0])
		}

		// Some assertions on record values in columns.
		columns := strings.Split(lines[1], ",")

		if length := len(columns); length != 6 {
			t.Fatalf("worng number of columns in line: %d", length)
		}

		if columns[0] != "998" {
			t.Errorf("wrong device id: %s", columns[0])
		}

		if columns[2] != "SSETestEvent" {
			t.Errorf("wrong event name: %s", columns[2])
		}

		if columns[5] != "25000" {
			t.Errorf("wrong mnc value: %s", columns[5])
		}
	})
}
