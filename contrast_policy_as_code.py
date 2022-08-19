import argparse
import json
import logging
import os
import pathlib
import shutil
import sys
from typing import Any

from contrast_api import contrast_instance_from_json, load_config

logging.basicConfig(level=logging.INFO, format="%(levelname)s: %(message)s")
logger = logging.getLogger(__file__)

try:
    import yaml
except ImportError:
    logger.fatal("pyyaml module is not installed (see README)")
    exit(1)

args_parser = argparse.ArgumentParser(
    description="Export Assess policy defaults and overrides."
)
args_parser.add_argument(
    "-f",
    "--folder",
    help="Output folder.",
    type=pathlib.Path,
    default="output",
)
args_parser.add_argument(
    "-t",
    "--type",
    help="Output type.",
    choices=["JSON", "YAML"],
    type=str.upper,
    default="YAML",
)
args_parser.add_argument(
    "-o",
    "--org-id",
    "--organization-id",
    help="ID of the organization to retrieve Assess policy from.",
    type=str,
    default=os.getenv("CONTRAST_ORG_ID"),
    required="CONTRAST_ORG_ID" not in os.environ.keys(),
)
args = args_parser.parse_args()

# do this first as it'll fail for any missing creds/connection details
config = load_config()
contrast = contrast_instance_from_json(config)

output_path = args.folder
if not output_path.exists():
    output_path.mkdir(parents=True)
elif not output_path.is_dir():
    logger.error(f"{output_path} already exists but is not a directory")
    sys.exit(1)

ORG_UUID = args.org_id
OUTPUT_YAML = args.type == "YAML"
OUTPUT_PATH = str(output_path)
EXTENSION = "yaml" if OUTPUT_YAML else "json"

contrast.test_connection()
if not contrast.test_org_access(ORG_UUID):
    logger.error(f"Unable to access org {ORG_UUID} - check org ID or credentials.")
    sys.exit(1)


def serialize(data: Any) -> str:
    """Return data as a string in the configured format (yaml/json)."""
    if OUTPUT_YAML:
        return yaml.dump(data)
    else:
        return json.dumps(
            data, indent=2, sort_keys=True
        )  # sort_keys should help maintain output stability


# Write out organization defaults for rules (applicable to new applications, and acting as our baseline)
rules_resp = contrast.api_request(f"{ORG_UUID}/rules?expand=app_assess_rules")

rules = rules_resp["rules"]
output = {
    "rules": {
        rule["name"]: {
            "title": rule["title"],
            "dev": rule["enabled_dev"],
            "qa": rule["enabled_qa"],
            "prod": rule["enabled_prod"],
        }
        for rule in rules
    }
}

defaults_file = f"{OUTPUT_PATH}/defaults.{EXTENSION}"
with open(defaults_file, "w") as defaults:
    defaults.write(serialize(output))

logger.info(f"Wrote rule defaults to {defaults_file}")

# Write out information on any rule overrides for specific applications
default_envs = ["dev", "qa", "prod"]
override_envs = ["development", "qa", "production"]  # sigh, naming consistency

overrides = {}

for rule in rules:
    app_overrides = {}
    for i, env in enumerate(default_envs):
        default_enabled = rule[f"enabled_{env}"]

        # find those apps who have adjusted vs the default for this rule
        overridden_apps = rule[f"{override_envs[i]}_breakdown"][
            f'{"off" if default_enabled else "on"}_applications'
        ]

        for app in overridden_apps:
            output_key = f"{app['appId']}/{app['appName']}"
            override = app_overrides.get(output_key, {})
            override[env] = not default_enabled
            app_overrides[output_key] = override

    if len(app_overrides) > 0:
        overrides[rule["name"]] = app_overrides


overrides_path = f"{OUTPUT_PATH}/overrides/"

# empty the overrides folder first in case some rules are no longer overridden
shutil.rmtree(overrides_path, True)
os.mkdir(overrides_path)

for rule in overrides.keys():
    with open(f"{overrides_path}{rule}.{EXTENSION}", "w") as override:
        override.write(serialize({"apps": overrides[rule]}))

logger.info(f"Wrote rule overrides to {overrides_path}")
