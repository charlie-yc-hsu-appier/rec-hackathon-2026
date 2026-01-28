*** Settings ***
Resource            ../res/init.robot

*** Variables ***
# Default subids for vendors (used when Config API doesn't have subid configured)
# This allows testing without waiting for Config API to have campaign data
&{DEFAULT_VENDOR_SUBIDS}
...    inl_corp_1=kkobizdw1
...    inl_corp_2=kkobiz1200x600
...    inl_corp_3=dlapp4968m1200x600
...    inl_corp_4=24kkobizdw1
...    inl_corp_5=pinpuzzleAPP
...    inl_corp_6=bluediarydw2
...    linkmine=FAXXi4vdOY
...    adpacker=103947SA
...    replace=KRpartner01
...    adpopcorn=CPSapapp1

*** Keywords ***
# Config API #
I have a config_api session
  [Documentation]  Create HTTP session for Config API with Basic Authentication.
  ...
  ...              *Purpose*
  ...              Initialize HTTP session for Config API with environment-specific credentials
  ...              and automatic retry configuration.
  ...
  ...              *Usage Example*
  ...              | I have a config_api session |
  ...
  ...              *Implementation*
  ...              - Base URL: https://${CONFIG_API_HOST}
  ...              - Authentication: Basic Auth (staging credentials)
  ...              - Auto-retry: 500, 502, 503, 504 status codes
  ...              - Timeout: 5 seconds
  ...              - Warnings disabled
  ...
  ...              *Prerequisites*
  ...              - ${CONFIG_API_HOST} must be set from valueset.dat
  
  [Tags]          robot:flatten

  # Basic Auth for Config API
  &{HEADERS} =    Create Dictionary       Content-Type=application/json  Authorization=Basic cmVjLXFhOmlsb3JmZ3Bva3J3aGl1bWtwb2FzZGZva2hwcWV3eWhi
  Create Session  ConfigAPISession        url=https://${CONFIG_API_HOST}  headers=&{HEADERS}  disable_warnings=1  retry_status_list=[500,502,503,504]  timeout=5


Validate config api response
  [Documentation]  Common validation for Config API responses.
  ...
  ...              *Purpose*
  ...              Validate Config API response status and ensure response body is not empty.
  ...
  ...              *Parameters*
  ...              - ${resp}: Response object from Config API
  ...              - ${error_context}: Context string for error messages
  ...
  ...              *Usage Example*
  ...              | ${resp} = | Get On Session | ConfigAPISession | url=/v0/campaigns/123 |
  ...              | Validate config api response | ${resp} | Campaign API error |
  ...
  ...              *Implementation*
  ...              - Validates status code is not 4xx or 5xx
  ...              - Ensures response JSON is not empty
  ...              - Fails with detailed error message if validation fails
  
  [Arguments]      ${resp}  ${error_context}
  
  ${resp.status_code} =   Convert To String       ${resp.status_code}
  Should Not Match Regexp  ${resp.status_code}  (^(4|5)..$)  ${error_context}: ${resp.status_code} - ${resp.content}
  Should Not Be Empty     ${resp.json()}          ${error_context} returned empty response


