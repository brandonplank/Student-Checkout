module brandonplank.org/checkout

go 1.17

require (
	github.com/gofiber/fiber/v2 v2.27.0
	github.com/gofiber/template v1.6.23
	github.com/joho/godotenv v1.4.0
	github.com/mileusna/crontab v1.2.0
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292
)

require (
	brandonplank.org/checkout/models v0.0.0-00010101000000-000000000000 // indirect
	github.com/gocarina/gocsv v0.0.0-20211203214250-4735fba0c1d9 // indirect
	github.com/jordan-wright/email v4.0.1-0.20210109023952-943e75fe5223+incompatible // indirect
)

replace (
	//brandonplank.org/checkout/global => ./Global
	brandonplank.org/checkout/models => ./Models
	brandonplank.org/checkout/routes => ./Routes
)

require (
	brandonplank.org/checkout/routes v0.0.0
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/klauspost/compress v1.14.2 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.33.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20220209214540-3681064d5158 // indirect
)
