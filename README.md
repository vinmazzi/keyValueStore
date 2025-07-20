For the app to work set the following environment variables:


FRONTEND_TYPE="rest"
HTTP_IDLE_TIMEOUT=1
HTTP_LISTEN_PORT=":8080"
TRANSACTION_LOGGER_BACKEND="postgres"
POSTGRES_HOST="localhost"
POSTGRES_USERNAME="postgres"
POSTGRES_PASSWORD="mysecretpassword"
POSTGRES_SSLMODE="disable"
POSTGRES_DATABASE="key_value_store"
ENCODER_TYPE="base64"
