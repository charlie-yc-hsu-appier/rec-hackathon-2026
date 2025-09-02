*** Keywords ***
# Config API #
I have a config_api session
  [Tags]          robot:flatten

  # Basic Auth for Config API
  IF  '${SERVER_ENV}' == 'stag' or '${SERVER_ENV}' == 'dev'
    &{HEADERS} =    Create Dictionary       Content-Type=application/json  Authorization=Basic cmVjLXFhOmlsb3JmZ3Bva3J3aGl1bWtwb2FzZGZva2hwcWV3eWhi
  ELSE IF  '${SERVER_ENV}' == 'prod'
    &{HEADERS} =    Create Dictionary       Content-Type=application/json  Authorization=Basic cmVjLXFhOkZ6VUVyVGRSTTI2d2RoMlIyMk5lbktmeXlya2t2VGhR
  ELSE
    Fail    No SERVER_ENV params
  END

  Create Session  ConfigAPISession        url=https://${CONFIG_API_HOST}  headers=&{HEADERS}  disable_warnings=1  retry_status_list=[500,502,503,504]  timeout=5


I would like to get campaign_ids with group vendor_inl_corp
  [Documentation]    Get campaign IDs that have group "vendor_inl_corp" from experiments endpoint
  [Arguments]        ${site_id}    ${service_type_id}
  
  # Start to post the request
  ${resp} =               Get On Session          ConfigAPISession        url=/v0/recommend/experiments?site_id=${site_id}&service_type_id=${service_type_id}

  # Extra the request/response-header to Testrial & console
  Set Test Message        \n                      append=yes
  Set Test Message        *** Config API Request URL(Get experiments for ${site_id}): ${resp.url}  append=yes

  # Set the assertion variable for the following assertion variables: status_code
  Set Test Variable       ${status_code}          ${resp.status_code}
  ${resp.status_code} =   Convert To String       ${resp.status_code}
  Should Not Match Regexp  ${resp.status_code}  (^(4|5)..$)  [Get experiments] Something wrong with Config API \n reqest_url: ${resp.url} \n http_status_code: ${resp.status_code} \n message: ${resp.content} \n
  ...                     values=False

  Should Not Be Empty     ${resp.json()}          [${site_id}] @{TEST TAGS} FAIL: Can't get the response (experiments) via the ConfigAPI response: ${resp.json()}, Request: ${resp.url}. please check the 'experiments' setting
  Set Test Variable       ${request.url}          ${resp.request.url}
  Set Test Variable       ${request.headers}      ${resp.request.headers}

  # Extract campaign_ids where any bucket has group "vendor_inl_corp"
  @{vendor_inl_corp_campaigns} =  Get Value From Json     ${resp.json()}    $.experiments[*].distributions[?(@.buckets[*].group == 'vendor_inl_corp')].campaign_id
  
  # Filter out empty campaign_ids
  @{filtered_campaigns} =  Create List
  FOR  ${campaign_id}  IN  @{vendor_inl_corp_campaigns}
    IF  '${campaign_id}' != '' and '${campaign_id}' != '${EMPTY}'
      Append To List      ${filtered_campaigns}    ${campaign_id}
    END
  END

  Set Test Message        \n                      append=yes
  Set Test Message        The extracted vendor_inl_corp campaign IDs from Config API are: ${filtered_campaigns}  append=yes
  Set Test Message        \n                      append=yes

  RETURN                  ${filtered_campaigns}


