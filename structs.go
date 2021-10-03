package main

type ContainerInstance struct {
	Ec2InstanceId string
	Status string
	RunningTasksCount int64
	IpAddress string
}

type TaskArns []struct {

  }

type describeTasks struct {
	Tasks []struct {
		Attachments          []interface{} `json:"attachments"`
		Availabilityzone     string        `json:"availabilityZone"`
		Clusterarn           string        `json:"clusterArn"`
		Connectivity         string        `json:"connectivity"`
		Connectivityat       float64       `json:"connectivityAt"`
		Containerinstancearn string        `json:"containerInstanceArn"`
		Containers           []struct {
			Containerarn      string        `json:"containerArn"`
			Taskarn           string        `json:"taskArn"`
			Name              string        `json:"name"`
			Image             string        `json:"image"`
			Runtimeid         string        `json:"runtimeId"`
			Laststatus        string        `json:"lastStatus"`
			Networkbindings   []interface{} `json:"networkBindings"`
			Networkinterfaces []interface{} `json:"networkInterfaces"`
			Healthstatus      string        `json:"healthStatus"`
			CPU               string        `json:"cpu"`
			Memoryreservation string        `json:"memoryReservation"`
			Imagedigest       string        `json:"imageDigest,omitempty"`
		} `json:"containers"`
		CPU           string  `json:"cpu"`
		Createdat     float64 `json:"createdAt"`
		Desiredstatus string  `json:"desiredStatus"`
		Group         string  `json:"group"`
		Healthstatus  string  `json:"healthStatus"`
		Laststatus    string  `json:"lastStatus"`
		Launchtype    string  `json:"launchType"`
		Memory        string  `json:"memory"`
		Overrides     struct {
			Containeroverrides []struct {
				Name string `json:"name"`
			} `json:"containerOverrides"`
			Inferenceacceleratoroverrides []interface{} `json:"inferenceAcceleratorOverrides"`
		} `json:"overrides"`
		Pullstartedat     float64       `json:"pullStartedAt"`
		Pullstoppedat     float64       `json:"pullStoppedAt"`
		Startedat         float64       `json:"startedAt"`
		Startedby         string        `json:"startedBy"`
		Tags              []interface{} `json:"tags"`
		Taskarn           string        `json:"taskArn"`
		Taskdefinitionarn string        `json:"taskDefinitionArn"`
		Version           int           `json:"version"`
	} `json:"tasks"`
	Failures []interface{} `json:"failures"`
}

