package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"math"
	"strconv"
	"time"
)

func LoadConfig(path string) IConfig {
	envMap, err := godotenv.Read(path)
	if err != nil {
		log.Fatalf("load dotenv failed: %v", err)
	}

	return &config{
		app: &app{
			host: envMap["APP_HOST"],
			port: func() int {
				p, err := strconv.Atoi(envMap["APP_PORT"])
				if err != nil {
					log.Fatalf("load app port fialed %v", err)
				}

				return p
			}(),
			name:    envMap["APP_NAME"],
			version: envMap["APP_VERSION"],
			readTimeout: func() time.Duration {
				p, err := strconv.Atoi(envMap["APP_READ_TIMEOUT"])
				if err != nil {
					log.Fatalf("load app read timeout fialed %v", err)
				}

				return time.Duration(int64(p) * int64(math.Pow10(9)))
			}(),
			writeTimeout: func() time.Duration {
				p, err := strconv.Atoi(envMap["APP_WRITE_TIMEOUT"])
				if err != nil {
					log.Fatalf("load app write timeout fialed %v", err)
				}

				return time.Duration(int64(p) * int64(math.Pow10(9)))
			}(),
			bodyLimit: func() int {
				p, err := strconv.Atoi(envMap["APP_BODY_LIMIT"])
				if err != nil {
					log.Fatalf("load app body limit fialed %v", err)
				}

				return p
			}(),
			fileLimit: func() int {
				p, err := strconv.Atoi(envMap["APP_FILE_LIMIT"])
				if err != nil {
					log.Fatalf("load app file limit fialed %v", err)
				}

				return p
			}(),
			gcpBucket: envMap["APP_GCP_BUCKET"],
		},
		db: &db{
			host: envMap["DB_HOST"],
			port: func() int {
				p, err := strconv.Atoi(envMap["DB_PORT"])
				if err != nil {
					log.Fatalf("load db port fialed %v", err)
				}

				return p
			}(),
			protocol: envMap["DB_PROTOCOL"],
			username: envMap["DB_USERNAME"],
			password: envMap["DB_PASSWORD"],
			database: envMap["DB_DATABASE"],
			sslMode:  envMap["DB_SSL_MODE"],
			maxConnections: func() int {
				p, err := strconv.Atoi(envMap["DB_MAX_CONNECTIONS"])
				if err != nil {
					log.Fatalf("load db max connections fialed %v", err)
				}

				return p
			}(),
		},
		jwt: &jwt{
			adminKey:  envMap["JWT_ADMIN_KEY"],
			secretKey: envMap["JWT_SECRET_KEY"],
			apiKey:    envMap["JWT_API_KEY"],
			accessExpiresAt: func() int {
				t, err := strconv.Atoi(envMap["JWT_ACCESS_EXPIRES"])
				if err != nil {
					log.Fatalf("load jwt access expires at fialed %v", err)
				}

				return t
			}(),
			refreshExpiresAt: func() int {
				t, err := strconv.Atoi(envMap["JWT_REFRESH_EXPIRES"])
				if err != nil {
					log.Fatalf("load jwt refresh expires at fialed %v", err)
				}

				return t
			}(),
		},
	}
}

type IConfig interface {
	App() IAppConfig
	Db() IDbConfig
	Jwt() IJwtConfig
}

type config struct {
	app *app
	db  *db
	jwt *jwt
}

type IAppConfig interface {
	Url() string // host:port
	Name() string
	Version() string
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	BodyLimit() int
	FileLimit() int
	GCPBucket() string
}

type app struct {
	host         string
	port         int
	name         string
	version      string
	readTimeout  time.Duration
	writeTimeout time.Duration
	bodyLimit    int // bytes
	fileLimit    int // bytes
	gcpBucket    string
}

func (a *app) Url() string {
	return fmt.Sprintf("%s:%d", a.host, a.port)
}

func (a *app) Name() string { return a.name }

func (a *app) Version() string { return a.version }

func (a *app) ReadTimeout() time.Duration { return a.readTimeout }

func (a *app) WriteTimeout() time.Duration { return a.writeTimeout }

func (a *app) BodyLimit() int { return a.bodyLimit }

func (a *app) FileLimit() int { return a.fileLimit }

func (a *app) GCPBucket() string { return a.gcpBucket }

func (c *config) App() IAppConfig { return c.app }

type IDbConfig interface {
	Url() string
	MaxOpenConn() int
}

type db struct {
	host           string
	port           int
	protocol       string
	username       string
	password       string
	database       string
	sslMode        string
	maxConnections int
}

func (d *db) Url() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.host,
		d.port,
		d.username,
		d.password,
		d.database,
		d.sslMode,
	)
}

func (d *db) MaxOpenConn() int { return d.maxConnections }

func (c *config) Db() IDbConfig {
	return c.db
}

type IJwtConfig interface {
	SecretKey() []byte
	AdminKey() []byte
	ApiKey() []byte
	AccessExpiresAt() int
	RefreshExpiresAt() int
	SetJwtAccessExpires(t int)
	SetJwtRefreshExpires(t int)
}

type jwt struct {
	adminKey         string
	secretKey        string
	apiKey           string
	accessExpiresAt  int // sec
	refreshExpiresAt int // sec
}

func (j *jwt) SecretKey() []byte { return []byte(j.secretKey) }

func (j *jwt) AdminKey() []byte { return []byte(j.adminKey) }

func (j *jwt) ApiKey() []byte { return []byte(j.apiKey) }

func (j *jwt) AccessExpiresAt() int { return j.accessExpiresAt }

func (j *jwt) RefreshExpiresAt() int { return j.refreshExpiresAt }

func (j *jwt) SetJwtAccessExpires(t int) { j.accessExpiresAt = t }

func (j *jwt) SetJwtRefreshExpires(t int) { j.refreshExpiresAt = t }

func (c *config) Jwt() IJwtConfig {
	return c.jwt
}
