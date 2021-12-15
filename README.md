# Custom fork of Portainer
Reason to create this customized version of Portainer is that IMO it switched to incorrect track on [portainer/portainer#2137](https://github.com/portainer/portainer/pull/2137) when they switched to configuration where it is **not possible** to control resources created outside of Portainer without admin permissions.

And what is even worse they decided to **not accept** my pull request which would allow user choosing between old and new behavior [portainer/portainer#2424](https://github.com/portainer/portainer/pull/2424)

As result of that I ended up to creating customized version from latest state without that change and it have been available on [here](https://github.com/olljanat/portainer/tree/1.19.1-custom4) which I have been using on Docker environments (both standalone and swarm).


Now years later Portainer added Kubernetes support which I'm interested to use because it's UI is useful for visualization and troubleshooting purposes.

Unfortunately they are decided that Portainer will totally ignore Kubernetes internal RBAC which means that anyone who want to use it are forced to use Portainer as only/primary management tool or alternatively maintain access rights on two places (=> Portainer is master, Kubernetes is slave).

Also [my proposal](https://github.com/portainer/portainer/issues/6100) about sharing kubeconfig with other tools got rejected.


So here we are once again with customized version of Portainer which totally reverses roles on way that Kubernetes RBAC is master and Portainer is forced to slave mode.

## How it is done?
This fork is based on Portainer 2.9.3 version and it contains following customizations:
* ~Auto create LDAP users as admin instead of standard user (=> totally skip Portainer internal RBAC)~ (will be removed)
* Auto create OAuth users as admin instead of standard user (=> totally skip Portainer's internal RBAC)
* Disabled all create and modify functionalities (=> make it read-only)
* Hardcoded `Scopes` value to `id,email,name` on way that it works with Azure AD so we can re-use that field without modifying database.
* Added Azure AD group ID check on way that only users which are part of Azure AD group specified on `Scopes` field are allowed to login.

NOTE!!! Only Kubernetes endpoints are really read-only as that that is done with Kubernetes RBAC.
Docker/Docker swarm APIs currently allow all functionalities and buttons are just hided from UI.

You can find docker images from: https://hub.docker.com/r/ollijanatuinen/portainer

## Usage (on Kubernetes)
1. Deploy Portainer-CE 2.9.3 (normal version)
2. Configure endpoints.
3. Configure LDAP / OAuth with auto user create
4. Drop Portainer but do NOT remove its data
5. Create cluster role for using [this](https://github.com/olljanat/portainer/blob/2.9.3-customization/README.md) YAML
6. Deploy custom version of Portainer on way that service account `portainer-sa-clusteradmin` and will use portainer role created above
7. Have fun. All users login with LDAP and/or OAuth should be automatically created and see all resources on read-only mode.

Alternatively if you are familiar with Portainer API you can deploy customized version directly and do all configuration from API.
