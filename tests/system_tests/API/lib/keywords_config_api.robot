*** Keywords ***
# Config API #
I have a config_api session
  [Tags]          robot:flatten

  # Basic Auth for Config API
  &{HEADERS} =    Create Dictionary       Content-Type=application/json  Authorization=Basic cmVjLXFhOmlsb3JmZ3Bva3J3aGl1bWtwb2FzZGZva2hwcWV3eWhi
  Create Session  ConfigAPISession        url=https://${CONFIG_API_HOST}  headers=&{HEADERS}  disable_warnings=1  retry_status_list=[500,502,503,504]  timeout=5


Validate config api response
  [Documentation]  Common validation for Config API responses
  [Arguments]      ${resp}  ${error_context}
  
  ${resp.status_code} =   Convert To String       ${resp.status_code}
  Should Not Match Regexp  ${resp.status_code}  (^(4|5)..$)  ${error_context}: ${resp.status_code} - ${resp.content}
  Should Not Be Empty     ${resp.json()}          ${error_context} returned empty response


I would like to get campaign_ids with group vendor_inl_corp
  [Documentation]  Get campaign IDs that have group "vendor_inl_corp" from experiments endpoint
  [Arguments]             ${site_id}              ${service_type_id}

  # Get experiments data
  ${resp} =               Get On Session          ConfigAPISession        url=/v0/recommend/experiments?site_id=${site_id}&service_type_id=${service_type_id}
  Validate config api response  ${resp}  Experiments API error for ${site_id}

  # Extract campaign_ids where any bucket has group "vendor_inl_corp"
  @{vendor_inl_corp_campaigns} =  Get Value From Json  ${resp.json()}  $.experiments[*].distributions[?(@.buckets[*].group == 'vendor_inl_corp')].campaign_id

  # Filter out empty campaign_ids
  @{filtered_campaigns} =  Create List
  FOR  ${campaign_id}  IN  @{vendor_inl_corp_campaigns}
    IF  '${campaign_id}' != '' and '${campaign_id}' != '${EMPTY}'
      Append To List  ${filtered_campaigns}   ${campaign_id}
    END
  END

  ${campaign_count} =     Get Length              ${filtered_campaigns}
  Log                     Found ${campaign_count} vendor_inl_corp campaigns: ${filtered_campaigns}
  RETURN                  ${filtered_campaigns}





I would like to check campaign status
  [Documentation]  Check campaign status and extract active inl_rec_api_group_ratio indices
  [Arguments]             ${campaign_id}

  # Get campaign details
  ${resp} =               Get On Session          ConfigAPISession        url=/v0/campaigns/${campaign_id}
  Validate config api response  ${resp}  Campaign ${campaign_id} API error

  # Extract campaign status
  ${campaign_status} =    Get Value From Json     ${resp.json()}          $.campaign.status_code

  # Extract active inl_rec_api_group_ratio indices
  @{active_indices} =     Create List
  ${has_inl_ratio} =      Run Keyword And Return Status
  ...                     Should Have Value In Json  ${resp.json()}  $.campaign.configs.inl_rec_api_group_ratio

  IF  ${has_inl_ratio}
    ${inl_ratios} =     Get Value From Json     ${resp.json()}      $.campaign.configs.inl_rec_api_group_ratio
    ${ratio_list} =     Set Variable            ${inl_ratios}[0]

    FOR  ${index}  ${ratio}  IN ENUMERATE  @{ratio_list}
      ${ratio_float} =    Convert To Number   ${ratio}
      IF  ${ratio_float} != 0.0
        Append To List  ${active_indices}   ${index}
      END
    END
  END

  Log                     Campaign ${campaign_id}: status=${campaign_status}[0], active_indices=${active_indices}

  &{result} =             Create Dictionary       status=${campaign_status}[0]  active_indices=${active_indices}
  RETURN                  ${result}


