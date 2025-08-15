*** Keywords ***
# Config API #
I have a config_api session
  [Tags]          robot:flatten

  # Basic Auth for Config API
  IF  '${SERVER_ENV}' == 'stag' or '${SERVER_ENV}' == 'dev'
    &{HEADERS} =    Create Dictionary       Content-Type=application/json  Authorization=Basic cmVjLXFhOmlsb3JmZ3Bva3J3aGl1bWtwb2FzZGZva2hwcWV3eWhi
  ELSE IF  '${SERVER_ENV}' == 'prod'
    &{HEADERS} =    Create Dictionary       Content-Type=application/json  Authorization=Basic cmVjLXFhOkZ6VUVyVGRSTTI2d2RoMlIyMk5lbktmeXlya2t2VGhR
  ELSE
    Fail    No SERVER_ENV params
  END

  Create Session  ConfigAPISession        url=https://${CONFIG_API_HOST}  headers=&{HEADERS}  disable_warnings=1  retry_status_list=[500,502,503,504]  timeout=5


I would like to check the status when the datafeed_id="${df}"
  [Tags]                  robot:flatten

  ${resp} =               Get On Session          ConfigAPISession        url=/v0/datafeeds/${df}
  Should Not Be Empty  ${resp.json()}  [${df}] @{TEST TAGS} FAIL: Can't get the response (datafeeds) via the ConfigAPI response: ${resp.json()} , Request: ${resp.request.url}. please check the 'datafeeds' setting
  Set Test Variable       ${request.url}          ${resp.request.url}
  Set Test Variable       ${request.headers}      ${resp.request.headers}

  ${df_status} =          Get Value From Json     ${resp.json()}          $.status
  RETURN                  ${df_status}


I would like to get the datafeed_id when site_id="${sid}"
  # Start to post the request
  ${resp} =               Get On Session          ConfigAPISession        url=/v0/sites/${sid}

  # Extra the request/response-header to Testrial & console
  Set Test Message        \n                      append=yes
  Set Test Message        *** Config API Request URL: ${resp.url}  append=yes

  # Set the assertion variable for the following assertion variables: status_code
  Set Test Variable       ${status_code}          ${resp.status_code}
  ${resp.status_code} =   Convert To String       ${resp.status_code}
  Should Not Match Regexp  ${resp.status_code}  (^(4|5)..$)  [${sid}] Something wrong with Config API \n reqest_url: ${resp.request.url} \n http_status_code: ${resp.status_code} \n message: ${resp.content} \n
  ...                     values=False

  Should Not Be Empty     ${resp.json()}          [${sid}] @{TEST TAGS} FAIL: Can't get the response (sites) via the ConfigAPI response: ${resp.json()}, Request: ${resp.request.url}. please check the 'sites' setting
  Set Test Variable       ${request.url}          ${resp.request.url}
  Set Test Variable       ${request.headers}      ${resp.request.headers}
  @{extracted_df_list} =  Get Value From Json     ${resp.json()}          $.site.datafeeds[*]

  # Set Test Variable  ${extracted_df}  ${extracted_df}
  FOR  ${df}  IN  @{extracted_df_list}
    ${df_status} =      I would like to check the status when the datafeed_id="${df}"
    IF  "${df_status[0]}" == "active"
      Set Test Variable   ${extracted_df}     ${df}
      IF  "${df_status[0]}" == "active"
        BREAK
      END
    END
  END

  IF  '${extracted_df}' == '${Empty}'
    Fail    No datafeed_id is active in this list: @{extracted_df_list}
  END

  Set Test Message        \n                      append=yes
  Set Test Message        The extracted values from Config API are: \n  append=yes
  Set Test Message        extracted_df: ${extracted_df}  append=yes
  Set Test Message        \n                      append=yes


