package qiniu

import "time"

type ErrorResponse struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e *ErrorResponse) Error() string {
	return e.Err
}

type CDNDomainListResponse struct {
	Domains []struct {
		Name               string      `json:"name"`
		PareDomain         string      `json:"pareDomain"`
		Type               string      `json:"type"`
		Cname              string      `json:"cname"`
		TestURLPath        string      `json:"testURLPath"`
		Protocol           string      `json:"protocol"`
		Platform           string      `json:"platform"`
		Product            string      `json:"product"`
		GeoCover           string      `json:"geoCover"`
		QiniuPrivate       bool        `json:"qiniuPrivate"`
		OperationType      string      `json:"operationType"`
		OperatingState     string      `json:"operatingState"`
		OperatingStateDesc string      `json:"operatingStateDesc"`
		FreezeType         string      `json:"freezeType"`
		CreateAt           time.Time   `json:"createAt"`
		ModifyAt           time.Time   `json:"modifyAt"`
		CouldOperateBySelf bool        `json:"couldOperateBySelf"`
		UidIsFreezed       bool        `json:"uidIsFreezed"`
		OemMail            string      `json:"oemMail"`
		TagList            interface{} `json:"tagList"`
		KvTagList          interface{} `json:"kvTagList"`
		IpTypes            int         `json:"ipTypes"`
		DeliveryBucket     string      `json:"deliveryBucket"`
		DeliveryBucketType string      `json:"deliveryBucketType"`
		DeliveryBucketFop  struct {
			Enable           bool        `json:"enable"`
			SufyDeliveryHost string      `json:"sufyDeliveryHost"`
			NewStyle         interface{} `json:"newStyle"`
			DeleteStyleNames interface{} `json:"deleteStyleNames"`
			NewSeparator     interface{} `json:"newSeparator"`
		} `json:"deliveryBucketFop"`
		IsSufy           bool   `json:"isSufy"`
		IsPcdnBackup     bool   `json:"isPcdnBackup"`
		IsPcdnBackup302  bool   `json:"isPcdnBackup302"`
		ChargeBackupType string `json:"chargeBackupType"`
		OperTaskId       string `json:"operTaskId"`
		OperTaskType     string `json:"operTaskType"`
		OperTaskErrCode  int    `json:"operTaskErrCode"`
	} `json:"domains"`
	Marker string `json:"marker"`
}