Get active vendor_inl_corp campaign ids from config api
  [Documentation]  Get vendor_inl_corp campaign IDs that are not in "Finished" status with their active inl_rec_api_group_ratio indices
  [Arguments]             ${site_id}=android--com.coupang.mobile_s2s_v3  ${service_type_id}=crossx_recommend

  I have a config_api session

  # Get all vendor_inl_corp campaign IDs from experiments
  ${all_campaign_ids} =   I would like to get campaign_ids with group vendor_inl_corp  ${site_id}  ${service_type_id}

  # Check each campaign's status and filter out "Finished" ones
  &{active_campaigns_info} =  Create Dictionary
  @{active_campaign_ids} =  Create List

  FOR  ${campaign_id}  IN  @{all_campaign_ids}
    IF  '${campaign_id}' != '' and '${campaign_id}' != '${EMPTY}'
      ${campaign_info} =      I would like to check campaign status  ${campaign_id}
      ${campaign_status} =    Set Variable            ${campaign_info}[status]
      ${active_indices} =     Set Variable            ${campaign_info}[active_indices]

      IF  '${campaign_status}' != 'Finished'
        Append To List      ${active_campaign_ids}  ${campaign_id}
        Set To Dictionary   ${active_campaigns_info}  ${campaign_id}=${campaign_info}
        Log                 Campaign ${campaign_id} is active (status: ${campaign_status}, active indices: ${active_indices})
      ELSE
        Log     Campaign ${campaign_id} is finished, skipping
      END
    END
  END

  ${campaign_count} =     Get Length              ${active_campaign_ids}
  Log                     Found ${campaign_count} active vendor_inl_corp campaigns: ${active_campaign_ids}

  Set Test Variable       ${extracted_active_vendor_inl_corp_campaigns}  ${active_campaign_ids}
  Set Test Variable       ${extracted_vendor_inl_corp_campaigns_info}  ${active_campaigns_info}

  &{result} =             Create Dictionary       campaign_ids=${active_campaign_ids}  campaigns_info=${active_campaigns_info}
  RETURN                  ${result}


Get vendor inl_corp indices for yaml configuration
  [Documentation]  Get all active vendor inl_corp indices from config api for yaml configuration testing
  [Arguments]             ${site_id}=android--com.coupang.mobile_s2s_v3  ${service_type_id}=crossx_recommend  ${validate_yaml_safety}=${True}

  # Get all active vendor_inl_corp campaign information
  ${campaign_result} =    Get active vendor_inl_corp campaign ids from config api  ${site_id}  ${service_type_id}
  ${campaigns_info} =     Set Variable            ${campaign_result}[campaigns_info]

  # Collect all active indices from all campaigns
  @{all_active_indices} =  Create List
  FOR  ${campaign_id}  IN  @{campaign_result}[campaign_ids]
    ${campaign_info} =      Get From Dictionary     ${campaigns_info}   ${campaign_id}
    ${active_indices} =     Set Variable            ${campaign_info}[active_indices]

    FOR  ${index}  IN  @{active_indices}
      Append To List  ${all_active_indices}   ${index}
    END
  END

  # Remove duplicates and sort the indices
  ${unique_indices} =     Remove Duplicates       ${all_active_indices}
  ${sorted_indices} =     Evaluate                sorted([int(x) for x in $unique_indices])

  # Create vendor mapping
  @{vendor_names} =       Create List
  @{yaml_vendor_entries} =  Create List
  FOR  ${index}  IN  @{sorted_indices}
    ${vendor_name} =    Set Variable            inl_corp_${index}
    ${yaml_entry} =     Set Variable            ${SPACE*2}${vendor_name}:
    Append To List      ${vendor_names}         ${vendor_name}
    Append To List      ${yaml_vendor_entries}  ${yaml_entry}
  END

  # Generate YAML configuration snippet
  ${yaml_snippet} =       Catenate                SEPARATOR=\n
  ...                     vendors:
  ...                       # Auto-generated inl_corp vendors from Config API
  FOR  ${entry}  IN  @{yaml_vendor_entries}
    ${yaml_snippet} =   Catenate    SEPARATOR=\n    ${yaml_snippet}     ${entry}
  END
  ${yaml_snippet} =       Catenate                SEPARATOR=\n            ${yaml_snippet}
  ...                       # End of auto-generated inl_corp vendors

  # Validation warnings
  @{safety_warnings} =    Create List
  ${max_index} =          Evaluate                max($sorted_indices) if $sorted_indices else -1
  ${vendor_count} =       Get Length              ${sorted_indices}

  Log                     Active vendors: ${vendor_count} (${vendor_names}), indices: ${sorted_indices}

  # Generate formatted YAML Config Report for test output
  Set Test Message        \n=== YAML Config Report ===  append=yes
  Set Test Message        Active vendors: ${vendor_count} (${vendor_names})  append=yes
  Set Test Message        Indices: ${sorted_indices}  append=yes
  Set Test Message        === End Report ===\n    append=yes

  Set Test Variable       ${extracted_vendor_indices_for_yaml}  ${sorted_indices}
  Set Test Variable       ${extracted_vendor_names_for_yaml}  ${vendor_names}
  Set Test Variable       ${extracted_yaml_snippet_for_vendors}  ${yaml_snippet}

  &{result} =             Create Dictionary
  ...                     indices=${sorted_indices}
  ...                     vendor_names=${vendor_names}
  ...                     yaml_snippet=${yaml_snippet}
  ...                     max_index=${max_index}
  ...                     vendor_count=${vendor_count}
  RETURN                  ${result}


