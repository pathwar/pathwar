# Pathwar Auth0 Configuration :lock:

Pathwar uses Auth0 to manage authentication to its platform.
Through the OpenID Connect protocol, Auth0 communicates to Pathwar a token that identifies and authenticates the user.

#### Here are the steps to follow to create your Auth0 Tenant

- Go to: https://auth0.com/
- Identify yourself or log in
- From the dashboard, create a "Single Page Application".
- In the application settings, define "Allowed Callbacks URLs", "Allowed Logout URLs" and "Allowed Web Origins" with the front-end address
- Accept "Allow Cross-Origin Authentication".
- Put the front-end URL in "Allow Origins (CORS)".
- Set the token lifetime as desired
- Then create an API on Auth0, allowing the back-end to interact with our Application
- Activate RBAC authorization policies for the API and the addition of Permissions in the access token
- Also activate "Allow Skipping User Content" and "Allow Offline Access".
- Now create "agent" and "admin" permissions for this API
- This way you can give these permissions to the users you want, "admin" allows you to control everything and "agent" is given to the token used by the pathwar agent to authenticate itself.
- Once you have done this you should fill the environments files in `~/web` and default data into ``go/pkg/pwsso/testing.go``