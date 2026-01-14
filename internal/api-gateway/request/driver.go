package request

// NewDriver 绑定 JSON
type NewDriver struct {
	UserId                  int64  `form:"user_id" json:"user_id" xml:"user_id"  binding:"required"`
	RealName                string `form:"real_name" json:"real_name" xml:"real_name"  binding:"required"`
	IdCardNo                string `form:"id_card_no" json:"id_card_no" xml:"id_card_no"  binding:"required"`
	IdCardFrontUrl          string `form:"id_card_front_url" json:"id_card_front_url" xml:"id_card_front_url"  binding:"required"`
	IdCardBackUrl           string `form:"id_card_back_url" json:"id_card_back_url" xml:"id_card_back_url"  binding:"required"`
	IdCardExpireTime        string `form:"id_card_expire_time" json:"id_card_expire_time" xml:"id_card_expire_time"  binding:"required"`
	DriverLicenseNo         string `form:"driver_license_no" json:"driver_license_no" xml:"driver_license_no"  binding:"required"`
	DriverLicenseUrl        string `form:"driver_license_url" json:"driver_license_url" xml:"driver_license_url"  binding:"required"`
	DriverLicenseGetTime    string `form:"driver_license_get_time" json:"driver_license_get_time" xml:"driver_license_get_time"  binding:"required"`
	DriverLicenseExpireTime string `form:"driver_license_expire_time" json:"driver_license_expire_time" xml:"driver_license_expire_time"  binding:"required"`
	DriverAge               int64  `form:"driver_age" json:"driver_age" xml:"driver_age"  binding:"required"`
	HealthCertUrl           string `form:"health_cert_url" json:"health_cert_url" xml:"health_cert_url"`
	ResidencePermitUrl      string `form:"residence_permit_url" json:"residence_permit_url" xml:"residence_permit_url"`
	CityCode                string `form:"city_code" json:"city_code" xml:"city_code"  binding:"required"`
	VehicleNo               string `form:"vehicle_no" json:"vehicle_no" xml:"vehicle_no"  binding:"required"`
	VehicleType             string `form:"vehicle_type" json:"vehicle_type" xml:"vehicle_type"  binding:"required"`
	VehicleBrand            string `form:"vehicle_brand" json:"vehicle_brand" xml:"vehicle_brand"  binding:"required"`
	VehicleModel            string `form:"vehicle_model" json:"vehicle_model" xml:"vehicle_model"  binding:"required"`
	VehicleColor            string `form:"vehicle_color" json:"vehicle_color" xml:"vehicle_color"  binding:"required"`
	Vin                     string `form:"vin" json:"vin" xml:"vin"  binding:"required"`
	EngineNo                string `form:"engine_no" json:"engine_no" xml:"engine_no"`
	LicenseNo               string `form:"license_no" json:"license_no" xml:"license_no"  binding:"required"`
	LicenseExpireDate       string `form:"license_expire_date" json:"license_expire_date" xml:"license_expire_date"  binding:"required"`
	InsuranceExpireDate     string `form:"insurance_expire_date" json:"insurance_expire_date" xml:"insurance_expire_date"  binding:"required"`
}

// UpdDriver 绑定 JSON
type UpdDriver struct {
	RealName                string `form:"real_name" json:"real_name" xml:"real_name"`
	IdCardNo                string `form:"id_card_no" json:"id_card_no" xml:"id_card_no"`
	IdCardFrontUrl          string `form:"id_card_front_url" json:"id_card_front_url" xml:"id_card_front_url"`
	IdCardBackUrl           string `form:"id_card_back_url" json:"id_card_back_url" xml:"id_card_back_url"`
	IdCardExpireTime        string `form:"id_card_expire_time" json:"id_card_expire_time" xml:"id_card_expire_time"`
	DriverLicenseNo         string `form:"driver_license_no" json:"driver_license_no" xml:"driver_license_no"`
	DriverLicenseUrl        string `form:"driver_license_url" json:"driver_license_url" xml:"driver_license_url"`
	DriverLicenseGetTime    string `form:"driver_license_get_time" json:"driver_license_get_time" xml:"driver_license_get_time"`
	DriverLicenseExpireTime string `form:"driver_license_expire_time" json:"driver_license_expire_time" xml:"driver_license_expire_time"`
	DriverAge               int64  `form:"driver_age" json:"driver_age" xml:"driver_age"`
	HealthCertUrl           string `form:"health_cert_url" json:"health_cert_url" xml:"health_cert_url"`
	ResidencePermitUrl      string `form:"residence_permit_url" json:"residence_permit_url" xml:"residence_permit_url"`
	CityCode                string `form:"city_code" json:"city_code" xml:"city_code"`
	VehicleNo               string `form:"vehicle_no" json:"vehicle_no" xml:"vehicle_no"`
	VehicleType             string `form:"vehicle_type" json:"vehicle_type" xml:"vehicle_type"`
	VehicleBrand            string `form:"vehicle_brand" json:"vehicle_brand" xml:"vehicle_brand"`
	VehicleModel            string `form:"vehicle_model" json:"vehicle_model" xml:"vehicle_model"`
	VehicleColor            string `form:"vehicle_color" json:"vehicle_color" xml:"vehicle_color"`
	Vin                     string `form:"vin" json:"vin" xml:"vin"`
	EngineNo                string `form:"engine_no" json:"engine_no" xml:"engine_no"`
	LicenseNo               string `form:"license_no" json:"license_no" xml:"license_no"`
	LicenseExpireDate       string `form:"license_expire_date" json:"license_expire_date" xml:"license_expire_date"`
	InsuranceExpireDate     string `form:"insurance_expire_date" json:"insurance_expire_date" xml:"insurance_expire_date"`
}
