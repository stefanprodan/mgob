package config

type ModuleConfig struct {
	AWSClient    bool `envconfig:"en_aws_cli" default:"false"`
	AzureClient  bool `envconfig:"en_azure" default:"false"`
	GCloudClient bool `envconfig:"en_gcloud" default:"false"`
	GnuPG        bool `envconfig:"en_gpg" default:"true"`
	MinioClient  bool `envconfig:"en_minio" default:"false"`
	RCloneClient bool `envconfig:"en_rclone" default:"false"`
}
