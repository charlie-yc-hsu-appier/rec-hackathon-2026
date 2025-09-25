*** Keywords ***
# Vendor API testing utility keywords #
Auto select test dimensions
  [Arguments]         ${request_url}=${Empty}  ${vendor_name}=${Empty}
  [Documentation]  Auto-select test dimensions from predefined sizes
  ...              Returns dimensions dictionary with width, height, and additional vendor parameters
  ...              Available test sizes: 300x300, 1200x627, 1200x600
  ...              For linkmine vendor: Also generates app_bundleId (empty string), imp_adType (2 or 3)
  ...              Note: request_url parameter is kept for compatibility but not used

  # Predefined test dimensions
  @{test_sizes} =     Create List
  ...                 300x300
  ...                 1200x627
  ...                 1200x600

  # Select a random size from the predefined list for testing
  ${list_length} =    Get Length          ${test_sizes}
  ${random_index} =   Evaluate            __import__('random').randint(0, ${list_length}-1)
  ${selected_size} =  Set Variable        ${test_sizes}[${random_index}]

  # Parse the selected size
  ${size_parts} =     Split String        ${selected_size}        x
  ${width} =          Set Variable        ${size_parts}[0]
  ${height} =         Set Variable        ${size_parts}[1]

  # Base dimensions dictionary
  &{dimensions} =     Create Dictionary
  ...                 width=${width}
  ...                 height=${height}

  Log                 📏 Selected test dimensions: ${width}x${height}

  # Generate additional vendor parameters only for linkmine
  ${is_linkmine} =    Run Keyword And Return Status
  ...                 Should Be Equal     ${vendor_name}          linkmine

  IF  ${is_linkmine}
    # Linkmine-specific parameters with updated requirements
    @{ad_types} =       Create List     2                       3

    # Random selection for adtype only
    ${adtype} =         Evaluate        __import__('random').choice($ad_types)

    # Set bundle_id as empty string and remove web_host (site_domain)
    ${bundle_id} =      Set Variable    ${Empty}

    # Add to dimensions dictionary (no web_host anymore)
    Set To Dictionary   ${dimensions}   bundle_id=${bundle_id}  adtype=${adtype}

    Log                 Linkmine params - bundle_id: ${bundle_id} (empty), adtype: ${adtype}
  END

  RETURN              &{dimensions}


Generate UUID4
  [Documentation]  Generate a random UUID4 for user_id
  ${uuid} =   Evaluate    str(__import__('uuid').uuid4())
  RETURN      ${uuid}


Encode Base64
  [Arguments]     ${text}
  [Documentation]  Encode text to base64
  ${encoded} =    Evaluate    __import__('base64').b64encode('${text}'.encode()).decode()
  RETURN          ${encoded}


Parse tracking config
  [Arguments]             ${tracking_queries}     ${vendor_name}=${Empty}
  [Documentation]  Parse tracking queries to extract parameter configuration
  ...              Example tracking queries: [{"key": "param1", "value": "{click_id_base64}"}]
  ...              Returns config dictionary with param_name, uses_base64, and has_group_id

  # Special handling for adpacker vendor
  ${is_adpacker_vendor} =  Run Keyword And Return Status
  ...                      Should Be Equal     ${vendor_name}      adpacker
  
  # Special handling for INL vendors
  ${is_inl_vendor} =      Run Keyword And Return Status
  ...                     Should Contain      ${vendor_name}      inl

  ${final_param_name} =   Set Variable        unknown
  ${uses_base64} =        Set Variable        ${FALSE}
  ${has_group_id} =       Set Variable        ${FALSE}

  # Find click_id related parameters in queries
  FOR  ${query}  IN  @{tracking_queries}
    ${key} =              Get From Dictionary  ${query}            key
    ${value} =            Get From Dictionary  ${query}            value
    
    # Check if this query contains click_id
    ${is_click_id_param} =  Run Keyword And Return Status
    ...                     Should Contain       ${value}            click_id
    
    IF  ${is_click_id_param}
      ${final_param_name} =  Set Variable       ${key}
      
      # Check if uses base64 encoding
      ${uses_base64} =    Run Keyword And Return Status
      ...                 Should Contain       ${value}            base64
      
      # Check if requires group_id (typically for INL vendors with subparam)
      ${has_group_id} =   Set Variable If     '${key}' == 'subparam'  ${TRUE}  ${FALSE}
      
      BREAK
    END
  END

  # Special handling for adpacker vendor
  IF  ${is_adpacker_vendor}
    ${final_param_name} =   Set Variable        ssp_click_id
    ${uses_base64} =        Set Variable        ${TRUE}
  END

  # Special handling for INL vendors without explicit tracking queries
  IF  ${is_inl_vendor} and '${final_param_name}' == 'unknown'
    Log                     INL vendor detected without tracking queries, using subparam base64 parameter for compatibility
    ${final_param_name} =   Set Variable        subparam
    ${uses_base64} =        Set Variable        ${TRUE}
    ${has_group_id} =       Set Variable        ${TRUE}
  END

  # Create config dictionary
  &{config} =             Create Dictionary
  ...                     param_name=${final_param_name}
  ...                     uses_base64=${uses_base64}
  ...                     has_group_id=${has_group_id}

  Log                     📋 Tracking config for ${vendor_name}: ${config}
  RETURN                  &{config}


Load vendor config from file
  [Arguments]         ${config_file_path}=${Empty}
  [Documentation]  Load vendor configuration from config.yaml file
  ...              Returns vendors section as YAML string for testing
  ...              Default path: deploy/rec-vendor-api/secrets/config.yaml (from project root)

  # Calculate default config file path if not provided
  IF  '${config_file_path}' == '${Empty}'
    # From keywords_util.robot location: lib -> API -> system_tests -> tests -> project_root
    ${project_root} =   Set Variable    ${CURDIR}/../../../../..
    ${config_file_path} = Set Variable  ${project_root}/deploy/rec-vendor-api/secrets/config.yaml
  END

  # Read the configuration file
  ${yaml_content} =   Get File                ${config_file_path}

  # Parse the YAML - vendors directly at top level
  ${config_data} =    Evaluate                yaml.safe_load('''${yaml_content}''')  yaml

  # Convert back to YAML string for testing framework
  ${vendor_yaml} =    Evaluate                yaml.dump($config_data)  yaml

  Log                 📄 Loaded vendor config from: ${config_file_path}
  RETURN              ${vendor_yaml}
