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
  ...               [
  ...                 {
  ...                   "product_id": "1703093047",
  ...                   "url": "https://...",
  ...                   "image": "https://..."
  ...                 }
  ...               ]
  
  # Response should be a list/array
  ${response_type} =    Evaluate    type($response_json).__name__
  Should Be Equal    ${response_type}    list
  ...    Response should be a list/array, but got: ${response_type}
  
  # Response should not be empty
  Should Not Be Empty    ${response_json}    
  ...    Response array should not be empty
  
  # Validate each product in the response
  FOR    ${product}    IN    @{response_json}
    Dictionary Should Contain Key    ${product}    product_id
    ...    Each product should contain 'product_id' key
    
    Dictionary Should Contain Key    ${product}    url
    ...    Each product should contain 'url' key
    
    Dictionary Should Contain Key    ${product}    image
    ...    Each product should contain 'image' key
    
    # Validate that values are not empty
    ${product_id} =    Get From Dictionary    ${product}    product_id
    Should Not Be Empty    ${product_id}    product_id should not be empty
    
    ${url} =    Get From Dictionary    ${product}    url
    Should Not Be Empty    ${url}    url should not be empty
    
    ${image} =    Get From Dictionary    ${product}    image
    Should Not Be Empty    ${image}    image should not be empty
    
    Log    ✅ Product ${product_id} structure validation passed
  END
  
  ${product_count} =    Get Length    ${response_json}
  Log    ✅ Response structure validation passed for ${product_count} products


Validate product patch contains product ids
  [Arguments]    ${response_json}    ${param_name}    ${expected_click_id_base64}    ${vendor_name}=${Empty}
  [Documentation]    Validate that each product contains the correct tracking parameter
  ...               with base64 encoded click_id in the URL
  ...               New response format: array of products with product_id, url, image
  ...               Special handling for INL vendors with URL encoded parameters

  # Response should be a list/array
  Should Not Be Empty    ${response_json}
  ...    Response array should not be empty

  # Check each product in the response
  FOR    ${product}    IN    @{response_json}
    ${product_id} =    Get From Dictionary    ${product}    product_id
    ${product_url} =    Get From Dictionary    ${product}    url
    ${product_image} =    Get From Dictionary    ${product}    image
    
    # Check if this is an INL vendor
    ${is_inl_vendor} =    Run Keyword And Return Status    Should Contain    ${vendor_name}    inl
    Log    Debug - vendor_name: ${vendor_name}, param_name: ${param_name}, is_inl_vendor: ${is_inl_vendor}
    
    IF    ${is_inl_vendor} and '${param_name}' == 'subparam'
      # Special handling for inl_corp_5 vendor
      ${is_inl_corp_5} =    Run Keyword And Return Status    Should Contain    ${vendor_name}    inl_corp_5
      
      IF    ${is_inl_corp_5}
        # For inl_corp_5: subParam=pier (P is uppercase, fixed value)
        ${encoded_param} =    Set Variable    %26subParam%3Dpier
        ${search_pattern} =    Set Variable    ${encoded_param}
        Log    INL_CORP_5 vendor detected - searching for fixed parameter: ${search_pattern}
      ELSE
        # For other INL vendors, the parameter appears URL encoded in the land parameter
        # Format: %26subparam%3DMTJhYS4xMmFh (URL encoded &subparam=base64)
        # We need to encode both the & and = characters
        ${encoded_param} =    Evaluate    "%26" + urllib.parse.quote("${param_name}") + "%3D" + "${expected_click_id_base64}"    urllib.parse
        ${search_pattern} =    Set Variable    ${encoded_param}
        Log    INL vendor detected - searching for URL encoded parameter: ${search_pattern}
      END
    ELSE
      # For non-INL vendors, use standard format
      ${search_pattern} =    Set Variable    ${param_name}=${expected_click_id_base64}
      Log    Standard vendor - searching for parameter: ${search_pattern}
    END
    
    # Verify the tracking parameter contains the base64 encoded click_id
    Should Contain    ${product_url}    ${search_pattern}
    ...    Product URL should contain ${search_pattern}, but got: ${product_url}
    
    Log    ✅ Product ${product_id} validation passed - URL contains correct tracking parameter: ${search_pattern}
  END
  
  ${product_count} =    Get Length    ${response_json}
  Log    ✅ All ${product_count} products validated successfully
  
  Log    ✅ All product_ids found in product_patch with correct tracking parameters