package settings

var CREATE_USER_FIELD_REQUIRED = []string{"Username", "Email", "Password"}

var UPDATE_USER_FIELD_REQUIRED = []string{"Email"}

var GET_USER_FIELD_REQUIRED = []string{"Email"}

var OBFUSCATED_FIELDS = map[string]string{"Password": "default"}

var JWTExpirationDelta = 3

var HMAC_SECRET = "v97iv7m0mi98BmPoGK81S7sKt1O1UBTV"

var CRYPTO_SECRET = "AES256Key-32Characters1234567890"

var NONCE = "bb8ef84243d2ee95a41c6c57"

var PROJECTID = "boardsandwater"

var EMAIL_PASSWORD = "######"

var EMAIL_SENDER = "locci.carlo.85@gmail.com"

var EMAIL_SMTP = "smtp.gmail.com"

var EMAIL_PORT = "587"
