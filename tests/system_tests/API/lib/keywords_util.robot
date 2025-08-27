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


Convert Response To String
  [Tags]              robot:flatten
  [Arguments]         ${resp.content}

  # Decode the byte to string
  # Set the assertion variable for the following assertion variables: resp_string
  ${string_resp_content} =  Decode Bytes To String  ${resp.content}  UTF-8
  Set Test Variable   ${resp_string}      ${string_resp_content}


Check "${expected_item}" key in the JSON response
  ${resp_json_value} =    Set Variable        ${resp_json}[${expected_item}]
  Should Not Be Empty     ${resp_json_value}  @{TEST TAGS} The "${expected_item}" is empty, Request:${request.url}, R_Header:${request.headers}
  Log                     Response of "${expected_item}" in the Json array: ${resp_json_value}  level=DEBUG


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


Check "${expected_item}" value doesn't exist in each JSON response items
  Variable Should Exist   ${resp_json['value']}
  ${resp_json_value} =    Set Variable            ${resp_json['value']}
  Should Not Be Empty     ${resp_json_value}      @{TEST TAGS} The value[] is empty, Request:${request.url}, R_Header:${request.headers}

  FOR  ${item}  IN  @{resp_json_value}
    Log                 ${item}                 level=DEBUG
    ${item_string} =    Convert JSON To String  ${item}
    ${match_list} =     Get Regexp Matches      ${item_string}  ${expected_item}
    Should Be Empty     ${match_list}           @{TEST TAGS} We're still able to see "${expected_item}" in response, Request:${request.url}, Response: ${item_string}
  END