I would like to get campaign_ids with group vendor_inl_corp
  [Documentation]  Get campaign IDs that have group "vendor_inl_corp" from experiments endpoint.
  ...
  ...              *Purpose*
  ...              Retrieve all campaign IDs configured with vendor_inl_corp group from
  ...              Config API experiments.
  ...
  ...              *Parameters*
  ...              - ${site_id}: Site identifier (e.g., android--com.coupang.mobile_s2s_v3)
  ...              - ${service_type_id}: Service type (e.g., crossx_recommend)
  ...
  ...              *Returns*
  ...              List of campaign IDs with vendor_inl_corp group (filtered, no empty values)
  ...
  ...              *Usage Example*
  ...              | @{campaigns} = | I would like to get campaign_ids with group vendor_inl_corp |
  ...              | ... | android--com.coupang.mobile_s2s_v3 | crossx_recommend |
  ...
  ...              *Implementation*
  ...              - Queries /v0/recommend/experiments endpoint
  ...              - Uses JSONPath to filter campaigns with vendor_inl_corp buckets
  ...              - Removes empty campaign IDs from result
  
  [Arguments]             ${site_id}              ${service_type_id}

  # Get experiments data
  ${resp} =               Get On Session          ConfigAPISession        url=/v0/recommend/experiments/custom_flow?site_id=${site_id}&service_type_id=${service_type_id}
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
  [Documentation]  Check campaign status and extract active inl_rec_api_group_ratio indices.
  ...
  ...              *Purpose*
  ...              Retrieve campaign status and identify which vendor group indices
  ...              have non-zero traffic allocation.
  ...
  ...              *Parameters*
  ...              - ${campaign_id}: Campaign identifier to check
  ...
  ...              *Returns*
  ...              Dictionary with:
  ...              - status: Campaign status code (e.g., "Running", "Finished")
  ...              - active_indices: List of indices where inl_rec_api_group_ratio != 0.0
  ...
  ...              *Usage Example*
  ...              | ${campaign_info} = | I would like to check campaign status | abc123def456 |
  ...              | ${status} = | Set Variable | ${campaign_info}[status] |
  ...              | ${indices} = | Set Variable | ${campaign_info}[active_indices] |
  ...
  ...              *Implementation*
  ...              - Queries /v0/campaigns/{campaign_id}
  ...              - Extracts status_code from campaign
  ...              - Parses inl_rec_api_group_ratio array for non-zero values
  ...              - Returns indices (0-based) where ratio > 0.0
  
  [Arguments]             ${campaign_id}

  # Get campaign details
  ${resp} =               Get On Session          ConfigAPISession        url=/v0/campaigns/${campaign_id}    expected_status=404
  
  # If campaign not found (404), treat as archived campaign
  IF  ${resp.status_code} == 404
    Log                   Campaign ${campaign_id} not found (404), treating as Archived
    &{empty_result} =     Create Dictionary       status=Archived  active_indices=${EMPTY}
    RETURN                ${empty_result}
  END
  
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
  [Documentation]  Get vendor_inl_corp campaign IDs that are not in "Finished" status with their active indices.
  ...
  ...              *Purpose*
  ...              Retrieve all active (non-Finished) vendor_inl_corp campaigns with their
  ...              traffic allocation information.
  ...
  ...              *Parameters*
  ...              - ${site_id}: Site identifier [default: android--com.coupang.mobile_s2s_v3]
  ...              - ${service_type_id}: Service type [default: crossx_recommend]
  ...
  ...              *Returns*
  ...              Dictionary with:
  ...              - campaign_ids: List of active campaign IDs
  ...              - campaigns_info: Dictionary mapping campaign_id to {status, active_indices}
  ...
  ...              *Side Effects*
  ...              - Sets ${extracted_active_vendor_inl_corp_campaigns} test variable
  ...              - Sets ${extracted_vendor_inl_corp_campaigns_info} test variable
  ...
  ...              *Usage Example*
  ...              | ${result} = | Get active vendor_inl_corp campaign ids from config api |
  ...              | @{campaign_ids} = | Set Variable | ${result}[campaign_ids] |
  ...              | &{campaigns_info} = | Set Variable | ${result}[campaigns_info] |
  ...
  ...              *Implementation*
  ...              1. Calls "I have a config_api session"
  ...              2. Gets all vendor_inl_corp campaigns from experiments
  ...              3. Checks each campaign's status
  ...              4. Filters out campaigns with status="Finished"
  ...              5. Returns active campaigns with their indices
  
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

      IF  '${campaign_status}' != 'Finished' and '${campaign_status}' != 'Archived'
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
  [Documentation]  Get all active vendor inl_corp indices from config api for yaml configuration testing.
  ...
  ...              *Purpose*
  ...              Generate YAML vendor configuration based on active inl_corp indices
  ...              from Config API campaigns.
  ...
  ...              *Parameters*
  ...              - ${site_id}: Site identifier [default: android--com.coupang.mobile_s2s_v3]
  ...              - ${service_type_id}: Service type [default: crossx_recommend]
  ...              - ${validate_yaml_safety}: Validation flag [default: ${True}]
  ...
  ...              *Returns*
  ...              Dictionary with:
  ...              - indices: Sorted list of unique active indices
  ...              - vendor_names: List of vendor names (inl_corp_0, inl_corp_1, ...)
  ...              - yaml_snippet: YAML configuration snippet
  ...              - max_index: Highest index value
  ...              - vendor_count: Number of active vendors
  ...
  ...              *Side Effects*
  ...              - Sets ${extracted_vendor_indices_for_yaml} test variable
  ...              - Sets ${extracted_vendor_names_for_yaml} test variable
  ...              - Sets ${extracted_yaml_snippet_for_vendors} test variable
  ...              - Logs YAML Config Report to test output
  ...
  ...              *Usage Example*
  ...              | ${config} = | Get vendor inl_corp indices for yaml configuration |
  ...              | @{indices} = | Set Variable | ${config}[indices] |
  ...              | Log | Active indices: ${indices} |
  ...
  ...              *Implementation*
  ...              1. Gets all active vendor_inl_corp campaigns
  ...              2. Collects all active indices from campaigns
  ...              3. Removes duplicates and sorts indices
  ...              4. Generates vendor names (inl_corp_{index})
  ...              5. Creates YAML snippet for vendors section
  
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
  [Documentation]  Get subid mapping for all vendors from Config API campaigns.
  ...
  ...              *Purpose*
  ...              Retrieve subid values for each vendor defined in YAML configuration
  ...              from Config API campaigns with vendor-related default_group.
  ...
  ...              *Parameters*
  ...              - ${yaml_content}: YAML configuration content containing vendor definitions
  ...              - ${site_id}: Site identifier [default: android--com.coupang.mobile_s2s_v3]
  ...              - ${service_type_id}: Service type [default: crossx_recommend]
  ...
  ...              *Returns*
  ...              Dictionary mapping vendor_name to subid value
  ...              (Empty string for vendors without subid)
  ...
  ...              *Usage Example*
  ...              | ${yaml} = | Get File | config.yaml |
  ...              | ${subid_mapping} = | Get vendor subids from config api | ${yaml} |
  ...              | ${inl_subid} = | Get From Dictionary | ${subid_mapping} | inl_corp_1 |
  ...
  ...              *Implementation*
  ...              1. Parses YAML to extract vendor list
  ...              2. Queries experiments API for campaigns with "vendor_" in default_group
  ...              3. For each vendor, searches campaigns for matching subid
  ...              4. If Config API doesn't have subid, uses default subid (from DEFAULT_VENDOR_SUBIDS)
  ...              5. Returns vendor_name → subid dictionary
  ...
  ...              *Fallback Strategy*
  ...              - Primary: Use subid from Config API campaigns
  ...              - Fallback: Use default subid defined in DEFAULT_VENDOR_SUBIDS variable
  ...              - Last resort: Use empty string if no default defined
  ...
  ...              *Prerequisites*
  ...              - Config API session must be initialized
  
  [Arguments]             ${yaml_content}         ${site_id}=android--com.coupang.mobile_s2s_v3  ${service_type_id}=crossx_recommend

  I have a config_api session

  # Parse YAML to get vendor names
  ${yaml_data} =          Evaluate                yaml.safe_load('''${yaml_content}''')  yaml
  ${vendors} =            Evaluate                $yaml_data['vendor_config']['vendors']
  
  # Get all campaign IDs from experiments and filter by default_group
  ${resp} =               Get On Session          ConfigAPISession        url=/v0/recommend/experiments/custom_flow?site_id=${site_id}&service_type_id=${service_type_id}
  Validate config api response  ${resp}  Experiments API error
  
  # Filter campaign IDs by checking if any bucket has default_group containing "vendor_"
  @{vendor_campaign_ids} =  Get Value From Json  ${resp.json()}  $.experiments[*].distributions[?(@.default_group =~ 'vendor_.*')].campaign_id

  # Remove duplicates and empty values
  ${unique_campaign_ids} =  Remove Duplicates      ${vendor_campaign_ids}
  @{filtered_campaign_ids} =  Create List
  FOR  ${campaign_id}  IN  @{unique_campaign_ids}
    IF  '${campaign_id}' != '' and '${campaign_id}' != '${EMPTY}'
      Append To List  ${filtered_campaign_ids}  ${campaign_id}
    END
  END
  
  ${campaign_count} =     Get Length              ${filtered_campaign_ids}
  ${vendor_count} =       Get Length              ${vendors}
  Log                     Searching subids across ${campaign_count} vendor campaigns for ${vendor_count} vendors
  Log                     Vendor campaigns found: ${filtered_campaign_ids}
  
  # Create vendor subid mapping
  &{vendor_subid_mapping} =  Create Dictionary
  
  # For each vendor in YAML, try to find its subid in filtered campaigns
  FOR  ${vendor_config}  IN  @{vendors}
    ${vendor_name} =      Get From Dictionary     ${vendor_config}    name
    ${vendor_subid} =     Find vendor subid in campaigns  ${filtered_campaign_ids}  ${vendor_name}
    
    # If Config API didn't provide subid, try to use default subid
    IF  '${vendor_subid}' == '${EMPTY}'
      ${has_default} =    Run Keyword And Return Status
      ...                 Dictionary Should Contain Key  ${DEFAULT_VENDOR_SUBIDS}  ${vendor_name}
      
      IF  ${has_default}
        ${vendor_subid} =  Get From Dictionary  ${DEFAULT_VENDOR_SUBIDS}  ${vendor_name}
        Set Test Message  ⚙️ Using default subid for ${vendor_name}: ${vendor_subid}  append=yes
        Log               Using default subid for ${vendor_name}: ${vendor_subid}
      ELSE
        Set Test Message  ℹ️ No subid found for ${vendor_name} - will use empty subid  append=yes
      END
    ELSE
      Set Test Message    ✅ Found subid for ${vendor_name}: ${vendor_subid} (from Config API)  append=yes
    END
    
    Set To Dictionary     ${vendor_subid_mapping}  ${vendor_name}=${vendor_subid}
  END
  
  Log                     Vendor subid mapping: ${vendor_subid_mapping}
  RETURN                  ${vendor_subid_mapping}


