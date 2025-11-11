*** Settings ***
Resource            ../res/init.robot

*** Keywords ***
# Vendor API testing utility keywords #
Auto select test dimensions
  [Documentation]  Auto-select random test dimensions and generate vendor-specific parameters.
  ...
  ...              *Purpose*
  ...              Randomly selects advertising dimensions from predefined sizes and generates
  ...              additional parameters required by specific vendors (linkmine, adpacker, adforus).
  ...
  ...              *Parameters*
  ...              - ${request_url}: (Legacy, not used) Kept for compatibility [default: ${Empty}]
  ...              - ${vendor_name}: Vendor identifier for parameter generation [default: ${Empty}]
  ...
  ...              *Returns*
  ...              Dictionary containing:
  ...              - width: Ad width in pixels
  ...              - height: Ad height in pixels
  ...              - adtype: Ad type (2 or 3) - for linkmine/adpacker vendors only
  ...              - bundle_id: Bundle identifier (empty string) - for linkmine vendor only
  ...              - os: Operating system (android/ios) - for adforus vendor only
  ...
  ...              *Usage Example*
  ...              | ${dimensions} = | Auto select test dimensions | vendor_name=linkmine |
  ...              | ${width} = | Get From Dictionary | ${dimensions} | width |
  ...              | ${height} = | Get From Dictionary | ${dimensions} | height |
  ...              | ${adtype} = | Get From Dictionary | ${dimensions} | adtype |
  ...
  ...              *Implementation*
  ...              1. Randomly selects from predefined sizes: 300x300, 1200x627, 1200x600
  ...              2. Parses selected size into width and height
  ...              3. For linkmine vendor: adds bundle_id (empty) and adtype (2 or 3)
  ...              4. For adpacker vendor: adds adtype (2 or 3)
  ...              5. For adforus vendor: adds os (android, will test both OS in template)
  ...
  ...              *Vendor-Specific Parameters*
  ...              - Linkmine: bundle_id=${Empty}, adtype=(2|3)
  ...              - Adpacker: adtype=(2|3)
  ...              - Adforus: os=android (both android/ios tested in calling template)
  
  [Arguments]         ${request_url}=${Empty}  ${vendor_name}=${Empty}

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

  # Generate additional vendor parameters for vendors that need them
  ${is_linkmine} =    Run Keyword And Return Status
  ...                 Should Be Equal     ${vendor_name}          linkmine
  ${is_adpacker} =    Run Keyword And Return Status
  ...                 Should Be Equal     ${vendor_name}          adpacker
  ${is_adforus} =     Run Keyword And Return Status
  ...                 Should Be Equal     ${vendor_name}          adforus

  # Generate adtype parameter for both linkmine and adpacker
  ${needs_adtype} =   Set Variable If     ${is_linkmine} or ${is_adpacker}     ${TRUE}     ${FALSE}
  
  IF  ${needs_adtype}
    # Common parameters for vendors that need adtype
    @{ad_types} =       Create List     2                       3
    ${adtype} =         Evaluate        __import__('random').choice($ad_types)
    
    # Add adtype to dimensions
    Set To Dictionary   ${dimensions}   adtype=${adtype}
    
    Log                 Generated adtype for ${vendor_name}: ${adtype}
  END

  # Linkmine-specific parameters (bundle_id)
  IF  ${is_linkmine}
    # Set bundle_id as empty string for linkmine
    ${bundle_id} =      Set Variable    ${Empty}
    Set To Dictionary   ${dimensions}   bundle_id=${bundle_id}
    
    Log                 Linkmine-specific params - bundle_id: ${bundle_id} (empty)
  END

  # Adforus-specific parameters (os for user_id transformation)
  IF  ${is_adforus}
    # For adforus vendor, we need to test both android and ios
    # This will be handled in the calling template to ensure both OS are tested
    # Default to android for single dimension generation (will be overridden in template)
    ${os} =             Set Variable    android
    Set To Dictionary   ${dimensions}   os=${os}
    
    Log                 Adforus-specific params - os: ${os} (will test both android and ios)
  END

  RETURN              &{dimensions}


Generate UUID4
  [Documentation]  Generate random UUID4 string for user_id parameter.
  ...
  ...              *Purpose*
  ...              Creates universally unique identifier for tracking user requests
  ...              in vendor API testing.
  ...
  ...              *Returns*
  ...              UUID4 string (e.g., "a1b2c3d4-e5f6-7890-abcd-ef1234567890")
  ...
  ...              *Usage Example*
  ...              | ${user_id} = | Generate UUID4 |
  ...              | Log | Generated user_id: ${user_id} |
  ...
  ...              *Implementation*
  ...              Uses Python's uuid.uuid4() to generate random UUID.
  
  ${uuid} =   Evaluate    str(__import__('uuid').uuid4())
  RETURN      ${uuid}


