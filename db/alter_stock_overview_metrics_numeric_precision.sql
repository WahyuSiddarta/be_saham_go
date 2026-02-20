-- Expands stock_overview_metrics numeric precision to reduce overflow risk.
-- Run once on existing databases.

DO $$
DECLARE
    col record;
BEGIN
    FOR col IN
        SELECT column_name, numeric_scale
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'stock_overview_metrics'
          AND data_type = 'numeric'
          AND numeric_precision = 20
          AND numeric_scale IN (2, 6)
    LOOP
        EXECUTE format(
            'ALTER TABLE public.stock_overview_metrics ALTER COLUMN %I TYPE NUMERIC(30,%s);',
            col.column_name,
            col.numeric_scale
        );
    END LOOP;
END $$;
