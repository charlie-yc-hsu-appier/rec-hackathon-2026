*** Keywords ***
# Utility #
Replacing the special characters for urlencode
  [Tags]              robot:flatten
  [Arguments]         ${string}
  ${temp_string} =    Replace String      ${string}       ,   %2C
  ${temp_string} =    Replace String      ${temp_string}  :   %3A
  ${temp_string} =    Replace String      ${temp_string}  /   %2F
  ${temp_string} =    Replace String      ${temp_string}  (   %28
  ${final_string} =   Replace String      ${temp_string}  )   %29
  Return              ${final_string}


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


Check "${expected_item}" in vast video JSON response
  Log                     Response of ${expected_item} in the vast JSON response.  level=DEBUG

  Variable Should Exist   ${resp_json['${expected_item}']}

  Log                     ${vast_video_item}  level=DEBUG

  ${is_video} =           Get Regexp Matches  ${vast_video_item}  mp4
  IF  ${is_video}
    Download the video and check the size > 1  ${vast_video_item}
  END


Check "${expected_item}" in spark JSON response
  [Arguments]             ${spark_json}           ${expected_object}  ${expected_rule}=${Empty}
  Log                     Start checking "${expected_item}" in the Spark JSON response.  level=DEBUG

  Should Have Value In Json  ${spark_json}  $.${expected_item}
  ${values} =             Get Value From Json     ${spark_json}       $.${expected_item}

  IF  '${expected_item}' == 'experiments[*]'
    ${group1} =         Evaluate        [e["group"] for e in ${values} if e["name"] == "matcher_group"]
    ${group2} =         Evaluate        [e["group"] for e in ${values} if e["name"] == "ranker"]
    ${merged_str} =     Evaluate        "reindeer_mg_%s_ranker_%s" % ('${group1}[0]', '${group2}[0]')
    ${values} =         Create List     ${merged_str}
  END

  Should Not Be Empty     ${values}
  ${object_type} =        Evaluate                type($expected_object).__name__

  # Compared the bidobjid from REC VS spark log (bidobjid)
  IF  "${object_type}" == "str" and "${expected_item}" == "bidobjid"
    Should Be Equal As Strings  ${values}[0]  ${expected_object}
  END

  # Compared the ufo list from REC VS spark log (ufos_user_features)
  IF  "${object_type}" == "DotDict" and "${expected_item}" == "ufos_user_features"
    Lists Should Be Equal   ${values}[0]    ${expected_object}
  END

  # Compared the sku list from REC VS spark log (recommended_skus[*].product_id) and the sku order in rank_result ($.rank_result.products[*].product_id)
  IF  "${object_type}" == "list" and "${expected_item}" == "recommended_skus[*].product_id"
    Lists Should Be Equal   ${values}           ${expected_object}
    ${rank_result_sku_list} =  Get Value From Json  ${spark_json}  $.rank_result.products[*].product_id

    # Based on the passing value ${expected_rule} to decide the order should be shuffled or not
    ${order_result} =       Check two lists are subsequence  ${rank_result_sku_list}  ${values}
    IF  "${expected_rule}" == "shuffled"
      Should Not Be True  ${order_result}     The lists match in order, REC: ${values}, Rank_Result: ${rank_result_sku_list}
    ELSE
      Should Be True  ${order_result}     The lists do not match in order, REC: ${values}, Rank_Result: ${rank_result_sku_list}
    END

    Log                     [Multi-Model Spark Log Info Assertion] Response of ${expected_item} in the spark JSON rank_result response (${expected_rule}): ${rank_result_sku_list}  level=DEBUG
  END

  # Compared the algo list from REC VS spark log (check_matcher_ranker_results)
  IF  "${object_type}" == "list" and "${expected_item}" == "experiments[*]"
    Should Contain Match    ${expected_object}  *${values}[0]*
  END

  Log                     [Spark Log Info Assertion] Response of ${expected_item} in the spark JSON response: ${values}  level=DEBUG


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


Check two lists are subsequence
  [Arguments]             ${long_list}        ${short_list}
  [Documentation]  This ensures that the index does not exceed the range of short_list while correctly determining
  ...              whether short_list is a subsequence of long_list.

  ${index} =              Set Variable        0
  ${index} =              Convert To Integer  ${index}
  ${short_list_len} =     Evaluate            len(${short_list})

  FOR  ${item}  IN  @{long_list}
    IF  ${index} < ${short_list_len}
      ${short_item} =     Set Variable    ${short_list}[${index}]
      IF  '${item}' == '${short_item}'
        ${index} =  Evaluate    ${index} + 1
      END
    END
  END

  ${result} =             Evaluate            ${index} == ${short_list_len}
  RETURN                  ${result}


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


