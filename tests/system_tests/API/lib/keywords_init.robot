*** Keywords ***
# Tearup/Down #
# Warm Up the EC-REC Service
Check Service Is Ready
  [Arguments]             ${datafeed_id}=android--com.hnsmall  ${sid}=android--com.hnsmall

  # Post a sample request after the deployment to check if the EC-REC service is ready or not
  Create Session          User2ItemSession    url=http://${VENDOR_HOST}  disable_warnings=1
  &{HEADERS} =            Create Dictionary   Content-Type=application/x-www-form-urlencoded
  ${resp} =               Get On Session      User2ItemSession        url=/w/${datafeed_id}/r?sid=${sid}&idfa=qa-check-service-ready&num_items=1\  headers=&{HEADERS}

  # Status Code should < 300
  ${resp.status_code} =   Convert To String   ${resp.status_code}
  Should Not Match Regexp  ${resp.status_code}  (^(3|4|5)..$)  [WarmUp][${sid}] Something wrong with EC_REC Service \n request_url: ${resp.url} \n http_status_code: ${resp.status_code} \n message: ${resp.content} \n
  ...                     values=False


# Get the Server Info (source from: ../res/valueset.dat)
Get Test Value
  [Documentation]  Need to pass the value of "tags" in valueset.dat file to get a specific Server Info
  [Arguments]             ${server}

  ${valuesetname} =       Acquire Value Set       ${server}

  ${HTTP_METHOD} =        Get Value From Set      HTTP_METHOD

  # VENDOR API
  ${VENDOR_HOST} =        Get Value From Set      VENDOR_HOST
  # Config API
  ${CONFIG_API_HOST} =    Get Value From Set      CONFIG_API_HOST

  # Set all variable
  Set Suite Variable      ${SERVER_ENV}           ${server}
  Set Suite Variable      ${HTTP_METHOD}          ${HTTP_METHOD}
  Set Suite Variable      ${VENDOR_HOST}          ${VENDOR_HOST}
  Set Suite Variable      ${CONFIG_API_HOST}      ${CONFIG_API_HOST}


Release Test Value
  Release Value Set
