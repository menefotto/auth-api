# YACP 

[![GoDoc](https://godoc.org/github.com/wind85/confparse?status.svg)](https://godoc.org/github.com/wind85/confparse)
[![Build Status](https://travis-ci.org/wind85/confparse.svg?branch=master)](https://travis-ci.org/wind85/confparse)
[![Coverage Status](https://coveralls.io/repos/github/wind85/confparse/badge.svg?branch=master)](https://coveralls.io/github/wind85/confparse?branch=master)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
### Yet Another Configuration Parser
This is a small package the provides a ini style configuration parser. This is 
what is allowed:

- Comments start, either with the "#" or ":" anything after it, till newline is ignored
- Sections are written like the following [default] and contain a map of key values,
  anything between square brackets is a valid section.
- Key and values are like "ip=192.168.10.1" ,the separator is "=" otherwise will
  not be considered a key value.
- The Parser can handle bool, int and floats (both 64bit), strings and string slices,
  as long as they are divided by columns.
- Empty lines are ignored, white spaces are ignored as well.

### How to use it
Pretty simple, there is only one method to create a new parser just call 
```
  ini, err := confparse.New("config-name.whatever")
```
It isn't name sensible any valid name can be passed. Then any of the valid supported 
values can be retrieved like so:
```
  value ,err := init.GetInt("sectionname.valuename")
  value ,err := init.GetFloat("sectionname.valuename")
  value ,err := init.GetSlice("sectionname.valuename")
  value ,err := init.GetString("sectionname.valuename")
  value ,err := init.GetBool("sectionname.valuename")
  section ,err := init.GetSection("sectionname")
```
There is also a Watch function that listen if any changes are made to the configuration
file, if it does find some, the configuration get reloaded and parsed every time a change
occurs. You can call it like so:
```
  ini.Watch()
```
If you need to run a function every time an event occurs on the file watched there the:
```
 ini.OnConfChange(run func(ev fsnotify.Event))
```
And that's pretty much about it.

#### Philosophy
This software is developed following the "mantra" keep it simple, stupid or better known as
KISS. Something so simple like configuration files should not required over engineered solutions.
Though it provides most of the functionality needed by generic configuration files, and most
important of all meaning full error messages.

#### Disclaimer
This software in alpha quality, don't use it in a production environment, it's not even
completed.

#### Thank You Notes
I should thanks the gopheraccademy since I wrote this little lib after I read an article about
configuration and tokenisers.
