*** Settings ***
# For more detail init setting, please refer to ../res/init.robot
Resource            ../res/init.robot
Test Timeout        ${TEST_API_TIMEOUT}

# For more detail init setting, please refer to ../res/valueset.robot
# Timeout Period: 10 times, Retry Strict Period: 1 sec,
Suite Setup         Get Test Value  ${ENV}
Suite Teardown      Release Test Value

*** Test Cases ***
# Automated vendor testing from YAML - Complete automation
[] [RAT] [VENDOR] [AUTO] Automated vendor testing from YAML
  [Tags]              testrailid=     RAT             VENDOR
  [Documentation]  Fully automated vendor testing using actual config.yaml file
  ...              Loads vendor configuration from deploy/rec-vendor-api/secrets/config.yaml
  ...              Tests all vendors in config with complete automation:
  ...              - Dynamic endpoint generation: /r/{vendor_name}
  ...              - Auto-selected test dimensions: 300x300, 1200x627, 1200x600
  ...              - Base64 encoding validation for click_id parameters
  ...              - Response structure validation (array of products)
  ...              - Product patch verification with tracking parameters

  # Load YAML configuration from actual config file
  # Note: size={width}x{height} parameters are auto-selected from predefined sizes (300x300, 1200x627, 1200x600)
  ${config_path} =        Set Variable    ${CURDIR}/../../../../deploy/rec-vendor-api/secrets/config.yaml
  ${yaml_content} =       Load vendor config from file    ${config_path}

  # Run complete automated testing for all vendors in YAML
  Test vendors from yaml configuration  ${yaml_content}


# Basic vendor API healthz test
[C4977913] [RAT] [VENDOR] [HEALTHZ] Test vendor API healthz endpoint
  [Tags]  testrailid=4977913  RAT     VENDOR  HEALTHZ
  [Documentation]  Test the vendor API healthz endpoint to ensure it returns proper status and message

  Given I have an vendor session
  When I would like to set the session under vendor endpoint with  endpoint=/healthz
  Then I would like to check status_code should be "200" within the current session