I would like to check campaign status
  [Documentation]    Check the status_code of a specific campaign and extract inl_rec_api_group_ratio active indices
  [Arguments]        ${campaign_id}
  
  # Start to post the request
  ${resp} =               Get On Session          ConfigAPISession        url=/v0/campaigns/${campaign_id}

  # Extra the request/response-header to Testrial & console
  Set Test Message        \n                      append=yes
  Set Test Message        *** Config API Request URL(Get campaign status for ${campaign_id}): ${resp.url}  append=yes

  # Set the assertion variable for the following assertion variables: status_code
  Set Test Variable       ${status_code}          ${resp.status_code}
  ${resp.status_code} =   Convert To String       ${resp.status_code}
  Should Not Match Regexp  ${resp.status_code}  (^(4|5)..$)  [Check campaign status] Something wrong with Config API \n reqest_url: ${resp.url} \n http_status_code: ${resp.status_code} \n message: ${resp.content} \n
  ...                     values=False

  Should Not Be Empty     ${resp.json()}          [${campaign_id}] @{TEST TAGS} FAIL: Can't get the response (campaign) via the ConfigAPI response: ${resp.json()}, Request: ${resp.url}. please check the 'campaign' setting
  Set Test Variable       ${request.url}          ${resp.request.url}
  Set Test Variable       ${request.headers}      ${resp.request.headers}

  # Extract campaign status_code
  ${campaign_status} =    Get Value From Json     ${resp.json()}    $.campaign.status_code
  
  # Extract inl_rec_api_group_ratio and find active indices (non-zero values)
  ${has_inl_ratio} =      Run Keyword And Return Status    
  ...                     Should Have Value In Json    ${resp.json()}    $.campaign.configs.inl_rec_api_group_ratio
  
  @{active_indices} =     Create List
  IF  ${has_inl_ratio}
    ${inl_ratios} =       Get Value From Json     ${resp.json()}    $.campaign.configs.inl_rec_api_group_ratio
    ${ratio_list} =       Set Variable            ${inl_ratios}[0]
    
    # Find indices where value is not 0.0 or 0
    FOR  ${index}  ${ratio}  IN ENUMERATE  @{ratio_list}
      ${ratio_float} =    Convert To Number       ${ratio}
      IF  ${ratio_float} != 0.0
        Append To List    ${active_indices}       ${index}
        Log               Index ${index}: ${ratio} (active)
      ELSE
        Log               Index ${index}: ${ratio} (inactive)
      END
    END
  ELSE
    Log                   No inl_rec_api_group_ratio found for campaign ${campaign_id}
  END
  
  Set Test Message        \n                      append=yes
  Set Test Message        Campaign ${campaign_id} status: ${campaign_status}[0]  append=yes
  Set Test Message        Active inl_rec_api_group_ratio indices: ${active_indices}  append=yes
  Set Test Message        \n                      append=yes

  # Return both status and active indices
  &{result} =             Create Dictionary       status=${campaign_status}[0]    active_indices=${active_indices}
  RETURN                  ${result}


Get active vendor_inl_corp campaign ids from config api
  [Documentation]    Get vendor_inl_corp campaign IDs that are not in "Finished" status with their active inl_rec_api_group_ratio indices
  ...                This keyword combines two steps:
  ...                1. Extract campaign IDs with vendor_inl_corp group from experiments
  ...                2. Check each campaign's status and inl_rec_api_group_ratio, filter out "Finished" campaigns
  ...                Returns: Dictionary with campaign_id as key and dictionary containing status and active_indices as value
  [Arguments]        ${site_id}=android--com.coupang.mobile_s2s_v3    ${service_type_id}=crossx_recommend
  
  I have a config_api session
  
  # Step 1: Get all vendor_inl_corp campaign IDs from experiments
  ${all_campaign_ids} =   I would like to get campaign_ids with group vendor_inl_corp    ${site_id}    ${service_type_id}
  
  # Step 2: Check each campaign's status and inl_rec_api_group_ratio, filter out "Finished" ones
  &{active_campaigns_info} =  Create Dictionary
  @{active_campaign_ids} =    Create List
  
  FOR  ${campaign_id}  IN  @{all_campaign_ids}
    IF  '${campaign_id}' != '' and '${campaign_id}' != '${EMPTY}'
      ${campaign_info} =      I would like to check campaign status    ${campaign_id}
      ${campaign_status} =    Set Variable    ${campaign_info}[status]
      ${active_indices} =     Set Variable    ${campaign_info}[active_indices]
      
      IF  '${campaign_status}' != 'Finished'
        Append To List        ${active_campaign_ids}    ${campaign_id}
        Set To Dictionary     ${active_campaigns_info}    ${campaign_id}=${campaign_info}
        Log                   Campaign ${campaign_id} is active (status: ${campaign_status}, active indices: ${active_indices})
      ELSE
        Log                   Campaign ${campaign_id} is finished, skipping
      END
    END
  END

  ${campaign_count} =     Get Length              ${active_campaign_ids}
  Set Test Message        \n                      append=yes
  Set Test Message        Found ${campaign_count} active vendor_inl_corp campaigns: ${active_campaign_ids}  append=yes
  
  # Log detailed information about each active campaign
  FOR  ${campaign_id}  IN  @{active_campaign_ids}
    ${info} =             Get From Dictionary     ${active_campaigns_info}    ${campaign_id}
    Set Test Message      Campaign ${campaign_id}: status=${info}[status], active_indices=${info}[active_indices]  append=yes
  END
  Set Test Message        \n                      append=yes
  
  Set Test Variable       ${extracted_active_vendor_inl_corp_campaigns}      ${active_campaign_ids}
  Set Test Variable       ${extracted_vendor_inl_corp_campaigns_info}        ${active_campaigns_info}
  
  # Return both the list of campaign IDs and detailed information
  &{result} =             Create Dictionary       campaign_ids=${active_campaign_ids}    campaigns_info=${active_campaigns_info}
  RETURN                  ${result}


