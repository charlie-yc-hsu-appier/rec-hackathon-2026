*** Keywords ***
# Automated YAML-driven vendor testing keywords #
Test vendors from yaml configuration
  [Arguments]    ${yaml_content}
  [Documentation]    Automated test for all vendors defined in YAML configuration
  ...               Tests the /r endpoint with dynamic parameters extracted from YAML
  ...               URL format: /r/{vendor_name}?user_id={uuid}&click_id={value}&w={width}&h={height}
  ...               Tests exactly the number of vendors defined in the YAML
  
  # Parse YAML once at the beginning
  ${yaml_data} =    Evaluate    yaml.safe_load('''${yaml_content}''')    yaml
  ${vendors} =    Get From Dictionary    ${yaml_data}    vendors
  
  ${vendor_count} =    Get Length    ${vendors}
  Log    Found ${vendor_count} vendor(s) in YAML configuration
  
  # Test each vendor defined in YAML
  FOR    ${vendor_config}    IN    @{vendors}
    ${vendor_name} =    Get From Dictionary    ${vendor_config}    name
    Log    Testing vendor: ${vendor_name}
    
    # Extract parameters from vendor configuration
    ${request_url} =    Get From Dictionary    ${vendor_config}    request_url
    ${tracking_url} =    Get From Dictionary    ${vendor_config}    tracking_url
    
    # Parse dimensions from request_url if available
    ${dimensions} =    Extract dimensions from request url    ${request_url}
    ${width} =    Get From Dictionary    ${dimensions}    width
    ${height} =    Get From Dictionary    ${dimensions}    height
    
    # Parse tracking URL to get parameter info
    ${tracking_config} =    Parse yaml tracking url template    ${tracking_url}
    ${param_name} =    Get From Dictionary    ${tracking_config}    param_name
    ${uses_base64} =    Get From Dictionary    ${tracking_config}    uses_base64
    
    # Generate test data
    ${user_id} =    Generate UUID4
    ${click_id} =    Set Variable    12aa
    ${click_id_base64} =    Run Keyword If    ${uses_base64}    Encode Base64    ${click_id}
    ...                     ELSE    Set Variable    ${click_id}
    
    # Test the vendor endpoint
    Given I have an vendor session
    When I would like to set the session under vendor endpoint with    endpoint=r/${vendor_name}    user_id=${user_id}    click_id=${click_id}    w=${width}    h=${height}
    
    # Verify response
    Then I would like to check status_code should be "200" within the current session
    
    # Validate response structure and content
    Validate vendor response structure    ${resp_json}
    Validate product patch contains product ids    ${resp_json}    ${param_name}    ${click_id_base64}
    
    Log    ✅ Vendor ${vendor_name} test PASSED
  END
  
  Log    ✅ Completed testing ${vendor_count} vendor(s) from YAML configuration


# BDD Template For Checking customized vendor Group #
Check the keeta vendor group
  [Arguments]         ${sid}        ${df}          ${oid}          ${expected_url_params}  ${expected_img_domain}
  [Documentation]     We should use the df="android--com.sankuai.sailor.afooddelivery_2" to verify the Keeta-api group

  Given I have a config_api session
  ${oid_list} =  I would like to get the not Finished & keeta_api enabled oid list when datafeed_id=android--com.sankuai.sailor.afooddelivery_2
  Should Not Be Empty   ${oid_list}  @{TEST TAGS} FAIL: Can't get any enabled keeta_api oid list via the ConfigAPI response: ${oid_list}

  # Since Keeta will limit specific range for lat/lon,we need to use the lat=22.2800 lon=114.1600
  Given I have an ecrec session and domain is "dync-stg.c.appier.net"
  When I would like to set the session under user2item endpoint with  sid=${sid}  df=${df}  oid=${oid_list[0]}
  ...            lat=22.2800    lon=114.1600   _debug_creative=false   idfa=55660000-0000-4C18-AAAA-556624AF0001  cid=RFTEST-Keeta  bidobjid=RFTEST-Keeta
  ...            num_items=14   no_cache=true

  Then I would like to check status_code should be "200" within the current session
  Then In the user2item response payload, the "url" should exist (EC-REC Common Value)
  Then Check the regex pattern in the given list  ${extracted_url}  (${expected_url_params})
  Then In the user2item response payload, the "img" should exist (EC-REC Common Value)
  Then Check the regex pattern in the given list  ${extracted_img_url}  (${expected_img_domain})
  Then In the user2item response payload, the "custom_label" should not exist

