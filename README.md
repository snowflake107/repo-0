# Contrast Assess Policy as Code

Script to help output Contrast Assess rules in JSON or YAML.
Rule defaults and application overrides are both written to disk.

Rule configuration is retrieved _from_ Contrast only, changes to the output contents will not be reflected back in Contrast.

## Requirements
- Python 3.10 (other versions _may_ work but are untested)
- Ability to install Python libraries from `requirements.txt`

## Setup
You can run this script locally with a Python install, or, in a container with the provided `Dockerfile`

### Container use

#### Pre-built
```bash
docker run -it --env-file=contrast.env -v $PWD/output:/usr/src/app/output ghcr.io/contrast-security-oss/assess-policy-as-code:main <...args...>
```

#### Local build
```bash
docker build . --tag contrast-policy-as-code # Build the container
docker run -it --env-file=contrast.env -v $PWD/output:/usr/src/app/output contrast-policy-as-code <...args...> # Run the container
```

### Local use
Use of a virtual environment is encouraged
```bash
python3 -m venv venv # Create the virtual environment
. venv/bin/activate # Activate the virtual environment
pip3 install -r requirements.txt # Install dependencies
. contrast.env # Setup environment
python3 contrast_policy_as_code.py <args> # Run script
```

## Connection and Authentication

The script **requires** the following environment variables to be defined:
- `CONTRAST__API__URL` - the URL to your Contast instance, e.g.: `https://contrast_instance.your_domain.tld/Contrast`
- `CONTRAST__API__API_KEY` - an API key with permission to access that instance
- `CONTRAST__API__AUTH_HEADER` - authorization header for a user with permission to access that instance (base 64 of `username:service_key`)
- `CONTRAST_ORG_ID` - organization ID - may also be passed with the `-o` command line argument

There are also the following optional environment variables:
- `INSECURE_SKIP_CERT_VALIDATION` - set to `true` or `1` to skip TLS certificate validation on network requests
- `HTTP_PROXY` - set to your proxy URL if a proxy is needed to reach Contrast

## Running

Full usage information:

```
usage: contrast_policy_as_code.py [-h] [-f FOLDER] [-t {JSON,YAML}] -o ORG_ID

Export Assess policy defaults and overrides.

options:
  -h, --help            show this help message and exit
  -f FOLDER, --folder FOLDER
                        Output folder.
  -t {JSON,YAML}, --type {JSON,YAML}
                        Output type.
  -o ORG_ID, --org-id ORG_ID, --organization-id ORG_ID
                        ID of the organization to retrieve Assess policy from.
```

## Output

Examples of the output - in both YAML and JSON - can be seen in the [`demo_output`](demo_output) folder.

The top-level `defaults.[json|yaml]` file lists the organization policy for each rule across the 3 environments.

The `overrides` folder contains a file for each rule where an application has overridden the defaults, providing detail on the application and environment overrides in place.

## Development Setup
Various tools enforce code standards, and are run as a pre-commit hook. This must be setup before committing changes with the following commands:
```bash
python3 -m venv venv # setup a virtual environment
. venv/bin/activate # activate the virtual environment
pip3 install -r requirements-dev.txt # install development dependencies (will also include app dependencies)
pre-commit install # setup the pre-commit hook which handles formatting
```