Using the given random bidobjid or generating a new one
  [Tags]              robot:flatten
  [Arguments]         ${args}
  # Since the bidobjid will be cached for 4 hours in the DNA service, we give a random value to avoid getting the same item
  ${has_bidobjid} =   Run Keyword And Return Status
  ...                 Dictionary Should Contain Key  ${args}  bidobjid
  IF  ${has_bidobjid}
    Set Test Variable   ${bidobjid}             ${args}[bidobjid]
  ELSE
    ${bidobjid} =       Generate Random String  15              [NUMBERS]abcdefghi
    Set Test Variable   ${bidobjid}             ${bidobjid}
  END
  RETURN              ${bidobjid}


Extract the required params for predict item
  [Arguments]             ${resp_json}

  # If the resp_json is from redis, then it should contain the key --> "p"
  ${is_from_redis} =      Run Keyword And Return Status
  ...                     Dictionary Should Contain Key  ${resp_json}  p

  IF  ${is_from_redis}
    ${sku} =                Get Value From Json     ${resp_json}    $.p
    Set Test Variable       ${extracted_sku}        ${sku[0]}

    ${custom_dna_layout} =  Get Value From Json     ${resp_json}    $.i
    Set Test Variable       ${extracted_custom_dna_layout}  ${custom_dna_layout[0]}

    ${predictid} =          Get Value From Json     ${resp_json}    $.id
    Set Test Variable       ${extracted_predictid}  ${predictid}
    Set Test Variable       ${extracted_predict_default}  ${0}
  ELSE
    ${sku} =                Get Value From Json     ${resp_json}    $.value[0].sku
    Set Test Variable       ${extracted_sku}        ${sku[0]}

    ${custom_dna_layout} =  Get Value From Json     ${resp_json}    $.value[0].custom_dna_layout
    ${custom_dna_layout} =  Get Regexp Matches      ${custom_dna_layout[0]}  v0\/(.+)  1
    Set Test Variable       ${extracted_custom_dna_layout}  ${custom_dna_layout[0]}

    ${predictid} =          Generate Random String  7               [NUMBERS]poqstuvwxyz
    Set Test Variable       ${extracted_predictid}  ${predictid}
    Set Test Variable       ${extracted_predict_default}  ${0}
  END

  Log                     sku: ${extracted_sku}, img: ${extracted_custom_dna_layout}, predict_id: ${extracted_predictid}, predict_default: ${extracted_predict_default}  level=DEBUG


Extract the required assertion data for multiple_URL
  [Tags]                  robot:flatten
  Variable Should Exist   ${resp_json['value']}
  ${resp_json_value} =    Set Variable            ${resp_json['value']}
  Should Not Be Empty     ${resp_json_value}      @{TEST TAGS} The value[] is empty

  ${dict_get_item} =      Create Dictionary

  FOR  ${item}  IN  @{resp_json_value}
    Log                 ${item}             level=DEBUG
    Dictionary Should Contain Key  ${item}  sku  @{TEST TAGS} No "sku" Key within the response, Request:${request.url}, R_Header:${request.headers}
    Dictionary Should Contain Key  ${item}  url  @{TEST TAGS} No "url" Key within the response, Request:${request.url}, R_Header:${request.headers}

    Set Local Variable  ${get_sku}          ${item['sku']}
    Set Local Variable  ${get_url}          ${item['url']}

    Set To Dictionary   ${dict_get_item}    ${get_sku}=${get_url}
  END

  Set Test Variable       ${mul_url_dict}         ${dict_get_item}


Get Current Epoch Time
  [Tags]              robot:flatten
  ${date} =           Get Current Date    exclude_millis=True
  ${epoch_date} =     Convert Date        ${date}                 epoch
  ${epoch_int} =      Convert To Integer  ${epoch_date}
  RETURN              ${epoch_int}


Generating all user ids
  ${idfa} =   Get A UUID In Upper
  ${auid} =   Generate Random String  12

  RETURN      ${idfa}                 ${auid}


Get A UUID In Upper
  [Tags]      robot:flatten

  ${uuid} =   Evaluate        uuid.uuid4()        modules=uuid
  ${UUID} =   Evaluate        "${uuid}".upper()
  RETURN      ${UUID}


