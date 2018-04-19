#!/bin/bash
ROOT=../../

# remove vendor directory
rm -rf $ROOT/vendor/*
cp ./vendor_empty.json $ROOT/vendor/vendor.json

echo "fetching cobra"
govendor fetch github.com/spf13/cobra

echo "fetching viper"
govendor fetch github.com/spf13/viper

echo "fetching validators"
govendor fetch github.com/asaskevich/govalidator@v9

echo "fetching gouuid"
govendor fetch github.com/nu7hatch/gouuid

echo "fetching logrus"
govendor fetch github.com/sirupsen/logrus@v1.0.4

echo "fetching echo"
govendor fetch github.com/labstack/echo@v3.2.1
govendor fetch github.com/labstack/echo/middleware@v3.2.1

echo "fetching gorm"
govendor fetch github.com/jinzhu/gorm
govendor fetch github.com/jinzhu/gorm/dialects/sqlite
govendor fetch github.com/jinzhu/gorm/dialects/postgres


