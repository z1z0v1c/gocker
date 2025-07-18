package registry

const URL = "registry-1.docker.io"

// AuthResponse represents the token response from the Docker Registry auth API.
type AuthResponse struct {
	Token string
}

// Manifest represents a platform-specific image manifest (schema v2).
// It includes metadata about the config blob and image layers.
type Manifest struct {
	SchemaVersion int
	MediaType     string
	Config        struct {
		MediaType string
		Size      int
		Digest    string
	}
	Layers []struct {
		MediaType string
		Size      int
		Digest    string
	}
}

// ManifestIndex represents a manifest list (multi-platform index).
// It maps platforms (e.g., linux/amd64) to specific manifest digests.
type ManifestIndex struct {
	SchemaVersion int
	MediaType     string
	Manifests     []struct {
		MediaType string
		Digest    string
		Platform  struct {
			Architecture string
			OS           string
		}
	}
}

type Config struct {
	Hostname     string              `json:"Hostname,omitempty"`
	Domainname   string              `json:"Domainname,omitempty"`
	User         string              `json:"User,omitempty"`
	AttachStdin  bool                `json:"AttachStdin,omitempty"`
	AttachStdout bool                `json:"AttachStdout,omitempty"`
	AttachStderr bool                `json:"AttachStderr,omitempty"`
	Tty          bool                `json:"Tty,omitempty"`
	OpenStdin    bool                `json:"OpenStdin,omitempty"`
	StdinOnce    bool                `json:"StdinOnce,omitempty"`
	Env          []string            `json:"Env,omitempty"`
	Cmd          []string            `json:"Cmd,omitempty"`
	Image        string              `json:"Image,omitempty"`
	Volumes      map[string]struct{} `json:"Volumes,omitempty"`
	WorkingDir   string              `json:"WorkingDir,omitempty"`
	Entrypoint   []string            `json:"Entrypoint,omitempty"`
	OnBuild      []string            `json:"OnBuild,omitempty"`
	Labels       map[string]string   `json:"Labels,omitempty"`
}

// ImageConfig represents the full image configuration.
type ImageConfig struct {
	Architecture    string `json:"architecture,omitempty"`
	Config          Config
	Container       string `json:"container,omitempty"`
	ContainerConfig Config `json:"container_config"`
	Created         string `json:"created,omitempty"`
	DockerVersion   string `json:"docker_version,omitempty"`
	History         []struct {
		Created    string `json:"created,omitempty"`
		CreatedBy  string `json:"created_by,omitempty"`
		EmptyLayer bool   `json:"empty_layer,omitempty"`
	}
	Os     string `json:"os,omitempty"`
	Rootfs struct {
		Type    string   `json:"type,omitempty"`
		DiffIds []string `json:"diff_ids,omitempty"`
	}
}
