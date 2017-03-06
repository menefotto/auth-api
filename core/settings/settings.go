package settings

// Required field by operation

var CREATE_USER_FIELD_REQUIRED = []string{"Username", "Email", "Password"}

var UPDATE_USER_FIELD_REQUIRED = []string{"Email"}

var GET_USER_FIELD_REQUIRED = []string{"Email"}

// Variuos

var PROJECTID = "boardsandwater" // related to the deployment on google cloud

var OBFUSCATED_FIELDS = map[string]string{"Password": "default"}

// Crypto settings

var JWT_LOGIN_DELTA = 3

var JWT_ACTIVATION_DELTA = 24 * 7

var JWT_PASSWORD_DELTA = 12

var HMAC_SECRET = "v97iv7m0mi98BmPoGK81S7sKt1O1UBTV"

var CRYPTO_SECRET = "AES256Key-32Characters1234567890"

var NONCE = "bb8ef84243d2ee95a41c6c57"

// Email settings

var EMAIL_PASSWORD = "Stovari1985"

var EMAIL_SENDER = "locci.carlo.85@gmail.com"

var EMAIL_SMTP = "smtp.gmail.com"

var EMAIL_PORT = "587"

var EMAIL_TEMPLATE_DIR = "templates/emails"

// Api specific settings

var API_URL = "http://localhost:8080/"