Find vendor subid in campaigns
  [Documentation]  Search for vendor subid across all campaigns.
  ...
  ...              *Purpose*
  ...              Iterate through campaign IDs to find subid for specific vendor.
  ...
  ...              *Parameters*
  ...              - ${campaign_ids}: List of campaign IDs to search
  ...              - ${vendor_name}: Vendor name to find subid for
  ...
  ...              *Returns*
  ...              Subid value if found, ${EMPTY} otherwise
  ...
  ...              *Usage Example*
  ...              | @{campaigns} = | Create List | camp1 | camp2 | camp3 |
  ...              | ${subid} = | Find vendor subid in campaigns | ${campaigns} | inl_corp_1 |
  ...
  ...              *Implementation*
  ...              - Loops through each campaign ID
  ...              - Calls "Get vendor subid from campaign" for each
  ...              - Returns first non-empty subid found
  
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
  [Documentation]  Get subid for a specific vendor from campaign configs.
  ...
  ...              *Purpose*
  ...              Extract subid value for vendor from campaign's subid_buckets configuration.
  ...
  ...              *Parameters*
  ...              - ${campaign_id}: Campaign identifier
  ...              - ${vendor_name}: Vendor name to match in subid_buckets
  ...
  ...              *Returns*
  ...              Subid value if found, ${EMPTY} otherwise
  ...
  ...              *Usage Example*
  ...              | ${subid} = | Get vendor subid from campaign | abc123 | inl_corp_1 |
  ...
  ...              *Implementation*
  ...              1. Queries /v0/campaigns/{campaign_id}
  ...              2. Checks if subid_buckets exists in configs
  ...              3. Searches for exact match of vendor_name in bucket keys
  ...              4. Falls back to partial match if no exact match
  ...              5. Randomly selects from partial matches if multiple found
  ...              6. Returns subid from first bucket of matched key
  ...
  ...              *Matching Strategy*
  ...              - Exact: "inl_corp_1" matches "inl_corp_1"
  ...              - Partial: "inl_corp_1" matches "inl_corp_1_1200x600"
  
  [Arguments]             ${campaign_id}          ${vendor_name}

  # Get campaign details
  ${resp} =               Get On Session          ConfigAPISession        url=/v0/campaigns/${campaign_id}
  Validate config api response  ${resp}  Campaign ${campaign_id} API error

  # Check if campaign has subid_buckets configuration
  ${has_subid_buckets} =  Run Keyword And Return Status
  ...                     Should Have Value In Json  ${resp.json()}  $.campaign.configs.subid_buckets

  ${subid} =              Set Variable            ${EMPTY}
  IF  ${has_subid_buckets}
    ${subid_buckets_obj} =  Get Value From Json   ${resp.json()}      $.campaign.configs.subid_buckets
    ${subid_buckets_list} =  Set Variable         ${subid_buckets_obj}[0]
    
    # Search through subid_buckets list for matching vendor names
    @{matching_keys} =    Create List
    
    # Collect all keys that start with vendor_name
    FOR  ${bucket_config}  IN  @{subid_buckets_list}
      ${bucket_key} =     Get From Dictionary     ${bucket_config}    key
      
      # Try exact match first
      IF  '${bucket_key}' == '${vendor_name}'
        ${buckets} =      Get From Dictionary     ${bucket_config}    buckets
        ${subid} =        Get From Dictionary     ${buckets[0]}       subid
        Log               Found exact match for ${vendor_name}: ${subid}
        BREAK
      END
      
      # Try partial match (e.g., inl_corp_1 matches inl_corp_1_1200x600)
      ${is_partial_match} =  Run Keyword And Return Status
      ...                    Should Start With   ${bucket_key}       ${vendor_name}
      IF  ${is_partial_match}
        Append To List    ${matching_keys}      ${bucket_config}
      END
    END
    
    # If no exact match found but partial matches exist, randomly select one
    IF  '${subid}' == '${EMPTY}' and @{matching_keys}
      ${matching_count} =  Get Length            ${matching_keys}
      ${random_index} =   Evaluate               random.randint(0, ${matching_count}-1)  random
      ${selected_config} =  Get From List       ${matching_keys}    ${random_index}
      ${selected_key} =   Get From Dictionary   ${selected_config}  key
      ${buckets} =        Get From Dictionary   ${selected_config}  buckets
      ${subid} =          Get From Dictionary   ${buckets[0]}       subid
      Log                 Found partial match for ${vendor_name} (randomly selected key: ${selected_key}): ${subid}
    END
  END

  RETURN                  ${subid}


