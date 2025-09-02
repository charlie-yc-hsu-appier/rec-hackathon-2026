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
[C4983449] [RAT] [VENDOR] [AUTO] Automated vendor testing from YAML
  [Tags]              testrailid=4983449  RAT     VENDOR
  [Documentation]  Fully automated vendor testing using actual config.yaml file with Config API integration
  ...              Loads vendor configuration from deploy/rec-vendor-api/secrets/config.yaml
  ...              Integrates with Config API to filter only active inl_corp vendors for testing
  ...              Tests all vendors in config with complete automation:
  ...              - Dynamic endpoint generation: /r/{vendor_name}
  ...              - Auto-selected test dimensions: 300x300, 1200x627, 1200x600
  ...              - Vendor-specific parameter handling:
  ...              * Standard vendors: user_id, click_id, w, h
  ...              * Linkmine vendor: adds web_host, bundle_id, adtype
  ...              * INL vendors: URL-encoded subparam with base64 encoding
  ...              * INL_corp_5: Special handling with subParam=pier
  ...              - Base64 encoding validation for click_id parameters
  ...              - Response structure validation (array of products)
  ...              - Product patch verification with tracking parameters
  ...              - Config API integration: Only tests active inl_corp_X vendors (backward compatible)
  ...              - Product patch verification with tracking parameters

  # Load YAML configuration from actual config file
  # Note: size={width}x{height} parameters are auto-selected from predefined sizes (300x300, 1200x627, 1200x600)
  ${config_path} =    Set Variable        ${CURDIR}/../../../../deploy/rec-vendor-api/secrets/config.yaml
  ${yaml_content} =   Load vendor config from file  ${config_path}

  # YAML maintains complete vendor configurations
  # Config API dynamically determines which inl_corp_X to test
  # Backward compatible, non-inl vendors unaffected
  # Only tests actually active inl vendors
  ${safe_vendor_config} =  Validate and generate safe vendor yaml configuration  ${yaml_content}

  # Run complete automated testing for all vendors in YAML (with validated safe configuration)
  Test vendors from yaml configuration  ${safe_vendor_config}


# Basic vendor API healthz test
[C4977913] [RAT] [VENDOR] [HEALTHZ] Test vendor API healthz endpoint
  [Tags]  testrailid=4977913  RAT     VENDOR  HEALTHZ
  [Documentation]  Test the vendor API healthz endpoint to ensure it returns proper status and message

  Given I have an vendor session
  When I would like to set the session under vendor endpoint with  endpoint=/healthz
  Then I would like to check status_code should be "200" within the current session
