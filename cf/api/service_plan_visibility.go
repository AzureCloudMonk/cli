package api

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry/cli/cf/api/resources"
	"github.com/cloudfoundry/cli/cf/configuration/coreconfig"
	"github.com/cloudfoundry/cli/cf/models"
	"github.com/cloudfoundry/cli/cf/net"
)

//go:generate counterfeiter . ServicePlanVisibilityRepository

type ServicePlanVisibilityRepository interface {
	Create(string, string) error
	List() ([]models.ServicePlanVisibilityFields, error)
	Delete(string) error
	Search(map[string]string) ([]models.ServicePlanVisibilityFields, error)
}

type CloudControllerServicePlanVisibilityRepository struct {
	config  coreconfig.Reader
	gateway net.Gateway
}

func NewCloudControllerServicePlanVisibilityRepository(config coreconfig.Reader, gateway net.Gateway) CloudControllerServicePlanVisibilityRepository {
	return CloudControllerServicePlanVisibilityRepository{
		config:  config,
		gateway: gateway,
	}
}

func (repo CloudControllerServicePlanVisibilityRepository) Create(serviceGUID, orgGUID string) error {
	url := "/v2/service_plan_visibilities"
	data := fmt.Sprintf(`{"service_plan_guid":"%s", "organization_guid":"%s"}`, serviceGUID, orgGUID)
	return repo.gateway.CreateResource(repo.config.ApiEndpoint(), url, strings.NewReader(data))
}

func (repo CloudControllerServicePlanVisibilityRepository) List() (visibilities []models.ServicePlanVisibilityFields, err error) {
	err = repo.gateway.ListPaginatedResources(
		repo.config.ApiEndpoint(),
		"/v2/service_plan_visibilities",
		resources.ServicePlanVisibilityResource{},
		func(resource interface{}) bool {
			if spv, ok := resource.(resources.ServicePlanVisibilityResource); ok {
				visibilities = append(visibilities, spv.ToFields())
			}
			return true
		})
	return
}

func (repo CloudControllerServicePlanVisibilityRepository) Delete(servicePlanGUID string) error {
	path := fmt.Sprintf("/v2/service_plan_visibilities/%s", servicePlanGUID)
	return repo.gateway.DeleteResourceSynchronously(repo.config.ApiEndpoint(), path)
}

func (repo CloudControllerServicePlanVisibilityRepository) Search(queryParams map[string]string) ([]models.ServicePlanVisibilityFields, error) {
	var visibilities []models.ServicePlanVisibilityFields
	err := repo.gateway.ListPaginatedResources(
		repo.config.ApiEndpoint(),
		combineQueryParametersWithUri("/v2/service_plan_visibilities", queryParams),
		resources.ServicePlanVisibilityResource{},
		func(resource interface{}) bool {
			if sp, ok := resource.(resources.ServicePlanVisibilityResource); ok {
				visibilities = append(visibilities, sp.ToFields())
			}
			return true
		})
	return visibilities, err
}
