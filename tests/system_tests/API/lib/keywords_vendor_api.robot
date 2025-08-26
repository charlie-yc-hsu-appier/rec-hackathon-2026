*** Keywords ***
# vendor API endpoint #
I have an vendor session

  Create Session  VendorSession    url=${HTTP_METHOD}://${VENDOR_HOST}  disable_warnings=1  retry_status_list=[500,502,503,504]


I would like to set the session under vendor endpoint with
  [Documentation]  Send GET request to vendor API endpoint.
  ...              Available parameters:
  ...              endpoint - The endpoint path (default: /healthz)
  ...              For /r endpoint format: endpoint=r/{vendor_name}
  ...              For example:
  ...              Given I would like to set the session under vendor endpoint with  endpoint=/healthz
  ...              Given I would like to set the session under vendor endpoint with  endpoint=r/linkmine  user_id=uuid  click_id=value  w=300  h=300
  [Arguments]             &{args}

  # Set default endpoint if not provided
  ${endpoint_has_value} =  Run Keyword And Return Status
  ...                     Dictionary Should Contain Key  ${args}  endpoint
  IF  ${endpoint_has_value}
    Set Local Variable  ${endpoint}  ${args}[endpoint]
  ELSE
    Set Local Variable  ${endpoint}  /healthz
  END

  # Handle optional parameters - Method 2: Direct dictionary filtering
  ${query_params} =         Create Dictionary
  @{param_names} =          Create List  vendor_key  user_id  click_id  w  h
  
  FOR  ${param}  IN  @{param_names}
    ${param_exists} =       Run Keyword And Return Status
    ...                     Dictionary Should Contain Key  ${args}  ${param}
    IF  ${param_exists}
      Set To Dictionary     ${query_params}  ${param}=${args}[${param}]
    END
  END

  # Set the request header with Accept: */*
  &{HEADERS} =            Create Dictionary
  ...                     Accept=*/*

  # Handle endpoint path - add leading slash if not present
  ${endpoint_starts_with_slash} =    Run Keyword And Return Status    Should Start With    ${endpoint}    /
  IF    not ${endpoint_starts_with_slash}
    ${endpoint} =    Set Variable    /${endpoint}
  END

  # Start to send the GET request - handle empty params safely
  ${has_params} =         Get Length              ${query_params}
  IF  ${has_params} > 0
    ${resp} =             Get On Session          VendorSession           url=${endpoint}     headers=&{HEADERS}    params=&{query_params}
  ELSE
    ${resp} =             Get On Session          VendorSession           url=${endpoint}     headers=&{HEADERS}
  END

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


Validate vendor response structure
  [Arguments]    ${response_json}
  [Documentation]    Validate the basic structure of vendor response
  ...               Expected structure:
  ...               {
  ...                 "product_ids": [...],
  ...                 "product_patch": {...}
  ...               }
  
  Dictionary Should Contain Key    ${response_json}    product_ids
  ...    Response should contain 'product_ids' key
  
  Dictionary Should Contain Key    ${response_json}    product_patch
  ...    Response should contain 'product_patch' key
  
  ${product_ids} =    Get From Dictionary    ${response_json}    product_ids
  Should Not Be Empty    ${product_ids}    product_ids should not be empty
  
  ${product_patch} =    Get From Dictionary    ${response_json}    product_patch
  Should Not Be Empty    ${product_patch}    product_patch should not be empty
  
  Log    ✅ Response structure validation passed


Validate product patch contains product ids
  [Arguments]    ${response_json}    ${param_name}    ${expected_click_id_base64}
  [Documentation]    Validate that product_ids appear in product_patch
  ...               and verify the tracking parameter contains correct base64 encoded click_id
  
  ${product_ids} =    Get From Dictionary    ${response_json}    product_ids
  ${product_patch} =    Get From Dictionary    ${response_json}    product_patch
  
  # Check that each product_id appears in product_patch
  FOR    ${product_id}    IN    @{product_ids}
    Dictionary Should Contain Key    ${product_patch}    ${product_id}
    ...    Product ID ${product_id} should exist in product_patch
    
    # Get the product info and validate URL contains correct parameter
    ${product_info} =    Get From Dictionary    ${product_patch}    ${product_id}
    Dictionary Should Contain Key    ${product_info}    url
    ...    Product ${product_id} should have 'url' field
    
    ${product_url} =    Get From Dictionary    ${product_info}    url
    
    # Verify the tracking parameter contains the base64 encoded click_id
    Should Contain    ${product_url}    ${param_name}=${expected_click_id_base64}
    ...    Product URL should contain ${param_name}=${expected_click_id_base64}, but got: ${product_url}
    
    Log    ✅ Product ${product_id} validation passed - URL contains correct tracking parameter
  END
  
  Log    ✅ All product_ids found in product_patch with correct tracking parameters