Get vendor inl_corp indices for yaml configuration
  [Documentation]    Get all active vendor inl_corp indices from config api for yaml configuration testing
  ...                This keyword collects all active_indices from vendor_inl_corp campaigns and merges them
  ...                to create a comprehensive list of vendor indices (inl_corp_0, inl_corp_1, etc.)
  ...                IMPORTANT: This only affects inl_corp vendor types, other vendor types are preserved
  ...                Returns: Dictionary with indices, vendor names, and validation info for safe YAML configuration
  [Arguments]        ${site_id}=android--com.coupang.mobile_s2s_v3    ${service_type_id}=crossx_recommend    ${validate_yaml_safety}=${True}
  
  # Get all active vendor_inl_corp campaign information
  ${campaign_result} =    Get active vendor_inl_corp campaign ids from config api    ${site_id}    ${service_type_id}
  ${campaigns_info} =     Set Variable    ${campaign_result}[campaigns_info]
  
  # Collect all active indices from all campaigns
  @{all_active_indices} =  Create List
  FOR  ${campaign_id}  IN  @{campaign_result}[campaign_ids]
    ${campaign_info} =    Get From Dictionary     ${campaigns_info}    ${campaign_id}
    ${active_indices} =   Set Variable            ${campaign_info}[active_indices]
    
    # Add each active index to the master list
    FOR  ${index}  IN  @{active_indices}
      Append To List      ${all_active_indices}    ${index}
      Log                 Added index ${index} from campaign ${campaign_id}
    END
  END
  
  # Remove duplicates and sort the indices
  ${unique_indices} =     Remove Duplicates       ${all_active_indices}
  ${sorted_indices} =     Evaluate                sorted([int(x) for x in $unique_indices])
  
  # Create vendor mapping for logging
  @{vendor_names} =       Create List
  @{yaml_vendor_entries} =  Create List
  FOR  ${index}  IN  @{sorted_indices}
    ${vendor_name} =      Set Variable            inl_corp_${index}
    ${yaml_entry} =       Set Variable            ${SPACE*2}${vendor_name}:
    Append To List        ${vendor_names}         ${vendor_name}
    Append To List        ${yaml_vendor_entries}  ${yaml_entry}
  END
  
  # Generate safe YAML configuration snippet
  ${yaml_snippet} =       Catenate                SEPARATOR=\n
  ...                     vendors:
  ...                     # Auto-generated inl_corp vendors from Config API
  ...                     # Only affects inl_corp_* entries, other vendors are preserved
  FOR  ${entry}  IN  @{yaml_vendor_entries}
    ${yaml_snippet} =     Catenate                SEPARATOR=\n    ${yaml_snippet}    ${entry}
  END
  ${yaml_snippet} =       Catenate                SEPARATOR=\n    ${yaml_snippet}
  ...                     # End of auto-generated inl_corp vendors
  
  # Validation warnings for YAML safety
  @{safety_warnings} =    Create List
  ${max_index} =          Evaluate                max($sorted_indices) if $sorted_indices else -1
  IF  ${max_index} > 10
    Append To List        ${safety_warnings}      High index detected (${max_index}). Consider reviewing vendor allocation.
  END
  
  ${vendor_count} =       Get Length              ${sorted_indices}
  IF  ${vendor_count} == 0
    Append To List        ${safety_warnings}      No active inl_corp vendors found. YAML configuration may be unnecessary.
  END
  
  Set Test Message        \n                      append=yes
  Set Test Message        === YAML Config Report ===  append=yes
  Set Test Message        Active vendors: ${vendor_count} (${vendor_names})  append=yes
  Set Test Message        Indices: ${sorted_indices}  append=yes
  
  IF  ${safety_warnings}
    Set Test Message      ⚠️  Warnings: ${safety_warnings}  append=yes
  END
  
  Set Test Message        === End Report ===\n  append=yes
  
  Set Test Variable       ${extracted_vendor_indices_for_yaml}    ${sorted_indices}
  Set Test Variable       ${extracted_vendor_names_for_yaml}      ${vendor_names}
  Set Test Variable       ${extracted_yaml_snippet_for_vendors}   ${yaml_snippet}
  
  # Return comprehensive result for safe YAML configuration
  &{result} =             Create Dictionary       
  ...                     indices=${sorted_indices}    
  ...                     vendor_names=${vendor_names}
  ...                     yaml_snippet=${yaml_snippet}
  ...                     safety_warnings=${safety_warnings}
  ...                     max_index=${max_index}
  ...                     vendor_count=${vendor_count}
  RETURN                  ${result}


