module github.com/lishengxie/MobileInternet

go 1.15

require (
	github.com/hyperledger/fabric-sdk-go v1.0.0
	github.com/lishengxie/MobileInternet/web v0.0.0-20210325122409-9d20087f105b
	github.com/lishengxie/MobileInternet/web/controller v0.0.0-20210325122409-9d20087f105b
	github.com/pkg/errors v0.9.1
)

replace (
	github.com/lishengxie/MobileInternet/web v0.0.0-20210325122409-9d20087f105b => ./web
	github.com/lishengxie/MobileInternet/web/controller v0.0.0-20210325122409-9d20087f105b => ./web/controller
)