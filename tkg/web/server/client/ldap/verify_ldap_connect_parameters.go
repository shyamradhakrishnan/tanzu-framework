// Code generated by go-swagger; DO NOT EDIT.

package ldap

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/vmware-tanzu/tanzu-framework/tkg/web/server/models"
)

// NewVerifyLdapConnectParams creates a new VerifyLdapConnectParams object
// with the default values initialized.
func NewVerifyLdapConnectParams() *VerifyLdapConnectParams {
	var ()
	return &VerifyLdapConnectParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewVerifyLdapConnectParamsWithTimeout creates a new VerifyLdapConnectParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewVerifyLdapConnectParamsWithTimeout(timeout time.Duration) *VerifyLdapConnectParams {
	var ()
	return &VerifyLdapConnectParams{

		timeout: timeout,
	}
}

// NewVerifyLdapConnectParamsWithContext creates a new VerifyLdapConnectParams object
// with the default values initialized, and the ability to set a context for a request
func NewVerifyLdapConnectParamsWithContext(ctx context.Context) *VerifyLdapConnectParams {
	var ()
	return &VerifyLdapConnectParams{

		Context: ctx,
	}
}

// NewVerifyLdapConnectParamsWithHTTPClient creates a new VerifyLdapConnectParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewVerifyLdapConnectParamsWithHTTPClient(client *http.Client) *VerifyLdapConnectParams {
	var ()
	return &VerifyLdapConnectParams{
		HTTPClient: client,
	}
}

/*VerifyLdapConnectParams contains all the parameters to send to the API endpoint
for the verify ldap connect operation typically these are written to a http.Request
*/
type VerifyLdapConnectParams struct {

	/*Credentials
	  LDAP configuration

	*/
	Credentials *models.LdapParams

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the verify ldap connect params
func (o *VerifyLdapConnectParams) WithTimeout(timeout time.Duration) *VerifyLdapConnectParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the verify ldap connect params
func (o *VerifyLdapConnectParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the verify ldap connect params
func (o *VerifyLdapConnectParams) WithContext(ctx context.Context) *VerifyLdapConnectParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the verify ldap connect params
func (o *VerifyLdapConnectParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the verify ldap connect params
func (o *VerifyLdapConnectParams) WithHTTPClient(client *http.Client) *VerifyLdapConnectParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the verify ldap connect params
func (o *VerifyLdapConnectParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithCredentials adds the credentials to the verify ldap connect params
func (o *VerifyLdapConnectParams) WithCredentials(credentials *models.LdapParams) *VerifyLdapConnectParams {
	o.SetCredentials(credentials)
	return o
}

// SetCredentials adds the credentials to the verify ldap connect params
func (o *VerifyLdapConnectParams) SetCredentials(credentials *models.LdapParams) {
	o.Credentials = credentials
}

// WriteToRequest writes these params to a swagger request
func (o *VerifyLdapConnectParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Credentials != nil {
		if err := r.SetBodyParam(o.Credentials); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
