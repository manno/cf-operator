// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"code.cloudfoundry.org/cf-operator/pkg/bosh/bpm"
	"code.cloudfoundry.org/cf-operator/pkg/bosh/bpmconverter"
	"code.cloudfoundry.org/cf-operator/pkg/bosh/manifest"
	"code.cloudfoundry.org/cf-operator/pkg/kube/controllers/boshdeployment"
)

type FakeBPMConverter struct {
	ResourcesStub        func(string, bpmconverter.DomainNameService, string, *manifest.InstanceGroup, manifest.ReleaseImageProvider, bpm.Configs, string) (*bpmconverter.Resources, error)
	resourcesMutex       sync.RWMutex
	resourcesArgsForCall []struct {
		arg1 string
		arg2 bpmconverter.DomainNameService
		arg3 string
		arg4 *manifest.InstanceGroup
		arg5 manifest.ReleaseImageProvider
		arg6 bpm.Configs
		arg7 string
	}
	resourcesReturns struct {
		result1 *bpmconverter.Resources
		result2 error
	}
	resourcesReturnsOnCall map[int]struct {
		result1 *bpmconverter.Resources
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBPMConverter) Resources(arg1 string, arg2 bpmconverter.DomainNameService, arg3 string, arg4 *manifest.InstanceGroup, arg5 manifest.ReleaseImageProvider, arg6 bpm.Configs, arg7 string) (*bpmconverter.Resources, error) {
	fake.resourcesMutex.Lock()
	ret, specificReturn := fake.resourcesReturnsOnCall[len(fake.resourcesArgsForCall)]
	fake.resourcesArgsForCall = append(fake.resourcesArgsForCall, struct {
		arg1 string
		arg2 bpmconverter.DomainNameService
		arg3 string
		arg4 *manifest.InstanceGroup
		arg5 manifest.ReleaseImageProvider
		arg6 bpm.Configs
		arg7 string
	}{arg1, arg2, arg3, arg4, arg5, arg6, arg7})
	fake.recordInvocation("Resources", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6, arg7})
	fake.resourcesMutex.Unlock()
	if fake.ResourcesStub != nil {
		return fake.ResourcesStub(arg1, arg2, arg3, arg4, arg5, arg6, arg7)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.resourcesReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeBPMConverter) ResourcesCallCount() int {
	fake.resourcesMutex.RLock()
	defer fake.resourcesMutex.RUnlock()
	return len(fake.resourcesArgsForCall)
}

func (fake *FakeBPMConverter) ResourcesCalls(stub func(string, bpmconverter.DomainNameService, string, *manifest.InstanceGroup, manifest.ReleaseImageProvider, bpm.Configs, string) (*bpmconverter.Resources, error)) {
	fake.resourcesMutex.Lock()
	defer fake.resourcesMutex.Unlock()
	fake.ResourcesStub = stub
}

func (fake *FakeBPMConverter) ResourcesArgsForCall(i int) (string, bpmconverter.DomainNameService, string, *manifest.InstanceGroup, manifest.ReleaseImageProvider, bpm.Configs, string) {
	fake.resourcesMutex.RLock()
	defer fake.resourcesMutex.RUnlock()
	argsForCall := fake.resourcesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6, argsForCall.arg7
}

func (fake *FakeBPMConverter) ResourcesReturns(result1 *bpmconverter.Resources, result2 error) {
	fake.resourcesMutex.Lock()
	defer fake.resourcesMutex.Unlock()
	fake.ResourcesStub = nil
	fake.resourcesReturns = struct {
		result1 *bpmconverter.Resources
		result2 error
	}{result1, result2}
}

func (fake *FakeBPMConverter) ResourcesReturnsOnCall(i int, result1 *bpmconverter.Resources, result2 error) {
	fake.resourcesMutex.Lock()
	defer fake.resourcesMutex.Unlock()
	fake.ResourcesStub = nil
	if fake.resourcesReturnsOnCall == nil {
		fake.resourcesReturnsOnCall = make(map[int]struct {
			result1 *bpmconverter.Resources
			result2 error
		})
	}
	fake.resourcesReturnsOnCall[i] = struct {
		result1 *bpmconverter.Resources
		result2 error
	}{result1, result2}
}

func (fake *FakeBPMConverter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.resourcesMutex.RLock()
	defer fake.resourcesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBPMConverter) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ boshdeployment.BPMConverter = new(FakeBPMConverter)
