# Secret Store

Secret Store is designed to be As Simple As Possible. It provides a _reasonably_ secure API and WebUI for Secret CRUD.

The limitations of this project are:
- Single User with Full Access - No Organisations, Roles, Projects or complex IAM.
- Environment Managed Auth - no OIDC, LDAP, SSO.

## Motivation

[Infiscal](https://infisical.com/) is a great Self-Hostable Secret Manager. 

However the Enterprise features and complexity are out of scope for a simple User/HomeLab Deployment.

## Tech Stack

The Project is designed to be run as a single Docker Image.

Backend uses go and the Frontend uses Static Vue (served from the Backend). The data is stored in a SQLite DB.
