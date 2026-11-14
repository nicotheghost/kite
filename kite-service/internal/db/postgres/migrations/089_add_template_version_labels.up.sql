ALTER TABLE bundle_versions ADD COLUMN version_label TEXT NOT NULL DEFAULT '';
ALTER TABLE bundle_installs ADD COLUMN version_label TEXT NOT NULL DEFAULT '';
