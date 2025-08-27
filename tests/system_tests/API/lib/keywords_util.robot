*** Keywords ***
# Utility #
Replacing the special characters
  [Tags]              robot:flatten
  [Arguments]         ${string}
  ${temp_string} =    Replace String  ${temp_string}  @   %40
  ${temp_string} =    Replace String  ${temp_string}  &   %26
  ${temp_string} =    Replace String  ${temp_string}  $   %24
  ${temp_string} =    Replace String  ${temp_string}  '   %27
  ${temp_string} =    Replace String  ${temp_string}  "   %22
  Set Test Variable   ${temp_string}  ${temp_string}


Replacing the robotframework reserved words
  [Tags]              robot:flatten
  [Arguments]         ${string}
  Set Test Variable   ${temp_string}      ${string}

  # When is none = TRUE, then Replacing the special characters
  ${is none} =        Run Keyword And Return Status
  ...                 Should Be Equal     ${string}   ${None}
  IF  ${is none} != ${TRUE}
    Replacing the special characters  ${string}
  END
  RETURN              ${temp_string}


# These legacy functions are kept for potential future use with recommendation system testing
Check "${expected_item}" value in each JSON response items
  [Arguments]             ${layout_id}=${Empty}
  Log                     Start checking "${expected_item}" in the JSON response items.  level=DEBUG

  # Get the test tags
  ${is_psid_case} =       Run Keyword And Return Status
  ...                     List Should Contain Value  ${TEST TAGS}  psid

  Variable Should Exist   ${resp_json['value']}
  ${resp_json_value} =    Set Variable            ${resp_json['value']}
  Should Not Be Empty     ${resp_json_value}      @{TEST TAGS} The value[] is empty, Request:${request.url}, R_Header:${request.headers}

  @{list_get_item} =      Create List
  log                     ${Check_IAB_With_partner_id}

  FOR  ${item}  IN  @{resp_json_value}
    Log                     ${item}             level=DEBUG
    Dictionary Should Contain Key  ${item}  ${expected_item}
    ...                     @{TEST TAGS} No '${expected_item}' Key within the response, Request:${request.url}, R_Header:${request.headers}

    Set Test Variable       ${get_item}         ${item}[${expected_item}]
    # The "&", "@", "$" will let robot become confused and treat it as a dictionary, list or var
    # replacing the string to solve it
    ${get_item} =           Replacing the robotframework reserved words  ${get_item}

    # Check the category of each sku via the Produt API when ${check_IAB_with_partner_id} is not Empty
    IF  '${check_IAB_with_partner_id}' != '${Empty}' and '${expected_item}' == 'sku'
      Check the IAB rule with product API  ${extracted_df}  ${get_item}  ${check_IAB_with_partner_id}
    ELSE IF  ${is_psid_case} and '${expected_item}' == 'sku'
      Check the PSID with product API  ${extracted_df}  ${get_item}  ${extracted_an_actived_psid}
    END

    Variable Should Exist   ${get_item}         No value within the Key "${expected_item}"
    Append To List          ${list_get_item}    ${get_item}
  END

  # Log message
  IF  '${check_IAB_with_partner_id}' != '${Empty}' and '${expected_item}' == 'sku'
    No Operation
  ELSE
    Log     ${list_get_item}    level=DEBUG
  END

  # Check the kakao_kr_v2 path
  IF  '${expected_item}' == 'custom_dna_layout' or '${expected_item}' == 'video_link'
    ${with_layout_id} =     Run Keyword And Return Status
    ...                     Variable Should Exist   ${layout_id}

    IF  ${with_layout_id}
      Check the regex pattern in the given list  ${list_get_item}  (${layout_id})
    ELSE
      Fail    Please provide the valid layout_id for verify the value of path within the custom_dna_layout
    END
  END

  # Download and check the size of each imgs and set the ${extracted_img_url}
  IF  '${expected_item}' == 'img' or '${expected_item}' == 'custom_dna_layout' or '${expected_item}' == 'ImageDNA_1200x627'
    Download the image and check the size > 1  ${list_get_item}
    Set Test Variable   ${extracted_img_url}    ${list_get_item}
  ELSE IF  '${expected_item}' == 'video_link'
    Download the video and check the size > 1  ${list_get_item}[0]
    Set Test Variable   ${extracted_video_url}  ${list_get_item}[0]
  END

  # Set the ${extracted_url} for the following assertion
  IF  '${expected_item}' == 'url'
    Set Test Variable   ${extracted_url}    ${list_get_item}
  END

  # Set the ${extracted_sku} for the following assertion
  IF  '${expected_item}' == 'sku'
    Set Test Variable   ${extracted_sku}    ${list_get_item}
  END


