package incapsula

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIncapRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceIncapRuleCreate,
		Read:   resourceIncapRuleRead,
		Update: resourceIncapRuleUpdate,
		Delete: resourceIncapRuleDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idSlice := strings.Split(d.Id(), "/")
				if len(idSlice) != 2 || idSlice[0] == "" || idSlice[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected site_id/rule_id", d.Id())
				}

				siteID := idSlice[0]
				d.Set("site_id", siteID)

				ruleID := idSlice[1]
				d.SetId(ruleID)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"site_id": {
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "Rule name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"action": {
				Description: "Rule action. See the detailed descriptions in the API documentation. Possible values: `RULE_ACTION_REDIRECT`, `RULE_ACTION_SIMPLIFIED_REDIRECT`, `RULE_ACTION_REWRITE_URL`, `RULE_ACTION_REWRITE_HEADER`, `RULE_ACTION_REWRITE_COOKIE`, `RULE_ACTION_DELETE_HEADER`, `RULE_ACTION_DELETE_COOKIE`, `RULE_ACTION_RESPONSE_REWRITE_HEADER`, `RULE_ACTION_RESPONSE_DELETE_HEADER`, `RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE`, `RULE_ACTION_FORWARD_TO_DC`, `RULE_ACTION_ALERT`, `RULE_ACTION_BLOCK`, `RULE_ACTION_BLOCK_USER`, `RULE_ACTION_BLOCK_IP`, `RULE_ACTION_RETRY`, `RULE_ACTION_INTRUSIVE_HTML`, `RULE_ACTION_CAPTCHA`, `RULE_ACTION_RATE`, `RULE_ACTION_CUSTOM_ERROR_RESPONSE`",
				Type:        schema.TypeString,
				Required:    true,
			},
			// Optional Arguments
			"filter": {
				Description: "The filter defines the conditions that trigger the rule action. For action `RULE_ACTION_SIMPLIFIED_REDIRECT` filter is not relevant. For other actions, if left empty, the rule is always run.",
				Type:        schema.TypeString,
				Optional:    true,
				DiffSuppressFunc: func(k, remoteState, desiredState string, d *schema.ResourceData) bool {
					return strings.TrimSpace(remoteState) == strings.TrimSpace(desiredState)
				},
			},
			"response_code": {
				Description: "For `RULE_ACTION_REDIRECT` or `RULE_ACTION_SIMPLIFIED_REDIRECT` rule's response code, valid values are `302`, `301`, `303`, `307`, `308`. For `RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE` rule's response code, valid values are all 3-digits numbers. For `RULE_ACTION_CUSTOM_ERROR_RESPONSE`, valid values are `400`, `401`, `402`, `403`, `404`, `405`, `406`, `407`, `408`, `409`, `410`, `411`, `412`, `413`, `414`, `415`, `416`, `417`, `419`, `420`, `422`, `423`, `424`, `500`, `501`, `502`, `503`, `504`, `505`, `507`.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"add_missing": {
				Description: "Add cookie or header if it doesn't exist (Rewrite cookie rule only).",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"rewrite_existing": {
				Description: "Rewrite cookie or header if it exists.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"from": {
				Description: "Pattern to rewrite. For `RULE_ACTION_REWRITE_URL` - Url to rewrite. For `RULE_ACTION_REWRITE_HEADER` and `RULE_ACTION_RESPONSE_REWRITE_HEADER` - Header value to rewrite. For `RULE_ACTION_REWRITE_COOKIE` - Cookie value to rewrite.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"to": {
				Description: "Pattern to change to. `RULE_ACTION_REWRITE_URL` - Url to change to. `RULE_ACTION_REWRITE_HEADER` and `RULE_ACTION_RESPONSE_REWRITE_HEADER` - Header value to change to. `RULE_ACTION_REWRITE_COOKIE` - Cookie value to change to.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"rewrite_name": {
				Description: "Name of cookie or header to rewrite. Applies only for `RULE_ACTION_REWRITE_COOKIE`, `RULE_ACTION_REWRITE_HEADER` and `RULE_ACTION_RESPONSE_REWRITE_HEADER`.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"dc_id": {
				Description: "Data center to forward request to. Applies only for `RULE_ACTION_FORWARD_TO_DC`.",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"port_forwarding_context": {
				Description: "Context for port forwarding. \"Use Port Value\" or \"Use Header Name\". Applies only for `RULE_ACTION_FORWARD_TO_PORT`.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"port_forwarding_value": {
				Description: "Port number or header name for port forwarding. Applies only for `RULE_ACTION_FORWARD_TO_PORT`.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"rate_context": {
				Description: "The context of the rate counter. Possible values `IP` or `Session`. Applies only to rules using `RULE_ACTION_RATE`.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"rate_interval": {
				Description: "The interval in seconds of the rate counter. Possible values is a multiple of `10`; minimum `10` and maximum `300`. Applies only to rules using `RULE_ACTION_RATE`.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"error_type": {
				Description: "The error that triggers the rule. `error.type.all` triggers the rule regardless of the error type. Applies only for `RULE_ACTION_CUSTOM_ERROR_RESPONSE`. Possible values: `error.type.all`, `error.type.connection_timeout`, `error.type.access_denied`, `error.type.parse_req_error`, `error.type.parse_resp_error`, `error.type.connection_failed`, `error.type.deny_and_retry`, `error.type.ssl_failed`, `error.type.deny_and_captcha`, `error.type.2fa_required`, `error.type.no_ssl_config`, `error.type.no_ipv6_config`.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"error_response_format": {
				Description: "The format of the given error response in the error_response_data field. Applies only for `RULE_ACTION_CUSTOM_ERROR_RESPONSE`. Possible values: `json`, `xml`.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"error_response_data": {
				Description: "The response returned when the request matches the filter and is blocked. Applies only for `RULE_ACTION_CUSTOM_ERROR_RESPONSE`.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"multiple_deletions": {
				Description: "Delete multiple header occurrences. Applies only to rules using `RULE_ACTION_DELETE_HEADER` and `RULE_ACTION_RESPONSE_DELETE_HEADER`.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"override_waf_rule": {
				Description: "The setting to override. Possible values: SQL Injection, Remote File Inclusion, Cross Site Scripting, Illegal Resource Access.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"override_waf_action": {
				Description: "The action for the override rule. Possible values: Alert Only, Block Request, Block User, Block IP, Ignore.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"block_duration_type": {
				Description: "Block duration type.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"block_duration": {
				Description: "Value of the fixed block duration.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"block_duration_min": {
				Description: "The lower limit for the randomized block duration.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"block_duration_max": {
				Description: "The upper limit for the randomized block duration.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"enabled": {
				Description: "Enable or disable rule.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"send_notifications": {
				Description: "Send an email notification whenever this rule is triggered",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "true" && v != "false" {
						errs = append(errs, fmt.Errorf("%q must be either 'true' or 'false', got: %s", key, v))
					}
					return
				},
			},
		},
	}
}

func resourceIncapRuleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	rewriteExisting := new(bool)
	action := d.Get("action").(string)

	//#MY-15714 rewrite_existing mustn't be set for other rule types
	if action == "RULE_ACTION_RESPONSE_REWRITE_HEADER" || action == "RULE_ACTION_REWRITE_HEADER" || action == "RULE_ACTION_REWRITE_COOKIE" {
		*rewriteExisting = d.Get("rewrite_existing").(bool)
	} else {
		rewriteExisting = nil
	}

	var sendNotifications *bool
	if v, ok := d.GetOk("send_notifications"); ok {
		valStr := v.(string)
		valBool, err := strconv.ParseBool(valStr)
		if err != nil {
			return fmt.Errorf("Error parsing send_notifications: %s", err)
		}
		sendNotifications = &valBool
	}

	blockDurationDetails := &BlockDurationDetails{
		BlockDurationType: d.Get("block_duration_type").(string),
		BlockDuration:     d.Get("block_duration").(int),
		BlockDurationMin:  d.Get("block_duration_min").(int),
		BlockDurationMax:  d.Get("block_duration_max").(int),
	}

	if blockDurationDetails.BlockDurationType == "" {
		blockDurationDetails = nil
	}

	rule := IncapRule{
		Name:                  d.Get("name").(string),
		Action:                action,
		Filter:                d.Get("filter").(string),
		ResponseCode:          d.Get("response_code").(int),
		AddMissing:            d.Get("add_missing").(bool),
		RewriteExisting:       rewriteExisting,
		From:                  d.Get("from").(string),
		To:                    d.Get("to").(string),
		RewriteName:           d.Get("rewrite_name").(string),
		DCID:                  d.Get("dc_id").(int),
		PortForwardingContext: d.Get("port_forwarding_context").(string),
		PortForwardingValue:   d.Get("port_forwarding_value").(string),
		RateContext:           d.Get("rate_context").(string),
		RateInterval:          d.Get("rate_interval").(int),
		ErrorType:             d.Get("error_type").(string),
		ErrorResponseFormat:   d.Get("error_response_format").(string),
		ErrorResponseData:     d.Get("error_response_data").(string),
		MultipleDeletions:     d.Get("multiple_deletions").(bool),
		OverrideWafRule:       d.Get("override_waf_rule").(string),
		OverrideWafAction:     d.Get("override_waf_action").(string),
		Enabled:               d.Get("enabled").(bool),
		SendNotifications:     sendNotifications,
		BlockDurationDetails:  blockDurationDetails,
	}

	ruleWithID, err := client.AddIncapRule(d.Get("site_id").(string), &rule)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(ruleWithID.RuleID))

	return resourceIncapRuleRead(d, m)
}

