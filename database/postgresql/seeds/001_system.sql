INSERT INTO apb.system_settings
    (setting_key, setting_value)
VALUES
    ('app.name', 'A-Radius'),
    ('app.architecture', 'APB-COMPLETE'),
    ('database.version', '0.1.0')
ON CONFLICT (setting_key)
DO UPDATE SET
    setting_value = EXCLUDED.setting_value,
    updated_at = NOW();
