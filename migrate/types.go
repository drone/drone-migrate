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
		ID          int64  `meddler:"repo_id"`
		UID         string `meddler:"repo_uid"`
		UserID      int64  `meddler:"repo_user_id"`
		Namespace   string `meddler:"repo_namespace"`
		Name        string `meddler:"repo_name"`
		Slug        string `meddler:"repo_slug"`
		SCM         string `meddler:"repo_scm"`
		HTTPURL     string `meddler:"repo_clone_url"`
		SSHURL      string `meddler:"repo_ssh_url"`
		Link        string `meddler:"repo_html_url"`
		Branch      string `meddler:"repo_branch"`
		Private     bool   `meddler:"repo_private"`
		Visibility  string `meddler:"repo_visibility"`
		Active      bool   `meddler:"repo_active"`
		Config      string `meddler:"repo_config"`
		Trusted     bool   `meddler:"repo_trusted"`
		Protected   bool   `meddler:"repo_protected"`
		IgnoreForks bool   `meddler:"repo_no_forks"`
		IgnorePulls bool   `meddler:"repo_no_pulls"`
		Timeout     int64  `meddler:"repo_timeout"`
		Counter     int64  `meddler:"repo_counter"`
		Synced      int64  `meddler:"repo_synced"`
		Created     int64  `meddler:"repo_created"`
		Updated     int64  `meddler:"repo_updated"`
		Version     int64  `meddler:"repo_version"`
		Signer      string `meddler:"repo_signer"`
		Secret      string `meddler:"repo_secret"`
	}

	// BuildV0 is a Drone 0.x build.
	BuildV0 struct {
		ID        int64  `meddler:"build_id"`
		RepoID    int64  `meddler:"build_repo_id"`
		ConfigID  int64  `meddler:"build_config_id"`
		Number    int64  `meddler:"build_number"`
		Parent    int64  `meddler:"build_parent"`
		Event     string `meddler:"build_event"`
		Status    string `meddler:"build_status"`
		Error     string `meddler:"build_error"`
		Enqueued  int64  `meddler:"build_enqueued"`
		Created   int64  `meddler:"build_created"`
		Started   int64  `meddler:"build_started"`
		Finished  int64  `meddler:"build_finished"`
		Deploy    string `meddler:"build_deploy"`
		Commit    string `meddler:"build_commit"`
		Branch    string `meddler:"build_branch"`
		Ref       string `meddler:"build_ref"`
		Refspec   string `meddler:"build_refspec"`
		Remote    string `meddler:"build_remote"`
		Title     string `meddler:"build_title"`
		Message   string `meddler:"build_message"`
		Timestamp int64  `meddler:"build_timestamp"`
		Sender    string `meddler:"build_sender"`
		Author    string `meddler:"build_author"`
		Avatar    string `meddler:"build_avatar"`
		Email     string `meddler:"build_email"`
		Link      string `meddler:"build_link"`
		Signed    bool   `meddler:"build_signed"`   // deprecate
		Verified  bool   `meddler:"build_verified"` // deprecate
		Reviewer  string `meddler:"build_reviewer"`
		Reviewed  int64  `meddler:"build_reviewed"`
	}

	// BuildV1 is a Drone 1.x build.
	BuildV1 struct {
		ID           int64             `meddler:"build_id"`
		RepoID       int64             `meddler:"build_repo_id"`
		Trigger      string            `meddler:"build_trigger"`
		Number       int64             `meddler:"build_number"`
		Parent       int64             `meddler:"build_parent"`
		Status       string            `meddler:"build_status"`
		Error        string            `meddler:"build_error"`
		Event        string            `meddler:"build_event"`
		Action       string            `meddler:"build_action"`
		Link         string            `meddler:"build_link"`
		Timestamp    int64             `meddler:"build_timestamp"`
		Title        string            `meddler:"build_title"`
		Message      string            `meddler:"build_message"`
		Before       string            `meddler:"build_before"`
		After        string            `meddler:"build_after"`
		Ref          string            `meddler:"build_ref"`
		Fork         string            `meddler:"build_source_repo"`
		Source       string            `meddler:"build_source"`
		Target       string            `meddler:"build_target"`
		Author       string            `meddler:"build_author"`
		AuthorName   string            `meddler:"build_author_name"`
		AuthorEmail  string            `meddler:"build_author_email"`
		AuthorAvatar string            `meddler:"build_author_avatar"`
		Sender       string            `meddler:"build_sender"`
		Params       map[string]string `meddler:"build_params,json"`
		Deploy       string            `meddler:"build_deploy"`
		Started      int64             `meddler:"build_started"`
		Finished     int64             `meddler:"build_finished"`
		Created      int64             `meddler:"build_created"`
		Updated      int64             `meddler:"build_updated"`
		Version      int64             `meddler:"build_version"`
	}

	// StageV0 is a Drone 0.x stage.
	StageV0 struct {
		ID       int64             `meddler:"proc_id"`
		BuildID  int64             `meddler:"proc_build_id"`
		PID      int               `meddler:"proc_pid"`
		PPID     int               `meddler:"proc_ppid"`
		PGID     int               `meddler:"proc_pgid"`
		Name     string            `meddler:"proc_name"`
		State    string            `meddler:"proc_state"`
		Error    string            `meddler:"proc_error"`
		ExitCode int               `meddler:"proc_exit_code"`
		Started  int64             `meddler:"proc_started"`
		Stopped  int64             `meddler:"proc_stopped"`
		Machine  string            `meddler:"proc_machine"`
		Platform string            `meddler:"proc_platform"`
		Environ  map[string]string `meddler:"proc_environ,json"`
	}

	// StageV1 is a Drone 1.x stage.
	StageV1 struct {
		ID        int64             `meddler:"stage_id"`
		RepoID    int64             `meddler:"stage_repo_id"`
		BuildID   int64             `meddler:"stage_build_id"`
		Number    int               `meddler:"stage_number"`
		Name      string            `meddler:"stage_name"`
		Kind      string            `meddler:"stage_kind"`
		Type      string            `meddler:"stage_type"`
		Status    string            `meddler:"stage_status"`
		Error     string            `meddler:"stage_error"`
		ErrIgnore bool              `meddler:"stage_errignore"`
		ExitCode  int               `meddler:"stage_exit_code"`
		Machine   string            `meddler:"stage_machine"`
		OS        string            `meddler:"stage_os"`
		Arch      string            `meddler:"stage_arch"`
		Variant   string            `meddler:"stage_variant"`
		Kernel    string            `meddler:"stage_kernel"`
		Limit     int               `meddler:"stage_limit"`
		Started   int64             `meddler:"stage_started"`
		Stopped   int64             `meddler:"stage_stopped"`
		Created   int64             `meddler:"stage_created"`
		Updated   int64             `meddler:"stage_updated"`
		Version   int64             `meddler:"stage_version"`
		OnSuccess bool              `meddler:"stage_on_success"`
		OnFailure bool              `meddler:"stage_on_failure"`
		DependsOn []string          `meddler:"stage_depends_on,json"`
		Labels    map[string]string `meddler:"stage_labels,json"`
	}

	// StepV0 is a Drone 0.x step.
	StepV0 struct {
		ID       int64             `meddler:"proc_id"`
		BuildID  int64             `meddler:"proc_build_id"`
		PID      int               `meddler:"proc_pid"`
		PPID     int               `meddler:"proc_ppid"`
		PGID     int               `meddler:"proc_pgid"`
		Name     string            `meddler:"proc_name"`
		State    string            `meddler:"proc_state"`
		Error    string            `meddler:"proc_error"`
		ExitCode int               `meddler:"proc_exit_code"`
		Started  int64             `meddler:"proc_started"`
		Stopped  int64             `meddler:"proc_stopped"`
		Machine  string            `meddler:"proc_machine"`
		Platform string            `meddler:"proc_platform"`
		Environ  map[string]string `meddler:"proc_environ,json"`
	}

	// StepV1 is a Drone 1.x step.
	StepV1 struct {
		ID        int64  `meddler:"step_id"`
		StageID   int64  `meddler:"step_stage_id"`
		Number    int    `meddler:"step_number"`
		Name      string `meddler:"step_name"`
		Status    string `meddler:"step_status"`
		Error     string `meddler:"step_error"`
		ErrIgnore bool   `meddler:"step_errignore"`
		ExitCode  int    `meddler:"step_exit_code"`
		Started   int64  `meddler:"step_started"`
		Stopped   int64  `meddler:"step_stopped"`
		Version   int64  `meddler:"step_version"`
	}

	// LogsV0 is a Drone 0.x logs.
	LogsV0 struct {
		ID     int64  `meddler:"log_id"`
		ProcID int64  `meddler:"log_job_id"`
		Data   []byte `meddler:"log_data"`
	}

	// LogsV1 is a Drone 1.x logs.
	LogsV1 struct {
		ID   int64  `meddler:"log_id"`
		Data []byte `meddler:"log_data"`
	}

	// SecretV0 is a Drone 0.x secret.
	SecretV0 struct {
		ID         int64    `meddler:"secret_id"`
		RepoID     int64    `meddler:"secret_repo_id"`
		Name       string   `meddler:"secret_name"`
		Value      string   `meddler:"secret_value"`
		Images     string   `meddler:"secret_images"`
		Events     []string `meddler:"secret_events,json"`
		SkipVerify bool     `meddler:"secret_skip_verify"`
		Conceal    bool     `meddler:"secret_conceal"`
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
	
	// RegistryV1 is a Drone 1.x registry -- note that in 1.x these are stored in the secrets table, hence the similar format
	RegistryV1 struct {
		ID              int64  `meddler:"secret_id,pk"`
		RepoID          int64  `meddler:"secret_repo_id"`
		Name            string `meddler:"secret_name"`
		Data            string `meddler:"secret_data"`
		PullRequest     bool   `meddler:"secret_pull_request"`
		PullRequestPush bool   `meddler:"secret_pull_request_push"`
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
