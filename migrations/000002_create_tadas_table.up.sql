-- Create tadas table
CREATE TABLE IF NOT EXISTS
  tadas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    name VARCHAR (255) NOT NULL,
    description TEXT,
    created_by UUID NOT NULL,
    assigned_to UUID,
    status VARCHAR (20) NOT NULL DEFAULT 'in_progress',
    due_at TIMESTAMP
    WITH
      TIME ZONE,
      completed_at TIMESTAMP
    WITH
      TIME ZONE,
      created_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW(),
      updated_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW(),
      deleted_at TIMESTAMP
    WITH
      TIME ZONE,
      CONSTRAINT fk_tadas_created_by FOREIGN KEY (created_by) REFERENCES users (id),
      CONSTRAINT fk_tadas_assigned_to FOREIGN KEY (assigned_to) REFERENCES users (id),
      CONSTRAINT chk_tadas_status CHECK (
        status IN ('in_progress', 'cancelled', 'completed')
      )
  );

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_tadas_created_by ON tadas (created_by);

CREATE INDEX IF NOT EXISTS idx_tadas_assigned_to ON tadas (assigned_to);

CREATE INDEX IF NOT EXISTS idx_tadas_status ON tadas (status);

CREATE INDEX IF NOT EXISTS idx_tadas_created_at ON tadas (created_at);

CREATE INDEX IF NOT EXISTS idx_tadas_deleted_at ON tadas (deleted_at);