Get A Predict Item From Redis By AUID OR IDFA
  [Arguments]         ${Option}

  ${UUID} =           Get A UUID In Upper
  # Set Test Variable  ${UUID}  ${UUID}
  Set Test Variable   ${extracted_UUID}       ${UUID}

  IF  '${Option}' == 'IDFA'
    # Insert params: ${kafka_host}  ${kafka_bootstarp}  ${kafka_topic}  ${site_id}=${Empty}  ${idfa}=${Empty}  ${auid}=${Empty}
    Kafka Insert Data   ${KAFKA_PREDICT_HOST}   ${KAFKA_PREDICT_BOOTSTRAP}  ${KAFKA_PREDICT_TOPIC}  android--com.coupang.mobile_s2s_v3  ${UUID}
  ELSE
    # There has been a pre-setting mapping in the fuzzy-user-redis-asia-east1-stg already (Due to the ACL rule)
    # AUID -->  IDFA:  T15TAEQ41TET -->  F9E980E0-470F-4A49-9C30-192FEF23CFAF
    Kafka Insert Data   ${KAFKA_PREDICT_HOST}   ${KAFKA_PREDICT_BOOTSTRAP}  ${KAFKA_PREDICT_TOPIC}  android--com.coupang.mobile_s2s_v3
    ...                 F9E980E0-470F-4A49-9C30-192FEF23CFAF
  END

  Redis Init
  IF  '${Option}' == 'IDFA'
    ${redis_vaule} =    Wait Until Keyword Succeeds  2 min  strict: 5s
    ...                 Redis Get Value By Key  bid:android--com.coupang.mobile_s2s_v3:${UUID}:0:r
  ELSE
    ${redis_vaule} =    Wait Until Keyword Succeeds  2 min  strict: 5s
    ...                 Redis Get Value By Key  bid:android--com.coupang.mobile_s2s_v3:T15TAEQ41TET:0:r
  END
  ${redis_vaule} =    Convert String to JSON  ${redis_vaule}

  IF  '${Option}' == 'IDFA'
    Log     The IDFA of the predicted item in the Kafka is : ${UUID}  level=DEBUG
  ELSE
    Log     The AUID of the predicted item in the Kafka is : F9E980E0-470F-4A49-9C30-192FEF23CFAF  level=DEBUG
  END
  Log                 The value of the predicted item in the Redis is: ${redis_vaule}  level=DEBUG
  RETURN              ${redis_vaule}


Get The Algo Groups By UUID
  [Arguments]         &{args}
  [Documentation]  The Usage of ./bin/print_exp_group.app:
  ...              -c string  Reyka config file
  ...              -cid string  Ad group id
  ...              -exp string  Experiment name
  ...              -group string  Target group of the experiment
  ...              -layout_id string  Layout id
  ...              -oid string  Campaign id
  ...              -partner_id string  Partner id
  ...              -sid string  Site id
  ...              -u string  UUID/IDFA

  # Check if the site_id,UUID or oid is passing by
  The key should exist in the params dictionary  ${args}  site_id
  The key should exist in the params dictionary  ${args}  UUID
  The key should exist in the params dictionary  ${args}  oid

  # More options, please refer to ./bin/print_exp_group.app --help
  ${result} =         Run Process             ./bin/print_exp_group.app  -c  ./deploy/rec-reyka/secrets/config.yaml  -sid  ${args}[site_id]  -u  ${args}[UUID]  -oid
  ...                 ${args}[oid]
  Log Many            stdout: ${result.stdout}  stderr: ${result.stderr}
  ${algo_json} =      Convert String To Json  ${result.stdout}

  # Using the oid to get the custom_inl info, the coupang_vendor_oid is listed in the init.robot
  IF  '${args}[oid]' in @{coupang_vendor_oid}
    ${algo_group} =     Get Value From Json     ${algo_json}    $.exp.algorithm.group
  ELSE
    ${algo_group} =     Get Value From Json     ${algo_json}    $.exp.matcher_group.group
  END

  Log                 The "${args}[UUID]" is selected by "${algo_group}[0]"  level=DEBUG
  RETURN              ${algo_group}[0]


