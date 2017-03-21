# AUTH-API authentication/autorization api
[![GoDoc](https://godoc.org/github.com/wind85/auth-api?status.svg)](https://godoc.org/github.com/wind85/auth-api)
[![Build Status](https://travis-ci.org/wind85/auth-api.svg?branch=master)](https://travis-ci.org/wind85/auth-api)
[![Coverage Status](https://coveralls.io/repos/github/wind85/auth-api/badge.svg?branch=master)](https://coveralls.io/github/wind85/auth-api?branch=master)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


## What is auth-api
Exactly what the name says is an api that aims to provide a multi-cloud backed 
api to manage user profiles, fully featured with registration flow, password 
reset flow, user activation flow, user profile update. 
Cloud providers currently supported google or amazon.

## How do I add another cloud vendor?
Under the core/databases there is defined the interface that the database should 
provide, to add another provider simply add in implementation the satisfies the Db 
interface and select that implementation under the database config file and also 
add to the core/managers constructor New your backend implementation name.

## How is it any different from others auth framework?
Auth-api aims to be more efficient and easier to scale horizontally ( just add 
more instances ), under the amazon it should be able to register roughly 15-20 
user par second and under google 5-10 user par seconds ( more official benchmarks 
need to be run this somehow a conservative consideration ). It also aims to be 
more efficient in terms of sessions, instead of the common approach of keeping 
all the users session in memory, sessions are avoided altogether, instead jwt 
tokens are used and invalidated by adding them to a blacklist is kept in a separate
database table, therefore reducing the memory foot-print to zero.
Authentication is done matching the cookie (encrypted with jwt inside) with the 
xrsf-token therefore not requiring the database to be hit. Auth-api also provides
throttling, and timeouts for request exceeding the configured timeout.
Everything is highly configurable through a toml based configuration, configuration
can be changed on the fly and the server will automatically reload itself, also 
the server is can is setup to reload itself on SIGHUP and shutdown itself on 
SIGTERM.

#### Under active developement not yet finished...
Go doc link is currently invalid, test coverage is roughly at 10%, and it's 
not even remotely usable...Third party registration,login provider to be added.

#### Philosophy
This software is developed following the "mantra" keep it simple, stupid or 
better known as KISS. Something so simple like configuration files should not
required over engineered solutions. Though it provides most of the functionality 
needed by generic configuration files, and most important of all meaning full 
error messages.

#### Disclaimer
This software in alpha quality, don't use it in a production environment, it's 
not even completed.

#### Thank You Notes
I should thanks the gopheraccademy since I wrote this little lib after I read an 
article about configuration and tokenisers.