Validate and generate safe vendor yaml configuration
  [Documentation]  Backward compatible YAML configuration validation - filters only active inl_corp vendors.
  ...
  ...              *Purpose*
  ...              Filter YAML vendor configuration to include only active inl_corp vendors
  ...              based on Config API traffic allocation.
  ...
  ...              *Parameters*
  ...              - ${original_yaml_content}: Original YAML configuration content
  ...
  ...              *Returns*
  ...              Filtered YAML configuration (original if no inl vendors found)
  ...
  ...              *Usage Example*
  ...              | ${original_yaml} = | Get File | config.yaml |
  ...              | ${safe_yaml} = | Validate and generate safe vendor yaml configuration | ${original_yaml} |
  ...
  ...              *Implementation*
  ...              1. Parses original YAML to get vendor list
  ...              2. Checks if any vendor contains "inl"
  ...              3. If no inl vendors: returns original YAML (backward compatibility)
  ...              4. If inl vendors exist:
  ...                 - Gets active indices from Config API
  ...                 - Keeps all non-inl vendors
  ...                 - Keeps only active inl_corp_{index} vendors
  ...                 - Rebuilds YAML with filtered vendor list
  ...              5. Logs filtering details to test output
  ...
  ...              *Side Effects*
  ...              - Logs kept/filtered vendors with ✅/❌ icons
  ...              - Sets test messages with filtering summary
  
  [Arguments]             ${original_yaml_content}

  # Parse original YAML
  ${yaml_data} =          Evaluate                yaml.safe_load('''${original_yaml_content}''')  yaml
  ${original_vendors} =   Evaluate                $yaml_data['vendor_config']['vendors']

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
  Set To Dictionary       ${yaml_data}[vendor_config]  vendors=${filtered_vendors}
  ${filtered_yaml} =      Evaluate                yaml.dump($yaml_data, default_flow_style=False)  yaml

  ${total_count} =        Get Length              ${filtered_vendors}
  Log                     Final configuration: ${total_count} vendors (active indices: ${active_indices})
  Set Test Message        📊 Final configuration: ${total_count} vendors (active indices: ${active_indices})  append=yes

  RETURN                  ${filtered_yaml}


