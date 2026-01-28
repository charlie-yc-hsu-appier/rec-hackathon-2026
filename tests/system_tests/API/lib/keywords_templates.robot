*** Settings ***
Resource            ../res/init.robot

*** Keywords ***
# Automated YAML-driven vendor testing keywords #
Test vendors from yaml configuration
  [Arguments]             ${yaml_content}
  [Documentation]  *Purpose:*
  ...              Automated end-to-end testing for all vendors defined in YAML configuration.
  ...              Orchestrates complete vendor testing workflow with dynamic parameter generation, vendor-specific handling, and comprehensive validation.
  ...              
  ...              *Parameters:*
  ...              - yaml_content: YAML configuration content containing vendor definitions, request/tracking configurations
  ...              
  ...              *Standard URL Format:*
  ...              /r/{vendor_name}?user_id={uuid}&click_id={value}&w={width}&h={height}
  ...              
  ...              *Vendor-specific Parameter Handling:*
  ...              - Standard vendors: user_id, click_id, w, h, subid (from Config API)
  ...              - Linkmine vendor: adds bundle_id (empty string), adtype parameters
  ...              - Adpacker vendor: adds adtype parameter
  ...              - INL vendors: URL-encoded subparam with base64 encoding
  ...              - INL_corp_5: Special subParam=pier (uppercase P, fixed value)
  ...              - Keeta vendor: adds lat=22.3264, lon=114.1661, k_campaign_id (from Config API)
  ...              - Adforus vendor: adds os (android/ios), auto-transforms user_id case based on OS
  ...              
  ...              *Usage Example:*
  ...              ```robotframework
  ...              # Load YAML config and test all vendors
  ...              ${config_path} =  Set Variable  ${CURDIR}/../../../../deploy/rec-vendor-api/secrets/config.yaml
  ...              ${yaml_content} =  Load vendor config from file  ${config_path}
  ...              ${safe_vendor_config} =  Validate and generate safe vendor yaml configuration  ${yaml_content}
  ...              Test vendors from yaml configuration  ${safe_vendor_config}
  ...              ```
  ...              
  ...              *Implementation:*
  ...              1. Parse YAML configuration to extract vendor list
  ...              2. Get vendor subid mapping from Config API for all vendors
  ...              3. For each vendor in YAML:
  ...                 a. Extract vendor name and configuration
  ...                 b. Handle Keeta vendor: Retrieve k_campaign_id from Config API
  ...                 c. Retrieve vendor-specific subid from Config API mapping
  ...                 d. Determine OS testing strategy (Adforus: dual-OS, others: single)
  ...                 e. For each OS iteration (1 for standard vendors, 2 for Adforus):
  ...                    - Auto-select test dimensions (300x300, 1200x627, 1200x600)
  ...                    - Extract vendor-specific parameters (bundle_id, adtype, os)
  ...                    - Parse tracking configuration (param_name, base64 encoding)
  ...                    - Generate test data (user_id UUID, click_id)
  ...                    - Encode click_id to base64 if required
  ...                    - Create vendor session
  ...                    - Build request parameters based on vendor type
  ...                    - Handle Keeta: Add GPS coordinates and campaign_id
  ...                    - Handle Adforus: Transform user_id case based on OS
  ...                    - Send GET request to /r/{vendor_name} endpoint
  ...                    - Validate response structure
  ...                    - Validate tracking parameters in product URLs
  ...                    - Log test results
  ...              4. Complete testing for all vendors
  ...              
  ...              *Keeta Integration Features:*
  ...              - Dynamic Config API integration: Searches running campaigns
  ...              - Campaign criteria: datafeed_id=android--com.sankuai.sailor.afooddelivery_2, status_code=Running
  ...              - Uses JSONPath filtering ($[?(@.datafeed_id==...)]) for efficient campaign discovery
  ...              - Extracts first campaign with non-empty keeta_campaign_name
  ...              - Skips image field validation (Keeta responses may not include image)
  ...              - Skips click_id tracking validation (Keeta uses k_campaign_id instead)
  ...              - GPS coordinates: lat=22.3264 (Hong Kong), lon=114.1661
  ...              
  ...              *Adforus Integration Features:*
  ...              - OS-specific adid case handling:
  ...                * Android: Converts user_id to lowercase (adid must be lowercase)
  ...                * iOS: Converts user_id to uppercase (adid must be uppercase)
  ...              - Comprehensive dual-OS testing: Automatically tests both Android and iOS in single run
  ...              - Dedicated adforus testing workflow ensuring complete OS coverage
  ...              - No subid requirement (similar to Keeta vendor)
  ...              - Validates adid case in product URLs matches OS requirement
  ...              
  ...              *Test Dimensions:*
  ...              Auto-selects from: 300x300, 1200x627, 1200x600
  ...              Different dimensions may be selected for different vendors in same run
  ...              
  ...              *Prerequisites:*
  ...              - YAML configuration must contain valid vendor definitions
  ...              - Config API must be accessible for subid/campaign retrieval
  ...              - Vendor API session must be creatable
  ...              - For Keeta: Must have running campaigns with keeta_campaign_name
  ...              - ${HTTP_METHOD} and ${VENDOR_HOST} must be set (via Get Test Value)
  ...              
  ...              *Validation Coverage:*
  ...              - Response structure validation (array of products with required fields)
  ...              - Tracking parameter validation (base64 encoding, URL encoding for INL)
  ...              - Product data validation (product_id, url, image fields)
  ...              - Vendor-specific validation (Keeta skips, Adforus adid case)
  ...              
  ...              *Special Cases:*
  ...              - Keeta vendor: Requires valid running campaign, skips standard validations
  ...              - Adforus vendor: Tests both Android and iOS, validates adid case transformation
  ...              - INL vendors: Uses URL-encoded subparam format in land parameter
  ...              - INL_corp_5: Uses fixed subParam=pier instead of dynamic base64 value
  ...              - Linkmine vendor: Requires bundle_id (can be empty) and adtype parameters
  ...              - Vendors without subid in Config API: Uses empty string for subid parameter
  ...              
  ...              *Error Handling:*
  ...              - Fails immediately if Keeta campaign configuration not found
  ...              - Validates YAML structure before testing
  ...              - Comprehensive logging for debugging vendor-specific issues
  ...              - Emoji indicators: 🎯 for special handling, ✅ for successful validation

  # Parse YAML once at the beginning
  ${yaml_data} =          Evaluate                yaml.safe_load('''${yaml_content}''')  yaml
  ${vendors} =            Evaluate                $yaml_data['vendor_config']['vendors']

  ${vendor_count} =       Get Length              ${vendors}
  Log                     Found ${vendor_count} vendor(s) in YAML configuration

  # Get subid mapping for all vendors from Config API
  ${vendor_subid_mapping} =  Get vendor subids from config api  ${yaml_content}

  # Test each vendor defined in YAML
  FOR  ${vendor_config}  IN  @{vendors}
    ${vendor_name} =        Get From Dictionary     ${vendor_config}    name
    
    # Temporarily skip Keeta vendor due to client API issues
    ${is_keeta} =           Run Keyword And Return Status
    ...                     Should Be Equal         ${vendor_name}      keeta
    IF  ${is_keeta}
      Log                   ⚠️ Skipping Keeta vendor due to client API issues  WARN
      Set Test Message      ⚠️ Skipped Keeta vendor (temporary - client API issues)  append=yes
      CONTINUE
    END
    
    Log                     Testing vendor: ${vendor_name}

    # Handle Keeta vendor specially
    ${is_keeta} =           Run Keyword And Return Status
    ...                     Should Be Equal         ${vendor_name}      keeta
    IF  ${is_keeta}
      # Get Keeta campaign configuration from Config API
      ${keeta_campaign_name} =  Get keeta campaign configuration
      
      IF  '${keeta_campaign_name}' == '${EMPTY}'
        Fail                Failed to get keeta_campaign_name from Config API - no valid Keeta campaign found
      END
      
      Log                   Testing Keeta vendor with campaign: ${keeta_campaign_name}
    END

    # Get subid for this vendor
    ${vendor_subid} =       Get From Dictionary     ${vendor_subid_mapping}  ${vendor_name}  default=${EMPTY}

    # Check if this is adforus vendor for dual-OS testing
    ${is_adforus} =         Run Keyword And Return Status
    ...                     Should Be Equal         ${vendor_name}      adforus

    # Prepare OS list - adforus tests both android/ios, others test once
    IF  ${is_adforus}
      @{os_list} =          Create List      android     ios
      Log                   Adforus vendor: will test both android and ios
    ELSE
      @{os_list} =          Create List      ${EMPTY}
      Log                   Standard vendor: single test execution
    END

    # Test vendor with each OS (single iteration for non-adforus vendors)
    FOR  ${current_os}  IN  @{os_list}
      IF  ${is_adforus}
        Log                 Testing Adforus vendor with OS: ${current_os}
      END

      # Extract parameters from vendor configuration
      ${request_config} =     Get From Dictionary     ${vendor_config}    request
      ${request_url} =        Get From Dictionary     ${request_config}   url
      ${request_queries} =    Get From Dictionary     ${request_config}   queries
      
      ${tracking_config} =    Get From Dictionary     ${vendor_config}    tracking
      ${tracking_url} =       Get From Dictionary     ${tracking_config}  url
      ${tracking_queries} =   Get From Dictionary     ${tracking_config}  queries

      # Auto-select test dimensions and vendor-specific parameters
      ${dimensions} =         Auto select test dimensions  ${request_url}  ${request_queries}  ${vendor_name}
      ${width} =              Get From Dictionary     ${dimensions}       width
      ${height} =             Get From Dictionary     ${dimensions}       height

      # Extract vendor-specific parameters if available
      ${has_bundle_id} =      Run Keyword And Return Status
      ...                     Dictionary Should Contain Key  ${dimensions}  bundle_id
      ${has_adtype} =         Run Keyword And Return Status
      ...                     Dictionary Should Contain Key  ${dimensions}  adtype
      
      IF  ${has_bundle_id}
        ${bundle_id} =        Get From Dictionary     ${dimensions}   bundle_id
      END
      
      IF  ${has_adtype}
        ${adtype} =           Get From Dictionary     ${dimensions}   adtype
      END
      
      # For adforus, use current_os from loop; for others, check dimensions or set empty
      IF  ${is_adforus}
        ${os} =               Set Variable            ${current_os}
        ${has_os} =           Set Variable            ${TRUE}
        Log                   Using adforus OS from loop: ${os}
      ELSE
        ${has_os} =           Run Keyword And Return Status
        ...                   Dictionary Should Contain Key  ${dimensions}  os
        IF  ${has_os}
          ${os} =             Get From Dictionary     ${dimensions}   os
        ELSE
          ${os} =             Set Variable            ${EMPTY}
        END
      END

      # Parse tracking configuration (check both request and tracking queries)
      ${tracking_config} =    Parse tracking config  ${request_queries}  ${tracking_queries}  ${vendor_name}
      ${param_name} =         Get From Dictionary     ${tracking_config}  param_name
      ${uses_base64} =        Get From Dictionary     ${tracking_config}  uses_base64

      # Generate test data
      ${user_id} =            Generate UUID4
      # Click ID is combination of cid.oid
      ${click_id} =           Set Variable            12aa.12aa
      IF  ${uses_base64}
        ${click_id_base64} =    Encode Base64   ${click_id}
      ELSE
        ${click_id_base64} =    Set Variable    ${click_id}
      END

      # Test the vendor endpoint
      Given I have an vendor session

      # Prepare common parameters
      ${common_params} =      Create Dictionary
      ...                     endpoint=r/${vendor_name}
      ...                     user_id=${user_id}
      ...                     click_id=${click_id}
      ...                     w=${width}
      ...                     h=${height}

      # Add subid if available (API will validate if required)
      IF  '${vendor_subid}' != '${EMPTY}'
        Set To Dictionary     ${common_params}        subid=${vendor_subid}
        Log                   Using subid for ${vendor_name}: ${vendor_subid}
      ELSE
        Log                   No subid available for ${vendor_name} - will use empty subid
      END

      # Add vendor-specific parameters if they exist
      # Add bundle_id for linkmine vendor
      IF  ${has_bundle_id}
        Set To Dictionary     ${common_params}        bundle_id=${bundle_id}
        Log                   Added bundle_id for ${vendor_name}: ${bundle_id}
      END
      
      # Add adtype for vendors that need it (linkmine, adpacker)
      IF  ${has_adtype}
        Set To Dictionary     ${common_params}        adtype=${adtype}
        Log                   Added adtype for ${vendor_name}: ${adtype}
      END

      # Add Keeta-specific parameters if needed
      IF  ${is_keeta}
        Set To Dictionary     ${common_params}
        ...                   lat=22.3264
        ...                   lon=114.1661
        ...                   k_campaign_id=${keeta_campaign_name}
        Log                   Making Keeta API call with campaign: ${keeta_campaign_name}
      END

      # Add Adforus-specific parameters if needed
      IF  ${is_adforus}
        Set To Dictionary     ${common_params}
        ...                   os=${os}
        Log                   Making Adforus API call with os: ${os}, user_id: ${user_id}
      END

      # Make the API call with all parameters
      When I would like to set the session under vendor endpoint with  &{common_params}

      # Verify response
      Then I would like to check status_code should be "200" within the current session

      # Validate response structure and content
      Validate vendor response structure  ${resp_json}  ${vendor_name}
      Validate product patch contains product ids  ${resp_json}  ${param_name}  ${click_id_base64}  ${vendor_name}  ${os}  ${user_id}

      IF  ${is_adforus}
        Log                   ✅ Adforus vendor ${vendor_name} test PASSED with OS: ${os}
      ELSE
        Log                   ✅ Vendor ${vendor_name} test PASSED
      END
    END
  END

  Log                     ✅ Completed testing ${vendor_count} vendor(s) from YAML configuration
