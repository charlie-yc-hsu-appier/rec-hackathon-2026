*** Settings ***
# For more detail init setting, please refer to ../res/init.robot
Resource            ../res/init.robot
Test Timeout        ${TEST_API_TIMEOUT}


# For more detail init setting, please refer to ../res/valueset.robot
# Timeout Period: 10 times, Retry Strict Period: 1 sec,
Suite Setup         Get Test Value  ${ENV}
Suite Teardown      Release Test Value

*** Test Cases ***
# Vendor API test cases section
[C4977913] [RAT] [VENDOR] [HEALTHZ] Test vendor API healthz endpoint
  [Tags]  testrailid=4977913      RAT             VENDOR          HEALTHZ
  [Documentation]  Test the vendor API healthz endpoint to ensure it returns proper status and message

  Given I have an vendor session
  When I would like to set the session under vendor endpoint with  endpoint=/healthz
  Then I would like to check status_code should be "200" within the current session


# Coupang INL Group
[C4202651] [RAT] [INL] Check the Coupang INL group
  [Tags]      testrailid=4202651  RAT-T     inl_group
  [Template]  Check the Coupang vendor group
  [Documentation]  We'll check if the subparam exists in the url and the img domain should be ads-partners.coupang.com
  ...              the subparam is encode (cid.bidobjid) in base64, you could refer to https://appier.atlassian.net/browse/AI-23485
  ...              when ${with_fix_click_id} is ${TRUE}, we'll use "RFTEST" as the cid and bidobjid
  ...              otherwise, we'll let the cid and bidobjid blank
  ...              In this template, we will also use the layout_id to check related code in INL api:
  ...              https://bitbucket.org/plaxieappier/rec-reyka/src/staging/config-template/config-prd.yaml
  ...              Turn off the group 3 at 2025/06/11: https://appier.slack.com/archives/C07TTKA4SHJ/p1749604656454839

  # Args: ${oid}  ${group_id}  ${layout_id}  ${expected_group_name}  ${expected_url_params}  ${expected_img_domain}  ${expected_layout_code}=${Empty}  ${with_fix_click_id}=${Empty}
  RghYkBdSRuyGkDcvIE6Eqg  0  300x300  adapi.inlcorp.com  subparam%3DLg  ads-partners.coupang.com  650alldb2
  RghYkBdSRuyGkDcvIE6Eqg  1  dna_1200x627  api.adreload.com  subparam%3DUkZURVNULlJGVEVTVA  ads-partners.coupang.com  dynamicntdw1  ${TRUE}
  RghYkBdSRuyGkDcvIE6Eqg  2  kakao_kr_v2  api.adiostech.com  subparam%3DLg  ads-partners.coupang.com  kkobiz1200x600
  #RghYkBdSRuyGkDcvIE6Eqg  3  video_720x1280_30s  cp.edl.co.kr  subparam%3DLg  ads-partners.coupang.com  appsspdw3
  RghYkBdSRuyGkDcvIE6Eqg  4  video_1280x720_30s  cdw.adsrv.co.kr  subparam%3DLg  ads-partners.coupang.com  24650alldb2

  4QSP2IQMQbahY5ipS88nnQ  0  video_1280x720_30s  adapi.inlcorp.com  subparam%3DLg  ads-partners.coupang.com  650alldb2
  4QSP2IQMQbahY5ipS88nnQ  1  300x300  api.adreload.com  subparam%3DUkZURVNULlJGVEVTVA  ads-partners.coupang.com  650alldb2  ${TRUE}
  4QSP2IQMQbahY5ipS88nnQ  2  dna_1200x627  api.adiostech.com  subparam%3DLg  ads-partners.coupang.com  dynamicntdw1
  #4QSP2IQMQbahY5ipS88nnQ  3  kakao_kr_v2  cp.edl.co.kr  subparam%3DLg  ads-partners.coupang.com  appsspdw2
  4QSP2IQMQbahY5ipS88nnQ  4  video_720x1280_30s  cdw.adsrv.co.kr  subparam%3DLg  ads-partners.coupang.com  24650alldb2

  ybZF8EpZQ86G76I2LnsMnA  0  video_720x1280_30s  adapi.inlcorp.com  subparam%3DLg  ads-partners.coupang.com  650alldb2
  ybZF8EpZQ86G76I2LnsMnA  1  video_1280x720_30s  api.adreload.com  subparam%3DUkZURVNULlJGVEVTVA  ads-partners.coupang.com  650alldb2  ${TRUE}
  ybZF8EpZQ86G76I2LnsMnA  2  300x300  api.adiostech.com  subparam%3DLg  ads-partners.coupang.com  650alldb2
  #ybZF8EpZQ86G76I2LnsMnA  3  dna_1200x627  cp.edl.co.kr  subparam%3DLg  ads-partners.coupang.com  appsspdw1
  ybZF8EpZQ86G76I2LnsMnA  4  kakao_kr_v2  cdw.adsrv.co.kr  subparam%3DLg  ads-partners.coupang.com  24kkobizdw1


