ALTER TABLE commands ADD COLUMN draft_flow_source JSONB;
ALTER TABLE event_listeners ADD COLUMN draft_flow_source JSONB;
