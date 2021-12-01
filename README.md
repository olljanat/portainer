# Custom fork of Portainer
This fork is based on Portainer 2.9.3 version and it contains following customizations:
* Auto create LDAP users as admin instead of standard user (=> totally skip Portainer internal RBAC)
* Auto create OAuth users as admin instead of standard user (=> totally skip Portainer internal RBAC)
* Disabled all create, modify functionalities (=> make it read-only)

NOTE!!! Only Kubernetes endpoints are really read-only as that that is done with Kubernetes RBAC.
Docker/Docker swarm APIs currently allow all functionalities and buttons are just hided from UI.

You can find docker images from: https://hub.docker.com/r/ollijanatuinen/portainer

# Usage (on Kubernetes)
1. Deploy Portainer-CE 2.9.3 (normal version)
2. Configure endpoints.
3. Configure LDAP / OAuth with auto user create
4. Drop Portainer but do NOT remove its data
5. Create cluster role for using [this](https://github.com/olljanat/portainer/blob/2.9.3-customization/README.md) YAML
6. Deploy custom version of Portainer on way that service account `portainer-sa-clusteradmin` and will use portainer role created above
7. Have fun. All users login with LDAP and/or OAuth should be automatically created and see all resources on read-only mode.
