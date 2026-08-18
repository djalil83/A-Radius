from pathlib import Path
import yaml

path = Path('internal/openapi/openapi.yaml')
doc = yaml.safe_load(path.read_text())
assert doc['openapi'] == '3.1.0'
assert '/healthz' in doc['paths']
assert '/api/v1/subscription-profiles' in doc['paths']
assert '/api/v1/subscription-profiles/{id}' in doc['paths']
assert '/api/v1/subscription-profiles/{id}/revisions' in doc['paths']
assert 'bearerAuth' in doc['components']['securitySchemes']
expected = {
    'subscription_profiles.read',
    'subscription_profiles.create',
    'subscription_profiles.update',
    'subscription_profiles.archive',
    'subscription_profiles.read_history',
}
actual = set()
for path_item in doc['paths'].values():
    for operation in path_item.values():
        if isinstance(operation, dict) and 'x-required-permission' in operation:
            actual.add(operation['x-required-permission'])
assert actual == expected, (actual, expected)
for name in ['Profile', 'CreateRequest', 'UpdateRequest', 'Revision', 'ErrorResponse']:
    assert name in doc['components']['schemas']
print(f'OpenAPI structure OK: {path}')
