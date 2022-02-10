module brandonplank.org/checkout

go 1.17

require (
	github.com/gofiber/fiber/v2 v2.26.0
	github.com/gofiber/template v1.6.22
)

//brandonplank.org/checkout/global => ./Global
//brandonplank.org/checkout/models => ./Models
replace brandonplank.org/checkout/routes => ./Routes

require (
	//brandonplank.org/checkout/global v0.0.0
	//brandonplank.org/checkout/models v0.0.0
	brandonplank.org/checkout/routes v0.0.0
	github.com/andybalholm/brotli v1.0.2 // indirect
	github.com/klauspost/compress v1.13.4 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.32.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20211205182925-97ca703d548d // indirect
)