Validate and generate safe vendor yaml configuration
  [Documentation]    Backward compatible YAML configuration validation
  ...                Non-inl vendors: Keep as-is
  ...                Inl vendors: Use Config API to filter only active inl_corp_X for testing
  [Arguments]        ${original_yaml_content}
  
  # Parse original YAML
  ${yaml_data} =          Evaluate                yaml.safe_load('''${original_yaml_content}''')  yaml
  ${original_vendors} =   Get From Dictionary     ${yaml_data}    vendors
  
  # Check if any vendor contains "inl" 
  ${has_inl} =            Set Variable            ${False}
  FOR  ${vendor}  IN  @{original_vendors}
    ${vendor_name} =      Get From Dictionary     ${vendor}    name
    ${contains_inl} =     Run Keyword And Return Status    Should Contain    ${vendor_name}    inl
    IF  ${contains_inl}
      ${has_inl} =        Set Variable            ${True}
      BREAK
    END
  END
  
  Set Test Message        Check if YAML contains inl vendors: ${has_inl}    append=yes
  
  # If no inl vendors found, return original YAML (backward compatibility)
  IF  not ${has_inl}
    Set Test Message      ✅ No inl vendors found, using original configuration    append=yes
    RETURN                ${original_yaml_content}
  END
  
  # Get active inl_corp indices from Config API
  Set Test Message        🔍 Found inl vendors, getting active indices from Config API...    append=yes
  ${config_result} =      Get vendor inl_corp indices for yaml configuration
  ${active_indices} =     Set Variable            ${config_result}[indices]
  
  # Filter vendors: keep all non-inl vendors + only active inl vendors
  @{filtered_vendors} =   Create List
  
  FOR  ${vendor}  IN  @{original_vendors}
    ${vendor_name} =      Get From Dictionary     ${vendor}    name
    ${contains_inl} =     Run Keyword And Return Status    Should Contain    ${vendor_name}    inl
    
    IF  not ${contains_inl}
      # Keep all non-inl vendors
      Append To List      ${filtered_vendors}     ${vendor}
      Set Test Message    ✅ Kept: ${vendor_name}    append=yes
    ELSE
      # For inl vendors, check if index is active
      ${vendor_index} =   Get Regexp Matches      ${vendor_name}    inl_corp_(\\d+)    1
      IF  ${vendor_index}
        ${index} =        Convert To Integer      ${vendor_index}[0]
        ${is_active} =    Run Keyword And Return Status    List Should Contain Value    ${active_indices}    ${index}
        IF  ${is_active}
          Append To List  ${filtered_vendors}     ${vendor}
          Set Test Message    ✅ Kept active: ${vendor_name} (index ${index})    append=yes
        ELSE
          Set Test Message    ❌ Filtered out inactive: ${vendor_name} (index ${index})    append=yes
        END
      ELSE
        Set Test Message    ⚠️  Unknown inl vendor format: ${vendor_name}    append=yes
      END
    END
  END
  
  # Rebuild YAML with filtered vendors
  Set To Dictionary       ${yaml_data}            vendors=${filtered_vendors}
  ${filtered_yaml} =      Evaluate                yaml.dump($yaml_data, default_flow_style=False)    yaml
  
  ${total_count} =        Get Length              ${filtered_vendors}
  Set Test Message        📊 Final configuration: ${total_count} vendors (active indices: ${active_indices})    append=yes
  
  RETURN                  ${filtered_yaml}

