## Software Developer (Backend) Task for Monk Commerce 2025

### Prerequisites for project

- [Go](https://go.dev/doc/install) version >= 1.24
- [GNU Make](https://search.brave.com/search?q=gnu+make) version 4.3 (Optional)

### How to run the project

You can run the project in multiple ways
1. `make start`
2. OR `go build -o bin/app cmd/*.go && ./bin/app`
3. OR simply `go run cmd/*.go`

NOTE: Run `go mod tidy` if you are going to run it using the 3rd method

### Project Overview

```sh
.
├── Makefile
├── README.md
├── bin ## for binaries
├── cart ## cart package handling the coupon apply and applicable apis
├── cmd ## entrypoint
├── coupon ## coupon package for the coupon CRUD
├── go.mod
├── go.sum
└── utils ## some common utilities
```

- The project uses the in memory db (map[int]entity) to handle the db ops
- The current implementation is prone to race conditions and runtime panics if multiple APIs were requested concurrently
- The concurrency issue could be fixed with a simple `sync.Mutex` or `sync.Map`. But for the purpose of the assignment it has been kept as is.
- For our cart I have gone with a static product list with limitations that product_id can be from 1 to 10 and price of product will be product_id * 10
- The current version implements the 3 coupons described in the requirement document, i.e.
    - BxGY
    - Cartwise
    - ProductWise
- The way the project tackles different coupon is leveraging Go's interface
- All coupon implement `CouponDetails` interface
- Whenever required we retrieve the actual concrete type from it and use it to calculate relevant coupon apply

### Additional Cases

- We can have upto limit in the CartWise and ProductWise coupons, as currently they are flat discounts. We can add a max cap on it. Aptly named Cart Wise Upto and Product Wise Upto Coupon
- We can add a product category and have a coupon with discount on product category. This too can have a up-to variant. Let's name them Product Category Wise and Product Category Wise upto coupons
