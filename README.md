
# Go authentication/autorization api

[![GoDoc](https://godoc.org/github.com/spf13/viper?status.svg)]()

## What is auth-api
Exactly what the name says is an api that aims to provide a multi-cloud backed api to manage user profiles,
fully featured with registration flow, password reset flow, user activation flow, user profile update. Cloud providers currently supported google or amazon.

## How do I add another cloud vendor?
Under the core/databases there is defined the interface that the database should provide, to add another provider simply add in implementation the satisfies the Db interface and select that implementation under the database config file and also add to the core/managers constructor New your backend implementation name.

## How is it any different from others auth framework?
Auth-api aims to be more efficient and easier to scale horizontally ( just add more instances ), under the amazon could it should be able to register roughly 15-20 user par second and under google 5-10 user par seconds ( more official benchmarks, this are the currently roughly measured performance ). It also aims to be more efficient in terms of open sessions, instead of the common approach of keeping all the users session in memory, sessions are avoided altogether, instead jwt-tokens are used and to invalidated them a blacklist is kept in a separate database table, therefore reducing the memory foot-print to zero.
Authentication is done matching the cookie validity with the xrsf-token therefore not requiring the database to be hit. Auth-api also provides throttling, and timeouts for request exceeding the configured timeout.
Everything is highly configurable except through a toml based configuration, configuration can be changed on the fly and server will automatically reload itself.

#### Under active developement not yet finished...
Go doc link is currently invalid, test coverage is roughly at 10%, and it's not even remotely usable...Third party registration,login provider to be added.
