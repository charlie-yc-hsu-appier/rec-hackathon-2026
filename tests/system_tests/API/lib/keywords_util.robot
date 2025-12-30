*** Settings ***
Resource            ../res/init.robot

*** Keywords ***
# Vendor API testing utility keywords #
Auto select test dimensions
  [Documentation]  Auto-select random test dimensions and generate vendor-specific parameters from configuration.
  ...
  ...              *Purpose*
  ...              Randomly selects advertising dimensions from predefined sizes and generates
  ...              additional parameters required by vendors based on their YAML configuration.
  ...              Parameters are automatically detected from request.queries configuration.
  ...
  ...              *Parameters*
  ...              - ${request_url}: (Legacy, not used) Kept for compatibility [default: ${Empty}]
  ...              - ${request_queries}: List of request query dictionaries from vendor.request.queries
  ...              - ${vendor_name}: Vendor identifier for special handling [default: ${Empty}]
  ...
  ...              *Returns*
  ...              Dictionary containing:
  ...              - width: Ad width in pixels
  ...              - height: Ad height in pixels
  ...              - adtype: Ad type (randomly selected: 2 or 3) - if vendor config uses {adtype}
  ...              - bundle_id: Bundle identifier (randomly selected from common apps) - if vendor config uses {bundle_id}
  ...              - os: Operating system (android/ios) - for adforus vendor only (special test requirement)
  ...
  ...              *Usage Example*
  ...              | ${dimensions} = | Auto select test dimensions | ${request_url} | ${request_queries} | ${vendor_name} |
  ...              | ${width} = | Get From Dictionary | ${dimensions} | width |
  ...
  ...              *Implementation*
  ...              1. Randomly selects from predefined sizes: 300x300, 1200x627, 1200x600
  ...              2. Parses selected size into width and height
  ...              3. Scans request.queries for {adtype} → generates random adtype (2 or 3) if found
  ...              4. Scans request.queries for {bundle_id} → generates random bundle_id if found
  ...              5. For adforus vendor: adds os (android, will test both OS in template)
  ...
  ...              *Configuration-Driven Approach*
  ...              - Automatically detects required parameters from YAML configuration
  ...              - No hardcoded vendor names needed (except adforus OS handling)
  ...              - New vendors are automatically supported if they use standard parameter names
  
  [Arguments]         ${request_url}=${Empty}  ${request_queries}=@{EMPTY}  ${vendor_name}=${Empty}

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

  # Configuration-driven parameter detection from request.queries
  ${needs_adtype} =   Set Variable        ${FALSE}
  ${needs_bundle_id} =  Set Variable      ${FALSE}
  
  # Check request_queries for required parameters
  FOR  ${query}  IN  @{request_queries}
    ${value} =        Get From Dictionary  ${query}  value
    
    # Check if this vendor uses {adtype}
    ${has_adtype} =   Run Keyword And Return Status
    ...               Should Contain      ${value}  {adtype}
    ${needs_adtype} =  Set Variable If    ${has_adtype}    ${TRUE}    ${needs_adtype}
    
    # Check if this vendor uses {bundle_id}
    ${has_bundle_id} =  Run Keyword And Return Status
    ...                 Should Contain   ${value}  {bundle_id}
    ${needs_bundle_id} =  Set Variable If    ${has_bundle_id}    ${TRUE}    ${needs_bundle_id}
  END
  
  # Generate adtype if needed (detected from config)
  IF  ${needs_adtype}
    @{ad_types} =       Create List     2                       3
    ${adtype} =         Evaluate        __import__('random').choice($ad_types)
    Set To Dictionary   ${dimensions}   adtype=${adtype}
    Log                 Generated adtype (from config): ${adtype}
  END

  # Generate bundle_id if needed (detected from config)
  IF  ${needs_bundle_id}
    @{bundles} =        Create List     com.coupang.mobile      kr.co.gmarket.mobile    com.elevenst        com.auction.mobile
    ${bundle_id} =      Evaluate        __import__('random').choice($bundles)
    Set To Dictionary   ${dimensions}   bundle_id=${bundle_id}
    Log                 Generated bundle_id (from config): ${bundle_id}
  END

  # Adforus-specific parameters (os for user_id transformation)
  # This is special test requirement, not in YAML config
  ${is_adforus} =     Run Keyword And Return Status
  ...                 Should Be Equal     ${vendor_name}          adforus
  
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
  [Documentation]  Extract tracking parameter configuration from vendor queries (request and tracking).
  ...
  ...              *Purpose*
  ...              Parses vendor YAML configuration to determine click_id parameter
  ...              name, encoding requirements, and group_id needs for URL validation.
  ...              Configuration is automatically detected from YAML, no hardcoding needed.
  ...
  ...              *Parameters*
  ...              - ${request_queries}: List of request query dictionaries from vendor.request.queries
  ...              - ${tracking_queries}: List of tracking query dictionaries from vendor.tracking.queries
  ...              - ${vendor_name}: Vendor identifier for logging [default: ${Empty}]
  ...
  ...              *Returns*
  ...              Dictionary containing:
  ...              - param_name: Click tracking parameter name (e.g., 'click_id', 'ssp_click_id', 'puid')
  ...              - uses_base64: Boolean - true if click_id requires base64 encoding
  ...              - has_group_id: Boolean - true if parameter requires group_id (INL vendors)
  ...
  ...              *Usage Example*
  ...              | ${request_queries} = | Get From Dictionary | ${vendor.request} | queries |
  ...              | ${tracking_queries} = | Get From Dictionary | ${vendor.tracking} | queries |
  ...              | ${config} = | Parse tracking config | ${request_queries} | ${tracking_queries} | vendor_name |
  ...              | ${param_name} = | Get From Dictionary | ${config} | param_name |
  ...
  ...              *Implementation*
  ...              1. Searches request_queries first for click_id-related parameters
  ...              2. Falls back to tracking_queries if not found in request
  ...              3. Checks if value contains 'base64' for encoding requirement
  ...              4. For 'subparam' key: sets has_group_id=True (INL vendors)
  ...              5. For INL vendors without explicit config: defaults to subparam+base64
  ...
  ...              *Configuration-Driven Approach*
  ...              - No hardcoded vendor names - all behavior driven by YAML config
  ...              - Automatically handles new vendors without code changes
  
  [Arguments]             ${request_queries}      ${tracking_queries}     ${vendor_name}=${Empty}

  ${final_param_name} =   Set Variable        unknown
  ${uses_base64} =        Set Variable        ${FALSE}
  ${has_group_id} =       Set Variable        ${FALSE}
  
  # Check if this is an INL vendor (special handling needed)
  ${is_inl_vendor} =      Run Keyword And Return Status
  ...                     Should Contain      ${vendor_name}      inl

  # Priority 1: Check tracking.queries first (these appear in returned URLs)
  ${tracking_count} =     Get Length          ${tracking_queries}
  IF  ${tracking_count} > 0
    FOR  ${query}  IN  @{tracking_queries}
      ${key} =            Get From Dictionary  ${query}            key
      ${value} =          Get From Dictionary  ${query}            value
      
      # Check if this query contains click_id
      ${is_click_id_param} =  Run Keyword And Return Status
      ...                     Should Contain       ${value}            click_id
      
      IF  ${is_click_id_param}
        ${final_param_name} =  Set Variable       ${key}
        
        # Check if uses base64 encoding
        ${uses_base64} =  Run Keyword And Return Status
        ...               Should Contain       ${value}            base64
        
        # Check if requires group_id (typically for INL vendors with subparam)
        ${has_group_id} =  Set Variable If     '${key}' == 'subparam'  ${TRUE}  ${FALSE}
        
        Log             Found click_id parameter in tracking.queries for ${vendor_name}: ${key} (base64=${uses_base64})
        BREAK
      END
    END
  END

  # Priority 2: If not found in tracking.queries, check request.queries
  # (for vendors like adpacker/binalab where click_id param is in request)
  # Skip this for INL vendors as they use subparam fallback
  IF  '${final_param_name}' == 'unknown' and not ${is_inl_vendor}
    ${request_count} =    Get Length          ${request_queries}
    IF  ${request_count} > 0
      FOR  ${query}  IN  @{request_queries}
        ${key} =          Get From Dictionary  ${query}            key
        ${value} =        Get From Dictionary  ${query}            value
        
        # Check if this query contains click_id
        ${is_click_id_param} =  Run Keyword And Return Status
        ...                     Should Contain       ${value}            click_id
        
        IF  ${is_click_id_param}
          ${final_param_name} =  Set Variable       ${key}
          
          # Check if uses base64 encoding
          ${uses_base64} =  Run Keyword And Return Status
          ...               Should Contain       ${value}            base64
          
          # Check if requires group_id (typically for INL vendors with subparam)
          ${has_group_id} =  Set Variable If     '${key}' == 'subparam'  ${TRUE}  ${FALSE}
          
          Log             Found click_id parameter in request.queries for ${vendor_name}: ${key} (base64=${uses_base64})
          BREAK
        END
      END
    END
  END

  # Priority 3: Fallback for INL vendors without explicit click_id in queries
  # INL vendors typically use 'subparam' in tracking URLs even if not explicitly configured
  IF  ${is_inl_vendor} and '${final_param_name}' == 'unknown'
    Log                     INL vendor detected without click_id in queries, using subparam base64 parameter for compatibility
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

