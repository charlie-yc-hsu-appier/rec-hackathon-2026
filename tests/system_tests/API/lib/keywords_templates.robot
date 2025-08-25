*** Keywords ***
# BDD Template For Checking Coupang vendor Group #
Check the Coupang vendor group
  [Arguments]
  ...                 ${oid}
  ...                 ${group_id}
  ...                 ${layout_id}
  ...                 ${expected_group_name}
  ...                 ${expected_url_params}
  ...                 ${expected_img_domain}
  ...                 ${expected_layout_code}=${Empty}
  ...                 ${with_fix_click_id}=${Empty}

  Given I have an vendor session

  # Parse layout_id to extract width and height
  ${w} =    Set Variable    300
  ${h} =    Set Variable    300
  
  # Check if layout_id contains dimensions (e.g., "300x300", "1200x627")
  ${dimension_match} =    Run Keyword And Return Status    Should Match Regexp    ${layout_id}    ^.*?(\\d+)x(\\d+).*$
  IF    ${dimension_match}
    ${matches} =    Get Regexp Matches    ${layout_id}    ^.*?(\\d+)x(\\d+).*$    1    2
    IF    ${matches}
      ${w} =    Set Variable    ${matches[0]}
      ${h} =    Set Variable    ${matches[1]}
    END
  END

  IF  '${with_fix_click_id}' != '${Empty}'
    Set Local Variable  ${_cid}         RFTEST
    Set Local Variable  ${_bidobjid}    RFTEST
  ELSE
    Set Local Variable  ${_cid}         ${Empty}
    Set Local Variable  ${_bidobjid}    ${Empty}
  END

  # When the vendor is from linkmine, we don't need the i_group param
  IF  '${group_id}' != '${Empty}'
    When I would like to set the session under vendor endpoint with  endpoint=r  vendor_key=${oid}  user_id=55660000-0000-4C18-AAAA-556624AF0000  click_id=${_cid}  w=${w}  h=${h}
  # When the vendor is from replace (oid = ig6jGmNbQvqiqpQ0XDqNpw), we'll use the xst to decide the layout
  ELSE IF  '${oid}' == 'ig6jGmNbQvqiqpQ0XDqNpw'
    When I would like to set the session under vendor endpoint with  endpoint=r  vendor_key=${oid}  user_id=55660000-0000-4C18-AAAA-556624AF0000  click_id=${_cid}  w=${w}  h=${h}
  ELSE
    When I would like to set the session under vendor endpoint with  endpoint=r  vendor_key=${oid}  user_id=55660000-0000-4C18-AAAA-556624AF0000  click_id=${_cid}  w=${w}  h=${h}
  END

  Then I would like to check status_code should be "200" within the current session
  Then In the user2item response payload, there is no duplicated SKU item in the response
  Then In the user2item response payload, the "url" should exist (EC-REC Common Value)
  Then Check the regex pattern in the given list  ${extracted_url}  (${expected_group_name})
  Then Check the regex pattern in the given list  ${extracted_url}  (${expected_url_params})
  # When the vendor is from INL, we will base on the group to check the layout code
  IF  '${group_id}' != '${Empty}'
    Then Check the regex pattern in the given list  ${extracted_url}  (${expected_layout_code})
  END
  Then In the user2item response payload, the "img" should exist (EC-REC Common Value)
  Then Check the regex pattern in the given list  ${extracted_img_url}  (${expected_img_domain})
  Then In the user2item response payload, the "custom_label" should not exist


# BDD Template For Checking customized vendor Group #
Check the keeta vendor group
  [Arguments]         ${sid}        ${df}          ${oid}          ${expected_url_params}  ${expected_img_domain}
  [Documentation]     We should use the df="android--com.sankuai.sailor.afooddelivery_2" to verify the Keeta-api group

  Given I have a config_api session
  ${oid_list} =  I would like to get the not Finished & keeta_api enabled oid list when datafeed_id=android--com.sankuai.sailor.afooddelivery_2
  Should Not Be Empty   ${oid_list}  @{TEST TAGS} FAIL: Can't get any enabled keeta_api oid list via the ConfigAPI response: ${oid_list}

  # Since Keeta will limit specific range for lat/lon,we need to use the lat=22.2800 lon=114.1600
  Given I have an ecrec session and domain is "dync-stg.c.appier.net"
  When I would like to set the session under user2item endpoint with  sid=${sid}  df=${df}  oid=${oid_list[0]}
  ...            lat=22.2800    lon=114.1600   _debug_creative=false   idfa=55660000-0000-4C18-AAAA-556624AF0001  cid=RFTEST-Keeta  bidobjid=RFTEST-Keeta
  ...            num_items=14   no_cache=true

  Then I would like to check status_code should be "200" within the current session
  Then In the user2item response payload, the "url" should exist (EC-REC Common Value)
  Then Check the regex pattern in the given list  ${extracted_url}  (${expected_url_params})
  Then In the user2item response payload, the "img" should exist (EC-REC Common Value)
  Then Check the regex pattern in the given list  ${extracted_img_url}  (${expected_img_domain})
  Then In the user2item response payload, the "custom_label" should not exist

