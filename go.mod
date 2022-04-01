module brandonplank.org/checkout

go 1.17

require (
	github.com/getsentry/sentry-go v0.12.0
	github.com/gofiber/fiber/v2 v2.27.0
	github.com/gofiber/template v1.6.23
	github.com/joho/godotenv v1.4.0
	github.com/mileusna/crontab v1.2.0
	go.mongodb.org/mongo-driver v1.8.4
)

require (
	brandonplank.org/checkout/models v0.0.0-00010101000000-000000000000 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/gocarina/gocsv v0.0.0-20211203214250-4735fba0c1d9 // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/jordan-wright/email v4.0.1-0.20210109023952-943e75fe5223+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.0.2 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/text v0.3.7 // indirect
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
