package provider

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &wxOneProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &wxOneProvider{
			version: version,
		}
	}
}

// wxOneProvider maps provider schema data to a Go type.
type wxOneProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// wxOneProvider is the provider implementation.
type wxOneProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *wxOneProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "wxone"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *wxOneProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func setCookiesMiddleware(next http.RoundTripper, cookie string) http.RoundTripper {
	return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		// Add the Cookie header to the request
		req.Header.Add("Cookie", cookie)
		return next.RoundTrip(req)
	})
}

type WxOneClients struct {
	httpClient    *resty.Client
	graphqlClient graphql.Client
}

// Configure prepares a WX-One API client for data sources and resources.
func (p *wxOneProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	// Retrieve provider data from configuration
	var config wxOneProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown WX-ONE API Host",
			"The provider cannot create the WX-ONE API client as there is an unknown configuration value for the WX-ONE API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the WX_ONE_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown WX-ONE API Username",
			"The provider cannot create the WX-ONE API client as there is an unknown configuration value for the WX-ONE API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the WX_ONE_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown WX-ONE API Password",
			"The provider cannot create the WX-ONE API client as there is an unknown configuration value for the WX-ONE API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the WX_ONE_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("WX_ONE_HOST")
	username := os.Getenv("WX_ONE_USERNAME")
	password := os.Getenv("WX_ONE_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing WX-ONE API Host",
			"The provider cannot create the WX-ONE API client as there is a missing or empty value for the WX-ONE API host. "+
				"Set the host value in the configuration or use the WX_ONE_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing WX-ONE API Username",
			"The provider cannot create the WX-ONE API client as there is a missing or empty value for the WX-ONE API username. "+
				"Set the username value in the configuration or use the WX_ONE_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing WX-ONE API Password",
			"The provider cannot create the WX-ONE API client as there is a missing or empty value for the WX-ONE API password. "+
				"Set the password value in the configuration or use the WX_ONE_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "wx_one_host", host)
	ctx = tflog.SetField(ctx, "wx_one_username", username)
	ctx = tflog.SetField(ctx, "wx_one_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "wx_one_password")

	tflog.Debug(ctx, "Creating WX-ONE client")

	restClient := resty.New()
	restClient.SetCookieJar(nil)
	restClient.SetTimeout(10 * time.Second)

	challengeResponse, err := restClient.R().SetBody(map[string]string{"username": username}).
		Post(fmt.Sprintf("%s/challenge", host))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create WX-ONE API Client",
			"An unexpected error occurred when creating the WX-ONE API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"WX-ONE Client Error: "+err.Error(),
		)
		return
	}

	var challenge map[string]interface{}
	err = json.Unmarshal(challengeResponse.Body(), &challenge)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create WX-ONE API Client",
			"An unexpected error occurred when creating the WX-ONE API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"WX-ONE Client Error: "+err.Error(),
		)
		return
	}

	initialHashValue := strings.ToUpper(password) + challenge["salt"].(string)
	hash := sha512.Sum512([]byte(initialHashValue))

	// Define the number of rounds
	rounds := int(challenge["rounds"].(float64))

	// Perform the hashing multiple rounds
	for i := 0; i < rounds; i++ {
		// Concatenate your strings
		hashValue := challenge["challenge"].(string) + challenge["date"].(string) + "wizardtales.com" + hex.EncodeToString(hash[:])
		hash = sha512.Sum512([]byte(hashValue))
	}

	// Convert the final hash to a hex string
	hashedPassword := hex.EncodeToString(hash[:])

	loginResponse, err := restClient.R().SetBody(map[string]string{"username": username, "password": hashedPassword}).
		Post(fmt.Sprintf("%s/login", host))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create WX-ONE API Client",
			"An unexpected error occurred when creating the WX-ONE API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"WX-ONE Client Error: "+err.Error(),
		)
		return
	}

	var login map[string]interface{}
	err = json.Unmarshal(loginResponse.Body(), &login)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create WX-ONE API Client",
			"An unexpected error occurred when creating the WX-ONE API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"WX-ONE Client Error: "+err.Error(),
		)
		return
	}

	// Extract cookies from the response
	cookies := loginResponse.Header().Get("Set-Cookie")
	cookieParts := strings.Split(cookies, ";")
	cookie := cookieParts[0]

	tflog.Info(ctx, "#######", map[string]interface{}{"login": login["auth"].(bool)})
	tflog.Info(ctx, "#######", map[string]interface{}{"cookie": cookie})

	httpClient := &http.Client{
		Transport: setCookiesMiddleware(http.DefaultTransport, cookie),
	}

	grqphqlClient := graphql.NewClient(host+"/graphql", httpClient)
	meResp, err := me(ctx, grqphqlClient)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create WX-ONE API Client",
			"An unexpected error occurred when creating the WX-ONE API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"WX-ONE Client Error: "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "#######", map[string]interface{}{"me": meResp.Me.Id})

	wxOneClients := WxOneClients{
		httpClient:    restClient,
		graphqlClient: grqphqlClient,
	}

	// Make the http and grqphql clients available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = &wxOneClients
	resp.ResourceData = &wxOneClients
}

// DataSources defines the data sources implemented in the provider.
func (p *wxOneProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *wxOneProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewKeyResource,
	}
}
