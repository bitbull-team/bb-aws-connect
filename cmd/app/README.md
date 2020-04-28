## Applications commands

This category contain commands used to install, run, build or configure applications.

### Application types

This CLI try to guess what kind of application is checking files existence or composer/npm dependencies installed. Supported application are:

Simple projects:

* Composer - a simple PHP application
* NPM - a simple NodeJS application
* YARN - a simple NodeJS application that use yarn
* GO - a GO project

Frameworks:

* Laravel
* Wordpress
* Magento
* Magento + Wordpress
* Magento2

### Application configuration

Configuration file (by default `.bb-cli.yml`) ha the following configurations:
```yml
app:
  type: "go" # Override application type
```
