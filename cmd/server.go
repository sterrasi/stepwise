package cmd

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"

	"github.com/sterrasi/stepwise/logging"
	"github.com/sterrasi/stepwise/users"
	"github.com/sterrasi/stepwise/util"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	applicationName = "stepwise"
)

type AutoTLSConfig struct {
	CertCacheDir  string   `mapstructure:"cert-cache-dir"`
	HostWhiteList []string `mapstructure:"host-white-list"`
}

type TLSConfig struct {
	Cert string `mapstructure:"cert"`
	Key  string `mapstructure:"key"`
}

type ServerConfig struct {
	Address string         `mapstructure:"address"`
	AutoTLS *AutoTLSConfig `mapstructure:"auto-tls"`
	TLS     *TLSConfig     `mapstructure:"tls"`
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start the Stepwise web application",
	Long:  "Starts an instance of the stepwise web application.",

	Run: func(cmd *cobra.Command, args []string) {

		// logging
		if err := initLogging(); err != nil {
			panic(err.Error())
		}
		defer logging.DeinitializeLogging()

		// database
		db, err := initDatabase()
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		// server
		e := echo.New()

		// Middleware
		e.Pre(middleware.RemoveTrailingSlash())
		e.Pre(middleware.HTTPSRedirect())

		e.Use(middleware.Gzip())
		e.Use(middleware.Secure())
		e.Use(logging.LoggerMiddleware())
		e.Use(middleware.Recover())

		// Register Users API
		usersConfig := &users.Config{}
		if err := viper.UnmarshalKey("users", usersConfig); err != nil {
			panic(err.Error())
		}
		users.Register(e.Group("/users"), db, usersConfig)

		e.GET("/", hello)

		// Start server
		e.Logger.Fatal(startServer(e))
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// initializes the GORM database
func initDatabase() (*gorm.DB, error) {
	databaseConfig := &util.DatabaseConfig{}

	var err error
	var db *gorm.DB

	if err = viper.UnmarshalKey("database", databaseConfig); err != nil {
		return nil, err
	}

	db, err = util.InitDatabase(databaseConfig)
	if err != nil {
		return nil, err
	}

	if databaseConfig.Migrate {
		db.AutoMigrate(&users.User{})
	}
	return db, nil
}

func startServer(e *echo.Echo) error {

	serverConfig := &ServerConfig{}

	if err := viper.UnmarshalKey("server", serverConfig); err != nil {
		return err
	}

	// check for autoTLS
	if serverConfig.AutoTLS != nil {
		logrus.Infof("Starting server using AutoTLS")

		// apply a white list policy
		if len(serverConfig.AutoTLS.HostWhiteList) > 0 {
			e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(serverConfig.AutoTLS.HostWhiteList...)
		}

		// Cache certificates in a local directory
		e.AutoTLSManager.Cache = autocert.DirCache(serverConfig.AutoTLS.CertCacheDir)

		// start server
		return e.StartAutoTLS(serverConfig.Address)

	} else if serverConfig.TLS != nil {
		logrus.Infof("Starting server using TLS")

		if serverConfig.TLS.Cert == "" {
			return fmt.Errorf("Certificate not provided for TLS enabled server")
		}
		if serverConfig.TLS.Key == "" {
			return fmt.Errorf("Key not provided for TLS enabled server")
		}
		return e.StartTLS(serverConfig.Address, serverConfig.TLS.Cert, serverConfig.TLS.Key)

	}

	return fmt.Errorf("Either the [server.auto-tls] or [server.tls] section must be filled out")
}

func initLogging() error {
	logConfig := &logging.LogConfig{}

	if err := viper.UnmarshalKey("logging", logConfig); err != nil {
		return err
	}
	if err := logging.InitLogging(logConfig); err != nil {
		return err
	}
	if logConfig.LogStartupInfo {
		logStartupInfo()
	}
	return nil
}

func logStartupInfo() {
	logrus.Info("starting Stepwise...")
}
