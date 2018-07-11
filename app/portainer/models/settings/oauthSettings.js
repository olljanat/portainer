function OAuthSettingsViewModel(data) {
  this.IdpUrl = data.IdpUrl;
  this.ClientID = data.ClientID;
  this.ClientSecret = data.ClientSecret;
  this.RedirectURL = data.RedirectURL;
  this.Scopes = data.Scopes;
}