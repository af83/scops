package config

import (
	"flag"
	"os"
	"time"

	"github.com/af83/scops/clock"
	"github.com/af83/scops/logger"
)

var Config = struct {
	Debug  bool
	Syslog bool
	Gzip   bool

	RemoteUrl string
	AuthToken string
	Plugin    string

	Cycle time.Duration
}{}

func LoadConfig() {
	logger.Log.Printf("Load configuration")

	loadEnvConfig()
	loadFlagConfig()

	logger.Log.Debug = Config.Debug
	logger.Log.Syslog = Config.Syslog
}

func loadEnvConfig() {
	Config.Debug = os.Getenv("SCOPS_DEBUG") == "TRUE"
	Config.Syslog = os.Getenv("SCOPS_SYSLOG") == "TRUE"
	Config.Gzip = os.Getenv("SCOPS_GZIP") == "TRUE"
	Config.RemoteUrl = checkEnv("SCOPS_REMOTE", "localhost/test/push")
	Config.AuthToken = checkEnv("SCOPS_TOKEN", "testToken")
	Config.Plugin = checkEnv("SCOPS_PLUGIN", "")

	cycle := checkEnv("SCOPS_CYCLE", "30s")
	d, err := time.ParseDuration(cycle)
	if err != nil {
		logger.Log.Panicf("Error with SCOPS_CYCLE environment variable: %v", err)
	}
	Config.Cycle = d
}

func loadFlagConfig() {
	clockPtr := flag.String("testclock", "", "Use a fake clock at time given. Format 20060102-1504")
	debugPtr := flag.Bool("debug", false, "Enable debug messages")
	sysPtr := flag.Bool("syslog", false, "Redirect messages to syslog")
	gzipPtr := flag.Bool("gzip", false, "Gzip requests")
	remotePtr := flag.String("remote", "", "Remote URL to send messages to")
	authPtr := flag.String("token", "", "Authorization token")
	pluginPtr := flag.String("plugin", "", "Plugin to use to get the data")
	cyclePtr := flag.Duration("cycle", 0, "Cycle duration")

	flag.Parse()

	flagset := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { flagset[f.Name] = true })

	if flagset["plugin"] {
		Config.Plugin = *pluginPtr
	}
	if Config.Plugin == "" {
		logger.Log.Panicf("No plugin set")
	}

	if flagset["testclock"] {
		testTime, err := time.Parse("20060102-1504", *clockPtr)
		if err != nil {
			logger.Log.Panicf("Error with testclock command line arguments: %v", err)
		}
		clock.SetDefaultClock(clock.NewFakeClockAt(testTime))
	}
	if flagset["debug"] {
		Config.Debug = *debugPtr
	}
	if flagset["syslog"] {
		Config.Syslog = *sysPtr
	}
	if flagset["gzip"] {
		Config.Gzip = *gzipPtr
	}
	if flagset["remote"] {
		Config.RemoteUrl = *remotePtr
	}
	if flagset["token"] {
		Config.AuthToken = *authPtr
	}
	if flagset["cycle"] {
		Config.Cycle = *cyclePtr
	}
}

func checkEnv(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}
