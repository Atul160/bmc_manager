ECC BMC manager

The ECC BMC manager is a centralized tool designed to manage hardware from multiple vendors such as Dell, HPE, Nutanix, and Lenovo. It provides a unified interface to interact with BMC (Baseboard Management Controllers) for performing tasks such as authentication, power management, retrieving system information, and firmware details. This document outlines the API usage, deployment, and additional details for setting up and using the BMC Manager.

Features
    Centralized management of multiple hardware BMCs.
    Vendor support for Dell, HPE, Nutanix, Lenovo, and more.
    RESTful APIs for seamless integration.
    Secure authentication using JWT.
    Comprehensive Swagger documentation for API usage.
    Lightweight and production-ready deployment using Docker.

Target Audience
    Internal Teams, such as Compute (Winserv/Linux), CCC, Stores Dispatch team, etc.

File Structure:

Overview:
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
│   ├── WalmartIntermediateCA01-SHA256-2021.crt
│   ├── WalmartIntermediateCA01-SHA256.crt
│   ├── WalmartIssuingCA-2FA-02-SHA256-2021.crt
│   ├── WalmartIssuingCA-2FA-03-SHA256-2021.crt
│   ├── WalmartIssuingCA-2FA-04-SHA256-2021.crt
│   ├── WalmartIssuingCA-TLS-01-SHA256-1.crt
│   ├── WalmartIssuingCA-TLS-01-SHA256-2.crt
│   ├── WalmartIssuingCA-TLS-01-SHA256-3.crt
│   ├── WalmartIssuingCA-TLS-01-SHA256-4.crt
│   ├── WalmartIssuingCA-TLS-01-SHA256-5.crt
│   └── WalmartRootCA-SHA256.crt
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

Key Directories
    api: Contains API endpoint implementations.
    bmc: Handles vendor-specific implementations of BMC interfaces.
    certs: Contains SSL certificates for secure communication.
    config: Application configuration management.
    docs: Swagger documentation.
    middleware: Implements authentication, logging, and other middleware utilities.
    services: Business logic for API endpoints.
    utils: Utility functions used throughout the application.


API Endpoints

1. Authentication
    Endpoint: POST /bmc/auth
    Description: Authenticate and obtain a JWT token for further API usage.

    Headers:
    {
    "authorization": "Basic <Basic AuthToken(HomeOffice Credentials)>"
    }

    Response:
    {
    "token": "<JWT Token>",
    "expiry": "<timestamp>"
    }

    Usage:
        Provide credentials for BMC authentication.
        The bmctype field is optional; if omitted, the system attempts to auto-detect the vendor.

    Example:
        curl -X POST http://localhost:8081/bmc/auth \
            -H "Content-Type: application/json" \
            -H 'authorization: Basic <Basic AuthToken>' \
            -d ''

2. Power Management
    Endpoint: POST /bmc/power
    Description: Perform power actions such as power-on, power-off, or reset.

    Body:
        {
        "bmctype": "<optional: string>",
        "ipaddress": "<string>",
        "action": "<string: on|off|reset>"
        }

    Response:
        {
        "status": "<string>",
        "message": "<string>"
        }

    Usage:
        Control the power state of a hardware unit.
    Example:
        curl -X POST http://localhost:8081/bmc/power \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer <JWT Token>" \
    -d '{"ipaddress": "192.168.1.10", "action": "on"}'

3. System Information
    Endpoint: POST /bmc/systeminfo
    Description: Retrieve system-level information from the BMC.

    Request Body:
        {
        "bmctype": "<optional: string>",
        "ipaddress": "<string>"
        }

    Response:
        {
        "manufacturer": "<string>",
        "model": "<string>",
        "serial": "<string>",
        "bios_version": "<string>"
        }

    Usage:
        Fetch details about the hardware system.

    Example:
        curl -X POST http://localhost:8081/bmc/systeminfo \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer <JWT Token>" \
    -d '{"ipaddress": "192.168.1.10"}'

4. Firmware Information
    Endpoint: POST /bmc/firmwareinfo
    Description: Retrieve firmware details from the BMC.

    Request Body:
        {
        "bmctype": "<optional: string>",
        "ipaddress": "<string>"
        }

    Response:
        {
        "firmware_version": "<string>",
        "release_date": "<string>",
        "components": [
            {
            "name": "<string>",
            "version": "<string>"
            }
        ]
        }

    Usage:
        Retrieve detailed firmware version and component information.

    Example:
        curl -X POST http://localhost:8081/bmc/firmwareinfo \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer <JWT Token>" \
    -d '{"ipaddress": "192.168.1.10"}'

Swagger Documentation
    Access detailed API documentation and testing interface at:
    http://server:8081/swagger/index.html

Deployment
1. Build the Application
    Use the provided Dockerfile to build and run the application:
    docker build -t ecc-bmc .

2. Run the Container
    docker run -p 8081:8081 ecc-bmc

3. Verify API is Running
    Access the Swagger docs or test the endpoints using tools like curl or Postman.

Logs
    Application logs are stored in the logs directory and can be monitored for debugging or audit purposes.

Contact
    For issues or contributions, please contact the development team or raise a GitHub issue.

This README provides a comprehensive overview of the Hardware Console Manager. Feel free to expand or modify it based on project needs