# /bin/sh

exec nginx -g "daemon off;" &
exec parrot serve --config=/config/config.yaml
