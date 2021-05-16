CREATE TABLE IF NOT EXISTS canvases (
    id UUID PRIMARY KEY,
    height INT NOT NULL,
    width INT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS rectangles (
    id UUID PRIMARY KEY,
    canvas_id UUID REFERENCES canvases(id),
    x INT NOT NULL,
    y INT NOT NULL,
    height INT NOT NULL,
    width INT NOT NULL,
    filler INT NOT NULL,
    outline INT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS fills (
    id UUID PRIMARY KEY,
    canvas_id UUID REFERENCES canvases(id),
    x INT NOT NULL,
    y INT NOT NULL,
    filler INT NOT NULL,
    created_at TIMESTAMP NOT NULL
);
