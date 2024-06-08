package config

import (
	"encoding/json"

	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	AppPort string `mapstructure:"APP_PORT"`

	DBHost                string `mapstructure:"DB_HOST"`
	DBUser                string `mapstructure:"DB_USER"`
	DBPassword            string `mapstructure:"DB_PASSWORD"`
	DBName                string `mapstructure:"DB_NAME"`
	DBPort                string `mapstructure:"DB_PORT"`
	ENV                   string `mapstructure:"APP_ENV"`
	RandomKey             string `mapstructure:"RANDOM_API_KEY"`
	RandomUrl             string `mapstructure:"RANDOM_RPC_URL"`
	SolanaRpcUrl          string `mapstructure:"TESTNET_RPC_URL"`
	Network               string `mapstructure:"NETWORK"`
	MasterWalletPubKey    string `mapstructure:"MASTER_WALLET_PUB_KEY"`
	MasterWalletPriKey    string `mapstructure:"MASTER_WALLET_PRI_KEY"`
	ProgramKey            string `mapstructure:"PROGRAM_KEY"`
	TokenPriceApi         string `mapstructure:"TOKEN_PRICE_API"`
	TatumBroadcastApi     string `mapstructure:"TATUM_BROADCAST_CONFIRM"`
	TatumApiKey           string `mapstructure:"TATUM_TESTNET_API_KEY"`
	TatumHmacSecret       string `mapstructure:"TATUM_HMAC_SECRET"`
	AWSAccessID           string `mapstructure:"AWS_ACCESS_ID"`
	AWSSecretKey          string `mapstructure:"AWS_SECRET_KEY"`
	AWSRegion             string `mapstructure:"AWS_REGION"`
	S3BucketName          string `mapstructure:"S3_BUCKET_NAME"`
	HyperspaceApiKey      string `mapstructure:"HYPERSPACE_API_KEY"`
	CloudWatchLogGroup    string `mapstructure:"CLOUD_WATCH_LOG_GROUP_NAME"`
	IpGeolocationApiKey   string `mapstructure:"IP_GEOLOCATION_API_KEY"`
	JupiterApiAccessKey   string `mapstructure:"JUPITER_AGGREGATER_API_ACCESS_KEY"`
	JupiterAggregater     string `mapstructure:"JUPITER_AGGREGATER_URL"`
	MixpanelToken         string `mapstructure:"MIXPANEL_TOKEN"`
	MixpanelServerUrl     string `mapstructure:"MIXPANEL_SERVER_URL"`
	MatricaApiAccessToken string `mapstructure:"MATRICA_API_ACCESS_TOKEN"`
	AdminApiAccessToken   string `mapstructure:"ADMIN_API_ACCESS_TOKEN"`
	RedisUrl              string `mapstructure:"REDIS_URL"`
	RedisPwd              string `mapstructure:"REDIS_PWD"`
	WeeklyRaffleRandomKey string `mapstructure:"WEEKLY_RAFFLE_RANDOM_KEY"`
}

var config Config

func Init() {
	var err error
	config, err = load(".")
	if err != nil {
		config, err = loadEnv()
	}
	if err != nil {
		logrus.Fatal("Could not load config: ", err)
	}
}

func Get() Config {
	return config
}

func load(path string) (conf Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&conf)
	return
}

func loadEnv() (conf Config, err error) {
	viper.SetEnvPrefix("duelana")
	viper.AutomaticEnv()

	conf = Config{
		AppPort:               viper.GetString("app_port"),
		DBHost:                viper.GetString("db_host"),
		DBUser:                viper.GetString("db_user"),
		DBPassword:            viper.GetString("db_password"),
		DBName:                viper.GetString("db_name"),
		DBPort:                viper.GetString("db_port"),
		ENV:                   viper.GetString("app_env"),
		RandomKey:             viper.GetString("random_api_key"),
		RandomUrl:             viper.GetString("random_rpc_url"),
		Network:               viper.GetString("network"),
		MasterWalletPubKey:    viper.GetString("master_wallet_pub_key"),
		MasterWalletPriKey:    viper.GetString("master_wallet_pri_key"),
		ProgramKey:            viper.GetString("program_key"),
		TokenPriceApi:         viper.GetString("TOKEN_PRICE_API"),
		TatumBroadcastApi:     viper.GetString("tatum_broadcast_confirm"),
		TatumApiKey:           viper.GetString("tatum_testnet_api_key"),
		TatumHmacSecret:       viper.GetString("tatum_hmac_secret"),
		AWSAccessID:           viper.GetString("AWS_ACCESS_ID"),
		AWSSecretKey:          viper.GetString("AWS_SECRET_KEY"),
		AWSRegion:             viper.GetString("AWS_REGION"),
		S3BucketName:          viper.GetString("S3_BUCKET_NAME"),
		HyperspaceApiKey:      viper.GetString("HYPERSPACE_API_KEY"),
		CloudWatchLogGroup:    viper.GetString("CLOUD_WATCH_LOG_GROUP_NAME"),
		IpGeolocationApiKey:   viper.GetString("IP_GEOLOCATION_API_KEY"),
		JupiterApiAccessKey:   viper.GetString("JUPITER_AGGREGATER_API_ACCESS_KEY"),
		JupiterAggregater:     viper.GetString("JUPITER_AGGREGATER_URL"),
		MixpanelToken:         viper.GetString("MIXPANEL_TOKEN"),
		MixpanelServerUrl:     viper.GetString("MIXPANEL_SERVER_URL"),
		MatricaApiAccessToken: viper.GetString("MATRICA_API_ACCESS_TOKEN"),
		AdminApiAccessToken:   viper.GetString("ADMIN_API_ACCESS_TOKEN"),
		RedisUrl:              viper.GetString("REDIS_URL"),
		WeeklyRaffleRandomKey: viper.GetString("WEEKLY_RAFFLE_RANDOM_KEY"),
	}
	if conf.Network == "mainnet" {
		conf.SolanaRpcUrl = viper.GetString("mainnet_rpc_url")
		conf.TatumApiKey = viper.GetString("tatum_mainnet_api_key")
	} else {
		conf.SolanaRpcUrl = viper.GetString("testnet_rpc_url")
		conf.TatumApiKey = viper.GetString("tatum_testnet_api_key")
	}
	err = nil

	return
}

var serverConfig models.ServerConfig = models.ServerConfig{
	BaseRakeBackRate:       5,
	AdditionalRakeBackRate: 0,
}
var apiRateConfig map[string]types.RateLimit
var websocketRateConfig map[string]types.RateLimit

func SetServerConfig(_serverConfig models.ServerConfig) {
	serverConfig = _serverConfig
	if err := json.Unmarshal(
		[]byte(
			serverConfig.ApiRateLimitConfiguration,
		),
		&apiRateConfig,
	); err != nil {
		return
	}
	if err := json.Unmarshal(
		[]byte(
			serverConfig.WebsocketRateLimitConfiguration,
		),
		&websocketRateConfig,
	); err != nil {
		return
	}
}

func GetServerConfig() models.ServerConfig {
	return serverConfig
}

func GetAPIRateConfig(key string) types.RateLimit {
	return apiRateConfig[key]
}

func GetWebsocketRateConfig(key string) types.RateLimit {
	return websocketRateConfig[key]
}
