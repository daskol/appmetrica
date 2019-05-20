package appmetrica

import (
	"errors"
	"io"
	"strconv"
	"strings"
)

const requiredCols = `application_id,event_name,event_timestamp`

var allowedHeader map[string]struct{}

func init() {
	var header = []string{
		"app_package_name", "app_version_name", "connection_type",
		"device_ipv6", "device_locale", "device_manufacturer", "device_model",
		"device_type", "event_json", "google_aid", "ios_ifa", "ios_ifv", "mcc",
		"mnc", "operator_name", "os_name", "os_version", "session_type",
		"windows_aid",
	}

	allowedHeader = make(map[string]struct{}, len(header))

	for _, label := range header {
		allowedHeader[label] = struct{}{}
	}
}

// EventIdentifierType encodes identifier which could be used in event import.
type EventIdentifierType int

const (
	DeviceID EventIdentifierType = iota
	ProfileID
)

// EventImporterState codes inner state of reader in event importer.
type EventImporterState int

const (
	Initial EventImporterState = iota
	HeaderWriting
	EventWriting
	EndOfEvents
)

// EventImporter
type EventImporter struct {
	err    error
	state  EventImporterState
	kind   EventIdentifierType
	offset int
	buffer []byte // line buffer
	header []string
	events []*ImportEvent
}

// NewEventImporter creates new instance of EventImporter. Developer has to
// specify type of user identifer which is either DeviceID or ProfileID and
// optional list of attribution columns. The list of mandatory columns should
// not to be specified since EventImporter adds them automatically. Mandatory
// columns are ApplicationID, either ProfileID or DeviceId, EventName, and
// EventTimestamp.
func NewEventImporter(kind EventIdentifierType, header ...string) *EventImporter {
	var importer = new(EventImporter)
	importer.state = Initial
	importer.buffer = make([]byte, 0, 4096) // 4kb
	importer.SetEventIdentifierType(kind)
	importer.SetHeader(header...)
	return importer
}

func (e *EventImporter) Import(events ...*ImportEvent) {
	if e.state != EndOfEvents {
		e.events = append(e.events, events...)
	} else {
		var msg = "failed to add event(s) to completed importer"
		e.err = errors.New(prefix + msg)
	}
}

func (e *EventImporter) ImportOne(event *ImportEvent) {
	e.Import(event)
}

func (e *EventImporter) ImportMany(events []*ImportEvent) {
	e.Import(events...)
}

func (e *EventImporter) Reset() {
	e.err = nil
	e.state = Initial
	e.offset = 0
	e.buffer = e.buffer[:0]
	e.header = e.header[:0]
	e.events = e.events[:0]
}

func (e *EventImporter) SetEventIdentifierType(kind EventIdentifierType) {
	if e.state == Initial {
		e.kind = kind
	} else {
		var msg = "failed to set identifier type during reading"
		e.err = errors.New(prefix + msg)
	}
}

func (e *EventImporter) SetHeader(header ...string) {
	if e.state != Initial {
		var msg = "failed to set identifier type during reading"
		e.err = errors.New(prefix + msg)
		return
	}

	e.header = make([]string, 0, len(header))

	for _, col := range header {
		if _, ok := allowedHeader[col]; ok {
			e.header = append(e.header, col)
		}
	}

}

func (e *EventImporter) Read(buffer []byte) (int, error) {
	// Short circut for EOF state.
	if e.state == EndOfEvents {
		return 0, io.EOF
	}

	// If there is nothing to flush or no events to read then EOF.
	if len(e.buffer) == 0 && len(e.events) == 0 {
		e.state = EndOfEvents
		return 0, io.EOF
	}

	switch e.state {
	case Initial:
		e.state = HeaderWriting
		fallthrough
	case HeaderWriting:
		return e.readHeader(buffer)
	case EventWriting:
		return e.readEvents(buffer)
	default:
		panic(prefix + "unexpected execution branch")
	}
}

func (e *EventImporter) encodeEvent(event *ImportEvent) {
	// Encode identifier device identifier.
	if e.kind == DeviceID {
		e.buffer = strconv.AppendUint(e.buffer, event.DeviceID, 10)
	} else {
		e.buffer = append(e.buffer, event.ProfileID...)
	}

	e.encodeRequiredColumns(event)
	e.encodeOptionalColumns(event)
	e.buffer = append(e.buffer, "\n"...)
}