func resourceIncapRuleRead(d *schema.ResourceData, m interface{}) error {
	// Implement by reading the SiteResponse for the site
	client := m.(*Client)

	ruleID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	rule, statusCode, err := client.ReadIncapRule(d.Get("site_id").(string), ruleID)

	// If the rule is deleted on the server, blow it out locally and run through the normal TF cycle
	if statusCode == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	// Update all of the properties
	d.Set("name", rule.Name)
	d.Set("action", rule.Action)
	d.Set("filter", rule.Filter)
	d.Set("response_code", rule.ResponseCode)
	d.Set("add_missing", rule.AddMissing)
	d.Set("from", rule.From)
	d.Set("to", rule.To)
	d.Set("rewrite_name", rule.RewriteName)
	d.Set("dc_id", rule.DCID)
	d.Set("port_forwarding_context", rule.PortForwardingContext)
	d.Set("port_forwarding_value", rule.PortForwardingValue)
	d.Set("rate_context", rule.RateContext)
	d.Set("rate_interval", rule.RateInterval)
	d.Set("error_type", rule.ErrorType)
	d.Set("error_response_format", rule.ErrorResponseFormat)
	d.Set("error_response_data", rule.ErrorResponseData)
	d.Set("multiple_deletions", rule.MultipleDeletions)
	d.Set("override_waf_rule", rule.OverrideWafRule)
	d.Set("override_waf_action", rule.OverrideWafAction)
	d.Set("enabled", rule.Enabled)
	if rule.SendNotifications != nil {
		d.Set("send_notifications", strconv.FormatBool(*rule.SendNotifications))
	}
	if rule.BlockDurationDetails != nil {
		d.Set("block_duration_type", rule.BlockDurationDetails.BlockDurationType)
		d.Set("block_duration", rule.BlockDurationDetails.BlockDuration)
		d.Set("block_duration_min", rule.BlockDurationDetails.BlockDurationMin)
		d.Set("block_duration_max", rule.BlockDurationDetails.BlockDurationMax)
	}

	action := d.Get("action").(string)

	if action == "RULE_ACTION_RESPONSE_REWRITE_HEADER" || action == "RULE_ACTION_REWRITE_HEADER" || action == "RULE_ACTION_REWRITE_COOKIE" {
		if rule.RewriteExisting != nil {
			d.Set("rewrite_existing", *rule.RewriteExisting)
		}
	} else {
		//align with schema default to avoid diff when importing resources
		d.Set("rewrite_existing", true)
	}

	return nil
}

func resourceIncapRuleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	rewriteExisting := new(bool)
	action := d.Get("action").(string)

	//#MY-15714 rewrite_existing mustn't be set for other rule types
	if action == "RULE_ACTION_RESPONSE_REWRITE_HEADER" || action == "RULE_ACTION_REWRITE_HEADER" || action == "RULE_ACTION_REWRITE_COOKIE" {
		*rewriteExisting = d.Get("rewrite_existing").(bool)
	} else {
		rewriteExisting = nil
	}

	var sendNotifications *bool
	if v, ok := d.GetOk("send_notifications"); ok {
		valStr := v.(string)
		valBool, err := strconv.ParseBool(valStr)
		if err != nil {
			return fmt.Errorf("Error parsing send_notifications: %s", err)
		}
		sendNotifications = &valBool
	}

	blockDurationDetails := &BlockDurationDetails{
		BlockDurationType: d.Get("block_duration_type").(string),
		BlockDuration:     d.Get("block_duration").(int),
		BlockDurationMin:  d.Get("block_duration_min").(int),
		BlockDurationMax:  d.Get("block_duration_max").(int),
	}

	if blockDurationDetails.BlockDurationType == "" {
		blockDurationDetails = nil
	}

	rule := IncapRule{
		Name:                  d.Get("name").(string),
		Action:                action,
		Filter:                d.Get("filter").(string),
		ResponseCode:          d.Get("response_code").(int),
		AddMissing:            d.Get("add_missing").(bool),
		RewriteExisting:       rewriteExisting,
		From:                  d.Get("from").(string),
		To:                    d.Get("to").(string),
		RewriteName:           d.Get("rewrite_name").(string),
		DCID:                  d.Get("dc_id").(int),
		PortForwardingContext: d.Get("port_forwarding_context").(string),
		PortForwardingValue:   d.Get("port_forwarding_value").(string),
		RateContext:           d.Get("rate_context").(string),
		RateInterval:          d.Get("rate_interval").(int),
		ErrorType:             d.Get("error_type").(string),
		ErrorResponseFormat:   d.Get("error_response_format").(string),
		ErrorResponseData:     d.Get("error_response_data").(string),
		MultipleDeletions:     d.Get("multiple_deletions").(bool),
		OverrideWafRule:       d.Get("override_waf_rule").(string),
		OverrideWafAction:     d.Get("override_waf_action").(string),
		Enabled:               d.Get("enabled").(bool),
		SendNotifications:     sendNotifications,
		BlockDurationDetails:  blockDurationDetails,
	}

	ruleID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	_, err = client.UpdateIncapRule(d.Get("site_id").(string), ruleID, &rule)

	if err != nil {
		return err
	}

	return nil
}

func resourceIncapRuleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	ruleID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteIncapRule(d.Get("site_id").(string), ruleID)
	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	return nil
}
