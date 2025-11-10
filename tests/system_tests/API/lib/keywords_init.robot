*** Settings ***
Resource            ../res/init.robot

*** Keywords ***
# Tearup/Down #
# Get the Server Info (source from: ../res/valueset.dat)
Get Test Value
  [Documentation]  Initialize test environment variables from valueset.dat configuration.
  ...
  ...              *Purpose*
  ...              Retrieves and sets server configuration variables for Vendor API and Config API
  ...              testing based on environment tag (stag-01, stag-02, prod).
  ...
  ...              *Parameters*
  ...              - ${server}: Environment tag from valueset.dat (e.g., 'stag-01', 'stag-02')
  ...
  ...              *Side Effects*
  ...              - Sets ${SERVER_ENV} suite variable (server environment)
  ...              - Sets ${HTTP_METHOD} suite variable (http/https)
  ...              - Sets ${VENDOR_HOST} suite variable (Vendor API host)
  ...              - Sets ${CONFIG_API_HOST} suite variable (Config API host)
  ...
  ...              *Usage Example*
  ...              | Get Test Value | stag-01 |
  ...
  ...              *Implementation*
  ...              1. Acquires value set from valueset.dat by server tag
  ...              2. Extracts HTTP_METHOD, VENDOR_HOST, CONFIG_API_HOST
  ...              3. Sets all values as suite-level variables for test execution
  ...
  ...              *Prerequisites*
  ...              - valueset.dat must exist in ../res/ directory
  ...              - Server tag must be defined in valueset.dat
  
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
  [Documentation]  Release value set resources acquired during test setup.
  ...
  ...              *Purpose*
  ...              Releases pabot value set lock to allow other parallel test executions
  ...              to access the same value set resources.
  ...
  ...              *Usage Example*
  ...              | Release Test Value |
  ...
  ...              *Implementation*
  ...              Calls pabot.PabotLib's Release Value Set keyword to unlock resources.
  ...
  ...              *Prerequisites*
  ...              - Must be called after Get Test Value
  ...              - Used in test teardown to ensure proper resource cleanup
  
  Release Value Set
