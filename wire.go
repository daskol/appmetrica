package appmetrica

type Application struct {
	APIKey128             string `json:"api_key128,omitempty"`
	CreateDate            string `json:"create_date,omitempty"`
	GDPRAgreementAccepted bool   `json:"gdpr_agreement_accepted,omitempty"`
	HideAddress           bool   `json:"hide_address,omitempty"`
	ID                    uint64 `json:"id,omitempty"`
	Label                 string `json:"label,omitempty"`
	LabelID               uint64 `json:"label_id,omitempty"`
	Name                  string `json:"name,omitempty"`
	OwnerLogin            string `json:"owner_login,omitempty"`
	Permission            string `json:"permission,omitempty"`
	PermissionDate        string `json:"permission_date,omitempty"`
	Status                string `json:"status,omitempty"`
	TimeZoneName          string `json:"time_zone_name,omitempty"`
	TimeZoneOffset        int    `json:"time_zone_offset,omitempty"`
	UID                   uint64 `json:"uid,omitempty"`
	UseUniversalLinks     bool   `json:"use_universal_links,omitempty"`
}

type Applications []Application

type Error struct {
	Type    string `json:"error_type"`
	Message string `json:"message"`
}

type Response struct {
	Application  *Application  `json:"application,omitempty"`
	Applications []Application `json:"applications,omitempty"`

	Errors       []Error `json:"errors,omitempty"`
	ErrorCode    int     `json:"code,omitempty"`
	ErrorMessage string  `json:"message,omitempty"`
}

type ConnectionType int

const (
	CT_Unknown ConnectionType = iota
	CT_WiFi
	CT_Cell
)

type ImportEvent struct {
	ApplicationID      int         `json:"application_id"`
	ProfileID          string      `json:"profile_id"`
	DeviceID           uint64      `json:"appmetrica_device_id"`
	SessionType        string      `json:"session_type,omitempty"`
	IFA                string      `json:"ios_ifa,omitempty"`
	IFV                string      `json:"ios_ifv,omitempty"`
	GoogleAID          string      `json:"google_aid,omitempty"`
	WindowsAID         string      `json:"windows_aid,omitempty"`
	OSName             string      `json:"os_name,omitempty"`
	OSVersion          string      `json:"os_version,omitempty"`
	DeviceManufacturer string      `json:"device_manufacturer,omitempty"`
	DeviceModel        string      `json:"device_model,omitempty"`
	DeviceType         string      `json:"device_type,omitempty"`
	DeviceLocale       string      `json:"device_locale,omitempty"`
	AppVersionName     string      `json:"app_version_name,omitempty"`
	AppPackageName     string      `json:"app_package_name,omitempty"`
	EventName          string      `json:"event_name"`
	EventJSON          interface{} `json:"event_json,omitempty"`
	EventTimestamp     int64       `json:"event_timestamp"`
	ConnectionType     string      `json:"connection_type,omitempty"`
	OperatorName       string      `json:"operator_name,omitempty"`
	MCC                int         `json:"mcc,omitempty"`
	MNC                int         `json:"mnc,omitempty"`
	DeviceIPv6         string      `json:"device_ipv6,omitempty"`
}