type CDNDomainDetailResponse struct {
	Name               string      `json:"name"`
	PareDomain         string      `json:"pareDomain"`
	Type               string      `json:"type"`
	Cname              string      `json:"cname"`
	TestURLPath        string      `json:"testURLPath"`
	Protocol           string      `json:"protocol"`
	Platform           string      `json:"platform"`
	Product            string      `json:"product"`
	GeoCover           string      `json:"geoCover"`
	QiniuPrivate       bool        `json:"qiniuPrivate"`
	OperationType      string      `json:"operationType"`
	OperatingState     string      `json:"operatingState"`
	OperatingStateDesc string      `json:"operatingStateDesc"`
	FreezeType         string      `json:"freezeType"`
	CreateAt           time.Time   `json:"createAt"`
	ModifyAt           time.Time   `json:"modifyAt"`
	CouldOperateBySelf bool        `json:"couldOperateBySelf"`
	UidIsFreezed       bool        `json:"uidIsFreezed"`
	OemMail            string      `json:"oemMail"`
	TagList            interface{} `json:"tagList"`
	KvTagList          interface{} `json:"kvTagList"`
	IpTypes            int         `json:"ipTypes"`
	DeliveryBucket     string      `json:"deliveryBucket"`
	DeliveryBucketType string      `json:"deliveryBucketType"`
	DeliveryBucketFop  struct {
		Enable           bool        `json:"enable"`
		SufyDeliveryHost string      `json:"sufyDeliveryHost"`
		NewStyle         interface{} `json:"newStyle"`
		DeleteStyleNames interface{} `json:"deleteStyleNames"`
		NewSeparator     interface{} `json:"newSeparator"`
	} `json:"deliveryBucketFop"`
	IsSufy           bool   `json:"isSufy"`
	IsPcdnBackup     bool   `json:"isPcdnBackup"`
	IsPcdnBackup302  bool   `json:"isPcdnBackup302"`
	ChargeBackupType string `json:"chargeBackupType"`
	Source           struct {
		SourceCname                  string      `json:"sourceCname"`
		SourceType                   string      `json:"sourceType"`
		SourceHost                   string      `json:"sourceHost"`
		TestSourceHost               string      `json:"testSourceHost"`
		SourceIPs                    []string    `json:"sourceIPs"`
		SourceDomain                 string      `json:"sourceDomain"`
		SourceQiniuBucket            string      `json:"sourceQiniuBucket"`
		SourceURLScheme              string      `json:"sourceURLScheme"`
		AdvancedSources              interface{} `json:"advancedSources"`
		SkipCheckSource              bool        `json:"skipCheckSource"`
		TransferConf                 interface{} `json:"transferConf"`
		SourceTimeACL                bool        `json:"sourceTimeACL"`
		SourceTimeACLKeys            interface{} `json:"sourceTimeACLKeys"`
		MaxSourceRate                int         `json:"maxSourceRate"`
		MaxSourceConcurrency         int         `json:"maxSourceConcurrency"`
		AddRespHeader                interface{} `json:"addRespHeader"`
		UrlRewrites                  interface{} `json:"urlRewrites"`
		SourceRetryCodes             interface{} `json:"sourceRetryCodes"`
		FollowRedirect               bool        `json:"followRedirect"`
		Redirect30X                  interface{} `json:"redirect30x"`
		SourceRequestHeaderControls  interface{} `json:"sourceRequestHeaderControls"`
		SourceResponseHeaderControls interface{} `json:"sourceResponseHeaderControls"`
		MaxSourceRatePerIDC          int         `json:"maxSourceRatePerIDC"`
		MaxSourceConcurrencyPerIDC   int         `json:"maxSourceConcurrencyPerIDC"`
		TestURLPath                  string      `json:"testURLPath"`
		SourceIgnoreParams           interface{} `json:"sourceIgnoreParams"`
		SourceIgnoreAllParams        bool        `json:"sourceIgnoreAllParams"`
		Range                        struct {
			Enable string `json:"enable"`
		} `json:"range"`
		EnableSourceAuth bool `json:"enableSourceAuth"`
		SourceAuthInfo   struct {
			OssProvider     string `json:"ossProvider"`
			AccessKeyId     string `json:"accessKeyId"`
			AccessKeySecret string `json:"accessKeySecret"`
		} `json:"sourceAuthInfo"`
	} `json:"source"`
	Bsauth struct {
		Path                           interface{} `json:"path"`
		Method                         string      `json:"method"`
		Parameters                     interface{} `json:"parameters"`
		TimeLimit                      int         `json:"timeLimit"`
		UserAuthUrl                    string      `json:"userAuthUrl"`
		Strict                         bool        `json:"strict"`
		Enable                         bool        `json:"enable"`
		SuccessStatusCode              int         `json:"successStatusCode"`
		FailureStatusCode              int         `json:"failureStatusCode"`
		IsQiniuPrivate                 bool        `json:"isQiniuPrivate"`
		BackSourceWithResourcePath     bool        `json:"backSourceWithResourcePath"`
		BackSourceWithoutClientHeaders bool        `json:"backSourceWithoutClientHeaders"`
		ResponseWithSourceAuthCode     bool        `json:"responseWithSourceAuthCode"`
		ResponseWithSourceAuthBody     bool        `json:"responseWithSourceAuthBody"`
		UserAuthUrlIpLimitConf         struct {
			Enable   bool `json:"enable"`
			Limit    int  `json:"limit"`
			TimeSlot int  `json:"timeSlot"`
		} `json:"userAuthUrlIpLimitConf"`
		UserAuthReqConf struct {
			Body                       interface{} `json:"body"`
			Header                     interface{} `json:"header"`
			Urlquery                   interface{} `json:"urlquery"`
			IncludeClientHeadersInBody bool        `json:"includeClientHeadersInBody"`
		} `json:"userAuthReqConf"`
		UserAuthContentType  string `json:"userAuthContentType"`
		UserAuthRespBodyConf struct {
			Enable                 bool        `json:"enable"`
			ContentType            string      `json:"contentType"`
			SuccessConditions      interface{} `json:"successConditions"`
			SuccessLogicalOperator string      `json:"successLogicalOperator"`
			FailureConditions      interface{} `json:"failureConditions"`
			FailureLogicalOperator string      `json:"failureLogicalOperator"`
		} `json:"userAuthRespBodyConf"`
		UserBsauthResultCacheConf struct {
			CacheEnable     bool        `json:"cacheEnable"`
			CacheSingleType string      `json:"cacheSingleType"`
			CacheKeyElems   interface{} `json:"cacheKeyElems"`
			CacheShareHost  string      `json:"cacheShareHost"`
			CacheDuration   int         `json:"cacheDuration"`
		} `json:"userBsauthResultCacheConf"`
		UserAuthMatchRuleConf struct {
			Type string `json:"type"`
			Rule string `json:"rule"`
		} `json:"userAuthMatchRuleConf"`
	} `json:"bsauth"`
	External struct {
		EnableFop bool `json:"enableFop"`
		ImageSlim struct {
			EnableImageSlim  bool          `json:"enableImageSlim"`
			PrefixImageSlims []interface{} `json:"prefixImageSlims"`
			RegexpImageSlims []interface{} `json:"regexpImageSlims"`
		} `json:"imageSlim"`
	} `json:"external"`
	Cache struct {
		CacheControls []struct {
			Time     int    `json:"time"`
			Timeunit int    `json:"timeunit"`
			Type     string `json:"type"`
			Rule     string `json:"rule"`
		} `json:"cacheControls"`
		IgnoreParam   bool          `json:"ignoreParam"`
		IgnoreParams  []interface{} `json:"ignoreParams"`
		IncludeParams []interface{} `json:"includeParams"`
	} `json:"cache"`
	Referer struct {
		RefererType   string        `json:"refererType"`
		RefererValues []interface{} `json:"refererValues"`
		NullReferer   bool          `json:"nullReferer"`
	} `json:"referer"`
	TimeACL struct {
		Enable                bool        `json:"enable"`
		TimeACLKeys           interface{} `json:"timeACLKeys"`
		AuthType              string      `json:"authType"`
		AuthDelta             int         `json:"authDelta"`
		SufyTimeACLKeys       interface{} `json:"sufyTimeACLKeys"`
		SufyCallbackBody      interface{} `json:"sufyCallbackBody"`
		CheckUrl              string      `json:"checkUrl"`
		AdvanceFunctionEnable bool        `json:"advanceFunctionEnable"`
		RuleType              string      `json:"ruleType"`
		Rules                 interface{} `json:"rules"`
		Params                interface{} `json:"params"`
		ParamStr              string      `json:"paramStr"`
		ToLowerCase           string      `json:"toLowerCase"`
		UrlEncode             string      `json:"urlEncode"`
		HashMethod            string      `json:"hashMethod"`
		Verification          struct {
			Name     string `json:"name"`
			Locate   string `json:"locate"`
			FailCode int    `json:"failCode"`
		} `json:"verification"`
	} `json:"timeACL"`
	IpACL struct {
		IpACLType   string        `json:"ipACLType"`
		IpACLValues []interface{} `json:"ipACLValues"`
	} `json:"ipACL"`
	UaACL                  interface{}   `json:"uaACL"`
	RequestHeaders         interface{}   `json:"requestHeaders"`
	ResponseHeaderControls []interface{} `json:"responseHeaderControls"`
	Https                  struct {
		CertId      string `json:"certId"`
		ForceHttps  bool   `json:"forceHttps"`
		Http2Enable bool   `json:"http2Enable"`
		FreeCert    bool   `json:"freeCert"`
	} `json:"https"`
	RegisterNo         string    `json:"registerNo"`
	ConfigProcessRatio int       `json:"configProcessRatio"`
	HurryUpFreecert    bool      `json:"hurryUpFreecert"`
	HttpsOPTime        time.Time `json:"httpsOPTime"`
	Range              struct {
		Enable string `json:"enable"`
	} `json:"range"`
	OperTaskId      string `json:"operTaskId"`
	OperTaskType    string `json:"operTaskType"`
	OperTaskErrCode int    `json:"operTaskErrCode"`
}

type UploadCertRequest struct {
	Name       string `json:"name"`
	CommonName string `json:"common_name"`
	Pri        string `json:"pri"`
	Ca         string `json:"ca"`
}

type UploadCertResponse struct {
	CertID string `json:"certID"`
}

type SetCDNHttpsRequest struct {
	Certid      string `json:"certid"`
	ForceHttps  bool   `json:"forceHttps"`
	Http2Enable bool   `json:"http2Enable"`
}