func (e *EventImporter) encodeOptionalColumns(event *ImportEvent) {
	for _, label := range e.header {
		// Append column delimiter.
		e.buffer = append(e.buffer, ","...)

		// Append column value.
		switch label {
		case "app_package_name":
			e.buffer = append(e.buffer, event.AppPackageName...)
		case "app_version_name":
			e.buffer = append(e.buffer, event.AppVersionName...)
		case "connection_type":
			e.buffer = append(e.buffer, event.ConnectionType...)
		case "device_ipv6":
			e.buffer = append(e.buffer, event.DeviceIPv6...)
		case "device_locale":
			e.buffer = append(e.buffer, event.DeviceLocale...)
		case "device_manufacturer":
			e.buffer = append(e.buffer, event.DeviceManufacturer...)
		case "device_model":
			e.buffer = append(e.buffer, event.DeviceModel...)
		case "device_type":
			e.buffer = append(e.buffer, event.DeviceType...)
		case "event_json":
			println(prefix + "field `event_json` in event is not supported")
		case "google_aid":
			e.buffer = append(e.buffer, event.GoogleAID...)
		case "ios_ifa":
			e.buffer = append(e.buffer, event.IFA...)
		case "ios_ifv":
			e.buffer = append(e.buffer, event.IFV...)
		case "mcc":
			e.buffer = strconv.AppendInt(e.buffer, int64(event.MCC), 10)
		case "mnc":
			e.buffer = strconv.AppendInt(e.buffer, int64(event.MNC), 10)
		case "operator_name":
			e.buffer = append(e.buffer, event.OperatorName...)
		case "os_name":
			e.buffer = append(e.buffer, event.OSName...)
		case "os_version":
			e.buffer = append(e.buffer, event.OSVersion...)
		case "session_type":
			e.buffer = append(e.buffer, event.SessionType...)
		case "windows_aid":
			e.buffer = append(e.buffer, event.WindowsAID...)
		default:
			panic(prefix + "unexpected execution branch")
		}
	}
}

func (e *EventImporter) encodeRequiredColumns(event *ImportEvent) {
	e.buffer = append(e.buffer, ","...)
	e.buffer = strconv.AppendInt(e.buffer, int64(event.ApplicationID), 10)
	e.buffer = append(e.buffer, ","...)
	e.buffer = append(e.buffer, event.EventName...)
	e.buffer = append(e.buffer, ","...)
	e.buffer = strconv.AppendInt(e.buffer, event.EventTimestamp, 10)
}

func (e *EventImporter) readHeader(buffer []byte) (int, error) {
	// Initialize buffer with header line.
	if len(e.buffer) == 0 {
		if e.kind == DeviceID {
			e.buffer = append(e.buffer, `appmetrica_device_id,`...)
		} else {
			e.buffer = append(e.buffer, `profile_id,`...)
		}

		e.buffer = append(e.buffer, requiredCols...)

		if len(e.header) > 0 {
			e.buffer = append(e.buffer, ","...)
			e.buffer = append(e.buffer, strings.Join(e.header, ",")...)
		}

		e.buffer = append(e.buffer, "\n"...)
	}

	// Write rest of header line to output buffer.
	var read = copy(buffer, e.buffer[e.offset:])
	e.offset += read

	// If we wrote all header line then change state.
	if e.offset == len(e.buffer) {
		e.offset = 0
		e.buffer = e.buffer[:e.offset]
		e.state = EventWriting
	}

	// There are no any errors but we force similarity to Reader interface.
	return read, nil
}

func (e *EventImporter) readEvents(buffer []byte) (int, error) {
	// Initialize buffer with import event line.
	if len(e.buffer) == 0 {
		var last = len(e.events) - 1 // Index of last event.
		var event = e.events[last]   // Get last element.
		e.events = e.events[:last]   // Remove last element.
		e.encodeEvent(event)
	}

	// Write rest of header line to output buffer.
	var read = copy(buffer, e.buffer[e.offset:])
	e.offset += read

	// If we wrote all header line then change state.
	if e.offset == len(e.buffer) {
		e.offset = 0
		e.buffer = e.buffer[:e.offset]
	}

	// If there is no events by now then transit to termination state.
	if len(e.events) == 0 {
		e.state = EndOfEvents
	}

	// There are no any errors but we force similarity to Reader interface.
	return read, nil
}