Get The UUID By Exp
  [Arguments]             &{args}
  [Documentation]  The Usage of ./bin/print_exp_group.app:
  ...              -c string  Reyka config file
  ...              -cid string  Ad group id
  ...              -exp string  Experiment name
  ...              -group string  Target group of the experiment
  ...              -layout_id string  Layout id
  ...              -oid string  Campaign id
  ...              -partner_id string  Partner id
  ...              -sid string  Site id
  ...              -u string  UUID/IDFA
  ...              You could apply the extra_exp & extra_group as the 2nd filter if you need.

  # Check if the site_id,UUID or oid is passing by
  The key should exist in the params dictionary  ${args}  site_id
  The key should exist in the params dictionary  ${args}  exp
  The key should exist in the params dictionary  ${args}  group
  ${has_extra_exp} =      Run Keyword And Return Status
  ...                     Dictionary Should Contain Key  ${args}  extra_exp
  ${has_extra_group} =    Run Keyword And Return Status
  ...                     Dictionary Should Contain Key  ${args}  extra_group

  IF  ${has_extra_group} and ${has_extra_group}
    FOR  ${i}  IN RANGE  50
      ${result} =         Run Process             ./bin/print_exp_group.app  -c  ./deploy/rec-reyka/secrets/config.yaml  -sid  ${args}[site_id]  -exp  ${args}[exp]  -group
      ...                 ${args}[group]
      TRY
        ${algo_json} =      Convert String To Json  ${result.stdout}
        ${exped_idfa} =     Get Value From Json     ${algo_json}        $.idfa
        ${exp_info} =       Get Value From Json     ${algo_json}        $.exp.${args}[exp]
        ${current_extra_group} =  Get Value From Json  ${algo_json}  $.exp.${args}[extra_exp].group

        IF  '${current_extra_group}[0]' == '${args}[extra_group]'
          Log     [Multiple Filters Exp Setting] Under ${args}[exp], we'd like to check $.exp.${args}[extra_exp] setting: ${current_extra_group} should be same as we expect: ${args}[extra_group]  level=DEBUG
          BREAK
        END
      EXCEPT
        FAIL    🛑 Fail to get the IDFA via print_exp_group tool after 50 tries...
      END
    END
  ELSE
    # More options, please refer to ./bin/print_exp_group.app --help
    ${result} =         Run Process             ./bin/print_exp_group.app  -c  ./deploy/rec-reyka/secrets/config.yaml  -sid  ${args}[site_id]  -exp  ${args}[exp]  -group
    ...                 ${args}[group]
    Log Many            stdout: ${result.stdout}  stderr: ${result.stderr}
    ${algo_json} =      Convert String To Json  ${result.stdout}
    ${exped_idfa} =     Get Value From Json     ${algo_json}        $.idfa
    ${exp_info} =       Get Value From Json     ${algo_json}        $.exp.${args}[exp]
  END

  Log                     The UUID "${exped_idfa}[0]" is under the config setting: "${algo_json}"  level=DEBUG
  RETURN                  ${exped_idfa}[0]


Check The Algo Distribution Is 100%
  [Arguments]     ${json_obj}

  ${all_prob} =   Get Value From Json     ${json_obj}     $.[*].prob
  ${sum_prob} =   Evaluate                sum(map(float, ${all_prob}))
  Should Be Equal As Numbers  ${sum_prob}  1.0
  ...             msg=@{TEST TAGS} The sum of the algo distribution setting isn't "1.0", Request:${request.url}, R_Header:${request.headers}
  RETURN          ${json_obj}


Check the amount of recommended products should be equal to "${expected_num_items}"
  ${expected_num_items} =  Convert To Integer  ${expected_num_items}
  ${sku_list} =       Get Value From Json     ${resp_json}    $.value[*].sku
  ${cnt_skus} =       Get length              ${sku_list}
  Should Be Equal     ${expected_num_items}   ${cnt_skus}
  ...                 msg=The amount of sku from REC isn't we expected (From REC: ${cnt_skus}, We expected: ${expected_num_items})  values=False


Check the rec_req_id should be within the response
  ${rec_req_id_list} =    Get Value From Json     ${resp_json}    $.rec_req_id
  Should Not Be Empty     ${rec_req_id_list}      We couldn't get value of rec_req_id from the rec response: ${resp_json}
  Set Test Variable       ${extracted_rec_req_id}  ${rec_req_id_list}[0]
  Log                     Value of "rec_req_id" in the Json response: ${rec_req_id_list}[0]  level=DEBUG


Check the land_append should be within the response
  ${land_append_list} =   Get Value From Json     ${resp_json}    $.land_append
  Should Not Be Empty     ${land_append_list}     We couldn't get value of land_append from the rec response: ${resp_json}
  Should Match Regexp     ${land_append_list}[0]  (<% this.item.pos %>~<% this.item.sku %>)
  ...                     The fix string(dynlnkg) of land_append isn't we expected
  Log                     Value of "land_append" in the Json response: ${land_append_list}[0]  level=DEBUG


Check the rmn_required keys should be within the response
  ${product_set} =        Get Value From Json     ${resp_json}    $.value[*].[product_id,score]
  Should Not Be Empty     ${product_set}          We couldn't get value of rmn_required keys from the rec_rmn response: ${resp_json}
  Log                     Value of "product_id & score" in the Json response(product_id,score): ${product_set}  level=DEBUG


Check the value of preview should be "${expected_status}" within the response
  ${preview_list} =   Get Value From Json     ${mimic_resp_json}  $.preview
  Should Be Equal As Strings  ${preview_list}[0]  ${expected_status}
  ...                 The value of preview from the rec response (${preview_list}[0]) isn't we expect('${expected_status}')
  Log                 Value of "preview" in the Json response: ${preview_list}[0]  level=DEBUG


Check no duplicated SKU items in the response
  ${sku_list} =   Get Value From Json     ${resp_json}    $.value[*].sku
  List Should Not Contain Duplicates  ${sku_list}  msg=@{TEST TAGS} Duplicated recommended items(sku) in reyka response: ${sku_list}

