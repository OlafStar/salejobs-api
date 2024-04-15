package api

import (
	"net/http"

	"github.com/olafstar/salejobs-api/internal/types"
)

func newHTTPError(message string) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusInternalServerError, 
		Message:    message,
	}
}

func validateCompany(company types.Company) *HTTPError {
	if company.Name == "" {
		return newHTTPError("Company name cannot be empty")
	}
	if ContainsSQLInjection(company.Name) {
		return newHTTPError("Internal server error")
	}
	if company.Size <= 0 {
		return newHTTPError("Company size must be greater than zero")
	}
	if company.Website == "" {
		return newHTTPError("Company website cannot be empty")
	}
	if ContainsSQLInjection(company.Website) {
		return newHTTPError("Internal server error")
	}
	return nil
}

func validateSalaryRange(salary []types.SalaryRange) *HTTPError {
	for _, s := range salary {
		if s.MinSalary < 0 || s.MaxSalary < 0 {
			return newHTTPError("Salary values must be non-negative")
		}
		if s.MinSalary > s.MaxSalary {
			return newHTTPError("MinSalary cannot be greater than MaxSalary")
		}
		if s.EmploymentType == "" {
			return newHTTPError("Employment type cannot be empty")
		}
		if ContainsSQLInjection(s.EmploymentType) {
			return newHTTPError("Internal server error")
		}
		if s.Currency == "" {
			return newHTTPError("Currency cannot be empty")
		}
		if ContainsSQLInjection(s.Currency) {
			return newHTTPError("Internal server error")
		}
	}
	return nil
}

func validateLocation(location types.Location) *HTTPError {
	if location.Country == "" {
		return newHTTPError("Country cannot be empty")
	}
	if ContainsSQLInjection(location.Country) {
		return newHTTPError("Internal server error")
	}
	if location.City == "" {
		return newHTTPError("City cannot be empty")
	}
	if ContainsSQLInjection(location.City) {
		return newHTTPError("Internal server error")
	}
	if location.Address == "" {
		return newHTTPError("Adress cannot be empty")
	}
	if ContainsSQLInjection(location.Address) {
		return newHTTPError("Internal server error")
	}
	return nil
}

func validateContact(contact types.ContactDetails) *HTTPError {
	if contact.Name == "" {
		return newHTTPError("Contact name cannot be empty")
	}
	if ContainsSQLInjection(contact.Name) {
		return newHTTPError("Internal server error")
	}
	if contact.Email == "" {
		return newHTTPError("Contact email cannot be empty")
	}
	if ContainsSQLInjection(contact.Email) {
		return newHTTPError("Internal server error")
	}
	if contact.Phone == "" {
		return newHTTPError("Contact phone cannot be empty")
	}
	if ContainsSQLInjection(contact.Phone) {
		return newHTTPError("Internal server error")
	}
	return nil
}

func validateApplyType(applyType types.ApplyTypeEnum) *HTTPError {
	if (applyType.Email == "" && applyType.Url == "") || (applyType.Email != "" && applyType.Url != "") {
		return newHTTPError("ApplyType must have either an Email or a URL, but not both")
	}
	if ContainsSQLInjection(applyType.Email) {
		return newHTTPError("Internal server error")
	}
	if ContainsSQLInjection(applyType.Url) {
		return newHTTPError("Internal server error")
	}
	return nil
}

func validateAdvertismentBody(body types.CreateAdvertisementBody) *HTTPError {
	if err := validateCompany(body.Company); err != nil {
		return err
	}
	if body.Title == "" {
		return newHTTPError("Title cannot be empty")
	}
	if ContainsSQLInjection(body.Title) {
		return newHTTPError("Internal server error")
	}

	if body.Experience == "" {
		return newHTTPError("Experience cannot be empty")
	}
	if ContainsSQLInjection(body.Experience) {
		return newHTTPError("Internal server error")
	}
	if body.Skill == "" {
		return newHTTPError("Skill cannot be empty")
	}
	if ContainsSQLInjection(body.Skill) {
		return newHTTPError("Internal server error")
	}
	if err := validateSalaryRange(body.Salary); err != nil {
		return err
	}
	if body.Description == "" {
		return newHTTPError("Description cannot be empty")
	}
	if ContainsSQLInjection(body.Description) {
		return newHTTPError("Internal server error")
	}
	if err := validateLocation(body.Location); err != nil {
		return err
	}
	if body.OperatingMode == "" {
		return newHTTPError("Operating mode cannot be empty")
	}
	if ContainsSQLInjection(body.OperatingMode) {
		return newHTTPError("Internal server error")
	}
	if body.TypeOfWork == "" {
		return newHTTPError("Type of work cannot be empty")
	}
	if ContainsSQLInjection(body.TypeOfWork) {
		return newHTTPError("Internal server error")
	}
	if err := validateApplyType(body.ApplyType); err != nil {
		return err
	}
	if !body.Consent {
		return newHTTPError("Consent must be given")
	}
	if err := validateContact(body.Contact); err != nil {
		return err
	}
	return nil
}