*** Keywords ***
# Tearup/Down #
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
