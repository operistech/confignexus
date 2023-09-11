# configNexus

## Table of Contents
- [Introduction](#introduction)
- [Requirements](#requirements)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running](#running)
- [Usage](#usage)
- [Example](#example)
- [Testing](#testing)
- [License](#license)

---

## Introduction

ConfigNexus serves as a single source of truth for configuration management in environments where multiple systems, like Ansible and Puppet, are in operation alongside additional custom tools. By centralizing configuration data, it facilitates a streamlined and coherent configuration strategy, mitigating potential conflicts and redundancies inherent in using diverse management systems concurrently.

This configuration management server is built to manage and serve configurations based on distinct domain patterns, offering a dynamic and adaptive approach to configuration management.

The project leverages the Go programming language, drawing on powerful libraries such as:
- **zerolog** for optimized logging
- **viper** for configuration handling
- **http** for web serving

By integrating ConfigNexus into your workflow, you empower your team to maintain consistency across various platforms, simplifying configuration processes and enhancing system reliability and efficiency.

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

### Running the configNexus Docker Container

To run the `configNexus` Docker container, you can use the following command which will run the container in the background and map port 9443 in the container to port 9443 on your host system:

```sh
docker run -d -p 9443:9443 operistech/confignexus
```


To run the container with environment variables, allowing you to configure the container at runtime, use a command like this:

```sh
docker run -d -p 9443:9443 -e "ENV_VAR_NAME=env_var_value" operistech/confignexus
```


### Notes
- Replace `ENV_VAR_NAME` and `env_var_value` with the actual name and value of the environment variable you want to set.

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

#### domains_regex.yaml

Hold any number of domain regex patterns that can match domains and named groups that can subsequently be used in template files.
The file contains a single key: `regex_patterns`, which is an array of objects. Each object has the following attributes:

- `name`: A unique identifier for the regular expression pattern.
- `regex`: The actual regular expression pattern.


```yaml
regex_patterns:
  - name: "Pattern1"
    regex: "^(?P<Datacenter>[a-z]{3})(?P<Function>[a-z]+)(?P<Instance>\\d+)\\.mgt\\.prod\\.example\\.com$"
  - name: "Pattern2"
    regex: "^(?P<Function>[a-z]+)(?P<Instance>\\d+)\\.(?P<Datacenter>[a-z]{3})\\.example\\.com$"
```

### Utilization in Templates
Parameterizing Named Groups

Named groups captured from these regular expressions can be parameterized and used in the ConfigNexus templates. For example, if a domain name matches Pattern1, you can use the captured Datacenter, Function, or Instance in a template like this:

```yaml

datacenter: {{ .Datacenter }}
function: {{ .Function }}
instance: {{ .Instance }}
```
By using these named groups in your templates, you can create highly dynamic configurations that adapt based on the domain name being processed.

#### Automatic Repo Monitoring

The application is configured to automatically monitor the associated repository for any changes. It will perform a `git pull` every 20 minutes to ensure that the latest configuration and code are always in sync with the deployed instance. This feature enables seamless updates without requiring manual intervention.


## Example
configNexus has an associated testConfigdata repository that when ran with confignexus will allow for some test domains to be fed through to generate a full json return of configuration data.

    CN_REPOADDRESS=https://github.com/operistech/testconfigdata.git go run cmd/server/main.go

Then server can be queried like:

    curl -k https://localhost:9443/details/slcpostgresql1.mgt.prod.example.com | python -m json.tool

which will return the following json:

    {
        "appserver": "appserver1.slc.example.com",
        "contact": {
            "phone_number": "555-555-1234",
            "security": "datacentersecurity@example.com",
            "support": "datacentersupport@example.com"
        },
        "datacenter": "slc",
        "environment": "staging",
        "features": {
            "firewall": true,
            "vpn": false
        },
        "function": "postgresql",
        "instance": 1,
        "ip_address": "192.168.0.7",
        "ip_range": "192.168.0.0/24",
        "ldap": "slcldap1.mgt.prod.example.com",
        "owner": "superappteam",
        "ports": [
            8102,
            8103
        ]
    }


## Testing

Run unit tests:


    go test ./...


## License

This project is licensed under the GPLv3 License - see the [LICENSE.md](LICENSE.md) file for details.