Encode Base64
  [Documentation]  Encode text string to base64 format for click_id parameters.
  ...
  ...              *Purpose*
  ...              Converts click_id or other tracking parameters to base64 encoding
  ...              as required by some vendor APIs.
  ...
  ...              *Parameters*
  ...              - ${text}: Text string to encode
  ...
  ...              *Returns*
  ...              Base64 encoded string
  ...
  ...              *Usage Example*
  ...              | ${click_id} = | Set Variable | 12aa.12aa |
  ...              | ${encoded} = | Encode Base64 | ${click_id} |
  ...              | Log | Encoded click_id: ${encoded} |
  ...
  ...              *Implementation*
  ...              Uses Python's base64.b64encode() to encode UTF-8 text.
  
  [Arguments]     ${text}
  ${encoded} =    Evaluate    __import__('base64').b64encode('${text}'.encode()).decode()
  RETURN          ${encoded}


Parse tracking config
  [Documentation]  Extract tracking parameter configuration from vendor tracking queries.
  ...
  ...              *Purpose*
  ...              Parses vendor YAML tracking configuration to determine click_id parameter
  ...              name, encoding requirements, and group_id needs for URL validation.
  ...
  ...              *Parameters*
  ...              - ${tracking_queries}: List of tracking query dictionaries from YAML
  ...                Example: [{"key": "subparam", "value": "{click_id_base64}"}]
  ...              - ${vendor_name}: Vendor identifier for special handling [default: ${Empty}]
  ...
  ...              *Returns*
  ...              Dictionary containing:
  ...              - param_name: Click tracking parameter name (e.g., 'subparam', 'ssp_click_id')
  ...              - uses_base64: Boolean - true if click_id requires base64 encoding
  ...              - has_group_id: Boolean - true if parameter requires group_id (INL vendors)
  ...
  ...              *Usage Example*
  ...              | ${tracking_queries} = | Get From Dictionary | ${vendor} | tracking.queries |
  ...              | ${config} = | Parse tracking config | ${tracking_queries} | inl_corp_1 |
  ...              | ${param_name} = | Get From Dictionary | ${config} | param_name |
  ...
  ...              *Implementation*
  ...              1. Searches tracking_queries for click_id-related parameters
  ...              2. Checks if value contains 'base64' for encoding requirement
  ...              3. For 'subparam' key: sets has_group_id=True (INL vendors)
  ...              4. Special handling for adpacker: param_name='ssp_click_id', uses_base64=True
  ...              5. Special handling for INL vendors without explicit config: defaults to subparam+base64
  ...
  ...              *Special Cases*
  ...              - Adpacker vendor: Always uses 'ssp_click_id' with base64
  ...              - INL vendors: Default to 'subparam' with base64 and group_id if not configured
  
  [Arguments]             ${tracking_queries}     ${vendor_name}=${Empty}

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
  [Documentation]  Load vendor configuration from YAML file for testing.
  ...
  ...              *Purpose*
  ...              Reads vendor configuration YAML file and returns content for
  ...              automated vendor testing workflows.
  ...
  ...              *Parameters*
  ...              - ${config_file_path}: Path to config.yaml file [default: ${Empty}]
  ...                If empty, uses default path: {project_root}/deploy/rec-vendor-api/secrets/config.yaml
  ...
  ...              *Returns*
  ...              YAML content as string containing vendor_config section
  ...
  ...              *Usage Example*
  ...              | ${yaml_content} = | Load vendor config from file |
  ...              | # Uses default path |
  ...              | ${yaml_content} = | Load vendor config from file | /custom/path/config.yaml |
  ...
  ...              *Implementation*
  ...              1. If config_file_path empty: calculates default path from project root
  ...              2. Reads YAML file content
  ...              3. Parses YAML to validate structure
  ...              4. Converts back to YAML string for testing framework
  ...              5. Returns complete YAML content
  ...
  ...              *Default Path Calculation*
  ...              From keywords_util.robot location:
  ...              lib → API → system_tests → tests → project_root
  ...              Then: project_root/deploy/rec-vendor-api/secrets/config.yaml
  ...
  ...              *Prerequisites*
  ...              - config.yaml file must exist at specified or default path
  ...              - YAML must contain valid vendor_config structure
  
  [Arguments]         ${config_file_path}=${Empty}

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

