*** Keywords ***
# Automated YAML-driven vendor testing keywords #
Test vendors from yaml configuration
  [Arguments]             ${yaml_content}
  [Documentation]  Automated test for all vendors defined in YAML configuration
  ...              Tests the /r endpoint with auto-selected test dimensions
  ...              URL format: /r/{vendor_name}?user_id={uuid}&click_id={value}&w={width}&h={height}
  ...              Auto-selects dimensions from: 300x300, 1200x627, 1200x600
  ...              Tests exactly the number of vendors defined in the YAML

  # Parse YAML once at the beginning
  ${yaml_data} =          Evaluate                yaml.safe_load('''${yaml_content}''')  yaml
  ${vendors} =            Get From Dictionary     ${yaml_data}        vendors

  ${vendor_count} =       Get Length              ${vendors}
  Log                     Found ${vendor_count} vendor(s) in YAML configuration

  # Test each vendor defined in YAML
  FOR  ${vendor_config}  IN  @{vendors}
    ${vendor_name} =        Get From Dictionary     ${vendor_config}    name
    Log                     Testing vendor: ${vendor_name}

    # Extract parameters from vendor configuration
    ${request_url} =        Get From Dictionary     ${vendor_config}    request_url
    ${tracking_url} =       Get From Dictionary     ${vendor_config}    tracking_url

    # Auto-select test dimensions from predefined sizes: 300x300, 1200x627, 1200x600
    ${dimensions} =         Auto select test dimensions  ${request_url}
    ${width} =              Get From Dictionary     ${dimensions}       width
    ${height} =             Get From Dictionary     ${dimensions}       height

    # Parse tracking URL to get parameter info
    ${tracking_config} =    Parse yaml tracking url template  ${tracking_url}
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
    When I would like to set the session under vendor endpoint with  endpoint=r/${vendor_name}  user_id=${user_id}  click_id=${click_id}  w=${width}  h=${height}

    # Verify response
    Then I would like to check status_code should be "200" within the current session

    # Validate response structure and content
    Validate vendor response structure  ${resp_json}
    Validate product patch contains product ids  ${resp_json}  ${param_name}  ${click_id_base64}

    Log                     ✅ Vendor ${vendor_name} test PASSED
  END

  Log                     ✅ Completed testing ${vendor_count} vendor(s) from YAML configuration
