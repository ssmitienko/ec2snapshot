package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func main() {
	instanceID := flag.String("i", "", "The ID of the instance to start")
	instanceDesc := flag.String("d", "", "Snapshot description")
	waitForComplete := flag.Bool("w", true, "Wait for snapshot to complete")
	flag.Parse()

	if *instanceID == "" {
		flag.Usage()
		return
	}

	if *instanceDesc == "" {
		fmt.Println("You must supply description for snapshot")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)
	if client == nil {
		panic("Failed to get client instance")
	}

	instanceSpecification := types.InstanceSpecification{
		ExcludeBootVolume: new(bool),
		InstanceId:        instanceID,
	}
	*instanceSpecification.ExcludeBootVolume = false

	createInstanceSnapshot := ec2.CreateSnapshotsInput{
		InstanceSpecification: &instanceSpecification,
		Description:           instanceDesc,
		DryRun:                new(bool),
	}
	*createInstanceSnapshot.DryRun = false

	createSnapshotsOutput, err := client.CreateSnapshots(context.TODO(), &createInstanceSnapshot)
	if err != nil {
		panic("CreateSnapshots error, " + err.Error())
	}

	if *waitForComplete {
		completed := false
		for !completed {
			completed = true
			time.Sleep(5 * time.Second)
			describeSnapshotsInput := ec2.DescribeSnapshotsInput{
				DryRun: new(bool),
			}
			for i := 0; i < len(createSnapshotsOutput.Snapshots); i++ {
				describeSnapshotsInput.SnapshotIds = append(describeSnapshotsInput.SnapshotIds, *createSnapshotsOutput.Snapshots[i].SnapshotId)
			}
			describeSnapshotsOutput, err := client.DescribeSnapshots(context.TODO(), &describeSnapshotsInput)
			if err != nil {
				panic("DescribeSnapshots error, " + err.Error())
			}
			for i := 0; i < len(describeSnapshotsOutput.Snapshots); i++ {
				if describeSnapshotsOutput.Snapshots[i].State == types.SnapshotStatePending {
					completed = false
				}
			}
		}
	}

	describeSnapshotsInput := ec2.DescribeSnapshotsInput{
		DryRun: new(bool),
	}
	for i := 0; i < len(createSnapshotsOutput.Snapshots); i++ {
		describeSnapshotsInput.SnapshotIds = append(describeSnapshotsInput.SnapshotIds, *createSnapshotsOutput.Snapshots[i].SnapshotId)
	}
	describeSnapshotsOutput, err := client.DescribeSnapshots(context.TODO(), &describeSnapshotsInput)
	if err != nil {
		panic("DescribeSnapshots error, " + err.Error())
	}
	for i := 0; i < len(describeSnapshotsOutput.Snapshots); i++ {
		fmt.Println(*describeSnapshotsOutput.Snapshots[i].VolumeId, *describeSnapshotsOutput.Snapshots[i].SnapshotId,
			*describeSnapshotsOutput.Snapshots[i].Progress,
			describeSnapshotsOutput.Snapshots[i].State)
	}
}
