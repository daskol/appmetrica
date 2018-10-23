package appmetrica

type Application struct {
	APIKey128             string `json:"api_key128"`
	CreateDate            string `json:"create_date"`
	GDPRAgreementAccepted bool   `json:"gdpr_agreement_accepted"`
	HideAddress           bool   `json:"hide_address"`
	ID                    uint64 `json:"id"`
	Label                 string `json:"label"`
	LabelID               uint64 `json:"label_id"`
	Name                  string `json:"name"`
	OwnerLogin            string `json:"owner_login"`
	Permission            string `json:"permission"`
	PermissionDate        string `json:"permission_date"`
	Status                string `json:"status"`
	TimeZoneName          string `json:"time_zone_name"`
	TimeZoneOffset        int    `json:"time_zone_offset"`
	UID                   uint64 `json:"uid"`
	UseUniversalLinks     bool   `json:"use_universal_links"`
}

type Applications []Application

type Error struct {
	Type    string `json:"error_type"`
	Message string `json:"message"`
}

type Response struct {
	Application  *Application  `json:"application"`
	Applications []Application `json:"applications"`

	Errors       []Error `json:"errors"`
	ErrorCode    int     `json:"code"`
	ErrorMessage string  `json:"message"`
}

type ConnectionType int

const (
	CT_Unknown ConnectionType = iota
	CT_WiFi
	CT_Cell
)

type ImportEvent struct {
	PostAPIKey         string      `json:"post_api_key"`
	ApplicationID      int         `json:"application_id"`
	ProfileID          string      `json:"profile_id"`
	DeviceID           int         `json:"appmetrica_device_id"`
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
