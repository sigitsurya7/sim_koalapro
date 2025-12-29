package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	DBDSN                string
	Port                 string
	MaxBodyBytes         int64
	ReadTimeout          time.Duration
	WriteTimeout         time.Duration
	IdleTimeout          time.Duration
	RateLimitRPS         float64
	RateLimitBurst       int
	LoginRateLimitRPS    float64
	LoginRateLimitBurst  int
	TrustedProxies       []string
	DBMaxOpenConns       int
	DBMaxIdleConns       int
	DBConnMaxLifetime    time.Duration
	JWTSecret            string
	MigrationPath        string
	StockityBaseURL      string
	CORSAllowedOrigins   []string
	CORSAllowCredentials bool
}

func Load() Config {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	migrationPath := os.Getenv("MIGRATION_PATH")
	if migrationPath == "" {
		migrationPath = "db/schema.sql"
	}

	stockityBaseURL := os.Getenv("STOCKITY_BASE_URL")
	if stockityBaseURL == "" {
		stockityBaseURL = "https://api.stockity.id"
	}
	corsAllowedOrigins := getSliceEnv("CORS_ALLOWED_ORIGINS")
	corsAllowCredentials := getBoolEnv("CORS_ALLOW_CREDENTIALS", false)

	maxBodyBytes := int64(getIntEnv("MAX_BODY_BYTES", 1<<20))
	readTimeout := time.Duration(getIntEnv("READ_TIMEOUT_SEC", 10)) * time.Second
	writeTimeout := time.Duration(getIntEnv("WRITE_TIMEOUT_SEC", 10)) * time.Second
	idleTimeout := time.Duration(getIntEnv("IDLE_TIMEOUT_SEC", 60)) * time.Second
	rateLimitRPS := getFloatEnv("RATE_LIMIT_RPS", 20)
	rateLimitBurst := getIntEnv("RATE_LIMIT_BURST", 40)
	loginRateLimitRPS := getFloatEnv("LOGIN_RATE_LIMIT_RPS", 5)
	loginRateLimitBurst := getIntEnv("LOGIN_RATE_LIMIT_BURST", 10)
	trustedProxies := getSliceEnv("TRUSTED_PROXIES")
	dbMaxOpenConns := getIntEnv("DB_MAX_OPEN_CONNS", 25)
	dbMaxIdleConns := getIntEnv("DB_MAX_IDLE_CONNS", 25)
	dbConnMaxLifetime := time.Duration(getIntEnv("DB_CONN_MAX_LIFETIME_MIN", 15)) * time.Minute

	return Config{
		DBDSN:                dsn,
		Port:                 port,
		JWTSecret:            jwtSecret,
		MigrationPath:        migrationPath,
		MaxBodyBytes:         maxBodyBytes,
		ReadTimeout:          readTimeout,
		WriteTimeout:         writeTimeout,
		IdleTimeout:          idleTimeout,
		RateLimitRPS:         rateLimitRPS,
		RateLimitBurst:       rateLimitBurst,
		LoginRateLimitRPS:    loginRateLimitRPS,
		LoginRateLimitBurst:  loginRateLimitBurst,
		TrustedProxies:       trustedProxies,
		DBMaxOpenConns:       dbMaxOpenConns,
		DBMaxIdleConns:       dbMaxIdleConns,
		DBConnMaxLifetime:    dbConnMaxLifetime,
		StockityBaseURL:      stockityBaseURL,
		CORSAllowedOrigins:   corsAllowedOrigins,
		CORSAllowCredentials: corsAllowCredentials,
	}
}

func getIntEnv(key string, def int) int {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	parsed, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("invalid %s, using default", key)
		return def
	}
	return parsed
}

func getFloatEnv(key string, def float64) float64 {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	parsed, err := strconv.ParseFloat(val, 64)
	if err != nil {
		log.Printf("invalid %s, using default", key)
		return def
	}
	return parsed
}

func getSliceEnv(key string) []string {
	val := os.Getenv(key)
	if val == "" {
		return nil
	}
	parts := strings.Split(val, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func getBoolEnv(key string, def bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	parsed, err := strconv.ParseBool(val)
	if err != nil {
		log.Printf("invalid %s, using default", key)
		return def
	}
	return parsed
}
