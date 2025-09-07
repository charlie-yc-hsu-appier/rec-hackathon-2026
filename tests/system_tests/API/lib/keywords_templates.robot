*** Keywords ***
# Automated YAML-driven vendor testing keywords #
Test vendors from yaml configuration
  [Arguments]             ${yaml_content}
  [Documentation]  Automated test for all vendors defined in YAML configuration
  ...              Tests the /r endpoint with auto-selected test dimensions and vendor-specific parameters
  ...              Standard URL format: /r/{vendor_name}?user_id={uuid}&click_id={value}&w={width}&h={height}
  ...              Linkmine URL format: adds &web_host={domain}&bundle_id={app_id}&adtype={2|3}
  ...              Auto-selects dimensions from: 300x300, 1200x627, 1200x600
  ...              Tests exactly the number of vendors defined in the YAML
  ...              Now includes subid parameter from Config API for each vendor

  # Parse YAML once at the beginning
  ${yaml_data} =          Evaluate                yaml.safe_load('''${yaml_content}''')  yaml
  ${vendors} =            Get From Dictionary     ${yaml_data}        vendors

  ${vendor_count} =       Get Length              ${vendors}
  Log                     Found ${vendor_count} vendor(s) in YAML configuration

  # Get subid mapping for all vendors from Config API
  ${vendor_subid_mapping} =  Get vendor subids from config api  ${yaml_content}

  # Test each vendor defined in YAML
  FOR  ${vendor_config}  IN  @{vendors}
    ${vendor_name} =        Get From Dictionary     ${vendor_config}    name
    Log                     Testing vendor: ${vendor_name}

    # Get subid for this vendor
    ${vendor_subid} =       Get From Dictionary     ${vendor_subid_mapping}  ${vendor_name}  default=${EMPTY}

    # Extract parameters from vendor configuration
    ${request_url} =        Get From Dictionary     ${vendor_config}    request_url
    ${tracking_url} =       Get From Dictionary     ${vendor_config}    tracking_url

    # Auto-select test dimensions and vendor-specific parameters
    # Sizes: 300x300, 1200x627, 1200x600
    # For linkmine: also generates web_host, bundle_id, adtype
    ${dimensions} =         Auto select test dimensions  ${request_url}  ${vendor_name}
    ${width} =              Get From Dictionary     ${dimensions}       width
    ${height} =             Get From Dictionary     ${dimensions}       height

    # Extract linkmine-specific parameters if available
    ${has_web_host} =       Run Keyword And Return Status
    ...                     Dictionary Should Contain Key  ${dimensions}  web_host
    IF  ${has_web_host}
      ${web_host} =   Get From Dictionary     ${dimensions}   web_host
      ${bundle_id} =  Get From Dictionary     ${dimensions}   bundle_id
      ${adtype} =     Get From Dictionary     ${dimensions}   adtype
    END

    # Parse tracking URL to get parameter info
    ${tracking_config} =    Parse yaml tracking url template  ${tracking_url}  ${vendor_name}
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

    # Check if this is linkmine vendor to include additional parameters
    ${is_linkmine} =        Run Keyword And Return Status
    ...                     Should Be Equal         ${vendor_name}      linkmine

    # Prepare common parameters
    ${common_params} =      Create Dictionary
    ...                     endpoint=r/${vendor_name}
    ...                     user_id=${user_id}
    ...                     click_id=${click_id}
    ...                     w=${width}
    ...                     h=${height}

    # Add subid if available
    IF  '${vendor_subid}' != '${EMPTY}'
      Set To Dictionary     ${common_params}        subid=${vendor_subid}
      Log                   Using subid for ${vendor_name}: ${vendor_subid}
    ELSE
      Log                   No subid available for ${vendor_name}
    END

    # Add linkmine-specific parameters if needed
    IF  ${is_linkmine}
      Set To Dictionary     ${common_params}
      ...                   web_host=${web_host}
      ...                   bundle_id=${bundle_id}
      ...                   adtype=${adtype}
    END

    # Make the API call with all parameters
    When I would like to set the session under vendor endpoint with  &{common_params}

    # Verify response
    Then I would like to check status_code should be "200" within the current session

    # Validate response structure and content
    Validate vendor response structure  ${resp_json}
    Validate product patch contains product ids  ${resp_json}  ${param_name}  ${click_id_base64}  ${vendor_name}  ${vendor_subid}

    Log                     ✅ Vendor ${vendor_name} test PASSED
  END

  Log                     ✅ Completed testing ${vendor_count} vendor(s) from YAML configuration