Get vendor subids from config api
  [Documentation]  Get subid mapping for all vendors from Config API campaigns
  ...              Returns a dictionary with vendor_name as key and subid as value
  [Arguments]             ${yaml_content}         ${site_id}=android--com.coupang.mobile_s2s_v3  ${service_type_id}=crossx_recommend

  I have a config_api session

  # Parse YAML to get vendor names
  ${yaml_data} =          Evaluate                yaml.safe_load('''${yaml_content}''')  yaml
  ${vendors} =            Get From Dictionary     ${yaml_data}        vendors
  
  # Get all campaign IDs
  ${resp} =               Get On Session          ConfigAPISession        url=/v0/recommend/experiments?site_id=${site_id}&service_type_id=${service_type_id}
  Validate config api response  ${resp}  Experiments API error
  
  # Extract and deduplicate campaign IDs
  @{all_campaign_ids} =   Get Value From Json     ${resp.json()}          $.experiments[*].distributions[*].campaign_id
  ${unique_campaign_ids} =  Remove Duplicates      ${all_campaign_ids}
  
  ${campaign_count} =     Get Length              ${unique_campaign_ids}
  ${vendor_count} =       Get Length              ${vendors}
  Log                     Searching subids across ${campaign_count} campaigns for ${vendor_count} vendors
  
  # Create vendor subid mapping
  &{vendor_subid_mapping} =  Create Dictionary
  
  # For each vendor in YAML, try to find its subid in campaigns
  FOR  ${vendor_config}  IN  @{vendors}
    ${vendor_name} =      Get From Dictionary     ${vendor_config}    name
    ${vendor_subid} =     Find vendor subid in campaigns  ${unique_campaign_ids}  ${vendor_name}
    
    Set To Dictionary     ${vendor_subid_mapping}  ${vendor_name}=${vendor_subid}
    
    IF  '${vendor_subid}' == '${EMPTY}'
      Set Test Message    ⚠️ No subid found for vendor: ${vendor_name}  append=yes
    ELSE
      Set Test Message    ✅ Found subid for ${vendor_name}: ${vendor_subid}  append=yes
    END
  END
  
  Log                     Vendor subid mapping: ${vendor_subid_mapping}
  RETURN                  ${vendor_subid_mapping}


Find vendor subid in campaigns
  [Documentation]  Search for vendor subid across all campaigns
  [Arguments]      ${campaign_ids}  ${vendor_name}
  
  FOR  ${campaign_id}  IN  @{campaign_ids}
    IF  '${campaign_id}' != '' and '${campaign_id}' != '${EMPTY}'
      ${subid_result} =   Get vendor subid from campaign  ${campaign_id}  ${vendor_name}
      IF  '${subid_result}' != '${EMPTY}'
        Log               Found subid for ${vendor_name}: ${subid_result} (from campaign ${campaign_id})
        RETURN            ${subid_result}
      END
    END
  END
  
  RETURN                  ${EMPTY}


Get vendor subid from campaign
  [Documentation]  Get subid for a specific vendor from campaign configs
  [Arguments]             ${campaign_id}          ${vendor_name}

  # Get campaign details
  ${resp} =               Get On Session          ConfigAPISession        url=/v0/campaigns/${campaign_id}
  Validate config api response  ${resp}  Campaign ${campaign_id} API error

  # Check if campaign has subids configuration
  ${has_subids} =         Run Keyword And Return Status
  ...                     Should Have Value In Json  ${resp.json()}  $.campaign.configs.subids

  ${subid} =              Set Variable            ${EMPTY}
  IF  ${has_subids}
    ${subids_obj} =       Get Value From Json     ${resp.json()}      $.campaign.configs.subids
    ${subids_dict} =      Set Variable            ${subids_obj}[0]
    
    # Try exact match first
    ${has_exact_match} =  Run Keyword And Return Status
    ...                   Dictionary Should Contain Key  ${subids_dict}  ${vendor_name}
    
    IF  ${has_exact_match}
      ${vendor_subid_array} =  Get From Dictionary  ${subids_dict}      ${vendor_name}
      ${subid} =          Extract subid from array  ${vendor_subid_array}
      Log                 Found exact match for ${vendor_name}: ${subid}
    ELSE
      # Try partial match (e.g., inl_corp_1 matches inl_corp_1_1200x600)
      ${subids_keys} =    Get Dictionary Keys      ${subids_dict}
      @{matching_keys} =  Create List
      
      # Collect all keys that start with vendor_name
      FOR  ${key}  IN  @{subids_keys}
        ${is_partial_match} =  Run Keyword And Return Status
        ...                    Should Start With   ${key}              ${vendor_name}
        IF  ${is_partial_match}
          Append To List  ${matching_keys}  ${key}
        END
      END
      
      # If we found matching keys, randomly select one
      ${matching_count} =  Get Length  ${matching_keys}
      IF  ${matching_count} > 0
        ${random_index} =  Evaluate  random.randint(0, ${matching_count}-1)  random
        ${selected_key} =  Get From List  ${matching_keys}  ${random_index}
        ${vendor_subid_array} =  Get From Dictionary  ${subids_dict}  ${selected_key}
        ${subid} =        Extract subid from array  ${vendor_subid_array}
        Log               Found partial match for ${vendor_name} (randomly selected key: ${selected_key} from ${matching_keys}): ${subid}
      END
    END
  END

  RETURN                  ${subid}


