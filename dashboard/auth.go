package dashboard

import (
	"net/http"
	"time"

	"github.com/levigross/grequests"
)

// Status is a global type so we can diff types of status in main pkg
type Status int

// Server struct to read information from private YAML.
type Server struct {
	Name         string
	URL          string
	Username     string
	Password     string
	Seed         string
	IsThirdParty bool
	Session      *grequests.Session
	Headers      *map[string]string
	Cookies      *[]http.Cookie
	Events       []string
	NewAlerts    Status
}

// Creds send credentials to login
type Creds struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Challenge string `json:"challenge"`
}

// Alerts struct as received from the dashboard
type Alerts struct {
	NextPage       string `json:"nextPage"`
	RemainingItems int    `json:"remainingItems"`
	Data           []struct {
		ID               string `json:"id"`
		LocalID          string `json:"localId"`
		EndpointID       string `json:"endpointId"`
		TriggerCondition int    `json:"triggerCondition"`
		Endpoint         struct {
			EndpointID         string    `json:"endpointId"`
			MachineID          string    `json:"machineId"`
			OsType             int       `json:"osType"`
			CPUVendor          int       `json:"cpuVendor"`
			Arch               int       `json:"arch"`
			CPUDescr           string    `json:"cpuDescr"`
			Kernel             string    `json:"kernel"`
			Os                 string    `json:"os"`
			Name               string    `json:"name"`
			Domain             string    `json:"domain"`
			State              int       `json:"state"`
			RegistrationTime   time.Time `json:"registrationTime"`
			AgentVersion       string    `json:"agentVersion"`
			ComponentsVersions []struct {
				Name    string `json:"name"`
				Version string `json:"version"`
				Build   string `json:"build"`
			} `json:"componentsVersions"`
			IsVirtualMachine    bool          `json:"isVirtualMachine"`
			IsDomainController  bool          `json:"isDomainController"`
			IsServer            bool          `json:"isServer"`
			SessionStart        time.Time     `json:"sessionStart"`
			SessionEnd          time.Time     `json:"sessionEnd"`
			LastSeenAt          time.Time     `json:"lastSeenAt"`
			DisconnectionReason int           `json:"disconnectionReason"`
			LocalAddr           string        `json:"localAddr"`
			HvStatus            int           `json:"hvStatus"`
			Macs                []string      `json:"macs"`
			Isolated            bool          `json:"isolated"`
			Connected           bool          `json:"connected"`
			Tags                []interface{} `json:"tags"`
			Groups              []struct {
				EndpointGroupID string `json:"endpointGroupId"`
				Name            string `json:"name"`
			} `json:"groups"`
		} `json:"endpoint"`
		TriggerEvents []struct {
			EventType     int  `json:"eventType"`
			ManuallyAdded bool `json:"manuallyAdded"`
			Process       struct {
				ID         string `json:"id"`
				ParentID   string `json:"parentId"`
				EndpointID string `json:"endpointId"`
				Program    struct {
					Path     string `json:"path"`
					Filename string `json:"filename"`
					Md5      string `json:"md5"`
					Sha1     string `json:"sha1"`
					Sha256   string `json:"sha256"`
					CertInfo struct {
						Signer  string `json:"signer"`
						Issuer  string `json:"issuer"`
						Trusted bool   `json:"trusted"`
						Expired bool   `json:"expired"`
					} `json:"certInfo"`
					Size   int    `json:"size"`
					Arch   string `json:"arch"`
					FsName string `json:"fsName"`
				} `json:"program"`
				User           string    `json:"user"`
				Pid            int       `json:"pid"`
				StartTime      time.Time `json:"startTime"`
				Ppid           int       `json:"ppid"`
				PstartTime     time.Time `json:"pstartTime"`
				UserSID        string    `json:"userSID"`
				PrivilegeLevel string    `json:"privilegeLevel"`
				NoGui          bool      `json:"noGui"`
				LogonID        string    `json:"logonId"`
			} `json:"process"`
			Data struct {
				TargetProcess struct {
					Program struct {
						Path     string `json:"path"`
						Filename string `json:"filename"`
						Md5      string `json:"md5"`
						Sha1     string `json:"sha1"`
						Sha256   string `json:"sha256"`
						CertInfo struct {
							Signer  string `json:"signer"`
							Issuer  string `json:"issuer"`
							Trusted bool   `json:"trusted"`
							Expired bool   `json:"expired"`
						} `json:"certInfo"`
						Size   int    `json:"size"`
						Arch   string `json:"arch"`
						FsName string `json:"fsName"`
					} `json:"program"`
					User           string `json:"user"`
					Pid            int    `json:"pid"`
					StartTime      int64  `json:"startTime"`
					Ppid           int    `json:"ppid"`
					PstartTime     int    `json:"pstartTime"`
					UserSID        string `json:"userSID"`
					PrivilegeLevel string `json:"privilegeLevel"`
					NoGui          bool   `json:"noGui"`
				} `json:"targetProcess"`
				T string `json:"_t"`
			} `json:"data"`
			EndpointID string    `json:"endpointId"`
			ID         string    `json:"id"`
			HappenedAt time.Time `json:"happenedAt"`
			Relevance  int       `json:"relevance"`
			Category   string    `json:"category"`
			ReceivedAt time.Time `json:"receivedAt"`
			Severity   string    `json:"severity"`
			LocalID    string    `json:"localId"`
			Trigger    bool      `json:"trigger"`
		} `json:"triggerEvents"`
		TotalEventCount  int `json:"totalEventCount"`
		ByTypeEventCount []struct {
			Type  int `json:"type"`
			Count int `json:"count"`
		} `json:"byTypeEventCount"`
		Impact            int           `json:"impact"`
		Severity          string        `json:"severity"`
		Closed            bool          `json:"closed"`
		State             int           `json:"state"`
		TerminationReason int           `json:"terminationReason,omitempty"`
		ReceivedAt        time.Time     `json:"receivedAt"`
		HappenedAt        time.Time     `json:"happenedAt"`
		Tags              []interface{} `json:"tags"`
		EndpointState     struct {
			OsType             int      `json:"osType"`
			CPUVendor          int      `json:"cpuVendor"`
			Arch               int      `json:"arch"`
			CPUDescr           string   `json:"cpuDescr"`
			Kernel             string   `json:"kernel"`
			Os                 string   `json:"os"`
			HvStatus           int      `json:"hvStatus"`
			Name               string   `json:"name"`
			Domain             string   `json:"domain"`
			Isolated           bool     `json:"isolated"`
			LocalAddr          string   `json:"localAddr"`
			Macs               []string `json:"macs"`
			ComponentsVersions []struct {
				Name    string `json:"name"`
				Version string `json:"version"`
				Build   string `json:"build"`
			} `json:"componentsVersions"`
			EndpointVersion string        `json:"endpointVersion"`
			Tags            []interface{} `json:"tags"`
			Groups          []struct {
				ID       string        `json:"id"`
				Name     string        `json:"name"`
				Parent   string        `json:"parent,omitempty"`
				Children []interface{} `json:"children"`
				License  struct {
					Limit struct {
						MaxEndpointCount       int `json:"maxEndpointCount"`
						MaxMobileEndpointCount int `json:"maxMobileEndpointCount"`
					} `json:"limit"`
					Expiration time.Time `json:"expiration"`
				} `json:"license,omitempty"`
			} `json:"groups"`
		} `json:"endpointState"`
		Notes string `json:"notes,omitempty"`
		Title string `json:"title,omitempty"`
	} `json:"data"`
}

