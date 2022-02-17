module brandonplank.org/checkout/routes

go 1.17

require (
	brandonplank.org/checkout/models v0.0.0-00010101000000-000000000000
	github.com/gocarina/gocsv v0.0.0-20211203214250-4735fba0c1d9
	github.com/gofiber/fiber/v2 v2.26.0
	github.com/jordan-wright/email v4.0.1-0.20210109023952-943e75fe5223+incompatible
)

replace brandonplank.org/checkout/models => ../Models

require (
	github.com/andybalholm/brotli v1.0.2 // indirect
	github.com/klauspost/compress v1.13.4 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.32.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
)
