module github.com/ethancls/cosmos-server

go 1.26.2

require (
	github.com/caddyserver/certmagic v0.21.3
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/coreos/go-oidc/v3 v3.18.0
	github.com/dexidp/dex v2.13.0+incompatible
	github.com/dexidp/dex/api/v2 v2.4.0
	github.com/eko/gocache/lib/v4 v4.2.0
	github.com/eko/gocache/store/go_cache/v4 v4.2.2
	github.com/eko/gocache/store/redis/v4 v4.2.2
	github.com/fsnotify/fsnotify v1.9.0
	github.com/gliderlabs/ssh v0.3.8
	github.com/go-jose/go-jose/v4 v4.1.4
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/golang/protobuf v1.5.4
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.0.2-0.20240212192251-757544f21357
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-secure-stdlib/base62 v0.1.2
	github.com/hashicorp/go-version v1.7.0
	github.com/jackc/pgx/v5 v5.5.5
	github.com/kardianos/service v1.2.3-0.20240613133416-becf2eb62b83
	github.com/lib/pq v1.12.3
	github.com/miekg/dns v1.1.72
	github.com/netbirdio/management-integrations/integrations v0.0.0-20260416123949-2355d972be42
	github.com/okta/okta-sdk-golang/v2 v2.18.0
	github.com/oschwald/maxminddb-golang v1.12.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pires/go-proxyproto v0.11.0
	github.com/prometheus/client_golang v1.23.2
	github.com/quic-go/quic-go v0.55.0
	github.com/redis/go-redis/v9 v9.7.3
	github.com/rs/cors v1.8.0
	github.com/rs/xid v1.3.0
	github.com/sirupsen/logrus v1.9.4
	github.com/spf13/cobra v1.10.2
	github.com/spf13/pflag v1.0.9
	github.com/stretchr/testify v1.11.1
	github.com/testcontainers/testcontainers-go v0.37.0
	github.com/testcontainers/testcontainers-go/modules/postgres v0.37.0
	github.com/vmihailenco/msgpack/v5 v5.4.1
	go.uber.org/mock v0.6.0
	go.uber.org/zap v1.27.0
	goauthentik.io/api/v3 v3.2023051.3
	golang.org/x/crypto v0.50.0
	golang.org/x/exp v0.0.0-20250620022241-b7579e27df2b
	golang.org/x/net v0.53.0
	golang.org/x/oauth2 v0.36.0
	golang.org/x/sync v0.20.0
	golang.org/x/sys v0.43.0
	golang.org/x/term v0.42.0
	golang.org/x/time v0.15.0
	google.golang.org/api v0.276.0
	google.golang.org/grpc v1.80.0
	google.golang.org/protobuf v1.36.11
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/postgres v1.5.7
	gorm.io/driver/sqlite v1.5.7
	gorm.io/gorm v1.25.12
	howett.net/plist v1.0.1
)

require (
	github.com/google/gopacket v1.1.19 // indirect
	github.com/huin/goupnp v1.2.0 // indirect
	github.com/jackpal/go-nat-pmp v1.0.2 // indirect
	github.com/koron/go-ssdp v0.0.4 // indirect
	github.com/libp2p/go-nat v0.2.0 // indirect
	github.com/libp2p/go-netroute v0.2.1 // indirect
)

replace github.com/dexidp/dex => github.com/netbirdio/dex v0.244.1-0.20260512110716-8d70ad8647c1

replace github.com/dexidp/dex/api/v2 => github.com/netbirdio/dex/api/v2 v2.0.0-20260512110716-8d70ad8647c1

replace github.com/kardianos/service => github.com/netbirdio/service v0.0.0-20240911161631-f62744f42502

replace github.com/mailru/easyjson => github.com/netbirdio/easyjson v0.9.0

replace github.com/pion/ice/v4 => github.com/netbirdio/ice/v4 v4.0.0-20250908184934-6202be846b51

replace github.com/cloudflare/circl => codeberg.org/cunicu/circl v0.0.0-20230801113412-fec58fc7b5f6

replace golang.zx2c4.com/wireguard => github.com/netbirdio/wireguard-go v0.0.0-20260628102922-2834bebf6c1a
