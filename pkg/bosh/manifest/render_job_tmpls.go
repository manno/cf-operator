package manifest

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	btg "github.com/viovanov/bosh-template-go"
	yaml "gopkg.in/yaml.v2"
)

// RenderJobTemplates will render templates for all jobs of the instance group
// https://bosh.io/docs/create-release/#job-specs
// TODO is this really a full boshManifestPath
func RenderJobTemplates(namespace string, boshManifestPath string, jobsDir string, jobsOutputDir string, instanceGroupName string, specIndex int) error {

	// Loading deployment manifest file
	resolvedYML, err := ioutil.ReadFile(boshManifestPath)
	if err != nil {
		return errors.Wrapf(err, "couldn't read manifest file %s", boshManifestPath)
	}
	boshManifest := Manifest{}
	err = yaml.Unmarshal(resolvedYML, &boshManifest)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal deployment manifest %s", boshManifestPath)
	}

	// Loop over instancegroups
	for _, instanceGroup := range boshManifest.InstanceGroups {

		// Filter based on the instance group name
		if instanceGroup.Name != instanceGroupName {
			continue
		}

		// Render all files for all jobs included in this instance_group.
		for _, job := range instanceGroup.Jobs {
			jobSpec, err := job.loadSpec(jobsDir)

			if err != nil {
				return errors.Wrapf(err, "failed to load job spec file %s", job.Name)
			}

			// Find job instance that's being rendered
			var currentJobInstance *JobInstance
			for _, instance := range instanceGroup.jobInstances(namespace, boshManifest.Name, job.Name, *jobSpec) {
				if instance.Index == specIndex {
					currentJobInstance = &instance
					break
				}
			}
			if currentJobInstance == nil {
				return fmt.Errorf("no instance found for spec index '%d'", specIndex)
			}

			// Loop over name and link
			// TODO only once per job instance?
			jobInstanceLinks := []Link{}
			for name, jobConsumersLink := range job.Properties.BOSHContainerization.Consumes {
				jobInstances := []JobInstance{}

				// Loop over instances of link
				// TODO calculate Instances from jobProviderLink info
				if len(jobConsumersLink.Instances) > 0 {
					jobConsumerLinkInstance := jobConsumersLink.Instances[0]
					for _, jobConsumerLinkInstance := range jobConsumersLink.Instances {
						jobInstances = append(jobInstances, JobInstance{
							Address: jobConsumerLinkInstance.Address,
							AZ:      jobConsumerLinkInstance.AZ,
							ID:      jobConsumerLinkInstance.ID,
							Index:   jobConsumerLinkInstance.Index,
							Name:    jobConsumerLinkInstance.Name,
						})
					}
				}

				jobInstanceLinks = append(jobInstanceLinks, Link{
					Name:       name,
					Instances:  jobInstances,
					Properties: jobConsumersLink.Properties,
				})
			}

			// Loop over templates for rendering files
			jobSrcDir := job.specDir(jobsDir)
			for source, destination := range jobSpec.Templates {
				absDest := filepath.Join(jobsOutputDir, job.Name, destination)
				os.MkdirAll(filepath.Dir(absDest), 0755)

				properties := job.Properties.ToMap()

				renderPointer := btg.NewERBRenderer(
					&btg.EvaluationContext{
						// TODO missing from btg, but calculated here
						// Instances:  jobInstanceLinks,
						Properties: properties,
					},

					&btg.InstanceInfo{
						Address: currentJobInstance.Address,
						AZ:      currentJobInstance.AZ,
						ID:      currentJobInstance.ID,
						Index:   string(currentJobInstance.Index),
						Name:    currentJobInstance.Name,
					},

					filepath.Join(jobSrcDir, JobSpecFilename),
				)

				// Create the destination file
				absDestFile, err := os.Create(absDest)
				if err != nil {
					return err
				}
				defer absDestFile.Close()
				if err = renderPointer.Render(filepath.Join(jobSrcDir, "templates", source), absDestFile.Name()); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
