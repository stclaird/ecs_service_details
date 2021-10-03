package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func print_results(processedContainerInstances []ContainerInstance, tasks []*ecs.Task, service_name string) {
	//Print the results
	fmt.Printf("Instances Running %v:\n", service_name)
	for _, ci := range processedContainerInstances {
		fmt.Printf("Status: %v | EC2 Instance ID: %v | Num Running Tasks: %v | EC2 Private IP: %v\n", ci.Status, ci.Ec2InstanceId, ci.RunningTasksCount, ci.IpAddress)
	}
	fmt.Printf("Tasks Running %v:\n", service_name)
	for _, task := range tasks {
		fmt.Printf("Task Arn: %v, Started At: %v\n", *task.TaskArn, task.StartedAt)
	}
}

func list_tasks(cluster_name string, service_name string, svc *ecs.ECS) *ecs.ListTasksOutput {
	//fetch the running tasks for given service
	input := &ecs.ListTasksInput{
		Cluster:     aws.String(cluster_name),
		ServiceName: aws.String(service_name),
	}

	result, err := svc.ListTasks(input)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func describe_tasks(svc *ecs.ECS, cluster_name string, task_arns []*string) ([]*ecs.Task, []ContainerInstance) {
	//fetch task details
	describe_tasks_params := &ecs.DescribeTasksInput{
		Cluster: aws.String(cluster_name),
		Tasks:   task_arns,
	}
	describe_tasks_resp, describe_tasks_err := svc.DescribeTasks(describe_tasks_params)
	if describe_tasks_err != nil {
		fmt.Println("Cannot retrieve ECS task details:", describe_tasks_err)
	}

	containerInstances := make([]string, 0)
	processedContainerInstances := []ContainerInstance{}

	//for every task in slice of task get the container instance ARN it is running on and store it in  containerInstances
	//Tasks can run on the same container so we only add the container instance once
	for _, task := range describe_tasks_resp.Tasks {
		strContainerInstanceArnValue := *task.ContainerInstanceArn
		exists := exists_in_slice(strContainerInstanceArnValue, containerInstances)

		if exists == false {
			containerInstances = append(containerInstances, strContainerInstanceArnValue)
		}
	}

	//for every container instance ARN in containerInstances
	//get the instance ID and add them to processedContainerInstances
	for _, container_arn := range containerInstances {
		ContainerInstance := describe_container_instances(container_arn, cluster_name, svc)
		processedContainerInstances = append(processedContainerInstances, ContainerInstance)
	}

	return describe_tasks_resp.Tasks, processedContainerInstances
}

func describe_container_instances(container_arn string, cluster_name string, svc *ecs.ECS) ContainerInstance {
	//Describe container instance running our service task
	var container ContainerInstance

	input := &ecs.DescribeContainerInstancesInput{
		Cluster: aws.String(cluster_name),
		ContainerInstances: []*string{
			aws.String(container_arn),
		},
	}
	describe_containers_resp, _ := svc.DescribeContainerInstances(input) //TODO Handle errors
	out := describe_containers_resp
	container.Ec2InstanceId = *out.ContainerInstances[0].Ec2InstanceId
	container.Status = *out.ContainerInstances[0].Status
	container.RunningTasksCount = *out.ContainerInstances[0].RunningTasksCount

	return container
}

func DescribeEC2Instances(ec2_svc *ec2.EC2) *ec2.DescribeInstancesOutput {
	//retrieve ECS Instance data
	filters := []*ec2.Filter{
		&ec2.Filter{},
	}

	input := ec2.DescribeInstancesInput{Filters: filters}
	ec2_resp, err := ec2_svc.DescribeInstances(&input)

	if err != nil {
		fmt.Println("Cannot retrieve EC2 Reservations:", err)
	}

	return ec2_resp
}

func ServiceToCluster(svc *ecs.ECS, serviceName string) string {
	//Takes a service name and searches for the cluster
	//input service and session, output ecs cluster name
	var cluster_to_return string

	fmt.Println("Service to Find:", serviceName)
	input := &ecs.ListClustersInput{}
	cluster_result, err := svc.ListClusters(input)
	if err != nil {
		fmt.Println("Cannot retrieve ECS Clusters:", err)
	}

	for _, cluster := range cluster_result.ClusterArns {
		cluster_name := strings.Split(*cluster, "/")[1]

		input := &ecs.ListServicesInput{
			Cluster: aws.String(cluster_name),
		}

		service_result, _ := svc.ListServices(input)

		for _, service := range service_result.ServiceArns {
			service_name := strings.Split(*service, "/")[2]
			if service_name == serviceName {
				cluster_to_return = cluster_name
				break
			}
		}
	}

	return cluster_to_return
}

func exists_in_slice(s string, strings []string) bool {
	//find a string in slice,
	//return boolean false for not found true for found

	for _, val := range strings {
		if val == s {
			return true
		}
	}
	return false
}
