*** Settings ***
# For more detail init setting, please refer to ../res/init.robot
Resource            ../res/init.robot
Test Timeout        ${TEST_API_TIMEOUT}

# For more detail init setting, please refer to ../res/valueset.robot
# Timeout Period: 10 times, Retry Strict Period: 1 sec,
Suite Setup         Get Test Value  ${ENV}
Suite Teardown      Release Test Value

*** Test Cases ***
# Automated vendor testing from YAML - Complete automation
[] [RAT] [VENDOR] [AUTO] Automated vendor testing from YAML
  [Tags]              testrailid=  RAT             VENDOR
  [Documentation]  Fully automated vendor testing using YAML configuration
  ...              Tests all vendors in YAML with complete automation:
  ...              - Dynamic endpoint generation: /r/{vendor_name}
  ...              - Parameter extraction from YAML (w, h from request_url)
  ...              - Base64 encoding validation
  ...              - Response structure validation
  ...              - Product patch verification

  # Complete YAML configuration for testing
  ${yaml_content} =   Catenate            SEPARATOR=\n
  ...                 vendors:
  ...                 ${SPACE}${SPACE}- name: linkmine
  ...                 ${SPACE}${SPACE}${SPACE}${SPACE}request_url: "https://api.adfork.kr/coupang_sch/?app_code=FAXXi4vdOY&limit=10&type=DNY&adid={user_id}&size=300x300"
  ...                 ${SPACE}${SPACE}${SPACE}${SPACE}tracking_url: "{product_url}&param1={click_id_base64}"
  ...                 ${SPACE}${SPACE}${SPACE}${SPACE}with_proxy: true

  # Run complete automated testing for all vendors in YAML
  Test vendors from yaml configuration  ${yaml_content}


# Basic vendor API healthz test
[C4977913] [RAT] [VENDOR] [HEALTHZ] Test vendor API healthz endpoint
  [Tags]  testrailid=4977913  RAT     VENDOR  HEALTHZ
  [Documentation]  Test the vendor API healthz endpoint to ensure it returns proper status and message

  Given I have an vendor session
  When I would like to set the session under vendor endpoint with  endpoint=/healthz
  Then I would like to check status_code should be "200" within the current session