I would like to get a list of oid when datafeed_id="${datafeed_id}" and status="${status}"
  [Arguments]             ${exclude_sa_group}=${Empty}

  @{available_status} =   Create List
  Append To List          ${available_status}     New                     Running     Paused  Stopped     Finished
  List Should Contain Value  ${available_status}  ${status}  The status: "${status}" doesn't exist in the ${available_status}(case-sensitive)

  # Start to post the request
  ${resp} =               Get On Session          ConfigAPISession        url=/v0/campaigns?status_code=${status}

  # Extra the request/response-header to Testrial & console
  Set Test Message        \n\r                      append=yes
  Set Test Message        *** Config API Request URL(Get oid based on "${status}" ${datafeed_id}):  append=yes
  Set Test Message        ${resp.url}             append=yes

  # Set the assertion variable for the following assertion variables: status_code
  Set Test Variable       ${status_code}          ${resp.status_code}
  ${resp.status_code} =   Convert To String       ${resp.status_code}
  Should Not Match Regexp  ${resp.status_code}  (^(4|5)..$)  [Get oid] Something wrong with Config API \n reqest_url: ${resp.url} \n http_status_code: ${resp.status_code} \n message: ${resp.content} \n  values=False

  @{extracted_oid_list} =  Get Value From Json  ${resp.json()}  $.campaigns[?(@.datafeed_id == '${datafeed_id}' & status_code == '${status}')].campaign_id

  IF  '${exclude_sa_group}' != '${Empty}'
    # Exclude those SA group oids
    Remove Values From List  ${extracted_oid_list}  4QSP2IQMQbahY5ipS88nnQ  RghYkBdSRuyGkDcvIE6Eqg  ybZF8EpZQ86G76I2LnsMnA
    Set Test Variable   ${extracted_running_oid_from_configapi}  ${extracted_oid_list}
  ELSE
    Set Test Variable   ${extracted_running_oid_from_configapi}  ${extracted_oid_list}
  END

  Set Test Message        \n                      append=yes
  Set Test Message        The extracted values from Config API are: \n  append=yes
  Set Test Message        extracted_running_oid_from_configapi: ${extracted_running_oid_from_configapi}  append=yes
  Set Test Message        \n                      append=yes


I would like to get the algorithm_seed when site_id="${sid}"
  # Start to post the request
  ${resp} =               Get On Session          ConfigAPISession        url=/v0/sites/${sid}

  # Extra the request/response-header to Testrial & console
  Set Test Message        \n                      append=yes
  Set Test Message        *** [Get algorithm_seed] Config API Request URL: ${resp.url}  append=yes

  # Set the assertion variable for the following assertion variables: status_code
  Set Test Variable       ${status_code}          ${resp.status_code}
  ${resp.status_code} =   Convert To String       ${resp.status_code}
  Should Not Match Regexp  ${resp.status_code}  (^(4|5)..$)
  ...                     [${sid}] [Get algorithm_seed] Something wrong with Config API \n reqest_url: ${resp.request.url} \n http_status_code: ${resp.status_code} \n message: ${resp.content} \n  values=False

  Should Not Be Empty     ${resp.json()}          [${sid}] @{TEST TAGS} FAIL: Can't get the response (sites) via the ConfigAPI response: ${resp.json()}, Request: ${resp.request.url}. please check the 'sites' setting
  ${algorithm_seed} =     Get Value From Json     ${resp.json()}          $.site.algorithm_seed

  Set Test Message        \n                      append=yes
  Set Test Message        The extracted algorithm_seed from Config API are: ${algorithm_seed}[0]  append=yes
  Set Test Message        \n                      append=yes

  RETURN                  ${algorithm_seed}[0]


I would like to update the metadata when site_id="${sid}" and status="${status}"
  &{json_meta_object} =   Create Dictionary   datafeed_id=${sid}  status=${status}
  ${resp} =               Put On Session      ConfigAPISession    url=/v0/datafeeds/${sid}  json=${json_meta_object}  expected_status=200


