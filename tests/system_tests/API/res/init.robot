*** Variables ***
${TEST_API_TIMEOUT}         4 minutes
${work_dir}                 ./tests/system_tests/API

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

