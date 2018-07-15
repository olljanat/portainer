function LDAPSettingsViewModel(data) {
  this.ReaderDN = data.ReaderDN;
  this.Password = data.Password;
  this.URL = data.URL;
  this.SearchSettings = data.SearchSettings;
}

function LDAPSearchSettings(data) {
  this.BaseDN = data.BaseDN;
  this.UsernameAttribute = data.UsernameAttribute;
  this.Filter = data.Filter;
  this.GroupBaseDN = data.GroupBaseDN;
  this.GroupAttribute = data.GroupAttribute;
  this.GroupFilter = data.GroupFilter;
}