I would like to get the not ${expected_status} & keeta_api enabled oid list when datafeed_id=${datafeed_id}
  # Start to post the request
  ${all_status} =  Create List  Running  Paused  Stopped  Finished
  Remove Values From List  ${all_status}  ${expected_status}
  ${query_string} =    Set Variable    status_code=${all_status[0]}
  
  ${length}=    Get Length    ${all_status}
  FOR    ${item}    IN RANGE    1    ${length}
    ${query_string} =    Catenate    SEPARATOR=&    ${query_string}    status_code=${all_status[${item}]}
  END

  ${resp} =               Get On Session          ConfigAPISession        url=/v0/campaigns?${query_string}

  # Extra the request/response-header to Testrial & console
  Set Test Message        \n                      append=yes
  Set Test Message        *** Config API Request URL(Get oid based on not "${expected_status}" ${datafeed_id}):  append=yes
  Set Test Message        ${resp.url}             append=yes

  # Set the assertion variable for the following assertion variables: status_code
  Set Test Variable       ${status_code}          ${resp.status_code}
  ${resp.status_code} =   Convert To String       ${resp.status_code}
  Should Not Match Regexp  ${resp.status_code}  (^(4|5)..$)  [Get oid] Something wrong with Config API \n reqest_url: ${resp.url} \n http_status_code: ${resp.status_code} \n message: ${resp.content} \n  values=False

  @{extracted_keeta_oid_list} =  Get Value From Json    ${resp.json()}    $.campaigns[?(@.datafeed_id == "${datafeed_id}" & @.status_code != "${expected_status}" & @.configs.keeta_enable_rta_api == true)].campaign_id

  Set Test Message        \n                      append=yes
  Set Test Message        The extracted values from Config API are: \n  append=yes
  Set Test Message        extracted_keeta_oid_list: ${extracted_keeta_oid_list}  append=yes
  Set Test Message        \n                      append=yes

  RETURN  ${extracted_keeta_oid_list}


# Util Keywords #
Get an actived datafeed id from sid
  [Arguments]     ${sid}
  I have a config_api session
  I would like to get the datafeed_id when site_id="${sid}"


Get an actived psid or use a given one
  [Arguments]             ${given_psid}=${Empty}
  IF  '${given_psid}' == '${Empty}'
    ${an_actived_psid} =    I would like to get an actived datafeed label when the  df=${extracted_df}  label_type=rtb-product-set
    Set Test Variable       ${extracted_an_actived_psid}  ${an_actived_psid}
  END
  Set Test Variable       ${extracted_an_actived_psid}  ${given_psid}


# Multiple URL #
I would like to get the multiple url setting when
  [Arguments]             &{args}
  # Required: datafeed_id, ad_solution, promotion_content, Optional: company_id
  The key should exist in the params dictionary  ${args}  datafeed_id
  The key should exist in the params dictionary  ${args}  ad_solution
  The key should exist in the params dictionary  ${args}  promotion_content

  ${has_company_id} =     Run Keyword And Return Status
  ...                     Dictionary Should Contain Key  ${args}  company_id
  IF  ${has_company_id}
    Set Test Variable   ${company_id}   ${args}[company_id]
  ELSE
    Set Test Variable   ${company_id}   ${Empty}
  END

  # Change the promotion_content to ios
  ${is_ios} =             Get Regexp Matches      ${args}[datafeed_id]    ios
  IF  ${is_ios}
    Set Test Variable   ${promotion_content}    ios
  END

  # Change the promotion_content to web
  @{is_web_list} =        Create List
  Append To List          ${is_web_list}          android--kr.co.ssg      ssg.com
  ${is_web_promotion} =   Run Keyword And Return Status
  ...                     List Should Contain Value  ${is_web_list}  ${args}[datafeed_id]
  IF  ${is_web_promotion}
    Set Test Variable   ${promotion_content}    web
  END

  # Start to post the request
  ${resp} =               Get On Session          ConfigAPISession        url=/v0/multiple-product-urls/datafeeds/${args}[datafeed_id]

  # Extra the request/response-header to Testrial & console
  Set Test Message        \n\r                      append=yes
  Set Test Message        *** Config API Request URL(Get multipleURL based on ${args}[datafeed_id]): ${resp.url}  append=yes

  # Set the assertion variable for the following assertion variables: status_code
  Set Test Variable       ${status_code}          ${resp.status_code}
  ${resp.status_code} =   Convert To String       ${resp.status_code}
  Should Not Match Regexp  ${resp.status_code}  (^(4|5)..$)  [Get multipleURL] Something wrong with Config API \n reqest_url: ${resp.url} \n http_status_code: ${resp.status_code} \n message: ${resp.content} \n
  ...                     values=False

  Should Not Be Empty  ${resp.json()}
  ...                  [${args}[datafeed_id]] @{TEST TAGS} FAIL: Can't get the response (multiple-product-urls) via the ConfigAPI response: ${resp.json()}, Request: ${resp.url}. please check the 'multiple-product-urls' setting
  Set Test Variable       ${request.url}          ${resp.request.url}
  Set Test Variable       ${request.headers}      ${resp.request.headers}

  Get the extracted variable for the multiple URL needs from config_api  ${resp.json()}  ${args}[datafeed_id]  ${args}[ad_solution]  ${args}[promotion_content]  ${company_id}


