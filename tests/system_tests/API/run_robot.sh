#!/bin/bash

# Shell Paramteres
# The work_dir on k8s
work_dir="./tests/system_tests/API"
run_time=$(date +%s)

# TestRail Parameters
TESTRAIL_EC_REC_PROJECT_ID="8"

# The Testrail API Parameters
TESTRAIL_HOST="appier.testrail.io"


usage () {
   cat <<HELP_USAGE

   $0 -u [user] -k [key] -c [case] -t [type] -d [dryrun]
   
   Example: ./run_robot.sh -u your@tesrail.account.com -k your_API_Key -c ec_rec -t rat -d no
   or
   Example: ./run_robot.sh --user your@tesrail.account.com --key your_API_Key --case ec_rec --type rat --dryrun no


   case:
   -u|--user      The testrail account 
   -k|--key       The testrail API key
   -c|--case      Choose the test cases of EC_REC project. (click,recommend)
   -t|--type      Choose the test type for your running test case. (rat,fast)
   -d|--dryrun    Run the cases without creating a test run in testrail (yse,no)
   -h|--help      Usage.
   

HELP_USAGE
}

get_user_info () {
       
        if [ -z $1 ]
          then
             echo -e "Please give the tesrail user account , within the -u option."
          else
             TESTRAIL_USER="$1"
        fi
}

get_key_info () {

        if [ -z $1 ]
          then
             echo -e "Please give the tesrail API key, within the -k option."
          else
             TESTRAIL_API_KEY="$1"
        fi
}

get_case () {

        if [ -z $1 ]
          then
             echo -e "Please give the cases value within the -c option."
          else
             TEST_CASE="$1"
        fi

}

get_case_type () {

        if [ -z $1 ]
          then
             echo -e "Please give the case type value within the -t option."
          else
             CASE_TYPE="$1"
        fi

}