Check the regex pattern in the given list
  [Arguments]             ${given_list}   ${regex}
  Log                     Start checking if the regex pattern: ${regex} exists in the given list: ${given_list}  level=DEBUG

  FOR  ${item}  IN  @{given_list}
    Should Match Regexp     ${item}         ${regex}    @{TEST TAGS} Can't get the regex:${regex} from list:${item}, Request:${request.url}, R_Header:${request.headers}
    ...                     values=false
  END


Download the image and check the size > 1
  [Arguments]         ${list_get_item}

  FOR  ${item}  IN  @{list_get_item}
    IF  '${item}' != 'novalue'
      Create Session      ImageSession            ${item}         disable_warnings=1
      ${image_resp} =     Get On Session          ImageSession    ${Empty}
      ${picture_info} =   Convert Binary To Image  ${image_resp.content}

      # picture_info = ['RGB', (1200, 627), 'JPEG', image_size]
      # log  ${picture_info}[1][0]  --> h
      # log  ${picture_info}[1][1]  --> W
      # log  ${picture_info}[3]
      Should Be True      ${picture_info}[3] > 1  @{TEST TAGS} The size of '${item}' is empty
    ELSE
      FAIL    @{TEST TAGS} We could not get the image since the path is missing
    END
  END


Download the video and check the size > 1
  [Arguments]             ${video_download_path}

  IF  '${video_download_path}' != '${Empty}'
    Create Session          VideoSession            ${video_download_path}  disable_warnings=1
    ${video_resp} =         Get On Session          VideoSession            ${Empty}
    Save Video              ${video_resp.content}   ${work_dir}/testsuite/output.mp4
    Wait Until Created      ${work_dir}/testsuite/output.mp4  timeout=15 seconds
    ${video_info} =         Get Video Info          ${work_dir}/testsuite/output.mp4

    # video_info = [1280, 720, 911, 233216]
    # log  ${video_info}[0]  --> width
    # log  ${video_info}[1]  --> height
    # log  ${video_info}[2]  --> num_frames
    # log  ${video_info}[3]  --> duration_ts
    # log  ${video_info}[4]  --> codec_long_name // H.265
    Should Be True          ${video_info}[3] > 1    @{TEST TAGS} The size of '${video_download_path}' is empty
    Should Match Regexp     ${video_info}[4]        (H.265)                 @{TEST TAGS} The codec_long_name isn't H265, what we get is ${video_info}[4]
    Log                     width: ${video_info}[0], height:${video_info}[1], size:${video_info}[3], codec_long_name:${video_info}[4]  level=DEBUG
  ELSE
    FAIL    @{TEST TAGS} We could not get the video since the path is missing
  END


# Vendor API testing utility keywords #
Auto select test dimensions
  [Arguments]         ${request_url}=${Empty}
  [Documentation]  Auto-select test dimensions from predefined sizes
  ...              Returns dimensions dictionary with width and height
  ...              Available test sizes: 300x300, 1200x627, 1200x600
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
  ${size_parts} =     Split String        ${selected_size}    x
  ${width} =          Set Variable        ${size_parts}[0]
  ${height} =         Set Variable        ${size_parts}[1]

  &{dimensions} =     Create Dictionary   width=${width}      height=${height}

  Log                 📏 Selected test dimensions: ${width}x${height} (from predefined test sizes)
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
  [Arguments]         ${tracking_url}
  [Documentation]  Parse YAML tracking_url to extract parameter configuration
  ...              Example: "{product_url}&param1={click_id_base64}" -> param_name=param1, uses_base64=true

  # Extract parameter name using regex
  ${param_matches} =  Get Regexp Matches  ${tracking_url}     [&?]([^=]+)=            1
  ${param_name} =     Set Variable If     ${param_matches}    ${param_matches[0]}     unknown

  # Check if uses base64 encoding
  ${uses_base64} =    Run Keyword And Return Status
  ...                 Should Contain      ${tracking_url}     base64

  # Check if requires group_id (typically for INL vendors)
  ${has_group_id} =   Set Variable If     '${param_name}' == 'subparam'  ${TRUE}  ${FALSE}

  # Create config dictionary
  &{config} =         Create Dictionary
  ...                 param_name=${param_name}
  ...                 uses_base64=${uses_base64}
  ...                 has_group_id=${has_group_id}

  RETURN              &{config}
