package oss

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

type (
	ClientOption = oss.ClientOption
	Bucket       = oss.Bucket
)

// client options
var (
	UseCname                 = oss.UseCname
	ForcePathStyle           = oss.ForcePathStyle
	Timeout                  = oss.Timeout
	MaxConns                 = oss.MaxConns
	SecurityToken            = oss.SecurityToken
	EnableMD5                = oss.EnableMD5
	MD5ThresholdCalcInMemory = oss.MD5ThresholdCalcInMemory
	EnableCRC                = oss.EnableCRC
	UserAgent                = oss.UserAgent
	Proxy                    = oss.Proxy
	AuthProxy                = oss.AuthProxy
	HttpClient               = oss.HTTPClient
	SetLogLevel              = oss.SetLogLevel
	SetLogger                = oss.SetLogger
	SetCredentialsProvider   = oss.SetCredentialsProvider
	SetLocalAddr             = oss.SetLocalAddr
	AuthVersion              = oss.AuthVersion
	AdditionalHeaders        = oss.AdditionalHeaders
	RedirectEnabled          = oss.RedirectEnabled
	InsecureSkipVerify       = oss.InsecureSkipVerify
	Region                   = oss.Region
	CloudBoxId               = oss.CloudBoxId
	Product                  = oss.Product
	VerifyObjectStrict       = oss.VerifyObjectStrict
)

type Client struct {
	*oss.Client
}

func New(cfg *Config, opts ...ClientOption) (*Client, error) {
	client, err := oss.New(cfg.Endpoint, cfg.AccessKeyId, cfg.AccessKeySecret, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{Client: client}, nil
}