# Coupang Linkmine Group
[C4608232] [RAT] [Linkmine] Check the Coupang Linkmine group
  [Tags]      testrailid=4608232  RAT-T     linkmine_group
  [Template]  Check the Coupang vendor group
  [Documentation]  We'll check if the param1 exists in the url and the img domain should be ads-partners.coupang.com
  ...              the param1 is encode (cid.bidobjid) in base64, you could refer to https://appier.atlassian.net/browse/AI-23485
  ...              when ${with_fix_click_id} is ${TRUE}, we'll use "RFTEST" as the cid and bidobjid
  ...              otherwise, we'll let the cid and bidobjid blank

  # Args: ${oid}  ${group_id}  ${layout_id}  ${expected_group_name}  ${expected_url_params}  ${expected_img_domain}  ${expected_layout_code}=${Empty}  ${with_fix_click_id}=${Empty}
  qyhK5S25SJqbWEMuJ1l4jA  ${Empty}  300x300  api.linkmine.co.kr  param1\=UkZURVNULlJGVEVTVA  ads-partners.coupang.com  ${Empty}  ${TRUE}
  qyhK5S25SJqbWEMuJ1l4jA  ${Empty}  300x300  api.linkmine.co.kr  param1\=Lg  ads-partners.coupang.com


# Coupang Replace Group
[C4962884] [RAT] [Replace] Check the Coupang Replace group
  [Tags]      testrailid=4962884  RAT-T     replace_group
  [Template]  Check the Coupang vendor group
  [Documentation]  We'll check if the param1 exists in the url and the img domain should be ads-partners.coupang.com
  ...              the param1 is encode (cid.bidobjid) in base64, you could refer to https://appier.atlassian.net/browse/AI-23485
  ...              when ${with_fix_click_id} is ${TRUE}, we'll use "RFTEST" as the cid and bidobjid
  ...              otherwise, we'll let the cid and bidobjid blank

  # Args: ${oid}  ${group_id}  ${layout_id}  ${expected_group_name}  ${expected_url_params}  ${expected_img_domain}  ${expected_layout_code}=${Empty}  ${with_fix_click_id}=${Empty}
  ig6jGmNbQvqiqpQ0XDqNpw  ${Empty}  300x300  click.adshot.network  click_id\=UkZURVNULlJGVEVTVA  ads-partners.coupang.com  ${Empty}  ${TRUE}
  ig6jGmNbQvqiqpQ0XDqNpw  ${Empty}  300x300  click.adshot.network  click_id\=Lg  ads-partners.coupang.com


# Coupang Adpopcorn Group
[C4965132] [RAT] [Adpopcorn] Check the Coupang Adpopcorn group
  [Tags]      testrailid=4965132  RAT-T     adpopcorn_group
  [Template]  Check the Coupang vendor group
  [Documentation]  We'll check if the param1 exists in the url and the img domain should be ads-partners.coupang.com
  ...              the param1 is encode (cid.bidobjid) in base64, you could refer to https://appier.atlassian.net/browse/AI-23485
  ...              when ${with_fix_click_id} is ${TRUE}, we'll use "RFTEST" as the cid and bidobjid
  ...              otherwise, we'll let the cid and bidobjid blank

  # Args: ${oid}  ${group_id}  ${layout_id}  ${expected_group_name}  ${expected_url_params}  ${expected_img_domain}  ${expected_layout_code}=${Empty}  ${with_fix_click_id}=${Empty}
  R0zFvq-OQSmVaJQ2XzOnTQ  ${Empty}  300x300  link.coupang.com  subparam\=UkZURVNULlJGVEVTVA  ads-partners.coupang.com  ${Empty}  ${TRUE}
  R0zFvq-OQSmVaJQ2XzOnTQ  ${Empty}  300x300  link.coupang.com  subparam\=Lg  ads-partners.coupang.com


# Coupang Adpacker Group
[C4975033] [RAT] [Adpacker] Check the Coupang Adpacker group
  [Tags]      testrailid=4975033  RAT-T     adpacker_group
  [Template]  Check the Coupang vendor group
  [Documentation]  We'll check if the param1 exists in the url and the img domain should be ads-partners.coupang.com
  ...              the param1 is encode (cid.bidobjid) in base64, you could refer to https://appier.atlassian.net/browse/AI-23485
  ...              when ${with_fix_click_id} is ${TRUE}, we'll use "RFTEST" as the cid and bidobjid
  ...              otherwise, we'll let the cid and bidobjid blank

  # Args: ${oid}  ${group_id}  ${layout_id}  ${expected_group_name}  ${expected_url_params}  ${expected_img_domain}  ${expected_layout_code}=${Empty}  ${with_fix_click_id}=${Empty}
  3TapVPIUS2u3zl4qmpDzPA  ${Empty}  300x300  ad.n-bridge.io  ssp_click_id\=UkZURVNULlJGVEVTVA  ads-partners.coupang.com  ${Empty}  ${TRUE}
  3TapVPIUS2u3zl4qmpDzPA  ${Empty}  300x300  ad.n-bridge.io  ssp_click_id\=Lg  ads-partners.coupang.com


# Keeta Customed Group
[C4911675] [RAT] [Keeta] Check the keeta group
  [Tags]      testrailid=4911675  RAT-T     keeta_group
  [Template]  Check the keeta vendor group
  [Documentation]  To verify the Keeta's recommendation
  ...              We should use the df="android--com.sankuai.sailor.afooddelivery_2" to verify the Keeta-api group

  # Args: ${sid}  ${df}  ${oid}  ${expected_url_params}  ${expected_img_domain}
  android--com.sankuai.sailor.afooddelivery  android--com.sankuai.sailor.afooddelivery_2  Xi34YugiQV6-rjI1g2ThSw  app.adjust.com  prod.adsappier.com


