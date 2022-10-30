# binance-trader

Make funny things by composing the binance-go API

- exchange & withdraw


## Setup

- add new file **.env**
```
BINANCE_API_KEY=
BINANCE_SECRET_KEY=
PROXY_URL=
```

## import 
```
go get github.com/pursonchen/binance-trader
```

## Caution
- self design pkg must be put in the GOPATH/src/github.com/pursonchen/xxx 
- **var \*Type** can't alloc mem but **param := new(Type)** And var Type