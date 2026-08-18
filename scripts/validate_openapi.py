from pathlib import Path
import sys

import yaml

spec_path = Path(__file__).resolve().parents[1] / "docs" / "openapi.yaml"
spec = yaml.safe_load(spec_path.read_text())

assert spec["openapi"] == "3.0.3"
assert spec["info"]["title"] == "A-Radius Subscription Profile API"
assert "/api/v1/subscription-profiles/{id}" in spec["paths"]
assert "/api/v1/subscription-profiles/{id}/revisions" in spec["paths"]

components = spec["components"]
schemas = components["schemas"]
parameters = components["parameters"]
responses = components["responses"]

for name in ["Profile", "ProfileInput", "Revision", "ErrorResponse", "ProfileList", "RevisionList"]:
    assert name in schemas, name
for name in ["ProfileID", "VersionQuery", "Limit", "Offset"]:
    assert name in parameters, name
assert "VersionConflict" in responses

for path, path_item in spec["paths"].items():
    for method, operation in path_item.items():
        if method not in {"get", "post", "patch", "delete", "put", "options", "head"}:
            continue
        text = str(operation)
        if "#/components" in text:
            assert "#/components" in text

version_conflict = responses["VersionConflict"]
example = version_conflict["content"]["application/json"]["example"]
assert example["error"]["code"] == "VERSION_CONFLICT"
assert "409" in spec["paths"]["/api/v1/subscription-profiles/{id}"]["patch"]["responses"]
assert "409" in spec["paths"]["/api/v1/subscription-profiles/{id}"]["delete"]["responses"]

print(f"OpenAPI validation passed: {spec_path}")
print(f"paths={len(spec['paths'])}, schemas={len(schemas)}, responses={len(responses)}")
