#!/bin/bash
while getopts ":c:s:" opt; do
  case $opt in
    c) CLUSTER="$OPTARG"
    ;;
    s) SERVICE="$OPTARG"
    ;;
  esac
done

if [[ -z ${CLUSTER+x} ]];
  then
    echo "Please set the ECS CLUSTER with the -c flag";
    exit 1
fi

if [[ -z ${SERVICE+x} ]];
  then
    echo "Please set the ECS Service with the -s flag";
    exit 1
fi

printEC2IP () {
    containerInstanceArn=$(aws ecs describe-tasks --cluster ${CLUSTER} --tasks ${taskID} | jq -r .tasks[].containerInstanceArn)
    ec2InstanceID=$(aws ecs describe-container-instances --cluster  ${CLUSTER} --container-instances ${containerInstanceArn} | jq -r .containerInstances[].ec2InstanceId)
    echo $(aws ec2 describe-instances --instance-ids ${ec2InstanceID} | jq -r .Reservations[].Instances[0].PrivateIpAddress)
}

tasks=$(aws ecs list-tasks --cluster ${CLUSTER} --service ${SERVICE} | jq -r .taskArns[0]) 

for taskID in ${tasks[@]}; do
    aws ecs describe-tasks --cluster ${CLUSTER} --tasks ${taskID} | jq
done