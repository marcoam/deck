package file

import (
	"io/ioutil"

	"github.com/kong/deck/state"
	yaml "gopkg.in/yaml.v2"
)

// KongStateToFile writes a state object to file with filename.
// It will omit timestamps and IDs while writing.
func KongStateToFile(kongState *state.KongState, filename string) error {
	var file fileStructure

	services, err := kongState.GetAllServices()
	if err != nil {
		return err
	}
	for _, s := range services {
		s := service{Service: s.Service}
		routes, err := kongState.GetAllRoutesByServiceID(*s.ID)
		if err != nil {
			return err
		}
		for _, r := range routes {
			r.Service = nil
			r.ID = nil
			r.CreatedAt = nil
			r.UpdatedAt = nil
			s.Routes = append(s.Routes, &route{Route: r.Route})
		}
		s.ID = nil
		s.CreatedAt = nil
		s.UpdatedAt = nil
		file.Services = append(file.Services, s)
	}

	c, err := yaml.Marshal(file)
	err = ioutil.WriteFile(filename, c, 0600)
	if err != nil {
		return err
	}
	return nil
}