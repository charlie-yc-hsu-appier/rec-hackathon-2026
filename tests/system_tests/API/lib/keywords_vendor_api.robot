*** Keywords ***
# vendor API endpoint #
I have an vendor session

  Create Session  VendorSession    url=${HTTP_METHOD}://${VENDOR_HOST}  disable_warnings=1  retry_status_list=[500,502,503,504]


I would like to set the session under vendor endpoint with
  [Documentation]  Send GET request to vendor API endpoint.
  ...              Available parameters:
  ...              endpoint - The endpoint path (default: /healthz)
  ...              For example:
  ...              Given I would like to set the session under vendor endpoint with  endpoint=/healthz
  [Arguments]             &{args}

  # Set default endpoint if not provided
  ${endpoint_has_value} =  Run Keyword And Return Status
  ...                     Dictionary Should Contain Key  ${args}  endpoint
  IF  ${endpoint_has_value}
    Set Local Variable  ${endpoint}  ${args}[endpoint]
  ELSE
    Set Local Variable  ${endpoint}  /healthz
  END

  # Set the request header with Accept: */*
  &{HEADERS} =            Create Dictionary
  ...                     Accept=*/*

  # Start to send the GET request
  ${resp} =               Get On Session          VendorSession           url=${endpoint}     headers=&{HEADERS}

  # Validate response is not empty
  Should Not Be Empty
  ...                     ${resp.json()}
  ...                     @{TEST TAGS} FAIL: The response from vendor API is empty: ${resp.json()}, Request:${resp.request.url}, R_Header:${resp.request.headers}

  Set Test Variable       ${request.url}          ${resp.request.url}
  Set Test Variable       ${request.headers}      ${resp.request.headers}

  # Extract the request/response-header to Testrial & console
  Set Test Message        \n                      append=yes
  Set Test Message        *** Request URL: ${resp.url}  append=yes
  Set Test Message        \n                      append=yes
  Set Test Message        *** Response Header: ${resp.headers}  append=yes

  # Set the assertion variable for the following assertion variables: status_code
  Set Test Variable       ${status_code}          ${resp.status_code}
  ${resp.status_code} =   Convert To String       ${resp.status_code}
  Should Not Match Regexp
  ...                     ${resp.status_code}
  ...                     (^(4|5)..$)
  ...                     Something wrong with Vendor API \n request_url: ${resp.url} \n http_status_code: ${resp.status_code} \n message: ${resp.content} \n
  ...                     values=False

  # Set JSON response variable
  Set Test Variable       ${resp_json}            ${resp.json()}

  # Validate healthz response message
  IF  '${endpoint}' == '/healthz'
    Should Be Equal As Strings
    ...                 ${resp_json}[message]
    ...                 ok
    ...                 @{TEST TAGS} FAIL: Healthz endpoint should return message: 'ok', but got: ${resp_json}[message]
  END


# Assertion Keywords #
I would like to check status_code should be "${expected_code}" within the current session
  Should Be Equal As Strings
  ...     ${status_code}
  ...     ${expected_code}
  ...     @{TEST TAGS} The status code isn't we expected:${expected_code}, we get:${status_code}, Request:${request.url}, R_Header:${request.headers}
  ...     values=False


I would like to check content_length should be "${expected_content_length}" within the tracker session
  Should Be Equal As Strings
  ...     ${content_length}
  ...     ${expected_content_length}
  ...     @{TEST TAGS} The content_length code isn't we expected:${expected_content_length}, we get:${content_length}, Request:${request.url}, R_Header:${request.headers}
  ...     values=False


In the user2item response payload, the value of alg should not be "${alg_type}"
  [Arguments]     ${rmn_ignore_num}=${FALSE}

  Check the "${alg_type}" should not in each JSON value items  rmn_ignore_num=${rmn_ignore_num}


In the user2item response payload, the "${expected_item}" should exist (EC-REC Common Key)
  Check "${expected_item}" key in the JSON response


In the user2item response payload, the "${expected_item}" should exist (EC-REC Common Value)
  IF  '${expected_item}' == 'rec_req_id'
    Check the rec_req_id should be within the response
  ELSE IF  '${expected_item}' == 'preview'
    Check the value of preview should be "${True}" within the response
  ELSE IF  '${expected_item}' == 'land_append'
    Check the land_append should be within the response
  ELSE IF  '${expected_item}' == 'rmn_required keys'
    Check the rmn_required keys should be within the response
  ELSE
    Check "${expected_item}" value in each JSON response items  layout_id=${extracted_layout_id}
  END


In the spark log, we can find the "${expected_idfa}" & "${expected_rec_req_id}" information
  [Arguments]             &{args}
  [Documentation]  You can add the check_xxx after the keywords to verify related cases

  # 1st find the spark_log information via idfa & rec_req_id
  ${spark_log_json} =     Spark Log Consumer      ${expected_idfa}    ${expected_rec_req_id}
  Should Not Be Empty     ${spark_log_json}
  log  ${spark_log_json}

  ${need_to_check_bidobjid} =  Run Keyword And Return Status
  ...                     Dictionary Should Contain Key  ${args}  check_bidobjid

  ${need_to_check_ufo} =  Run Keyword And Return Status
  ...                     Dictionary Should Contain Key  ${args}  check_ufo

  ${need_to_check_chosen_mg_and_ranker} =  Run Keyword And Return Status
  ...                     Dictionary Should Contain Key  ${args}  check_matcher_ranker_results

  ${need_to_check_rank_recommend_order} =  Run Keyword And Return Status
  ...                     Dictionary Should Contain Key  ${args}  check_rank_recommend_order

  # 2nd verify other features we want
  IF  ${need_to_check_bidobjid}
    Set Local Variable  ${expected_bidobjid}    ${args}[check_bidobjid]
    Check "bidobjid" in spark JSON response  ${spark_log_json}  ${expected_bidobjid}
  END

  IF  ${need_to_check_ufo}
    Set Local Variable  ${expected_ufo}     ${args}[check_ufo]
    Check "ufos_user_features" in spark JSON response  ${spark_log_json}  ${expected_ufo}
  END

  IF  ${need_to_check_chosen_mg_and_ranker}
    Set Local Variable  ${expected_chosen_mg_and_ranker}  ${args}[check_matcher_ranker_results]
    Check "experiments[*]" in spark JSON response  ${spark_log_json}  ${expected_chosen_mg_and_ranker}
  END

  IF  ${need_to_check_rank_recommend_order}
    Dictionary Should Contain Key  ${args}  check_order_should_be  msg=When the "need_to_check_rank_recommend_order" is "TRUE", please provide "shuffled" or not after the "check_order_should_be"
    Set Local Variable  ${expected_rank_recommend_order}  ${args}[check_rank_recommend_order]
    Set Local Variable  ${check_order_should_be}      ${args}[check_order_should_be]
    Check "recommended_skus[*].product_id" in spark JSON response  ${spark_log_json}  ${expected_rank_recommend_order}  ${check_order_should_be}
  END


In the user2item response payload, the "${expected_item}" should not exist
  Check "${expected_item}" value doesn't exist in each JSON response items


In the user2item response payload, the amount of recommended products should be as same as the vaule of num_items="${expected_num_items}"
  Check the amount of recommended products should be equal to "${expected_num_items}"


In the user2item response payload, there is no duplicated SKU item in the response
  Check no duplicated SKU items in the response


In all the user2item response payloads, the order of each recommend proudcts should be "${expected_result}"
  Should Match Regexp     ${expected_result}      (different|same)    msg= [Fail] Please check the expected_result, should be: (different|same)

  FOR  ${index}  IN RANGE  3
    ${random_1} =       Generate Random String  1                   0123456789
    ${random_1} =       Convert To Integer      ${random_1}

    IF  '${random_1}' == '8'
      Set Local Variable  ${random_2}     0
    ELSE IF  '${random_1}' == '9'
      Set Local Variable  ${random_2}     1
    ELSE
      ${random_2} =   Evaluate    ${random_1} + 2
    END

    ${status}           ${msg} =                Run Keyword And Ignore Error
    ...                 Lists Should Be Equal   ${${random_1}_products_list}  ${${random_2}_products_list}

    Set Test Message    \n\r                    append=yes
    Set Test Message
    ...                 [Shuffle Check] The order of list (${random_1} & ${random_2}) ${${random_1}_products_list} & ${${random_2}_products_list} is "${expected_result}"
    ...                 append=yes

    IF  "${status}" == "PASS" and "${expected_result}" == "different"
      FAIL
      ...     The order of two recommend (${${random_1}_products_list}|${${random_2}_products_list}) result is same --> 1st:${${random_1}_products_list} 2nd:${${random_2}_products_list}
    ELSE IF  "${status}" == "FAIL" and "${expected_result}" == "same"
      FAIL
      ...     The order of two recommend (${${random_1}_products_list}|${${random_2}_products_list}) result is different --> 1st:${${random_1}_products_list} 2nd:${${random_2}_products_list}
    END
  END


In all the user2item response payloads, the diversity of each recommend proudcts should >= "${expected_num}"
  ${expected_num} =   Convert To Integer  ${expected_num}
  ${cnt} =            Get length          ${u_all_product_list}
  Should Be True
  ...                 ${cnt} >= ${expected_num}
  ...                 The diversity of products isn't bigger than we expect: ${expected_num}

  Set Test Message    \n\r                append=yes
  Set Test Message    [Shuffle Check] The diversity(${cnt}) of list: ${u_all_product_list}  append=yes
