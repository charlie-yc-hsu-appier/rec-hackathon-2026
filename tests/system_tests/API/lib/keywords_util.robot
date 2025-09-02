*** Keywords ***
# Vendor API testing utility keywords #
Auto select test dimensions
  [Arguments]         ${request_url}=${Empty}  ${vendor_name}=${Empty}
  [Documentation]  Auto-select test dimensions from predefined sizes
  ...              Returns dimensions dictionary with width, height, and additional vendor parameters
  ...              Available test sizes: 300x300, 1200x627, 1200x600
  ...              For linkmine vendor: Also generates site_domain (web domain), app_bundleId (app store ID), imp_adType (2 or 3)
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
    # Linkmine-specific parameters with predefined options
    @{domains} =        Create List     coupang.com             gmarket.co.kr           11st.co.kr          auction.co.kr
    @{bundles} =        Create List     com.coupang.mobile      kr.co.gmarket.mobile    com.elevenst        com.auction.mobile
    @{ad_types} =       Create List     2                       3

    # Random selection
    ${web_host} =       Evaluate        __import__('random').choice($domains)
    ${bundle_id} =      Evaluate        __import__('random').choice($bundles)
    ${adtype} =         Evaluate        __import__('random').choice($ad_types)

    # Add to dimensions dictionary
    Set To Dictionary   ${dimensions}   web_host=${web_host}    bundle_id=${bundle_id}  adtype=${adtype}

    Log                 Linkmine params - web_host: ${web_host}, bundle_id: ${bundle_id}, adtype: ${adtype}
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


Parse yaml tracking url template
  [Arguments]             ${tracking_url}     ${vendor_name}=${Empty}
  [Documentation]  Parse YAML tracking_url to extract parameter configuration
  ...              Example: "{product_url}&param1={click_id_base64}" -> param_name=param1, uses_base64=true
  ...              Special handling for INL vendors: if tracking_url is "{product_url}" only,
  ...              automatically adds subparam={click_id_base64} for compatibility
  ...              Note: INL vendors use 'subparam' parameter in the tracking template

  # Special handling for INL vendors
  ${is_inl_vendor} =      Run Keyword And Return Status
  ...                     Should Contain      ${vendor_name}      inl
  ${modified_tracking_url} =  Set Variable  ${tracking_url}
  ${final_param_name} =   Set Variable        unknown

  IF  ${is_inl_vendor} and '${tracking_url}' == '{product_url}'
    Log                     INL vendor detected with simple tracking_url, adding subparam base64 parameter for compatibility
    ${modified_tracking_url} =  Set Variable  {product_url}&subparam={click_id_base64}
    ${final_param_name} =   Set Variable        subparam
  ELSE
    # Extract parameter name using regex for non-INL vendors
    ${param_matches} =      Get Regexp Matches  ${modified_tracking_url}  [&?]([^=]+)=  1
    ${final_param_name} =   Set Variable If     ${param_matches}    ${param_matches[0]}     unknown
  END

  # Check if uses base64 encoding
  ${uses_base64} =        Run Keyword And Return Status
  ...                     Should Contain      ${modified_tracking_url}  base64

  # Check if requires group_id (typically for INL vendors)
  ${has_group_id} =       Set Variable If     '${final_param_name}' == 'subparam'  ${TRUE}  ${FALSE}

  # Create config dictionary
  &{config} =             Create Dictionary
  ...                     param_name=${final_param_name}
  ...                     uses_base64=${uses_base64}
  ...                     has_group_id=${has_group_id}
  ...                     modified_tracking_url=${modified_tracking_url}

  RETURN                  &{config}


Load vendor config from file
  [Arguments]         ${config_file_path}=${Empty}
  [Documentation]  Load vendor configuration from config.yaml file
  ...              Returns vendor_config section as YAML string for testing
  ...              Default path: deploy/rec-vendor-api/secrets/config.yaml (from project root)

  # Calculate default config file path if not provided
  IF  '${config_file_path}' == '${Empty}'
    # From keywords_util.robot location: lib -> API -> system_tests -> tests -> project_root
    ${project_root} =   Set Variable    ${CURDIR}/../../../../..
    ${config_file_path} = Set Variable  ${project_root}/deploy/rec-vendor-api/secrets/config.yaml
  END

  # Read the configuration file
  ${yaml_content} =   Get File                ${config_file_path}

  # Parse the YAML to extract only vendor_config section
  ${config_data} =    Evaluate                yaml.safe_load('''${yaml_content}''')  yaml
  ${vendor_config} =  Get From Dictionary     ${config_data}          vendor_config

  # Convert vendor_config back to YAML string for testing framework
  ${vendor_yaml} =    Evaluate                yaml.dump($vendor_config)  yaml

  Log                 📄 Loaded vendor config from: ${config_file_path}
  RETURN              ${vendor_yaml}