Get the extracted variable for the multiple URL needs from config_api
  [Arguments]             ${json_obj}             ${datafeed_id}  ${ad_solution}  ${promotion_content}    ${company_id}

  IF  '${company_id}' != '${Empty}'
    # For testing the complex macro setting
    Set Local Variable      ${specific_company_id}  ${company_id}

    ${company_id} =         Get Value From Json     ${json_obj}     $.multiple_product_urls[?(company_id == '${specific_company_id}' & promotion_content == '${promotion_content}')].company_id
    Should Not Be Empty  ${company_id}
    ...                  [${datafeed_id}] @{TEST TAGS} FAIL: Can't get the mock data (specific_company_id) via the ConfigAPI response: ${specific_company_id}, Request: ${request.url}, please check the insert_mock script
    Set Test Variable       ${extracted_company_id_from_configapi}  ${company_id}

    ${url} =                Get Value From Json     ${json_obj}     $.multiple_product_urls[?(company_id == '${specific_company_id}' & ad_solution == '${ad_solution}' & promotion_content == '${promotion_content}')].url
    Set Test Variable       ${extracted_url_from_configapi}  ${url}

    ${mobile_url} =         Get Value From Json     ${json_obj}     $.multiple_product_urls[?(company_id == '${specific_company_id}' & ad_solution == '${ad_solution}' & promotion_content == '${promotion_content}')].mobile_url
    Set Test Variable       ${extracted_mobile_url_from_configapi}  ${mobile_url}

    ${id} =  Get Value From Json  ${json_obj}
    ...      $.multiple_product_urls[?(ad_solution == '${ad_solution}' & promotion_content == '${promotion_content}' & company_id == '${extracted_company_id_from_configapi[0]}')].campaigns[?(@.status_code != 'Finished')].id
    Set Test Variable       ${extracted_running_oid_from_configapi}  ${id}
  ELSE
    ${company_id_with_url_contain_http} =  Get Value From Json  ${json_obj}  $.multiple_product_urls[?(url =~ 'http.*' & ad_solution == '${ad_solution}' & promotion_content == '${promotion_content}')].company_id
    ${company_id} =         Get Value From Json     ${json_obj}     $.multiple_product_urls[?(ad_solution == '${ad_solution}' & promotion_content == '${promotion_content}')].company_id
    Should Not Be Empty  ${company_id}
    ...                  [${datafeed_id}] @{TEST TAGS} FAIL: Can't get the (company_id) when promotion_content is '${promotion_content}' & ad_solution is '${ad_solution}' via the ConfigAPI response: ${company_id}, Request: ${request.url}, please check the Config Setting.

    IF  ${company_id_with_url_contain_http}
      Remove Values From List  ${company_id}  ${company_id_with_url_contain_http[0]}
    END

    Set Test Variable       ${extracted_company_id_from_configapi}  ${company_id}

    ${url} =  Get Value From Json  ${json_obj}  $.multiple_product_urls[?(ad_solution == '${ad_solution}' & promotion_content == '${promotion_content}') & company_id == '${extracted_company_id_from_configapi[0]}'].url
    Set Test Variable       ${extracted_url_from_configapi}  ${url}

    ${mobile_url} =         Get Value From Json     ${json_obj}
    ...                     $.multiple_product_urls[?(ad_solution == '${ad_solution}' & promotion_content == '${promotion_content}' & company_id == '${extracted_company_id_from_configapi[0]}')].mobile_url
    Set Test Variable       ${extracted_mobile_url_from_configapi}  ${mobile_url}

    ${id} =  Get Value From Json  ${json_obj}
    ...      $.multiple_product_urls[?(ad_solution == '${ad_solution}' & promotion_content == '${promotion_content}' & company_id == '${extracted_company_id_from_configapi[0]}')].campaigns[?(@.status_code != 'Finished')].id
    Should Not Be Empty  ${id}
    ...                  [${datafeed_id}] @{TEST TAGS} FAIL: Can't get the running oid (!= Finished) when company_id is '${extracted_company_id_from_configapi[0]}' for the following testing, please check OID status in the iDash
    Set Test Variable       ${extracted_running_oid_from_configapi}  ${id}
  END

  Set Test Message        \n                      append=yes
  Set Test Message        The extracted values from Config API are: \n  append=yes
  Set Test Message
  ...                     Under extracted_company_id_from_configapi: "${extracted_company_id_from_configapi}", we got extracted_url_from_configapi: ${extracted_url_from_configapi}, extracted_mobile_url_from_configapi: ${extracted_mobile_url_from_configapi}, extracted_running_oid_from_configapi: ${extracted_running_oid_from_configapi}
  ...                     append=yes
  Set Test Message        \n                      append=yes