Check the "${alg_type}" should not in each JSON value items
  [Arguments]             ${rmn_ignore_num}=${FALSE}
  Log                     *** [Shouldn't equal to "${alg_type}"] Response of alg in the value array:  level=DEBUG

  Variable Should Exist   ${resp_json['value']}
  ${resp_json_value} =    Set Variable            ${resp_json['value']}
  Should Not Be Empty     ${resp_json_value}      @{TEST TAGS} The algo_type in value[] is empty, Request:${request.url}, R_Header:${request.headers}

  ${all_alg_list} =       Get Value From Json     ${resp_json}            $.value[*].alg

  IF  '${alg_type}' == 'hotitem'
    # Make sure the 1st item should be recommended by matcher/ranker
    Should Match Regexp     ${all_alg_list[0]}  ranker              @{TEST TAGS} The first algo type (pos[0]) isn't a matcher/ranker type: ${resp_json['value'][0]}
    IF  '${rmn_ignore_num}' == '${FALSE}'
      ${reindeer_num} =   Get Match Count     ${all_alg_list}     *ranker*
      Should Be True      ${expected_reindeer_num} <= ${reindeer_num}
      ...                 msg=The expected num of recommended items from reindeer (${reindeer_num}) isn't equal to we expected: ${expected_reindeer_num}
    END
  ELSE IF  '${alg_type}' == 'below_predefined'
    Should Contain Match    ${all_alg_list}     *predefined*    msg=The algo type isn't a below_predefined type, Request:${request.url}, R_Header:${request.headers}
  ELSE
    Should Not Contain Match  ${all_alg_list}  *${alg_type}*  msg=The algo type is a ${alg_type} type, Request:${request.url}, R_Header:${request.headers}
  END

  Set Test Variable       ${extracted_all_alg_list}  ${all_alg_list}
  Log                     ${all_alg_list}         level=DEBUG


Check "${expected_item}" key in the Header response
  [Tags]                  robot:flatten
  [Arguments]             ${should_be_same_as}=${Empty}
  Log                     Checking "${expected_item}" in the Header response.  level=DEBUG

  log                     ${resp_headers}
  Variable Should Exist   ${resp_headers['${expected_item}']}
  ...                     No '${expected_item}' key within the response headers, Request:${request.url}, R_Header:${request.headers}
  Set Test Variable       ${resp_header_expected_item}  ${resp_headers['${expected_item}']}

  ${is_video} =           Get Regexp Matches  ${resp_header_expected_item}  mp4
  ${is_image} =           Get Regexp Matches  ${resp_header_expected_item}  jpeg

  IF  ${is_video}
    Download the video and check the size > 1  ${resp_header_expected_item}
  ELSE IF  ${is_image}
    @{img_list} =   Create List     ${resp_header_expected_item}
    Download the image and check the size > 1  ${img_list}
  END

  IF  '${should_be_same_as}' != '${Empty}'
    Should Match Regexp     ${resp_header_expected_item}  ${should_be_same_as}
  END


Check if there is partner_id in the params for IAB filter
  [Tags]              robot:flatten
  [Arguments]         ${args}
  # When the ${check_IAB_with_partner_id} isn't empty, we'll check the category of each sku in the keyword "Check "${expected_item}" in each JSON value items"
  # The function will set ${partner_id} as ${check_IAB_with_partner_id} for the following steps
  ${partner_id_has_value} =  Run Keyword And Return Status
  ...                 Dictionary Should Contain Key  ${args}  partner_id
  IF  ${partner_id_has_value}
    Set Test Variable   ${check_IAB_with_partner_id}  ${args}[partner_id]
  ELSE
    Set Test Variable   ${check_IAB_with_partner_id}  ${Empty}
  END


Check the regex pattern in the given list
  [Arguments]             ${given_list}   ${regex}
  Log                     Start checking if the regex pattern: ${regex} exists in the given list: ${given_list}  level=DEBUG

  FOR  ${item}  IN  @{given_list}
    Should Match Regexp     ${item}         ${regex}    @{TEST TAGS} Can't get the regex:${regex} from list:${item}, Request:${request.url}, R_Header:${request.headers}
    ...                     values=false
  END


Check the "${alg_type}" should be in debug information
  ${algo_in_debugInfo} =  Get Value From Json     ${resp_json}        $.value[*].alg
  Log                     *** Algo Distribution Test Result --> From_Tool: ${alg_type} | Debug_Endpoint: ${algo_in_debugInfo}[0]  level=DEBUG

  ${matched_item} =       Get Regexp Matches      ${algo_in_debugInfo}[0]  (${alg_type})
  ${is_match} =           Get Length              ${matched_item}
  Should Be Equal As Numbers  ${is_match}  ${1}  msg=The match item Algo isn't what we want
  IF  ${is_match} == 0
    FAIL    @{TEST TAGS} (Regex: (${alg_type})) Can't get the algo info: "${alg_type}" from REC response:${resp_json}, Request: ${request.url}
  END


The key should exist in the params dictionary
  [Tags]              robot:flatten
  [Arguments]         ${args}         ${key}
  ${has_key} =        Run Keyword And Return Status
  ...                 Dictionary Should Contain Key  ${args}  ${key}
  IF  ${has_key}
    ${key_length} =     Get Length  ${args}[${key}]
    IF  ${key_length} <= 0
      FAIL    "[Empty ${key}] @{TEST TAGS} Please check the value of ${key} in params", Request:${request.url}, R_Header:${request.headers}
    END
  ELSE
    FAIL    "[No ${key}] @{TEST TAGS} Please check the key ${key} in params", Request:${request.url}, R_Header:${request.headers}
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
  ${list_length} =    Get Length    ${test_sizes}
  ${random_index} =   Evaluate    __import__('random').randint(0, ${list_length}-1)
  ${selected_size} =  Set Variable    ${test_sizes}[${random_index}]
  
  # Parse the selected size
  ${size_parts} =     Split String    ${selected_size}    x
  ${width} =          Set Variable    ${size_parts}[0]
  ${height} =         Set Variable    ${size_parts}[1]

  &{dimensions} =     Create Dictionary   width=${width}      height=${height}
  
  Log    📏 Selected test dimensions: ${width}x${height} (from predefined test sizes)
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
