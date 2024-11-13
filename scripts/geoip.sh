


FILE="internal/plugins/collectors/geo/GeoLite2-City.mmdb"

# # https://git.io/GeoLite2-City.mmdb
# # https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-City.mmdb

wget "https://git.io/GeoLite2-City.mmdb" -P internal/plugins/collectors/geo

if [ -f ! $FILE ]; then
   wget "https://git.io/GeoLite2-City.mmdb" -P internal/plugins/collectors/geo
else
   echo "File $FILE exists"
fi

if [ -f ! $FILE ]; then
   wget "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-City.mmdb" -P internal/plugins/collectors/geo
else
   echo "File $FILE exists"
fi

