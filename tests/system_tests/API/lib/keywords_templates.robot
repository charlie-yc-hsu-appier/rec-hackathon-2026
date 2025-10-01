*** Keywords ***
# Automated YAML-driven vendor testing keywords #
Test vendors from yaml configuration
  [Arguments]             ${yaml_content}
  [Documentation]  Automated test for all vendors defined in YAML configuration
  ...              Tests the /r endpoint with auto-selected test dimensions and vendor-specific parameters
  ...              Standard URL format: /r/{vendor_name}?user_id={uuid}&click_id={value}&w={width}&h={height}
  ...              
  ...              Vendor-specific parameter handling:
  ...              • Standard vendors: user_id, click_id, w, h, subid (from Config API)
  ...              • Linkmine vendor: adds bundle_id (empty string), adtype parameters
  ...              • Adpacker vendor: adds adtype parameter
  ...              • INL vendors: URL-encoded subparam with base64 encoding
  ...              • Keeta vendor: adds lat=22.3264, lon=114.1661, k_campaign_id (from Config API)
  ...              
  ...              Keeta integration features:
  ...              • Dynamic Config API integration: searches running campaigns
  ...              • Campaign criteria: datafeed_id=android--com.sankuai.sailor.afooddelivery_2
  ...              • Uses JSONPath filtering for efficient campaign discovery
  ...              • Skips image validation and click_id tracking validation for Keeta
  ...              
  ...              Auto-selects test dimensions from: 300x300, 1200x627, 1200x600
  ...              Tests exactly the number of vendors defined in the YAML
  ...              Includes subid parameter from Config API for each vendor
  ...              Comprehensive validation: response structure, tracking parameters, product data

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

    # Extract parameters from vendor configuration
    # Use structured request and tracking configuration
    ${request_config} =     Get From Dictionary     ${vendor_config}    request
    ${request_url} =        Get From Dictionary     ${request_config}   url
    
    ${tracking_config} =    Get From Dictionary     ${vendor_config}    tracking
    ${tracking_url} =       Get From Dictionary     ${tracking_config}  url
    ${tracking_queries} =   Get From Dictionary     ${tracking_config}  queries

    # Auto-select test dimensions and vendor-specific parameters
    # Sizes: 300x300, 1200x627, 1200x600
    # For linkmine and adpacker: also generates adtype parameter
    # For linkmine: additionally generates bundle_id (empty string)
    ${dimensions} =         Auto select test dimensions  ${request_url}  ${vendor_name}
    ${width} =              Get From Dictionary     ${dimensions}       width
    ${height} =             Get From Dictionary     ${dimensions}       height

    # Extract vendor-specific parameters if available
    ${has_bundle_id} =      Run Keyword And Return Status
    ...                     Dictionary Should Contain Key  ${dimensions}  bundle_id
    ${has_adtype} =         Run Keyword And Return Status
    ...                     Dictionary Should Contain Key  ${dimensions}  adtype
    
    IF  ${has_bundle_id}
      ${bundle_id} =  Get From Dictionary     ${dimensions}   bundle_id
    END
    
    IF  ${has_adtype}
      ${adtype} =     Get From Dictionary     ${dimensions}   adtype
    END

    # Parse tracking configuration
    ${tracking_config} =    Parse tracking config  ${tracking_queries}  ${vendor_name}
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

    # Add subid (required for all vendors except keeta)
    IF  ${is_keeta}
      Log                   Keeta vendor does not require subid
    ELSE IF  '${vendor_subid}' != '${EMPTY}'
      Set To Dictionary     ${common_params}        subid=${vendor_subid}
      Log                   Using subid for ${vendor_name}: ${vendor_subid}
    ELSE
      Fail                  No subid found for vendor ${vendor_name} - subid is required for all non-keeta vendors
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

    # Make the API call with all parameters
    When I would like to set the session under vendor endpoint with  &{common_params}

    # Verify response
    Then I would like to check status_code should be "200" within the current session

    # Validate response structure and content
    Validate vendor response structure  ${resp_json}  ${vendor_name}
    Validate product patch contains product ids  ${resp_json}  ${param_name}  ${click_id_base64}  ${vendor_name}

    Log                     ✅ Vendor ${vendor_name} test PASSED
  END

  Log                     ✅ Completed testing ${vendor_count} vendor(s) from YAML configuration
