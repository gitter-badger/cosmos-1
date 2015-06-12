rm -rf ./src/github.com/influxdb/influxdb
mkdir -p ./src/github.com/influxdb/influxdb
git clone git@github.com:influxdb/influxdb tmp
mv ./tmp/client ./src/github.com/influxdb/influxdb/
mv ./tmp/influxql ./src/github.com/influxdb/influxdb/
rm -rf ./tmp

rm -rf ./src/github.com/gorilla/context
mkdir -p ./src/github.com/gorilla
git clone git@github.com:gorilla/context ./src/github.com/gorilla/context
rm -rf ./src/github.com/gorilla/context/.git

rm -rf ./src/github.com/gorilla/mux
mkdir -p ./src/github.com/gorilla
git clone git@github.com:gorilla/mux ./src/github.com/gorilla/mux
rm -rf ./src/github.com/gorilla/mux/.git