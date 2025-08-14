*** Variables ***
${TEST_API_TIMEOUT}         4 minutes
${work_dir}                 ./tests/system_tests/API

# [Coupang Only] These oids below are related to INL,linkmine,replace and Adpopcorn and we'll use these oids to test the vendor pictures
@{coupang_vendor_oid}       RghYkBdSRuyGkDcvIE6Eqg  4QSP2IQMQbahY5ipS88nnQ  ybZF8EpZQ86G76I2LnsMnA  qyhK5S25SJqbWEMuJ1l4jA  ig6jGmNbQvqiqpQ0XDqNpw  R0zFvq-OQSmVaJQ2XzOnTQ

*** Settings ***
Library     BuiltIn
Library     Collections
Library     DateTime
Library     ImapLibrary
Library     OperatingSystem
Library     Process
Library     RequestsLibrary
Library     String
Library     XML
Library     SeleniumLibrary
Library     JSONLibrary
Library     pabot.PabotLib
Library     ConfluentKafkaLibrary
Library     RedisLibrary
Library     ../lib/image.py

Resource    ../lib/keywords_init.robot
Resource    ../lib/keywords_util.robot
Resource    ../lib/keywords_vendor_api.robot
Resource    ../lib/keywords_config_api.robot
Resource    ../lib/keywords_templates.robot