# Mapping the suite ID and run the test cases
run_case () {
    
    # when the $1=yes, means the dry-run 
    if [ $1 == "yes" ];then
        echo -e "\r\n==================== Start to Dry-Run the ${TEST_CASE^^} ${CASE_TYPE^^} Case ====================="
        echo -e "\r\n Dry Run  !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\r\n"
        sleep 1
        robot -i ${CASE_TYPE} --outputdir ${work_dir}/report ${work_dir}/testsuite/api_${TEST_CASE,,}_${CASE_TYPE,,}.robot
        echo -e "\r\n Dry-run is done ...\r\n"
        exit

    # when the $1=no, it will update the result 
    else
        if [ ${TEST_CASE,,} == "rec_vendor" ] && [ ${CASE_TYPE,,} == "rat" ];then
           SUITE_ID="279811"
        else 
         echo -e "\r\n !!! Something Wrong !!!\r\n"
         echo -e " !!! No mapping result could be found between test case & case type."
         echo -e " !!! please check suite_id in the testrail or the mapping relationship between case and type."
         echo -e "  --> Stop the testing ... \r\n"
         exit
        fi     

        echo -e "\r\n==================== Start to run the ${TEST_CASE^^} ${CASE_TYPE^^} Case ====================="
        sleep 1
        echo -e "The Test Case is: ${TEST_CASE^^}"
        echo -e "The Case Type is: ${CASE_TYPE^^}"
        echo -e "The Case Suite ID is: ${SUITE_ID}"

        if [ -z ${SUITE_ID} ];then 
            echo -e "\r\n !!! Something Wrong !!!\r\n"
            echo -e " !!! The suite_id is empty."
            echo -e " !!! Please help to check if the suite_id has been implemented in this run_robot.sh" 
            echo -e "  --> Stop the testing ... \r\n"
            exit
        fi

        run_id=$(http POST https://${TESTRAIL_HOST}/index.php?/api/v2/add_run/${TESTRAIL_EC_REC_PROJECT_ID} suite_id=${SUITE_ID} name=RobotFramework_${TEST_CASE^^}_${CASE_TYPE^^}_${run_time} include_all:=true  'Content-Type: application/json' -a ${TESTRAIL_USER}:${TESTRAIL_API_KEY} --ignore-stdin | jq '.id')
             
        if [ ${run_id} == "null" ];then 
                echo -e "\r\n !!! Something Wrong !!!\r\n"
                echo -e " !!! The run_id is null."
                echo -e " !!! Please help to check the paramters of testrail or API key."
                echo -e " !!! Also, please help to check if the test suit | test case | test type are existed in your project or testaril"
                echo -e " !!! TestRail User Account: ${TESTRAIL_USER}"
                echo -e "  --> Stop the testing ... \r\n"
                exit

        else
                echo -e "The Run Task ID is : ${run_id}\r\n"
                pabot --pabotlib --resourcefile ${work_dir}/res/valueset.dat --listener RetryFailed:3 --listener TestRailListener:${TESTRAIL_HOST}:${TESTRAIL_USER}:${TESTRAIL_API_KEY}:${run_id}:https:update --RemoveKeywords passed -v ENV:stag -i ${CASE_TYPE} -o ${TEST_CASE,,}_${CASE_TYPE,,}.xml --outputdir ${work_dir}/report ${work_dir}/testsuite/api_${TEST_CASE,,}_${CASE_TYPE,,}*.robot

                #robot --listener TestRailListener:${TESTRAIL_HOST}:${TESTRAIL_USER}:${TESTRAIL_API_KEY}:${run_id}:https:update -i ${CASE_TYPE} -o ${TEST_CASE,,}_${CASE_TYPE,,}.xml --outputdir ${work_dir}/report ${work_dir}/testsuite/api_${TEST_CASE,,}_${CASE_TYPE,,}.robot
                echo -e "\r\nFinishing the test and start to close the run task..."
                close_click_status=$(http POST https://${TESTRAIL_HOST}/index.php?/api/v2/close_run/${run_id} 'Content-Type: application/json' -a ${TESTRAIL_USER}:${TESTRAIL_API_KEY} | jq '.is_completed')
                echo -e "The Run Task - ${run_id} is completed : ${close_click_status}\r\n"

                pass_count=$(http GET https://${TESTRAIL_HOST}/index.php?/api/v2/get_run/${run_id} 'Content-Type: application/json' -a ${TESTRAIL_USER}:${TESTRAIL_API_KEY} --ignore-stdin | jq '.passed_count')
                failed_count=$(http GET https://${TESTRAIL_HOST}/index.php?/api/v2/get_run/${run_id} 'Content-Type: application/json' -a ${TESTRAIL_USER}:${TESTRAIL_API_KEY} --ignore-stdin | jq '.failed_count')
                pass_rate= awk -v x=$pass_count -v y=$(($pass_count + $failed_count)) 'BEGIN {printf "The Pass Rate - %.2f%%\n",x/y*100}' > ${work_dir}/testsuite/pass_rate.txt
                exit
        fi
    fi
}



options=$(getopt -o u:k:c:t:d:h --long user:,key:,case:,type:,dryrun:,help -- "$@")
[ $? -eq 0 ] || { 
    echo "Incorrect options provided"
    exit 1
}
eval set -- "$options"
while true; do
  case "$1" in

    -u|--user)
       get_user_info $2 ;;

    -k|--key)
       get_key_info $2 ;;
   
    -c|--case)
       get_case $2 ;;
   
    -t|--type)
       get_case_type $2 ;;

    -d|--dryrun)
        case "$2" in
           yes)
             run_case $2 ;;
           no)
             run_case no ;; 
           *)
             echo -e "\r\n The value of dryrun should be: yes|no.\r\n" 
             exit ;; 
        esac ;;      

    -h|--help)
       usage
       exit ;;

    esac
    shift
done
exit 0;