// GetEndpoints from ReaQta
type GetEndpoints struct {
	Data []struct {
		AgentVersion       string `json:"agentVersion"`
		Arch               int64  `json:"arch"`
		ComponentsVersions []struct {
			Build   string `json:"build"`
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"componentsVersions"`
		Connected           bool   `json:"connected"`
		CPUDescr            string `json:"cpuDescr"`
		CPUVendor           int64  `json:"cpuVendor"`
		DisconnectionReason int64  `json:"disconnectionReason"`
		Domain              string `json:"domain"`
		EndpointID          string `json:"endpointId"`
		Groups              []struct {
			Description     string `json:"description"`
			EndpointGroupID string `json:"endpointGroupId"`
			Name            string `json:"name"`
		} `json:"groups"`
		HvStatus           int64         `json:"hvStatus"`
		IsDomainController bool          `json:"isDomainController"`
		IsServer           bool          `json:"isServer"`
		IsVirtualMachine   bool          `json:"isVirtualMachine"`
		Isolated           bool          `json:"isolated"`
		Kernel             string        `json:"kernel"`
		LastSeenAt         string        `json:"lastSeenAt"`
		LocalAddr          string        `json:"localAddr"`
		MachineID          string        `json:"machineId"`
		Macs               []string      `json:"macs"`
		Name               string        `json:"name"`
		OpenIncidents      int64         `json:"openIncidents"`
		Os                 string        `json:"os"`
		OsType             int64         `json:"osType"`
		RegistrationTime   string        `json:"registrationTime"`
		SessionEnd         string        `json:"sessionEnd"`
		SessionStart       string        `json:"sessionStart"`
		State              int64         `json:"state"`
		Tags               []interface{} `json:"tags"`
	} `json:"data"`
	NextPage       string `json:"nextPage"`
	RemainingItems int64  `json:"remainingItems"`
}
