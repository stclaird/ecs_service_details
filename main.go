package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func main() {

	//TODO improve credential detection: Profiles, MFA etc
	//sess, err := session.NewSession()
	sess := session.Must(session.NewSession())

	ecs_svc := ecs.New(sess)
	ec2_svc := ec2.New(sess)

	var cluster_name string
	var service_name string

	var ecsCluster = flag.String("ecsCluster", "", "Specify the ECS CLUSTER where the ECS service is running")
	var serviceName = flag.String("serviceName", "", "Specify the ECS Service where the ECS service")
	flag.Parse()

	if *serviceName == "" {
		flag.Usage()
		os.Exit(1)
	}

	service_name = *serviceName

	if *ecsCluster == "" {
		clust := ServiceToCluster(ecs_svc, service_name)

		if clust == "" {
			fmt.Println("No ECS cluster specified and I could not detect it.\n Did you specify the correct service name or AWS Credentials")
			os.Exit(1)
		}
		cluster_name = clust

	} else {
		cluster_name = *ecsCluster
	}

	list_tasks_result := list_tasks(cluster_name, service_name, ecs_svc)
	tasks, processedContainerInstances := describe_tasks(ecs_svc, cluster_name, list_tasks_result.TaskArns)

	//examine the underlying EC2 instance of each container instance
	//append details of this such as IP address to the processedContainerInstances slice
	ec2_instances := DescribeEC2Instances(ec2_svc)

	for _, r := range ec2_instances.Reservations {
		for _, instance := range r.Instances {
			for i, p := range processedContainerInstances {
				if p.Ec2InstanceId == *instance.InstanceId {
					processedContainerInstances[i].IpAddress = *instance.PrivateIpAddress
				}
			}
		}
	}
	print_results(processedContainerInstances, tasks, service_name)
}
