*** Settings ***
Resource            ../res/init.robot

*** Keywords ***
# vendor API endpoint #
I have an vendor session
  [Documentation]  Create HTTP session for Vendor API endpoint.
  ...
  ...              *Purpose*
  ...              Initializes persistent HTTP session for Vendor API testing with
  ...              automatic retry on server errors.
  ...
  ...              *Configuration*
  ...              - Base URL: ${HTTP_METHOD}://${VENDOR_HOST}
  ...              - Auto-retry on server errors: 500, 502, 503, 504
  ...              - Warnings disabled
  ...
  ...              *Usage Example*
  ...              | I have an vendor session |
  ...
  ...              *Implementation*
  ...              Creates RequestsLibrary session named 'VendorSession' for all vendor
  ...              endpoint testing. ${HTTP_METHOD} and ${VENDOR_HOST} loaded from valueset.dat.
  ...
  ...              *Prerequisites*
  ...              - ${HTTP_METHOD} must be set (via Get Test Value)
  ...              - ${VENDOR_HOST} must be set (via Get Test Value)
  
  Create Session  VendorSession   url=${HTTP_METHOD}://${VENDOR_HOST}  disable_warnings=1  retry_status_list=[500,502,503,504]


I would like to set the session under vendor endpoint with
  [Documentation]  *Purpose:*
  ...              Send GET request to vendor API endpoint with flexible parameter support for different vendor types.
  ...              Handles healthz endpoint, standard vendor endpoints, and vendor-specific parameter requirements.
  ...              
  ...              *Parameters:*
  ...              - endpoint: The endpoint path (default: /healthz). For vendor endpoints use format r/{vendor_name}
  ...              - user_id: User identifier (automatically case-adjusted for vendor requirements)
  ...              - click_id: Click tracking identifier
  ...              - w, h: Image dimensions (width, height)
  ...              - bundle_id: Bundle identifier (linkmine-specific)
  ...              - adtype: Ad type (linkmine-specific)
  ...              - subid: Vendor subid from Config API
  ...              - lat, lon: GPS coordinates (keeta-specific)
  ...              - k_campaign_id: Campaign identifier (keeta-specific)
  ...              - os: Operating system (adforus-specific, android/ios)
  ...              
  ...              *Returns:*
  ...              Sets test variables: ${status_code}, ${resp_json}, ${request.url}, ${request.headers}
  ...              
  ...              *Usage Examples:*
  ...              ```robotframework
  ...              # Healthz endpoint
  ...              I would like to set the session under vendor endpoint with  endpoint=/healthz
  ...              
  ...              # Linkmine vendor with standard parameters
  ...              I would like to set the session under vendor endpoint with  endpoint=r/linkmine  user_id=${user_id}  click_id=${click_id}  w=300  h=300  bundle_id=${EMPTY}  adtype=banner  subid=vendor_subid
  ...              
  ...              # Keeta vendor with GPS and campaign
  ...              I would like to set the session under vendor endpoint with  endpoint=r/keeta  user_id=${user_id}  click_id=${click_id}  w=1200  h=627  lat=22.3264  lon=114.1661  k_campaign_id=${campaign_name}
  ...              
  ...              # Adforus vendor with OS parameter
  ...              I would like to set the session under vendor endpoint with  endpoint=r/adforus  user_id=${user_id}  click_id=${click_id}  w=300  h=300  os=android
  ...              ```
  ...              
  ...              *Implementation:*
  ...              1. Set default endpoint (/healthz) if not provided
  ...              2. Build query parameters dictionary from supported parameters (vendor_key, user_id, click_id, w, h, bundle_id, adtype, subid, lat, lon, k_campaign_id, os)
  ...              3. Add Accept: */* header to request
  ...              4. Normalize endpoint path (ensure leading slash)
  ...              5. Send GET request with VendorSession
  ...              6. Validate response is not empty
  ...              7. Extract request URL and headers
  ...              8. Set test variables (status_code, resp_json, request.url, request.headers)
  ...              9. Log request URL and response headers to test message
  ...              10. Validate status code is not 4xx or 5xx
  ...              
  ...              *Prerequisites:*
  ...              - VendorSession must be created via "I have an vendor session" keyword
  ...              - For vendor endpoints, appropriate parameters must be provided based on vendor type
  ...              
  ...              *Side Effects:*
  ...              - Sets ${status_code} test variable with HTTP response code
  ...              - Sets ${resp_json} test variable with parsed JSON response
  ...              - Sets ${request.url} test variable with full request URL
  ...              - Sets ${request.headers} test variable with request headers
  ...              - Appends request URL and response headers to test message log
  ...              
  ...              *Special Cases:*
  ...              - Adforus vendor: user_id parameter automatically case-adjusted (uppercase for iOS, lowercase for Android)
  ...              - Keeta vendor: Requires lat/lon GPS coordinates and k_campaign_id
  ...              - Linkmine vendor: Requires bundle_id (can be empty string) and adtype parameters
  ...              - INL vendors: Use subid parameter with URL-encoded subparam format
  ...              - Endpoint path automatically normalized (adds leading slash if missing)
  ...              - Empty query parameters dictionary handled gracefully
  ...              - Fails test if response is empty or status code is 4xx/5xx
  [Arguments]             &{args}

  # Set default endpoint if not provided
  ${endpoint_has_value} =  Run Keyword And Return Status
  ...                     Dictionary Should Contain Key  ${args}  endpoint
  IF  ${endpoint_has_value}
    Set Local Variable  ${endpoint}     ${args}[endpoint]
  ELSE
    Set Local Variable  ${endpoint}     /healthz
  END

  # Handle optional parameters - Method 2: Direct dictionary filtering
  ${query_params} =       Create Dictionary
  @{param_names} =        Create List             vendor_key              user_id             click_id            w                       h   bundle_id   adtype    subid   lat   lon   k_campaign_id   os

  FOR  ${param}  IN  @{param_names}
    ${param_exists} =   Run Keyword And Return Status
    ...                 Dictionary Should Contain Key  ${args}  ${param}
    IF  ${param_exists}
      Set To Dictionary   ${query_params}     ${param}=${args}[${param}]
    END
  END

  # Set the request header with Accept: */*
  &{HEADERS} =            Create Dictionary
  ...                     Accept=*/*

  # Handle endpoint path - add leading slash if not present
  ${endpoint_starts_with_slash} =  Run Keyword And Return Status
  ...                     Should Start With       ${endpoint}             /
  IF  not ${endpoint_starts_with_slash}
    ${endpoint} =   Set Variable    /${endpoint}
  END

  # Start to send the GET request - handle empty params safely
  ${has_params} =         Get Length              ${query_params}
  IF  ${has_params} > 0
    ${resp} =   Get On Session  VendorSession   url=${endpoint}     headers=&{HEADERS}  params=&{query_params}
  ELSE
    ${resp} =   Get On Session  VendorSession   url=${endpoint}     headers=&{HEADERS}
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
    ...     ${resp_json}[message]
    ...     ok
    ...     @{TEST TAGS} FAIL: Healthz endpoint should return message: 'ok', but got: ${resp_json}[message]
  END


# Assertion Keywords #
I would like to check status_code should be "${expected_code}" within the current session
  [Documentation]  Verify HTTP response status code matches expected value.
  ...
  ...              *Purpose*
  ...              Validates that API response status code equals the expected code.
  ...              Commonly used to verify successful responses (200) or error handling.
  ...
  ...              *Parameters*
  ...              - expected_code: Expected HTTP status code (e.g., "200", "404", "500")
  ...
  ...              *Usage Example*
  ...              | I would like to check status_code should be "200" within the current session |
  ...
  ...              *Implementation*
  ...              Compares ${status_code} test variable against expected value.
  ...              Fails with detailed error message including request URL and headers if mismatch.
  ...
  ...              *Prerequisites*
  ...              - ${status_code} must be set by previous API call
  ...              - ${request.url} and ${request.headers} for error reporting
  
  Should Be Equal As Strings
  ...     ${status_code}
  ...     ${expected_code}
  ...     @{TEST TAGS} The status code isn't we expected:${expected_code}, we get:${status_code}, Request:${request.url}, R_Header:${request.headers}
  ...     values=False


Validate vendor response structure
  [Arguments]             ${response_json}        ${vendor_name}=${EMPTY}
  [Documentation]  *Purpose:*
  ...              Validate the structural integrity of vendor API response to ensure it conforms to expected format.
  ...              Verifies response is a non-empty array of product objects with required fields.
  ...              
  ...              *Parameters:*
  ...              - response_json: JSON response from vendor API (should be list of product dictionaries)
  ...              - vendor_name: Vendor identifier (default: ${EMPTY}). Used to determine special validation rules (keeta/adforus skip image validation)
  ...              
  ...              *Expected Structure:*
  ...              ```json
  ...              [
  ...                {
  ...                  "product_id": "1703093047",
  ...                  "url": "https://...",
  ...                  "image": "https://..."  // optional for Keeta/Adforus vendors
  ...                }
  ...              ]
  ...              ```
  ...              
  ...              *Usage Examples:*
  ...              ```robotframework
  ...              # Standard vendor (requires image field)
  ...              Validate vendor response structure  ${resp_json}  linkmine
  ...              
  ...              # Keeta vendor (skips image validation)
  ...              Validate vendor response structure  ${resp_json}  keeta
  ...              
  ...              # Adforus vendor (skips image validation)
  ...              Validate vendor response structure  ${resp_json}  adforus
  ...              
  ...              # Generic validation (requires all fields)
  ...              Validate vendor response structure  ${resp_json}
  ...              ```
  ...              
  ...              *Implementation:*
  ...              1. Validate response type is list/array
  ...              2. Verify response is not empty
  ...              3. Determine if vendor is Keeta or Adforus (image validation skip flag)
  ...              4. Iterate through each product in response array
  ...              5. Validate required keys exist: product_id, url
  ...              6. Validate optional image key (skip for Keeta/Adforus)
  ...              7. Verify all field values are not empty
  ...              8. Log validation status for each product
  ...              9. Log final summary with product count
  ...              
  ...              *Validation Rules:*
  ...              - Response must be list type
  ...              - Response must contain at least one product
  ...              - Each product must have product_id (non-empty)
  ...              - Each product must have url (non-empty)
  ...              - Each product must have image (non-empty) EXCEPT Keeta and Adforus vendors
  ...              
  ...              *Special Cases:*
  ...              - Keeta vendor: Skips image field validation (image field not required)
  ...              - Adforus vendor: Skips image field validation (image field not required)
  ...              - Empty vendor_name: Requires all fields including image
  ...              - Logs emoji indicators: 🎯 for skipped validations, ✅ for passed validations

  # Response should be a list/array
  ${response_type} =      Evaluate                type($response_json).__name__
  Should Be Equal         ${response_type}        list
  ...                     Response should be a list/array, but got: ${response_type}

  # Response should not be empty
  Should Not Be Empty     ${response_json}
  ...                     Response array should not be empty

  # Check if this is Keeta or Adforus vendor (skip image validation)
  ${is_keeta} =           Run Keyword And Return Status
  ...                     Should Be Equal         ${vendor_name}      keeta
  ${is_adforus} =         Run Keyword And Return Status
  ...                     Should Be Equal         ${vendor_name}      adforus
  ${skip_image} =         Evaluate                ${is_keeta} or ${is_adforus}

  # Validate each product in the response
  FOR  ${product}  IN  @{response_json}
    Dictionary Should Contain Key  ${product}  product_id
    ...                     Each product should contain 'product_id' key

    Dictionary Should Contain Key  ${product}  url
    ...                     Each product should contain 'url' key

    # Skip image validation for Keeta and Adforus vendors
    IF  not ${skip_image}
      Dictionary Should Contain Key  ${product}  image
      ...                     Each product should contain 'image' key
    END

    # Validate that values are not empty
    ${product_id} =         Get From Dictionary     ${product}  product_id
    Should Not Be Empty     ${product_id}           product_id should not be empty

    ${url} =                Get From Dictionary     ${product}  url
    Should Not Be Empty     ${url}                  url should not be empty

    # Skip image validation for Keeta and Adforus vendors
    IF  not ${skip_image}
      ${image} =              Get From Dictionary     ${product}  image
      Should Not Be Empty     ${image}                image should not be empty
    ELSE
      Log                     🎯 ${vendor_name} vendor: skipping image validation
    END

    Log                     ✅ Product ${product_id} structure validation passed
  END

  ${product_count} =      Get Length              ${response_json}
  IF  ${skip_image}
    Log                   ✅ ${vendor_name} response structure validation passed for ${product_count} products (image validation skipped)
  ELSE
    Log                   ✅ Response structure validation passed for ${product_count} products
  END


Validate product patch contains product ids
  [Arguments]             ${response_json}        ${param_name}       ${expected_click_id_base64}  ${vendor_name}=${Empty}  ${os}=${Empty}  ${user_id}=${Empty}
  [Documentation]  *Purpose:*
  ...              Validate that product URLs contain correct tracking parameters with proper encoding.
  ...              Ensures click tracking is properly implemented across different vendor types with vendor-specific validation rules.
  ...              
  ...              *Parameters:*
  ...              - response_json: JSON response array from vendor API (list of product dictionaries)
  ...              - param_name: Name of the tracking parameter to validate (e.g., 'click_id', 'subparam')
  ...              - expected_click_id_base64: Base64-encoded click_id value to search for in product URLs
  ...              - vendor_name: Vendor identifier for special handling (default: ${Empty}). Options: keeta, adforus, inl_corp_X, linkmine
  ...              - os: Operating system for Adforus vendor (default: ${Empty}). Options: 'ios', 'android'
  ...              - user_id: User identifier for Adforus adid case validation (default: ${Empty})
  ...              
  ...              *Usage Examples:*
  ...              ```robotframework
  ...              # Standard vendor (linkmine)
  ...              Validate product patch contains product ids  ${resp_json}  click_id  ${encoded_click_id}  linkmine
  ...              
  ...              # INL vendor with subparam
  ...              Validate product patch contains product ids  ${resp_json}  subparam  ${encoded_click_id}  inl_corp_1
  ...              
  ...              # INL_corp_5 vendor with fixed subParam=pier (special case)
  ...              Validate product patch contains product ids  ${resp_json}  subparam  ${encoded_click_id}  inl_corp_5
  ...              
  ...              # Keeta vendor (skips click_id validation)
  ...              Validate product patch contains product ids  ${resp_json}  click_id  ${encoded_click_id}  keeta
  ...              
  ...              # Adforus vendor with OS-specific adid case validation
  ...              Validate product patch contains product ids  ${resp_json}  click_id  ${encoded_click_id}  adforus  os=android  user_id=${user_id}
  ...              Validate product patch contains product ids  ${resp_json}  click_id  ${encoded_click_id}  adforus  os=ios  user_id=${user_id}
  ...              ```
  ...              
  ...              *Implementation:*
  ...              1. Validate response is not empty
  ...              2. Check if vendor is Keeta (skip entire validation if true)
  ...              3. Iterate through each product in response array
  ...              4. Extract product_id, url, image from product dictionary
  ...              5. URL decode the product URL to normalize parameter format
  ...              6. Build search pattern:
  ...                 - INL_corp_5 with subparam: Fixed pattern 'subParam=pier' (uppercase P)
  ...                 - All other cases: Simple pattern '{param_name}={base64_value}'
  ...              7. Verify decoded URL contains expected tracking parameter pattern
  ...              8. If Adforus vendor: Validate adid case in product URL based on OS
  ...                 - iOS: Verify uppercase adid
  ...                 - Android: Verify lowercase adid
  ...              9. Log validation status for each product
  ...              10. Log final summary with total product count
  ...              
  ...              *Tracking Parameter Validation:*
  ...              - Product URLs are URL decoded before validation
  ...              - Standard format: {param_name}={base64_encoded_value}
  ...              - Special case: INL_corp_5 uses 'subParam=pier' (uppercase P, fixed value)
  ...              
  ...              *Prerequisites:*
  ...              - response_json must be valid JSON array
  ...              - Each product must contain url field
  ...              - For Adforus validation: os and user_id parameters must be provided
  ...              
  ...              *Special Cases:*
  ...              - Keeta vendor: Completely skips click_id validation (returns early)
  ...              - INL_corp_5 with subparam: Uses fixed parameter 'subParam=pier' (uppercase P) instead of dynamic base64 value
  ...              - Adforus vendor: Additional adid case validation after URL decode
  ...                * iOS: Converts user_id to uppercase and validates in decoded URL
  ...                * Android: Converts user_id to lowercase and validates in decoded URL
  ...                * Unknown OS: Logs warning and skips adid case validation
  ...              - All other vendors: Simple param_name=value format after URL decode
  ...              
  ...              *Validation Flow:*
  ...              1. Response non-empty check
  ...              2. Keeta vendor early return (skip all validation)
  ...              3. Per-product tracking parameter validation
  ...              4. Adforus adid case validation (if applicable)
  ...              5. Success logging with emoji indicators (🎯 ✅ ⚠️)

  # Response should be a list/array
  Should Not Be Empty     ${response_json}
  ...                     Response array should not be empty

  # Check if this is Keeta vendor (skip click_id validation)
  ${is_keeta} =           Run Keyword And Return Status
  ...                     Should Be Equal         ${vendor_name}      keeta

  IF  ${is_keeta}
    Log                   🎯 Keeta vendor detected - skipping click_id validation
    ${product_count} =    Get Length              ${response_json}
    Log                   ✅ Keeta vendor validation passed for ${product_count} products (click_id validation skipped)
    RETURN
  END

  # Check each product in the response
  FOR  ${product}  IN  @{response_json}
    ${product_id} =         Get From Dictionary     ${product}          product_id
    ${product_url} =        Get From Dictionary     ${product}          url
    ${product_image} =      Get From Dictionary     ${product}          image

    # URL decode the product URL to simplify pattern matching
    ${decoded_url} =        Evaluate                urllib.parse.unquote("${product_url}")  urllib.parse
    Log                     Decoded URL: ${decoded_url}
    
    # Check for INL_corp_5 special case (uses 'subParam=pier_{click_id}' format)
    ${is_inl_corp_5} =      Run Keyword And Return Status
    ...                     Should Contain  ${vendor_name}      inl_corp_5
    
    IF  ${is_inl_corp_5} and '${param_name}' == 'subparam'
      # For inl_corp_5: subParam=pier_{base64} (uppercase P, pier_ prefix)
      ${search_pattern} =     Set Variable    subParam=pier_${expected_click_id_base64}
      Log                     INL_CORP_5 vendor - searching for: ${search_pattern}
    ELSE
      # For all other vendors, use standard format: param_name=base64_value
      ${search_pattern} =     Set Variable    ${param_name}=${expected_click_id_base64}
      Log                     Searching for parameter: ${search_pattern}
    END

    # Verify the tracking parameter in decoded URL
    Should Contain          ${decoded_url}          ${search_pattern}
    ...                     Decoded URL should contain ${search_pattern}, but got: ${decoded_url}

    # Additional validation for Adforus vendor - check adid case in product_url
    ${is_adforus} =         Run Keyword And Return Status
    ...                     Should Be Equal         ${vendor_name}      adforus

    IF  ${is_adforus} and '${os}' != '${Empty}' and '${user_id}' != '${Empty}'
      Log                   🎯 Adforus vendor detected - validating adid case in product_url for OS: ${os}
      
      IF  '${os}' == 'ios'
        # For iOS, adid should be uppercase in product_url
        ${expected_adid_case} =  Convert To Uppercase  ${user_id}
        Should Contain      ${product_url}          ${expected_adid_case}
        ...                 Product URL should contain uppercase adid for iOS: adid=${expected_adid_case}, but got: ${product_url}
        Log                 ✅ iOS adid validation passed - found uppercase adid: ${expected_adid_case}
      ELSE IF  '${os}' == 'android'
        # For Android, adid should be lowercase in product_url
        ${expected_adid_case} =  Convert To Lowercase  ${user_id}
        Should Contain      ${product_url}          ${expected_adid_case}
        ...                 Product URL should contain lowercase adid for Android: adid=${expected_adid_case}, but got: ${product_url}
        Log                 ✅ Android adid validation passed - found lowercase adid: ${expected_adid_case}
      ELSE
        Log                 ⚠️ Unknown OS for adforus vendor: ${os}, skipping adid case validation
      END
    END

    Log                     ✅ Product ${product_id} validation passed - URL contains correct tracking parameter: ${search_pattern}
  END

  ${product_count} =      Get Length              ${response_json}
  Log                     ✅ All ${product_count} products validated successfully

  Log                     ✅ All product_ids found in product_patch with correct tracking parameters
