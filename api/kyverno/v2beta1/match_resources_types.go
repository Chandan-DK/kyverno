package v2beta1

import (
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// MatchResources is used to specify resource and admission review request data for
// which a policy rule is applicable.
// +kubebuilder:not:={required:{any,all}}
type MatchResources struct {
	// Any allows specifying resources which will be ORed
	// +optional
	Any kyvernov1.ResourceFilters `json:"any,omitempty"`

	// All allows specifying resources which will be ANDed
	// +optional
	All kyvernov1.ResourceFilters `json:"all,omitempty"`
}

// GetResourceFilters returns all resource filters
func (m *MatchResources) GetResourceFilters() kyvernov1.ResourceFilters {
	var filters kyvernov1.ResourceFilters
	filters = append(filters, m.All...)
	filters = append(filters, m.Any...)
	return filters
}

// GetKinds returns all kinds
func (m *MatchResources) GetKinds() []string {
	var kinds []string
	for _, value := range m.All {
		kinds = append(kinds, value.ResourceDescription.Kinds...)
	}
	for _, value := range m.Any {
		kinds = append(kinds, value.ResourceDescription.Kinds...)
	}
	return kinds
}

// ValidateNoUserInfo verifies that no user info is used
func (m *MatchResources) ValidateNoUserInfo(path *field.Path) (errs field.ErrorList) {
	anyPath := path.Child("any")
	for i, filter := range m.Any {
		errs = append(errs, filter.UserInfo.ValidateNoUserInfo(anyPath.Index(i))...)
	}
	allPath := path.Child("all")
	for i, filter := range m.All {
		errs = append(errs, filter.UserInfo.ValidateNoUserInfo(allPath.Index(i))...)
	}
	return errs
}

// ValidateResourceWithNoUserInfo implements programmatic validation and verifies that no user info is used
func (m *MatchResources) ValidateResourceWithNoUserInfo(path *field.Path, namespaced bool, clusterResources sets.Set[string]) (errs field.ErrorList) {
	if len(m.Any) > 0 && len(m.All) > 0 {
		errs = append(errs, field.Invalid(path, m, "Can't specify any and all together"))
	}
	anyPath := path.Child("any")
	for i, filter := range m.Any {
		errs = append(errs, filter.UserInfo.ValidateNoUserInfo(anyPath.Index(i))...)
		errs = append(errs, filter.ResourceDescription.Validate(anyPath.Index(i), namespaced, clusterResources)...)
	}
	allPath := path.Child("all")
	for i, filter := range m.All {
		errs = append(errs, filter.UserInfo.ValidateNoUserInfo(allPath.Index(i))...)
		errs = append(errs, filter.ResourceDescription.Validate(allPath.Index(i), namespaced, clusterResources)...)
	}
	return errs
}

// Validate implements programmatic validation
func (m *MatchResources) Validate(path *field.Path, namespaced bool, clusterResources sets.Set[string]) (errs field.ErrorList) {
	if len(m.Any) > 0 && len(m.All) > 0 {
		errs = append(errs, field.Invalid(path, m, "Can't specify any and all together"))
	}
	anyPath := path.Child("any")
	for i, filter := range m.Any {
		errs = append(errs, filter.UserInfo.Validate(anyPath.Index(i))...)
		errs = append(errs, filter.ResourceDescription.Validate(anyPath.Index(i), namespaced, clusterResources)...)
	}
	allPath := path.Child("all")
	for i, filter := range m.All {
		errs = append(errs, filter.UserInfo.Validate(allPath.Index(i))...)
		errs = append(errs, filter.ResourceDescription.Validate(allPath.Index(i), namespaced, clusterResources)...)
	}
	return errs
}
