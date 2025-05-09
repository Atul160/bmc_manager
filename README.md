```

```

# ECC BMC manager

The ECC BMC manager is a centralized tool designed to manage hardware from multiple vendors such as Dell, HPE, Nutanix, and Lenovo. It provides a unified interface to interact with BMC (Baseboard Management Controllers) for performing tasks such as authentication, power management, retrieving system information, and firmware details. This document outlines the API usage, deployment, and additional details for setting up and using the BMC Manager.

Features
    Centralized management of multiple hardware BMCs.
    Vendor support for Dell, HPE, Nutanix, Lenovo, and more.
    RESTful APIs for seamless integration.
    Secure authentication using JWT.
    Comprehensive Swagger documentation for API usage.
    Lightweight and production-ready deployment using Docker.

## Authors

- Atul Pahlazani

## Target Audience

    Internal Teams, such as Compute (Winserv/Linux), CCC, vertical1 Dispatch team, etc.

## File Structure:

```
├── Dockerfile
├── README.md
├── api
│   ├── firmware.go
│   ├── power.go
│   ├── system_info.go
│   ├── token_handler.go
│   └── types.go
├── bmc
│   ├── bmc_factory.go
│   ├── bmc_interface.go
│   ├── dell_idrac.go
│   ├── hpe_ilo.go
│   ├── lenovo_imm.go
│   ├── lenovo_xcc.go
│   └── nutanix_ipmi.go
├── certs
│   ├── ABCIntermediateCA01-SHA256-2021.crt
│   ├── ABCIntermediateCA01-SHA256.crt
│   ├── ABCIssuingCA-2FA-02-SHA256-2021.crt
│   ├── ABCIssuingCA-2FA-03-SHA256-2021.crt
│   ├── ABCIssuingCA-2FA-04-SHA256-2021.crt
│   ├── ABCIssuingCA-TLS-01-SHA256-1.crt
│   ├── ABCIssuingCA-TLS-01-SHA256-2.crt
│   ├── ABCIssuingCA-TLS-01-SHA256-3.crt
│   ├── ABCIssuingCA-TLS-01-SHA256-4.crt
│   ├── ABCIssuingCA-TLS-01-SHA256-5.crt
│   └── ABCRootCA-SHA256.crt
├── config
│   └── config.go
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── logs
│   └── app.log
├── main.go
├── middleware
│   ├── adauth.go
│   ├── jwtauth.go
│   └── logging.go
├── services
│   ├── firmware_service.go
│   ├── powerservices.go
│   ├── system_info_service.go
│   └── types.go
└── utils
    ├── async.go
    ├── auth_utils.go
    ├── bmc_utils.go
    ├── http_utils.go
    ├── logging.go
    └── other_utils.go
```

## Key Directories

    api: Contains API endpoint implementations.
    bmc: Handles vendor-specific implementations of BMC interfaces.
    certs: Contains SSL certificates for secure communication.
    config: Application configuration management.
    docs: Swagger documentation.
    middleware: Implements authentication, logging, and other middleware utilities.
    services: Business logic for API endpoints.
    utils: Utility functions used throughout the application.

## API Reference

#### Authenticate and obtain a JWT token for further API usage.

```http
    POST /bmc/auth 
```

    Usage:
        Provide credentials for BMC authentication.

    Headers:

```
    {         "authorization": "Basic <Basic AuthToken(Domain Credentials)>"         }
```

    Response:

```
    {   
        "token": "<JWT Token>",   
        "expiry": "<timestamp>"   
    }  
```

    Example:

```
    curl -X POST http://localhost:8086/bmc/auth \       
    -H "Content-Type: application/json" \       
    -H 'authorization: Basic <Basic AuthToken>' \       
    -d ''  
```

#### Perform power actions such as power-on, power-off, reset and bmcreset

```http
        POST /bmc/power   
```

    Usage:
        Control the power state of a hardware unit.
        The bmctype field is optional; if omitted, the system attempts to auto-detect the vendor.

```
Request Body:

    {   "bmctype": "<optional: string>",   
        "ipaddress": "<string>",   
        "action": "<string: on|off|reset>"   
    }
```

```
Response:

    {
    "result": [
        {
        "ip_address": "<ip_address><String>",
        "action": "<action><String>",
        "success": true | false
        }
    ]
    }
```

```
Example:

    curl -X POST http://localhost:8086/bmc/power \       
    -H "Content-Type: application/json" \       
    -H "Authorization: Bearer <JWT Token>" \       
    -d '{"ipaddress": "192.168.1.10", "action": "on"}'
```

#### Get system-level information from the BMC

```http
    POST /bmc/systeminfo   
```

    Usage:
        Fetch details about the hardware system.
        The bmctype field is optional; if omitted, the system attempts to auto-detect the vendor. 

```
Request Body:

    {   "bmctype": "<optional: string>",   
        "ipaddress": "<string>"   
    }  
```

```
Response:

    {
        "system_info": [
            {
            "device": "<Target Device><String>",
            "error": "<Error in case of any failure><String>"
            },
            {
            "device": "<Target Device><String>",
            "health": "<String>",
            "manufacturer": "<String>",
            "powerstate": "<String>",
            "model": "<String>",
            "biosversion": "<version>",
            "serialnumber": "<String>",
            "hostname": "<Server OS FQDN><String>",
            "memory": {},
            "cpu": {}
            },
        ]
    }
```

```
Example:

    curl -X POST http://localhost:8086/bmc/systeminfo \       
    -H "Content-Type: application/json" \      
    -H "Authorization: Bearer <JWT Token>" \       
    -d '{"ipaddress": "192.168.1.10"}'  
```

#### Get Firmware details from the BMC

```http
    POST /bmc/firmwareinfo   
```

    Usage:
        Retrieve detailed firmware version and component information.
        The bmctype field is optional; if omitted, the system attempts to auto-detect the vendor.

```
Request Body:

    {   
        "bmctype": "<optional: string>",   
        "ipaddress": "<string>"   
    }   
```

```
Response:

    {
    "firmware_info": [
        {
        "ip_address": "<Target Device><String>",
        "data":{<Firmware Inventory for all components>}
        }
    ]
    }
```

```
Example:

    curl -X POST http://localhost:8086/bmc/firmwareinfo \       
    -H "Content-Type: application/json" \       
    -H "Authorization: Bearer <JWT Token>" \       
    -d '{"ipaddress": "192.168.1.10"}'  
```

## Swagger Documentation

    Access detailed API documentation and testing interface at:
    http://localhost:8086/swagger/index.html

## Deployment

- Build the Application
  Use the provided Dockerfile to build and run the application:

```
   docker build  -t gtocdc/ec-bmc -f Dockerfile .
```

- Run the Container

```
   docker run -d -t --restart=always --name ec-bmc --net=web -p 8086:8086 -v <path to secrets directory>:/app/secrets gtocdc/ec-bmc:latest
```

- Verify API is Running
  Access the Swagger docs or test the endpoints using tools like curl or Postman.

## Logs

    Application logs are stored in the logs directory and can be monitored for debugging or audit purposes.

## Contact

    For issues or contributions, please contact the development team or raise a GitHub issue.

> This README provides a comprehensive overview of the Hardware Console Manager. Feel free to expand or modify it based on project needs
