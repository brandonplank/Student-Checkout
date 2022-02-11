module brandonplank.org/checkout/routes

go 1.17

require (
	github.com/gofiber/fiber/v2 v2.26.0
	github.com/xuri/excelize/v2 v2.5.0
)

replace (
	brandonplank.org/checkout/models => ../Models
)

require (
	brandonplank.org/checkout/models v0.0.0
)

require (
	github.com/andybalholm/brotli v1.0.2 // indirect
	github.com/klauspost/compress v1.13.4 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/richardlehane/mscfb v1.0.3 // indirect
	github.com/richardlehane/msoleps v1.0.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.32.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	github.com/xuri/efp v0.0.0-20210322160811-ab561f5b45e3 // indirect
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/net v0.0.0-20210726213435-c6fcb2dbf985 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/text v0.3.6 // indirect
)
