package injector

import (
	"reflect"
	"testing"

	"github.com/linkerd/linkerd2/controller/proxy-injector/fake"
	k8sPkg "github.com/linkerd/linkerd2/pkg/k8s"
	"k8s.io/api/core/v1"
)

func TestPatch(t *testing.T) {
	fixture := fake.NewFactory()

	trustAnchors, err := fixture.Volume("inject-trust-anchors-volume-spec.yaml")
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	secrets, err := fixture.Volume("inject-linkerd-secrets-volume-spec.yaml")
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	sidecar, err := fixture.Container("inject-sidecar-container-spec.yaml")
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	init, err := fixture.Container("inject-init-container-spec.yaml")
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	var (
		controllerNamespace = "linkerd"
		createdBy           = "linkerd/cli v18.8.4"
	)

	actual := NewPatch()
	actual.addContainer(sidecar)
	actual.addInitContainerRoot()
	actual.addInitContainer(init)
	actual.addVolumeRoot()
	actual.addVolume(trustAnchors)
	actual.addVolume(secrets)
	actual.addPodLabels(map[string]string{
		k8sPkg.ControllerNSLabel: controllerNamespace,
	})
	actual.addDeploymentLabels(map[string]string{
		k8sPkg.ControllerNSLabel: controllerNamespace,
	})
	actual.addPodAnnotations(map[string]string{
		k8sPkg.CreatedByAnnotation: createdBy,
	})

	expected := NewPatch()
	expected.patchOps = []*patchOp{
		{Op: "add", Path: patchPathContainer, Value: sidecar},
		{Op: "add", Path: patchPathInitContainerRoot, Value: []*v1.Container{}},
		{Op: "add", Path: patchPathInitContainer, Value: init},
		{Op: "add", Path: patchPathVolumeRoot, Value: []*v1.Volume{}},
		{Op: "add", Path: patchPathVolume, Value: trustAnchors},
		{Op: "add", Path: patchPathVolume, Value: secrets},
		{Op: "add", Path: patchPathPodLabels, Value: map[string]string{
			k8sPkg.ControllerNSLabel: controllerNamespace,
		}},
		{Op: "add", Path: patchPathDeploymentLabels, Value: map[string]string{
			k8sPkg.ControllerNSLabel: controllerNamespace,
		}},
		{Op: "add", Path: patchPathPodAnnotations, Value: map[string]string{k8sPkg.CreatedByAnnotation: createdBy}},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Content mismatch\nExpected: %s\nActual: %s", expected, actual)
	}
}