Extract subid from array
  [Documentation]  Helper to extract subID from vendor subid array
  [Arguments]      ${subid_array}
  
  IF  ${subid_array}
    ${subid} =    Get Value From Json     ${subid_array}  $[0].subID
    RETURN        ${subid}[0]
  ELSE
    RETURN        ${EMPTY}
  END


Validate and generate safe vendor yaml configuration
  [Documentation]  Backward compatible YAML configuration validation - filters only active inl_corp vendors
  [Arguments]             ${original_yaml_content}

  # Parse original YAML
  ${yaml_data} =          Evaluate                yaml.safe_load('''${original_yaml_content}''')  yaml
  ${original_vendors} =   Get From Dictionary     ${yaml_data}            vendors

  # Check if any vendor contains "inl"
  ${has_inl} =            Set Variable            ${False}
  FOR  ${vendor}  IN  @{original_vendors}
    ${vendor_name} =    Get From Dictionary     ${vendor}       name
    ${contains_inl} =   Run Keyword And Return Status
    ...                 Should Contain          ${vendor_name}  inl
    IF  ${contains_inl}
      ${has_inl} =    Set Variable    ${True}
      BREAK
    END
  END

  # If no inl vendors found, return original YAML (backward compatibility)
  IF  not ${has_inl}
    Log     No inl vendors found, using original configuration
    RETURN  ${original_yaml_content}
  END

  # Get active inl_corp indices from Config API
  Log                     Found inl vendors, getting active indices from Config API...
  ${config_result} =      Get vendor inl_corp indices for yaml configuration
  ${active_indices} =     Set Variable            ${config_result}[indices]

  # Filter vendors: keep all non-inl vendors + only active inl vendors
  @{filtered_vendors} =   Create List

  FOR  ${vendor}  IN  @{original_vendors}
    ${vendor_name} =    Get From Dictionary     ${vendor}           name
    ${contains_inl} =   Run Keyword And Return Status
    ...                 Should Contain          ${vendor_name}      inl

    IF  not ${contains_inl}
      # Keep all non-inl vendors
      Append To List      ${filtered_vendors}     ${vendor}
      Log                 Kept: ${vendor_name}
      Set Test Message    ✅ Kept: ${vendor_name}  append=yes
    ELSE
      # For inl vendors, check if index is active
      ${vendor_index} =   Get Regexp Matches      ${vendor_name}      inl_corp_(\\d+)     1
      IF  ${vendor_index}
        ${index} =          Convert To Integer      ${vendor_index}[0]
        ${is_active} =      Run Keyword And Return Status
        ...                 List Should Contain Value  ${active_indices}  ${index}
        IF  ${is_active}
          Append To List      ${filtered_vendors}     ${vendor}
          Log                 Kept active: ${vendor_name} (index ${index})
          Set Test Message    ✅ Kept active: ${vendor_name} (index ${index})  append=yes
        ELSE
          Log                 Filtered out inactive: ${vendor_name} (index ${index})
          Set Test Message    ❌ Filtered out inactive: ${vendor_name} (index ${index})  append=yes
        END
      ELSE
        Log                 Unknown inl vendor format: ${vendor_name}
        Set Test Message    ⚠️  Unknown inl vendor format: ${vendor_name}  append=yes
      END
    END
  END

  # Rebuild YAML with filtered vendors
  Set To Dictionary       ${yaml_data}            vendors=${filtered_vendors}
  ${filtered_yaml} =      Evaluate                yaml.dump($yaml_data, default_flow_style=False)  yaml

  ${total_count} =        Get Length              ${filtered_vendors}
  Log                     Final configuration: ${total_count} vendors (active indices: ${active_indices})
  Set Test Message        📊 Final configuration: ${total_count} vendors (active indices: ${active_indices})  append=yes

  RETURN                  ${filtered_yaml}
