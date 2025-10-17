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
  ...              
  ...              *Configuration Management:*
  ...              - Loads vendor configuration from deploy/rec-vendor-api/secrets/config.yaml
  ...              - Integrates with Config API to filter only active inl_corp vendors for testing  
  ...              - Backward compatible: non-inl vendors unaffected by Config API filtering
  ...              
  ...              *Vendor Testing Coverage:*
  ...              - Tests all vendors in config with complete automation including Keeta
  ...              - Dynamic endpoint generation: /r/{vendor_name}
  ...              - Auto-selected test dimensions: 300x300, 1200x627, 1200x600
  ...              
  ...              *Vendor-specific Parameter Handling:*
  ...              - Standard vendors: user_id, click_id, w, h, subid (from Config API)
  ...              - Linkmine vendor: adds bundle_id (empty string), adtype
  ...              - INL vendors: URL-encoded subparam with base64 encoding
  ...              - INL_corp_5: Special handling with subParam=pier
  ...              - Keeta vendor: Dynamic Config API integration with lat=22.3264, lon=114.1661, k_campaign_id
  ...              - Adforus vendor: OS-specific adid handling (lowercase Android, uppercase iOS)
  ...              
  ...              *Keeta Integration Features:*
  ...              - Searches running campaigns with JSONPath filtering
  ...              - Campaign criteria: status_code=Running, datafeed_id=android--com.sankuai.sailor.afooddelivery_2
  ...              - Uses first campaign with non-empty keeta_campaign_name
  ...              - Skips image validation and click_id tracking validation for Keeta responses
  ...              
  ...              *Adforus Integration Features:*
  ...              - OS-specific adid case handling: Android (lowercase), iOS (uppercase)
  ...              - Comprehensive dual-OS testing: automatically tests both Android and iOS
  ...              - Dedicated testing workflow ensuring complete OS coverage
  ...              - No subid requirement (similar to Keeta vendor)
  ...              
  ...              *Validation & Quality Assurance:*
  ...              - Base64 encoding validation for click_id parameters
  ...              - Response structure validation (array of products)
  ...              - Product patch verification with tracking parameters
  ...              - Comprehensive error handling and logging

  # Load YAML configuration from actual config file
  # Note: size={width}x{height} parameters are auto-selected from predefined sizes (300x300, 1200x627, 1200x600)
  ${config_path} =    Set Variable        ${CURDIR}/../../../../deploy/rec-vendor-api/secrets/config.yaml
  ${yaml_content} =   Load vendor config from file  ${config_path}

  # YAML maintains complete vendor configurations with structured format
  # Config API dynamically determines which inl_corp_X to test
  # Supports structured request/tracking with queries arrays
  # Only tests actually active inl vendors
  ${safe_vendor_config} =  Validate and generate safe vendor yaml configuration  ${yaml_content}

  # Run complete automated testing for all vendors in YAML (with validated safe configuration)
  Test vendors from yaml configuration  ${safe_vendor_config}


# Basic vendor API healthz test
[C4977913] [RAT] [VENDOR] [HEALTHZ] Test vendor API healthz endpoint
  [Tags]  testrailid=4977913  RAT     VENDOR  HEALTHZ
  [Documentation]  *Test the vendor API healthz endpoint to ensure it returns proper status and message*

  Given I have an vendor session
  When I would like to set the session under vendor endpoint with  endpoint=/healthz
  Then I would like to check status_code should be "200" within the current session
