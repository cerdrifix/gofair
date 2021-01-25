// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// OrderRunnerChange order runner change
//
// swagger:model OrderRunnerChange
type OrderRunnerChange struct {

	// full image
	FullImage bool `json:"fullImage,omitempty"`

	// Handicap - the handicap of the runner (selection) (null if not applicable)
	Hc float64 `json:"hc,omitempty"`

	// Selection Id - the id of the runner (selection)
	ID int64 `json:"id,omitempty"`

	// Matched Backs - matched amounts by distinct matched price on the Back side for this runner (selection)
	Mb [][]float64 `json:"mb"`

	// Matched Lays - matched amounts by distinct matched price on the Lay side for this runner (selection)
	Ml [][]float64 `json:"ml"`

	// Strategy Matches - Matched Backs and Matched Lays grouped by strategy reference
	Smc map[string]StrategyMatchChange `json:"smc,omitempty"`

	// Unmatched Orders - orders on this runner (selection) that are not fully matched
	Uo []*Order `json:"uo"`
}

// Validate validates this order runner change
func (m *OrderRunnerChange) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSmc(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUo(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OrderRunnerChange) validateSmc(formats strfmt.Registry) error {
	if swag.IsZero(m.Smc) { // not required
		return nil
	}

	for k := range m.Smc {

		if err := validate.Required("smc"+"."+k, "body", m.Smc[k]); err != nil {
			return err
		}
		if val, ok := m.Smc[k]; ok {
			if err := val.Validate(formats); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *OrderRunnerChange) validateUo(formats strfmt.Registry) error {
	if swag.IsZero(m.Uo) { // not required
		return nil
	}

	for i := 0; i < len(m.Uo); i++ {
		if swag.IsZero(m.Uo[i]) { // not required
			continue
		}

		if m.Uo[i] != nil {
			if err := m.Uo[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("uo" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this order runner change based on the context it is used
func (m *OrderRunnerChange) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateSmc(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateUo(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OrderRunnerChange) contextValidateSmc(ctx context.Context, formats strfmt.Registry) error {

	for k := range m.Smc {

		if val, ok := m.Smc[k]; ok {
			if err := val.ContextValidate(ctx, formats); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *OrderRunnerChange) contextValidateUo(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Uo); i++ {

		if m.Uo[i] != nil {
			if err := m.Uo[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("uo" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *OrderRunnerChange) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OrderRunnerChange) UnmarshalBinary(b []byte) error {
	var res OrderRunnerChange
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}