# Get an active label (psid)#
I would like to get an actived datafeed label when the
  [Arguments]             &{args}

  # Check if the df is passing by
  The key should exist in the params dictionary  ${args}  df
  Set Local Variable      ${df}                   ${args}[df]

  # Check if the label_type is passing by
  The key should exist in the params dictionary  ${args}  label_type
  Set Local Variable      ${label_type}           ${args}[label_type]

  # Fix-Me: AI-22711, We should use the Global Env Var, instead a local var
  &{HEADERS} =            Create Dictionary       Content-Type=application/json  Authorization=Basic cmVjLXFhOkZ6VUVyVGRSTTI2d2RoMlIyMk5lbktmeXlya2t2VGhR
  Create Session          ConfigAPISessionProd    url=https://recommendation-config.appier.co  headers=&{HEADERS}  disable_warnings=1  retry_status_list=[500,502,503,504]  timeout=1
  ${resp} =               Get On Session          ConfigAPISessionProd    url=/v0/datafeeds/${df}/labels?label_type=${label_type}  expected_status=200

  Should Not Be Empty     ${resp.json()}
  ...                     [${df}] @{TEST TAGS} FAIL: Can't get the response (labels) via the ConfigAPI response: ${resp.json()}, Request: , Request: ${resp.request.url}. please check the 'datafeed-labels' setting
  Set Test Variable       ${request.url}          ${resp.request.url}
  Set Test Variable       ${request.headers}      ${resp.request.headers}

  # Extra the request/response-header to Testrial & console
  Set Test Message        \n                      append=yes
  Set Test Message        *** Config API Request URL(Get datafeed labels on ${df}): ${resp.url}  append=yes

  ${actived_labels_list} =  Get Value From Json  ${resp.json()}  $.list[*].['label_id','enabled']

  # the actived_labels_list like: ['3QsodAuH', True, '8jYUfBj2OV', False, '9eilEeXD', False]
  # try to find the postition index when the boolean is True
  FOR  ${index}  ${item}  IN ENUMERATE  @{actived_labels_list}
    IF  ${index}%2 != 0
      IF  ${item}
        Set Test Variable   ${item_index}   ${index}
        IF  ${item}
          BREAK
        END
      END
    END
  END

  ${the_psid_index} =     Evaluate                ${item_index}-1
  ${extracted_actived_label} =  Get From List  ${actived_labels_list}  ${the_psid_index}

  Set Test Message        \n                      append=yes
  Set Test Message        [${df}] The extracted actived label item from Config API is: ${extracted_actived_label}  append=yes
  Set Test Message        \n                      append=yes
  RETURN                  ${extracted_actived_label}

