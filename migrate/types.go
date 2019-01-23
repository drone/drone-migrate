package migrate

import (
	"encoding/base64"
	"encoding/json"
)

type (
	// UserV0 is a Drone 0.x user.
	UserV0 struct {
		ID     int64  `meddler:"user_id"`
		Login  string `meddler:"user_login"`
		Token  string `meddler:"user_token"`
		Secret string `meddler:"user_secret"`
		Expiry int64  `meddler:"user_expiry"`
		Email  string `meddler:"user_email"`
		Avatar string `meddler:"user_avatar"`
		Active bool   `meddler:"user_active"`
		Admin  bool   `meddler:"user_admin"`
		Synced int64  `meddler:"user_synced"`
		Hash   string `meddler:"user_hash"`
	}

	// UserV1 is a Drone 1.x user.
	UserV1 struct {
		ID        int64  `meddler:"user_id"`
		Login     string `meddler:"user_login"`
		Email     string `meddler:"user_email"`
		Machine   bool   `meddler:"user_machine"`
		Admin     bool   `meddler:"user_admin"`
		Active    bool   `meddler:"user_active"`
		Avatar    string `meddler:"user_avatar"`
		Syncing   bool   `meddler:"user_syncing"`
		Synced    int64  `meddler:"user_synced"`
		Created   int64  `meddler:"user_created"`
		Updated   int64  `meddler:"user_updated"`
		LastLogin int64  `meddler:"user_last_login"`
		Token     string `meddler:"user_oauth_token"`
		Refresh   string `meddler:"user_oauth_refresh"`
		Expiry    int64  `meddler:"user_oauth_expiry"`
		Hash      string `meddler:"user_hash"`
	}

	// RepoV0 is a Drone 0.x repository.
	RepoV0 struct {
		ID          int64  `meddler:"repo_id"`
		UserID      int64  `meddler:"repo_user_id"`
		Owner       string `meddler:"repo_owner"`
		Name        string `meddler:"repo_name"`
		FullName    string `meddler:"repo_full_name"`
		Avatar      string `meddler:"repo_avatar"`
		Link        string `meddler:"repo_link"`
		Kind        string `meddler:"repo_scm"`
		Clone       string `meddler:"repo_clone"`
		Branch      string `meddler:"repo_branch"`
		Timeout     int64  `meddler:"repo_timeout"`
		Visibility  string `meddler:"repo_visibility"`
		IsPrivate   bool   `meddler:"repo_private"`
		IsTrusted   bool   `meddler:"repo_trusted"`
		IsGated     bool   `meddler:"repo_gated"`
		IsActive    bool   `meddler:"repo_active"`
		AllowPull   bool   `meddler:"repo_allow_pr"`
		AllowPush   bool   `meddler:"repo_allow_push"`
		AllowDeploy bool   `meddler:"repo_allow_deploys"`
		AllowTag    bool   `meddler:"repo_allow_tags"`
		Counter     int    `meddler:"repo_counter"`
		Config      string `meddler:"repo_config_path"`
		Hash        string `meddler:"repo_hash"`
	}

	// RepoV1 is a Drone 1.x repository.
	RepoV1 struct {
		ID         int64  `meddler:"repo_id"`
		UID        string `meddler:"repo_uid"`
		UserID     int64  `meddler:"repo_user_id"`
		Namespace  string `meddler:"repo_namespace"`
		Name       string `meddler:"repo_name"`
		Slug       string `meddler:"repo_slug"`
		SCM        string `meddler:"repo_scm"`
		HTTPURL    string `meddler:"repo_clone_url"`
		SSHURL     string `meddler:"repo_ssh_url"`
		Link       string `meddler:"repo_html_url"`
		Branch     string `meddler:"repo_branch"`
		Private    bool   `meddler:"repo_private"`
		Visibility string `meddler:"repo_visibility"`
		Active     bool   `meddler:"repo_active"`
		Config     string `meddler:"repo_config"`
		Trusted    bool   `meddler:"repo_trusted"`
		Protected  bool   `meddler:"repo_protected"`
		Timeout    int64  `meddler:"repo_timeout"`
		Counter    int64  `meddler:"repo_counter"`
		Synced     int64  `meddler:"repo_synced"`
		Created    int64  `meddler:"repo_created"`
		Updated    int64  `meddler:"repo_updated"`
		Version    int64  `meddler:"repo_version"`
		Signer     string `meddler:"repo_signer"`
		Secret     string `meddler:"repo_secret"`
	}

	// SecretV0 is a Drone 0.x secret.
	SecretV0 struct {
		ID           int64  `meddler:"secret_id"`
		RepoID       int64  `meddler:"secret_repo_id"`
		RepoFullname string `meddler:"repo_full_name"`
		Name         string `meddler:"secret_name"`
		Value        string `meddler:"secret_value"`
		Images       string `meddler:"secret_images"`
		Events       string `meddler:"secret_events"`
		SkipVerify   bool   `meddler:"secret_skip_verify"`
		Conceal      bool   `meddler:"secret_conceal"`
	}

	// SecretV1 is a Drone 1.x secret.
	SecretV1 struct {
		ID              int64  `meddler:"secret_id"`
		RepoID          int64  `meddler:"secret_repo_id"`
		Name            string `meddler:"secret_name"`
		Data            string `meddler:"secret_data"`
		PullRequest     bool   `meddler:"secret_pull_request"`
		PullRequestPush bool   `meddler:"secret_pull_request_push"`
	}

	// RegistryV0 is a Drone 0.x registry.
	RegistryV0 struct {
		ID           int64  `meddler:"registry_id"`
		RepoID       int64  `meddler:"registry_repo_id"`
		RepoFullname string `meddler:"repo_full_name"`
		Addr         string `meddler:"registry_addr"`
		Email        string `meddler:"registry_email"`
		Username     string `meddler:"registry_username"`
		Password     string `meddler:"registry_password"`
		Token        string `meddler:"registry_token"`
	}

	// DockerConfig defines required attributes from Docker registry credentials.
	DockerConfig struct {
		AuthConfigs map[string]AuthConfig `json:"auths"`
	}

	// AuthConfig contains authorization information for connecting to a Registry.
	AuthConfig struct {
		Email    string `json:"email,omitempty"`
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
		Auth     string `json:"auth,omitempty"`
	}
)

func (c AuthConfig) MarshalJSON() ([]byte, error) {
	result := struct {
		Auth  string `json:"auth,omitempty"`
		Email string `json:"email,omitempty"`
	}{
		Email: c.Email,
	}

	credentials := []byte(c.Username + ":" + c.Password)

	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(credentials)))
	base64.StdEncoding.Encode(encoded, credentials)

	result.Auth = string(encoded)

	return json.Marshal(result)
}