Get keeta campaign configuration
  [Documentation]  Get Keeta campaign configuration by searching for running campaigns.
  ...
  ...              *Purpose*
  ...              Retrieve keeta_campaign_name value from active Keeta campaigns
  ...              for use in Keeta vendor API testing.
  ...
  ...              *Returns*
  ...              Keeta campaign name (string) if found, ${EMPTY} otherwise
  ...
  ...              *Usage Example*
  ...              | ${k_campaign} = | Get keeta campaign configuration |
  ...              | Should Not Be Empty | ${k_campaign} |
  ...
  ...              *Implementation*
  ...              - Queries /v0/campaigns?status_code=Running
  ...              - Filters campaigns with:
  ...                * datafeed_id = android--com.sankuai.sailor.afooddelivery_2
  ...                * Non-empty keeta_campaign_name in configs
  ...              - Uses JSONPath for efficient filtering
  ...              - Returns keeta_campaign_name from first matching campaign
  ...
  ...              *Prerequisites*
  ...              - Config API session must be initialized
  
  I have a config_api session
  
  # Get all running campaigns
  ${resp} =               Get On Session          ConfigAPISession        url=/v0/campaigns?status_code=Running
  Validate config api response  ${resp}  Running campaigns API error
  
  # Use JSONPath to directly filter for Keeta campaigns
  # Filter: campaigns with datafeed_id=android--com.sankuai.sailor.afooddelivery_2 and non-empty keeta_campaign_name
  @{keeta_campaigns} =    Get Value From Json     ${resp.json()}          $.campaigns[?(@.datafeed_id == 'android--com.sankuai.sailor.afooddelivery_2' & @.configs.keeta_campaign_name != '')]
  
  ${keeta_campaign_count} =  Get Length           ${keeta_campaigns}
  Log                     Found ${keeta_campaign_count} Keeta campaigns matching criteria
  
  # If we found any matching campaigns, use the first one
  IF  ${keeta_campaign_count} > 0
    ${selected_campaign} =  Set Variable          ${keeta_campaigns}[0]
    ${campaign_id} =      Get Value From Json     ${selected_campaign}    $.campaign_id
    ${datafeed_id} =      Get Value From Json     ${selected_campaign}    $.datafeed_id
    ${keeta_campaign_name} =  Get Value From Json  ${selected_campaign}   $.configs.keeta_campaign_name
    
    Log                   ✅ Found Keeta campaign: ${campaign_id}[0]
    Log                   ✅ Keeta campaign name: ${keeta_campaign_name}[0]
    Log                   ✅ Datafeed ID: ${datafeed_id}[0]
    RETURN                ${keeta_campaign_name}[0]
  END
  
  # If we get here, no valid campaign was found
  Log                     ❌ No valid Keeta campaign found with criteria:
  Log                     - status_code=Running
  Log                     - datafeed_id=android--com.sankuai.sailor.afooddelivery_2
  Log                     - non-empty keeta_campaign_name
  
  RETURN                  ${EMPTY}
