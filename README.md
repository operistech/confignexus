# configNexus

## Table of Contents
- [Introduction](#introduction)
- [Requirements](#requirements)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running](#running)
- [Usage](#usage)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

---

## Introduction

`configNexus` is a configuration management server designed to manage and serve configurations based on domain patterns. The project is written in Go and utilizes libraries like `zerolog` for logging, `viper` for configuration, and `http` for web serving.

### Features
- HTTP to HTTPS redirection
- Support for YAML template processing
- Dynamic configurations based on domain patterns
- Environment-variable based settings
- Graceful shutdown and cleanup
- Debug logging
- Extensible design to integrate with other systems

## Requirements
- Go 1.16+
- git (For managing repositories)

## Installation

Clone the repository and build the application:


    git clone https://github.com/your-username/configNexus.git
    cd configNexus
    go build -o configNexus


## Configuration

The application looks for a `config.yaml` file at the root directory or `/etc/configNexus`. You can also use environment variables prefixed with `CN` to configure the application.

| Config Key    | Environment Variable | Default   | Description                           |
|---------------|----------------------|-----------|---------------------------------------|
| HTTPPort      | CN_HTTPPORT          | 9000      | HTTP port to listen on                |
| HTTPSPort     | CN_HTTPSPORT         | 9443      | HTTPS port to listen on               |
| ListenAddress | CN_LISTENADDRESS     | localhost | Address to listen on                  |
| HTTPEnabled   | CN_HTTPENABLED       | true      | Enable/Disable HTTP                   |
| HTTPRedirect  | CN_HTTPREDIRECT      | true      | Enable/Disable HTTP to HTTPS redirect |
| RepoAddress   | CN_REPOADDRESS       | false     | Set the path to the Config repository |
| RepoBranch    | CN_REPOBRANCH        | main      | The default git branch to monitor     |


For example, to set the HTTPS port:


    export CN_HTTPSPORT=9445


## Running

Run the application using:


    ./configNexus


## Usage

### Configuration Data

The project relies on various YAML configuration files to operate. These files are critical for the initialization and the dynamic behavior of the application. Here's a breakdown:

#### Required Files

- `all.yaml`: Holds global settings that apply to all instances, datacenters, and devices.
- `domains_regex.yaml`: Contains regular expressions for domain patterns.

#### Optional Files

- `<function>.yaml`: Holds settings specific to a particular function. E.g., `web.yaml`, `db.yaml`.
- `<datacenter>.yaml`: Holds settings specific to a particular datacenter. E.g., `us-east.yaml`, `eu-central.yaml`.
- `<hostname>.yaml`: Holds settings specific to a particular device identified by its hostname.

#### Order of Precedence

1. **Global Settings (`all.yaml`)**: These settings apply universally unless overridden.
2. **Function Settings (`<function>.yaml`)**: Specific to different functions like web, db, etc., and override global settings.
3. **Datacenter Settings (`<datacenter>.yaml`)**: These are specific to each datacenter and override both global and function settings.
4. **Device Settings (`<hostname>.yaml`)**: These are the most specific and will override all the above.

When the application is initialized or reconfigured, it merges settings in this order to derive the final settings. This way, specific configurations can be applied granularly, allowing for flexible system behavior.

#### Automatic Repo Monitoring

The application is configured to automatically monitor the associated repository for any changes. It will perform a `git pull` every 20 minutes to ensure that the latest configuration and code are always in sync with the deployed instance. This feature enables seamless updates without requiring manual intervention.



## Testing

Run unit tests:


    go test ./...


## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
